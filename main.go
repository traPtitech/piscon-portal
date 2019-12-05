package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/nagatea/piscon-portal/conoha"
	"golang.org/x/crypto/acme/autocert"
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

type Output struct {
	Pass     bool     `json:"pass"`
	Score    int64    `json:"score"`
	Campaign int64    `json:"campaign`
	Language string   `json:"language`
	Messages []string `json:"messages"`
}

type Message struct {
	gorm.Model
	Text string
}

type Team struct {
	gorm.Model
	Name     string      `gorm:"unique size:50" json:"name"`
	Instance []*Instance `json:"instance"`
	Results  []*Result   `json:"results"`
}

type User struct {
	gorm.Model
	Name       string `gorm:"unique size:50" json:"name"`
	ScreenName string `json:"screen_name"`
	TeamID     uint   `json:"team_id"`
}

type Instance struct {
	gorm.Model
	TeamID           uint   `json:"team_id"`
	GlobalIPAddress  string `json:"global_ip_address"`
	PrivateIPAddress string `json:"private_ip_address"`
	Password         string `json:"password"`
}

type Result struct {
	gorm.Model
	ID        int        `gorm:"AUTO_INCREMENT" json:"id"`
	TeamID    uint       `json:"team_id"`
	TaskID    uint       `json:"task_id"`
	Pass      bool       `json:"pass"`
	Score     int64      `json:"score"`
	Campaign  int64      `json:"campaign`
	Betterize string     `json:"betterize"`
	Messages  []*Message `json:"messages"`
	CreatedAt time.Time  `json:"created_at"`
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

	_db, err := gorm.Open("mysql", "isucon@/isucon?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	//_db.LogMode(true)
	defer _db.Close()
	db = _db

	db.AutoMigrate(&Task{}, &Message{}, &Result{}, &Instance{}, &Team{}, &User{}, &Question{})

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

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api := e.Group("/api")
	api.GET("/results", getAllResults)
	api.GET("/benchmark/queue", getBenchmarkQueue)
	api.GET("/newer", getNewer)
	api.GET("/questions", getQuestions)
	// api.POST("/instancelog", postInstanceLog)

	apiWithAuth := e.Group("/api", middlewareAuthUser)
	apiWithAuth.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	apiWithAuth.POST("/team", createTeam)
	apiWithAuth.POST("/user", createUser)
	// TODO: ユーザー名で認証してないので修正する必要がある
	apiWithAuth.GET("/team/:id", getTeam)
	apiWithAuth.GET("/user/:name", getUser)
	apiWithAuth.POST("/benchmark/:name/:id", queBenchmark)
	apiWithAuth.GET("/admin/team", getAllTeam)

	apiWithAuth.POST("/questions", postQuestions)
	apiWithAuth.PUT("/questions/:id", putQuestions)
	apiWithAuth.DELETE("/questions/:id", deleteQuestions)

	switch env {
	case "prod":
		e.StartAutoTLS(":443")
	default:
		e.Use(middleware.CORS())
		e.Start(":4000")
	}
	fmt.Println("end")
}

// traPかどうかの認証
// TODO: Fix ユーザーネーム認証
func middlewareAuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.NoContent(http.StatusForbidden)
		}
		req, _ := http.NewRequest("GET", "https://q.trap.jp/api/1.0/users/me", nil)
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		res, _ := client.Do(req)
		if res.StatusCode != 200 {
			return c.NoContent(http.StatusForbidden)
		}
		return next(c)
	}
}

func getNewer(c echo.Context) error {
	teams := []Team{}
	db.Raw("SELECT * FROM results AS PI LEFT JOIN teams ON PI.team_id = teams.id WHERE PI.id =( SELECT po.id FROM results AS po LEFT JOIN teams ON po.team_id = teams.id WHERE pass = 1 AND PI.team_id = po.team_id AND score > 0 ORDER BY po.score DESC LIMIT 1 ) AND (PI.created_at > (CURRENT_TIME() - INTERVAL 1 day))").Scan(&teams)
	return c.JSON(http.StatusOK, teams)
}

func getTeam(c echo.Context) error {
	id := c.Param("id")
	team := Team{}
	db.Where("id = ?", id).Find(&team)

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, Response{false, "登録されていません"})
	}

	db.Where("team_id = ?", &team.ID).Preload("Messages").Find(&team.Results)
	db.Model(&team).Related(&team.Instance)
	return c.JSON(http.StatusOK, team)
}

func getUser(c echo.Context) error {
	name := c.Param("name")
	user := User{}
	db.Where("name = ?", name).Find(&user)

	if user.Name == "" {
		return c.JSON(http.StatusNotFound, Response{false, "登録されていません"})
	}

	return c.JSON(http.StatusOK, user)
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

func createUser(c echo.Context) error {
	user := &User{}
	c.Bind(user)

	u := &User{}
	db.Where("name = ?", user.Name).Find(u)

	if u.Name != "" {
		return c.JSON(http.StatusNotFound, Response{false, "登録されています"})
	}

	db.Create(user)
	return c.JSON(http.StatusCreated, user)
}

func createTeam(c echo.Context) error {
	requestBody := &struct {
		Name             string `json:"name"`
		GlobalIPAddress  string `json:"global_ip_address"`
		PrivateIPAddress string `json:"private_ip_address"`
	}{}

	c.Bind(requestBody)

	if requestBody.Name == "" {
		return c.JSON(http.StatusBadRequest, Response{false, "リクエストボディの要素が足りません"})
	}

	t := &Team{}
	db.Where("name = ?", requestBody.Name).Find(t)

	if t.Name != "" {
		return c.JSON(http.StatusNotFound, Response{false, "登録されています"})
	}
	pass := genPassword()

	instance := Instance{
		GlobalIPAddress:  requestBody.GlobalIPAddress,
		PrivateIPAddress: requestBody.PrivateIPAddress,
		Password:         pass,
	}

	team := &Team{
		Name:     requestBody.Name,
		Instance: make([]*Instance, 3),
	}
	team.Instance[0] = &instance

	db.Create(team)
	return c.JSON(http.StatusCreated, team)
}

func getAllResults(c echo.Context) error {
	teams := []*Team{}
	db.Find(&teams)
	for _, team := range teams {
		db.Where("team_id = ?", &team.ID).Preload("Messages").Find(&team.Results)
	}
	return c.JSON(http.StatusOK, teams)
}

func queBenchmark(c echo.Context) error {
	id := c.Param("id")
	name := c.Param("name")

	req := struct {
		Betterize string `json:"betterize"`
	}{}

	c.Bind(&req)

	team := &Team{}
	db.Where("name = ?", name).Find(team)

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, Response{false, "登録されていません"})
	}

	db.Model(team).Related(&team.Instance)

	if team.Instance[0].PrivateIPAddress == "" {
		return c.JSON(http.StatusBadRequest, Response{false, "インスタンスが存在しません"})
	}

	ip := team.Instance[0].PrivateIPAddress

	if id == "2" {
		ip = team.Instance[1].PrivateIPAddress
	} else if id == "3" {
		ip = team.Instance[2].PrivateIPAddress
	}

	task := &Task{}

	db.Where("team_id = ?", team.ID).Not("state = 'done'").First(task)
	if task.CmdStr != "" {
		return c.JSON(http.StatusNotAcceptable, Response{false, "すでに登録されています"})
	}

	cmdStr := fmt.Sprintf("/home/xecua/go/src/github.com/isucon/isucon9-qualify/bin/benchmarker "+
		"-target-url %s", ip)
	t := &Task{
		CmdStr:    cmdStr,
		IP:        ip,
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

		res, err := exec.Command(command[0], command[1:]...).Output()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("end benchmark")

		fmt.Println(string(res))
		data := &Output{}
		err = json.Unmarshal(res, data)
		if err != nil {
			result := &Result{
				TeamID:    task.TeamID,
				TaskID:    task.ID,
				Pass:      false,
				Score:     0,
				Campaign:  0,
				Betterize: task.Betterize,
				Messages:  []*Message{{Text: err.Error()}},
			}
			db.Create(result)

			task.State = "done"
			db.Save(task)
			continue
		}
		messages := make([]*Message, len(data.Messages))
		for i, text := range data.Messages {
			messages[i] = &Message{Text: text}
		}

		result := &Result{
			TeamID:    task.TeamID,
			TaskID:    task.ID,
			Pass:      data.Pass,
			Score:     data.Score,
			Campaign:  data.Campaign,
			Betterize: task.Betterize,
			Messages:  messages,
		}

		db.Create(result)

		task.State = "done"
		db.Save(task)
	}
}
