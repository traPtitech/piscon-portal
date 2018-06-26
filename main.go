package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	shellwords "github.com/mattn/go-shellwords"
	"golang.org/x/crypto/acme/autocert"
)

var (
	checkTask  chan struct{}
	sendWorker chan *Task
	db         *gorm.DB
)

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
	Name    string    `gorm:"unique" json:"name"`
	Results []*Result `json:"results"`
}

type Result struct {
	ID        int        `gorm:"AUTO_INCREMENT" json:"id"`
	TeamID    uint       `json:"team_id"`
	TaskID    uint       `json:"task_id"`
	Pass      bool       `json:"pass"`
	Score     int64      `json:"score"`
	Suceess   int64      `json:"suceess"`
	Fail      int64      `json:"fail"`
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
	CmdStr string `json:"cmd_str"`
	IP     string `json:"ip"`
	State  string `json:"state"`
	TeamID uint   `json:"team_id"`
	Team   Team   `json:"team"`
}

func main() {
	sendWorker = make(chan *Task, 10)
	checkTask = make(chan struct{})

	go benchmarkWorker()

	_db, err := gorm.Open("sqlite3", "./db.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = _db

	db.AutoMigrate(&Message{}, &Result{}, &Team{}, &Task{})
	team := Team{Name: "to-hutohu"}
	if db.NewRecord(team) {
		db.Create(&team)
	}

	tasks := []*Task{}
	db.Not("state = 'done'").Find(&tasks)
	for _, t := range tasks {
		go func(task *Task) {
			sendWorker <- task
		}(t)
	}

	e := echo.New()
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.Pre(middleware.HTTPSNonWWWRedirect())

	signingKey := os.Getenv("JWTKey")
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(signingKey),
	}
	e.Use(middleware.JWTWithConfig(config))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*jwtCustomClaims)
			if claims.Name != "traP-showcase" {
				return c.NoContent(http.StatusForbidden)
			}
			return next(c)
		}
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/benchmark/run", queBenchmark)
	e.GET("/benchmark/queue", getBenchmarkQueue)

	api := e.Group("/api")
	api.GET("/results", getAllResults)

	e.StartAutoTLS(":443")
	fmt.Println("end")
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
	requestBody := &struct {
		IP string `json:"ip"`
	}{}
	c.Bind(requestBody)

	teamID := 1

	task := &Task{}

	//db.Where("team_id = ?", teamID).Not("state = 'done'").First(task)
	fmt.Println(teamID)
	if task.CmdStr != "" {
		return c.JSON(http.StatusNotAcceptable, map[string]string{"massage": "すでに登録されています"})
	}

	cmdStr := fmt.Sprintf("/home/benchmarker/bin/benchmarker -t http://%s -u /home/benchmarker/private-isu/benchmarker/userdata", requestBody.IP)
	t := &Task{
		CmdStr: cmdStr,
		IP:     requestBody.IP,
		State:  "waiting",
		TeamID: 1,
	}
	db.Create(t)

	go func() {
		sendWorker <- t
	}()

	return c.NoContent(http.StatusNoContent)
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
			panic(err)
		}
		result := &Result{
			TeamID:  task.TeamID,
			TaskID:  task.ID,
			Pass:    data.Pass,
			Score:   data.Score,
			Suceess: data.Suceess,
			Fail:    data.Fail,
		}

		for _, message := range data.Messages {
			result.Messages = append(result.Messages, &Message{Text: message})
		}
		db.Create(result)

		task.State = "done"
		db.Save(task)

	}
}
