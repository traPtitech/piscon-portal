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

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	shellwords "github.com/mattn/go-shellwords"
	plugin "github.com/traPtitech/piscon-portal/aws"
	"github.com/traPtitech/piscon-portal/model"
	"github.com/traPtitech/piscon-portal/router"
	"golang.org/x/crypto/acme/autocert"
	"gorm.io/driver/mysql"
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

type Config plugin.Config

func main() {
	sendWorker = make(chan *model.Task, 10)
	checkTask = make(chan struct{})
	checkInstance = make(chan *model.Instance)

	go benchmarkWorker()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	cfg, err := plugin.CreateDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err = plugin.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	go instanceInfo(*cfg)

	// _db, err := gorm.Open("mysql", "isucon@/isucon?charset=utf8&parseTime=True&loc=Local")
	_db, err := establishConnection()
	if err != nil {
		panic(err)
	}
	//_db.LogMode(true)
	_cl, err := _db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer _cl.Close()
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
	h := router.NewHandlers(client, db, checkInstance, sendWorker)
	h.SetUp(e)
	e.Use(middleware.CORS())
	e.Start(":4000")
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
func instanceInfo(cfg plugin.Config) {
	// 23時間ごとにtoken更新
	t := time.NewTicker(23 * time.Hour)
	for {
		select {
		case instance := <-checkInstance:
			fmt.Println("receive instance")
			go setupInstance(instance)
		case <-t.C:
			_client, err := plugin.New(cfg)
			if err != nil {
				log.Fatal(err)
			}
			client = _client
			fmt.Println("Client created")
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
			// networkID := os.Getenv("CONOHA_NETWORK_ID")
			// log.Printf("AttachPrivateNetwork name:%s networkID %s privateIP:%s\n", instance.Name, os.Getenv("CONOHA_NETWORK_ID"), instance.PrivateIPAddress)
			// client.AttachPrivateNetwork(instance.Name, networkID, instance.PrivateIPAddress)
			client.StartInstance(instance.Name)
			instance.Status = model.STARTING
		case model.STARTING:
			log.Println("wait starting")
			instance = waitStarting(instance)
		case model.ACTIVE:
			log.Println("write to db")
			db.Model(&model.Instance{Name: instance.Name}).Updates(instance)
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
		instance.GlobalIPAddress = _instance.GlobalIPAddress
		instance.Status = model.ACTIVE

		// // instanceのipv4のアドレスを抜き出そうとしてるけどもっといいやり方がありそう
		// for _, v := range _instance.Addresses {
		// 	for _, vv := range ([]interface{})(v.([]interface{})) {
		// 		if (vv.(map[string]interface{})["version"]).(float64) == 4 {
		// 			IPv4 = (vv.(map[string]interface{})["addr"]).(string)
		// 		}
		// 	}
		// }
		// if IPv4 != "" {
		// 	instance.GlobalIPAddress = IPv4
		// 	instance.Status = model.PRE_SHUTDOWN
		// }
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

func establishConnection() (*gorm.DB, error) {
	user := os.Getenv("MARIADB_USERNAME")
	if user == "" {
		user = "isucon"
	}

	pass := os.Getenv("MARIADB_PASSWORD")
	if pass == "" {
		pass = "isucon"
	}
	env := os.Getenv("ENV")
	host := os.Getenv("MARIADB_HOSTNAME")

	switch env {
	case "prod":
		if host == "" {
			host = "localhost"
		}
	default:
		if host == "" {
			host = "db"
		}
	}

	dbname := os.Getenv("MARIADB_DATABASE")
	if dbname == "" {
		dbname = "isucon"
	}
	dsn := fmt.Sprintf("%s:%s@(%s)/%s", user, pass, host, dbname) + "?charset=utf8mb4&parseTime=True&loc=Local"
	log.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// _db.BlockGlobalUpdate(true) <= GOrm v2でデフォルト有効になったらしい(要調査)
	return db, err
}
