package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/traPtitech/piscon-portal/aws"
	"github.com/traPtitech/piscon-portal/conoha"
	"github.com/traPtitech/piscon-portal/model"
	"golang.org/x/crypto/acme/autocert"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	checkTask     chan struct{}
	sendWorker    chan *model.Task
	checkInstance chan *model.Instance
	db            *gorm.DB
	client        model.ServerClient
)

const (
	MAX_INSTANCE_NUMBER = 2
)

func main() {
	sendWorker = make(chan *model.Task, 10)
	checkTask = make(chan struct{})
	checkInstance = make(chan *model.Instance)

	go benchmarkWorker()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	cfg, err := aws.CreateDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err = aws.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	go instanceInfo(opts)

	// _db, err := gorm.Open("mysql", "isucon@/isucon?charset=utf8&parseTime=True&loc=Local")
	_db, err := establishConnection()
	if err != nil {
		panic(err)
	}
	//_db.LogMode(true)
	defer _db.Close()
	db = _db
	// db.LogMode(true)

	db.AutoMigrate(&model.Task{}, &model.Message{}, &model.Result{}, &model.Instance{}, &model.Team{}, &model.User{}, &model.Question{})

	tasks := []*model.Task{}
	db.Not("state = 'done'").Find(&tasks)
	for _, t := range tasks {
		go func(task *model.Task) {
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
	apiWithAuth.POST("/instance/:team_id/:instance_number", createInstance)
	apiWithAuth.DELETE("/instance/:team_id/:instance_number", deleteInstance)
	// TODO: ユーザー名で認証してないので修正する必要がある
	apiWithAuth.GET("/team/:id", getTeam)
	apiWithAuth.GET("/user/:name", getUser)
	apiWithAuth.POST("/benchmark/:name/:instance_number", queBenchmark)
	apiWithAuth.GET("/admin/team", getAllTeam)

	apiWithAuth.POST("/questions", postQuestions)
	apiWithAuth.PUT("/questions/:id", putQuestions)
	apiWithAuth.DELETE("/questions/:id", deleteQuestions)

	// e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(os.Getenv("HOST"))
	// e.AutoTLSManager.Cache = autocert.DirCache("/etc/letsencrypt/live/piscon-portal.trap.jp/cert.pem")
	// e.Pre(middleware.HTTPSWWWRedirect())
	// switch env {
	// case "prod":
	// e.StartAutoTLS(":443")
	// default:
	e.Use(middleware.CORS())
	e.Start(":4000")
	// }
	fmt.Println("end")
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
		data := &model.Output{}
		err = json.Unmarshal(res, data)
		if err != nil {
			result := &model.Result{
				TeamID:    task.TeamID,
				TaskID:    task.ID,
				Pass:      false,
				Score:     0,
				Campaign:  0,
				Betterize: task.Betterize,
				Messages:  []*model.Message{{Text: err.Error()}},
			}
			db.Create(result)

			task.State = "done"
			db.Save(task)
			continue
		}
		messages := make([]*model.Message, len(data.Messages))
		for i, text := range data.Messages {
			messages[i] = &model.Message{Text: text}
		}

		result := &model.Result{
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

// activeになったらdbにipアドレスとかを含めて登録
func instanceInfo(opts gophercloud.AuthOptions) {
	// 23時間ごとにtoken更新
	t := time.NewTicker(23 * time.Hour)
	for {
		select {
		case instance := <-checkInstance:
			fmt.Println("receive instance")
			go setupInstance(instance)
		case <-t.C:
			client = conoha.New(opts)
			fmt.Println("Conoha Client created")
		}
	}
}

func setupInstance(_instance *model.Instance) {
	instance := _instance
L:
	for {
		switch instance.Status {
		case model.BUILDING:
			log.Println("wait building")
			instance = waitBuilding(instance)
		case model.PRE_SHUTDOWN:
			log.Println("pre shutdown")
			instance.Status = model.SHUTDOWNING
			time.Sleep(5 * time.Second)
			client.StopInstance(instance.InstanceId)
		case model.SHUTDOWNING:
			log.Println("shutdowning")
			instance = waitShutdown(instance)
		case model.SHUTOFF:
			log.Println("shutoff")
			networkID := os.Getenv("CONOHA_NETWORK_ID")
			log.Printf("AttachPrivateNetwork name:%s networkID %s privateIP:%s\n", instance.Name, os.Getenv("CONOHA_NETWORK_ID"), instance.PrivateIPAddress)
			client.AttachPrivateNetwork(instance.Name, networkID, instance.PrivateIPAddress)
			client.StartInstance(instance.Name)
			instance.Status = model.STARTING
		case model.STARTING:
			log.Println("wait starting")
			instance = waitStarting(instance)
		case model.ACTIVE:
			log.Println("write to db")
			db.Model(&model.Instance{Name: instance.Name}).Update(instance)
			break L
		}
	}
}

func waitBuilding(instance *model.Instance) *model.Instance {
	time.Sleep(10 * time.Second)

	_instance, err := client.GetInstanceInfo(instance.Name)
	if err != nil {
		fmt.Println(err)
	}

	if strings.ToUpper(_instance.Status) == model.ACTIVE {
		IPv4 := ""
		// instanceのipv4のアドレスを抜き出そうとしてるけどもっといいやり方がありそう
		for _, v := range _instance.Addresses {
			for _, vv := range ([]interface{})(v.([]interface{})) {
				if (vv.(map[string]interface{})["version"]).(float64) == 4 {
					IPv4 = (vv.(map[string]interface{})["addr"]).(string)
				}
			}
		}
		if IPv4 != "" {
			instance.GlobalIPAddress = IPv4
			instance.Status = model.PRE_SHUTDOWN
		}
	}
	return instance
}

func waitShutdown(instance *model.Instance) *model.Instance {
	time.Sleep(10 * time.Second)

	_instance, err := client.GetInstanceInfo(instance.Name)
	if err != nil {
		fmt.Println(err)
	}
	if strings.ToUpper(_instance.Status) == model.SHUTOFF {
		instance.Status = model.SHUTOFF
	}
	return instance
}

func waitStarting(instance *model.Instance) *model.Instance {
	time.Sleep(10 * time.Second)

	_instance, err := client.GetInstanceInfo(instance.Name)
	if err != nil {
		fmt.Println(err)
	}
	if strings.ToUpper(_instance.Status) == model.ACTIVE {
		instance.Status = model.ACTIVE
	}
	return instance
}
