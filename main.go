package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
)

var (
	queingCh     chan *Task
	taskQue      []*Task
	taskQueMutex *sync.Mutex
	endTask      chan struct{}
	sendWoker    chan *Task
	db           *gorm.DB
)

type Output struct {
	Pass     bool     `json:"pass"`
	Score    int64    `json:"score"`
	Suceess  int64    `json:"success"`
	Fail     int64    `json:"fail"`
	Messages []string `json:"messages"`
}

type Team struct {
	gorm.Model
	Name    string   `json:"name"`
	Results []Result `json:"results"`
}

type Result struct {
	ID        int       `gorm:"AUTO_INCREMENT" json:"id"`
	TeamID    uint      `json:"team_id"`
	Pass      bool      `json:"pass"`
	Score     int64     `json:"score"`
	Suceess   int64     `json:"suceess"`
	Fail      int64     `json:"fail"`
	Messages  []Message `gorm:"auto_preload" json:"messages"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID       uint   `json:"id"`
	ResultID int    `json:"result_id"`
	Text     string `json:"text"`
}

type Task struct {
	Cmd  *exec.Cmd `json:"-"`
	IP   string    `json:"ip"`
	Team Team      `json:"team"`
}

func main() {
	queingCh = make(chan *Task)
	taskQue = make([]*Task, 0)
	taskQueMutex = &sync.Mutex{}

	_db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = _db

	e := echo.New()

	e.POST("/benchmark/run", queBenchmark)
}

func queBenchmark(c echo.Context) error {
	requestBody := &struct {
		IP string
	}{}
	c.Bind(requestBody)

	go func() {
		cmd := exec.Command(fmt.Sprintf("/home/benchmarker/bin/benchmarker -t %s -d /home/benchmarker/private-isu/benchmarker/userdata", requestBody.IP))
		t := &Task{
			Cmd:  cmd,
			IP:   requestBody.IP,
			Team: Team{},
		}
		queingCh <- t
	}()

	return c.NoContent(http.StatusNoContent)
}

func getTaskQueInfo() []Team {
	taskQueMutex.Lock()
	teams := make([]Team, 0)
	for _, t := range taskQue {
		teams = append(teams, t.Team)
	}
	taskQueMutex.Unlock()
	return teams
}

func taskManager() {
	for {
		select {
		case task := <-queingCh:
			taskQueMutex.Lock()
			taskQue = append(taskQue, task)
			taskQueMutex.Unlock()

		case <-endTask:
			taskQueMutex.Lock()
			if len(taskQue) > 0 {
				sendWoker <- taskQue[0]
				taskQue = taskQue[1:]

			}
			taskQueMutex.Unlock()
		}
	}
}

func benchmarkWorker() {
	for {
		task := <-sendWoker
		err := task.Cmd.Run()
		if err != nil {
			panic(err)
		}

		output, err := task.Cmd.Output()
		if err != nil {
			panic(err)
		}

		data := &Output{}
		err = json.Unmarshal(output, data)
		if err != nil {
			panic(err)
		}
		result := &Result{
			TeamID:  task.Team.ID,
			Pass:    data.Pass,
			Score:   data.Score,
			Suceess: data.Suceess,
			Fail:    data.Fail,
		}

		for _, message := range data.Messages {
			result.Messages = append(result.Messages, Message{Text: message})
		}
		db.Create(result)
		endTask <- struct{}{}
	}
}
