package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/rackspace/gophercloud"
	"github.com/tohutohu/isucon-portal-go/conoha"
	"golang.org/x/crypto/acme/autocert"
	"github.com/joho/godotenv"
)

var (
	checkTask  chan struct{}
	sendWorker chan *Task
	db         *gorm.DB
	client     *conoha.ConohaClient
)

type Response struct {
	Suceess bool   `json:"suceess"`
	Message string `json:"message"`
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

type Output struct {
	Pass     bool     `json:"pass"`
	Score    int64    `json:"score"`
	Suceess  int64    `json:"success"`
	Fail     int64    `json:"fail"`
	Messages []string `json:"messages"`
}

type Team struct {
	gorm.Model
	Name     string    `gorm:"unique size:50" json:"name"`
	Instance Instance  `json:"instance"`
	Results  []*Result `json:"results"`
}

type Instance struct {
	gorm.Model
	TeamID       uint   `json:"team_id"`
	InstanceName string `json:"instance_id"`
	IPAddress    string `json:"ip_address"`
	Password     string `json:"password"`
}

type Result struct {
	ID        int        `gorm:"AUTO_INCREMENT" json:"id"`
	TeamID    uint       `json:"team_id"`
	TaskID    uint       `json:"task_id"`
	Pass      bool       `json:"pass"`
	Score     int64      `json:"score"`
	Suceess   int64      `json:"suceess"`
	Fail      int64      `json:"fail"`
	Betterize string     `json:"betterize"`
	Messages  []*Message `json:"messages"`
	CreatedAt time.Time  `json:"created_at"`
}

type Message struct {
	ID       uint   `json:"id"`
	ResultID int    `json:"result_id"`
	Text     string `json:"text"`
}

type Task struct {
	gorm.Model
	CmdStr    string `json:"cmd_str"`
	IP        string `json:"ip"`
	State     string `json:"state"`
	Betterize string `json:"betterize"`
	TeamID    uint   `json:"team_id"`
	Team      Team   `json:"team"`
}

type Question struct {
	gorm.Model
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func main() {
	sendWorker = make(chan *Task, 10)
	checkTask = make(chan struct{})

	go benchmarkWorker()

	err := godotenv.Load()
	if err != nil {
			fmt.Println("Error loading .env file")
	}

	_db, err := gorm.Open("mysql", "root@/isucon?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	_db.LogMode(true)
	defer _db.Close()
	db = _db

	db.AutoMigrate(&Message{}, &Task{}, &Result{}, &Instance{}, &Team{}, &Question{})

	tasks := []*Task{}
	db.Not("state = 'done'").Find(&tasks)
	for _, t := range tasks {
		go func(task *Task) {
			sendWorker <- task
		}(t)
	}

	e := echo.New()
	env := os.Getenv("ENV")
	if env == "prod" {
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		e.Pre(middleware.HTTPSNonWWWRedirect())
	}

	// signingKey := os.Getenv("JWTKey")
	// config := middleware.JWTConfig{
	// 	Claims:     &jwtCustomClaims{},
	// 	SigningKey: []byte(signingKey),
	// }
	// e.Use(middleware.JWTWithConfig(config))

	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		user := c.Get("user").(*jwt.Token)
	// 		claims := user.Claims.(*jwtCustomClaims)
	// 		if claims.Name != "traP-showcase" {
	// 			return c.NoContent(http.StatusForbidden)
	// 		}
	// 		return next(c)
	// 	}
	// })

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := e.Group("/api")
	api.GET("/results", getAllResults)
	api.POST("/team", updateTeam)
	api.GET("/team/:id", getTeam)
	api.GET("/admin/team", getAllTeam)
	api.POST("/benchmark/:id", queBenchmark)
	api.GET("/benchmark/queue", getBenchmarkQueue)
	api.GET("/newer", getNewer)
	api.GET("/questions", getQuestions)
	api.POST("/questions", postQuestions)
	api.PUT("/questions/:id", putQuestions)
	api.DELETE("/questions/:id", deleteQuestions)

	switch env {
	case "prod":
		e.StartAutoTLS(":443")
	default:
		e.Start(":4000")
	}
	fmt.Println("end")
}

func getNewer(c echo.Context) error {
	teams := []Team{}
	db.Raw("SELECT * FROM results AS PI LEFT JOIN teams ON PI.team_id = teams.id WHERE PI.id =( SELECT po.id FROM results AS po LEFT JOIN teams ON po.team_id = teams.id WHERE pass = 1 AND PI.team_id = po.team_id AND score > 0 ORDER BY po.score DESC LIMIT 1 ) AND (PI.created_at > (CURRENT_TIME() - INTERVAL 1 day))").Scan(&teams)
	return c.JSON(http.StatusOK, teams)
}

func getTeam(c echo.Context) error {
	id := c.Param("id")
	team := Team{}
	db.Where("name = ?", id).Find(&team)

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, Response{false, "登録されていません"})
	}
	db.Model(&team).Related(&team.Results)
	for _, result := range team.Results {
		db.Model(result).Related(&result.Messages)
	}
	db.Model(&team).Related(&team.Instance)
	return c.JSON(http.StatusOK, team)
}

func getAllTeam(c echo.Context) error {
	teams := []*Team{}
	db.Find(&teams)
	for _, team := range teams {
		db.Model(team).Related(&team.Results)
		db.Model(team).Related(&team.Instance)
		for _, result := range team.Results {
			db.Model(result).Related(&result.Messages)
		}
	}
	return c.JSON(http.StatusOK, teams)
}

func getQuestions(c echo.Context) error {
	questions := []*Question{}
	db.Find(&questions)
	// TODO: 改行対応

	return c.JSON(http.StatusOK, questions)
}

func postQuestions(c echo.Context) error {
	req := struct {
		Question string `json:"question"`
	}{}
	c.Bind(&req)
	question := &Question{
		Question: req.Question,
	}
	db.Create(question)
	return c.NoContent(http.StatusCreated)
}

func putQuestions(c echo.Context) error {
	id := c.Param("id")
	req := struct {
		Answer string `json:"answer"`
	}{}
	c.Bind(&req)
	question := &Question{
		Answer: req.Answer,
	}
	db.Model(question).Where("id = ?", id).Update(&question)

	return c.JSON(http.StatusOK, question)
}

func deleteQuestions(c echo.Context) error {
	id := c.Param("id")
	question := &Question{}
	db.Model(question).Where("id = ?", id).Delete(&question)

	return c.JSON(http.StatusOK, question)
}

func genPassword() string {
	pass := ""
	gen := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	for i := 0; i < 12; i++ {
		pass += string(gen[rand.Intn(len(gen))])
	}
	return pass
}

func updateTeam(c echo.Context) error {
	requestBody := &struct {
		Name string `json:"name"`
	}{}

	c.Bind(requestBody)

	if requestBody.Name == "" {
		return c.JSON(http.StatusBadRequest, Response{false, "Nameが空です"})
	}

	t := &Team{}
	db.Where("name = ?", requestBody.Name).Find(t)
	db.Model(t).Related(&t.Instance)

	if t.Name != "" && t.Instance.IPAddress != "" {
		return c.JSON(http.StatusNotFound, Response{false, "登録されています"})
	}
	pass := genPassword()
	client = conoha.New(gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		Username:         os.Getenv("OS_USERNAME"),
		TenantName:       os.Getenv("OS_TENANT_NAME"),
		Password:         os.Getenv("OS_PASSWORD"),
	})

	err := client.MakeInstance(requestBody.Name, pass)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, Response{false, "インスタンスの作成に失敗しました\n@to-hutohuに連絡してください"})
	}
	s, err := client.GetInstanceInfo(requestBody.Name)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, Response{false, "インスタンスの作成に失敗しました\n@to-hutohuに連絡してください"})
	}

	instance := Instance{
		InstanceName: s.ID,
		IPAddress:    strings.Replace(s.Name, "-", ".", -1),
		Password:     pass,
	}

	if t.Name != "" {
		t.Instance = instance
		db.Create(instance)
		db.Save(t)
		return c.JSON(http.StatusCreated, Response{true, "インスタンスを作成しました"})
	}

	team := &Team{
		Name:     requestBody.Name,
		Instance: instance,
	}
	db.Create(team)
	return c.JSON(http.StatusCreated, Response{true, "登録しました"})
}

func getAllResults(c echo.Context) error {
	teams := []*Team{}
	db.Find(&teams)
	for _, team := range teams {
		db.Model(team).Related(&team.Results)
		for _, result := range team.Results {
			db.Model(result).Related(&result.Messages)
		}
	}
	return c.JSON(http.StatusOK, teams)
}

func queBenchmark(c echo.Context) error {
	id := c.Param("id")

	req := struct {
		Betterize string `json:"betterize"`
	}{}

	c.Bind(&req)

	team := &Team{}
	db.Where("name = ?", id).Find(team)

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, Response{false, "登録されていません"})
	}

	db.Model(team).Related(&team.Instance)

	if team.Instance.IPAddress == "" {
		return c.JSON(http.StatusBadRequest, Response{false, "インスタンスが存在しません"})
	}

	task := &Task{}

	db.Where("team_id = ?", team.ID).Not("state = 'done'").First(task)
	if task.CmdStr != "" {
		return c.JSON(http.StatusNotAcceptable, Response{false, "すでに登録されています"})
	}

	cmdStr := fmt.Sprintf("/home/benchmarker/bin/benchmarker -t http://%s -u /home/benchmarker/private-isu/benchmarker/userdata", team.Instance.IPAddress)
	t := &Task{
		CmdStr:    cmdStr,
		IP:        team.Instance.IPAddress,
		State:     "waiting",
		TeamID:    team.ID,
		Betterize: req.Betterize,
	}
	db.Create(t)

	go func() {
		sendWorker <- t
	}()

	return c.JSON(http.StatusCreated, Response{true, "キューに追加しました"})
}

func getBenchmarkQueue(c echo.Context) error {
	tasks := getTaskQueInfo()
	for _, task := range tasks {
		db.Model(task).Related(&task.Team)
	}
	return c.JSON(http.StatusOK, tasks)
}

func getTaskQueInfo() []*Task {
	tasks := []*Task{}
	db.Table("tasks").Joins("LEFT JOIN teams ON `teams`.id = `tasks`.team_id").Not("state = 'done'").Find(&tasks)
	return tasks
}

func benchmarkWorker() {
	for {
		task := <-sendWorker
		fmt.Println("recieve task")
		task.State = "benchmark"
		db.Save(task)

		command, _ := shellwords.Parse(task.CmdStr)

		output, _ := exec.Command(command[0], command[1:]...).CombinedOutput()
		fmt.Println("end benchmark")

		fmt.Println(string(output))
		data := &Output{}
		err := json.Unmarshal(output, data)
		if err != nil {
			result := &Result{
				TeamID:    task.TeamID,
				TaskID:    task.ID,
				Pass:      false,
				Score:     0,
				Suceess:   0,
				Fail:      0,
				Betterize: task.Betterize,
				Messages:  []*Message{&Message{Text: err.Error()}},
			}
			db.Create(result)

			task.State = "done"
			db.Save(task)
			continue
		}
		result := &Result{
			TeamID:    task.TeamID,
			TaskID:    task.ID,
			Pass:      data.Pass,
			Score:     data.Score,
			Suceess:   data.Suceess,
			Fail:      data.Fail,
			Betterize: task.Betterize,
		}

		for _, message := range data.Messages {
			result.Messages = append(result.Messages, &Message{Text: message})
		}
		db.Create(result)

		task.State = "done"
		db.Save(task)
	}
}
