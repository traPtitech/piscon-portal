package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/traPtitech/piscon-portal/conoha"
	"golang.org/x/crypto/acme/autocert"
)

var (
	checkTask     chan struct{}
	sendWorker    chan *Task
	checkInstance chan *Instance
	db            *gorm.DB
	conohaClient  *conoha.ConohaClient
)

const (
	BUILD   = "BUILD"
	ACTIVE  = "ACTIVE"
	SHUTOFF = "SHUTOFF"

	SHUTDOWNING  = "SHUTDOWNING"
	NOT_EXIST    = "NOT_EXIST"
	STARTING     = "STARTING"
	PRE_SHUTDOWN = "PRE_SHUTDOWN"

	MAX_INSTANCE_NUMBER = 3
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
	ResultId uint
	Text     string
}

type Team struct {
	gorm.Model
	Name              string      `gorm:"unique size:50" json:"name"`
	Instance          []*Instance `json:"instance"`
	Results           []*Result   `json:"results"`
	MaxInstanceNumber int         `json:"max_instance_number"`
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
	InstanceNumber   uint   `json:"instance_number"`
	Status           string `json:"status"`
	Name             string `json:"name"`
}

type Result struct {
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
	CmdStr    string `json:"cmd_str" sql:"type:text;"`
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
	checkInstance = make(chan *Instance)

	go benchmarkWorker()
	go instanceInfo()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "https://identity.tyo2.conoha.io/v2.0",
		Username:         os.Getenv("CONOHA_USERNAME"),
		TenantName:       os.Getenv("CONOHA_TENANT_NAME"),
		TenantID:         os.Getenv("CONOHA_TENANT_ID"),
		Password:         os.Getenv("CONOHA_PASSWORD"),
	}

	conohaClient = conoha.New(opts)

	// _db, err := gorm.Open("mysql", "isucon@/isucon?charset=utf8&parseTime=True&loc=Local")
	_db, err := establishConnection()
	if err != nil {
		panic(err)
	}
	//_db.LogMode(true)
	defer _db.Close()
	db = _db
	// db.LogMode(true)

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
	log.Println((fmt.Sprintf("%s:%s@(%s)/%s", user, pass, host, dbname)+"?charset=utf8mb4&parseTime=True&loc=Local"))
	_db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", user, pass, host, dbname)+"?charset=utf8mb4&parseTime=True&loc=Local")
	_db.BlockGlobalUpdate(true)
	db = _db
	return db, err
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
	flag := false
	for i := 0; i < team.MaxInstanceNumber; i++ {
		flag = false
		for _, instance := range team.Instance {
			if instance.InstanceNumber == uint(i+1) {
				flag = true
			}
		}
		if !flag {
			emptyInstance := &Instance{}
			emptyInstance.InstanceNumber = uint(i + 1)
			emptyInstance.Status = NOT_EXIST
			team.Instance = append(team.Instance, emptyInstance)
		}
	}

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
		Name string `json:"name"`
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
	// pass := genPassword()

	team := &Team{
		Name:              requestBody.Name,
		MaxInstanceNumber: MAX_INSTANCE_NUMBER,
		// Instance: make([]*Instance, 0, 3),
		Instance: []*Instance{},
	}
	// team.Instance = append(team.Instance, &instance)
	// log.Printf("len %d,cap %d", len(team.Instance), cap(team.Instance))

	// log.Println(instance)

	// log.Println(team)
	db.Create(team)

	// db.Where("name = ?", team.Name).Find(t)

	// for i := 1; i <= 3; i++ {
	// 	db.Create(&Instance{
	// 		Status:         NOT_EXIST,
	// 		InstanceNumber: uint(i),
	// 		TeamID:         t.ID,
	// 	})
	// }
	return c.JSON(http.StatusCreated, team)
}

func createInstance(c echo.Context) error {

	instanceNumber, err := strconv.Atoi(c.Param("instance_number"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	teamId, err := strconv.Atoi(c.Param("team_id"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	if instanceNumber != 1 && instanceNumber != 2 && instanceNumber != 3 {
		return c.JSON(http.StatusBadRequest, Response{false, "instance number should be 1 or 2 or 3"})
	}

	name := fmt.Sprintf("%d-%d", teamId, instanceNumber)

	pass := genPassword()
	i := &Instance{}
	db.Where("name = ?", name).Find(i)
	if i.Name != "" {
		return c.JSON(http.StatusConflict, "既に登録されています")
	}

	privateIP := fmt.Sprintf("172.16.0.%d", teamId*10+instanceNumber)

	log.Printf("Makeinstance name:%s pass %s privateIP:%s\n", name, pass, privateIP)
	err = conohaClient.MakeInstance(name, pass, privateIP)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	instance := &Instance{
		Password:         pass,
		InstanceNumber:   uint(instanceNumber),
		TeamID:           uint(teamId),
		Name:             name,
		Status:           BUILD,
		GlobalIPAddress:  "",
		PrivateIPAddress: privateIP,
	}
	go func() {
		checkInstance <- instance
	}()
	db.Create(instance)

	return nil
}

func deleteInstance(c echo.Context) error {
	log.Println("delete command received")
	instanceNumber, err := strconv.Atoi(c.Param("instance_number"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	teamId, err := strconv.Atoi(c.Param("team_id"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	if instanceNumber != 1 && instanceNumber != 2 && instanceNumber != 3 {
		return c.JSON(http.StatusBadRequest, Response{false, "instance number should be 1 or 2 or 3"})
	}

	name := fmt.Sprintf("%d-%d", teamId, instanceNumber)
	i := &Instance{}
	if db.Where("name = ?", name).First(i).RecordNotFound() {
		return c.JSON(http.StatusNotFound, "指定したインスタンスが見つかりません")
	}

	err = conohaClient.DeleteInstance(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	i = &Instance{}
	db.Where("name = ?", name).Delete(i)

	return nil
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
	// id := strconv.atoi(c.Param("id"))
	instanceNumber, err := strconv.Atoi(c.Param("instance_number"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
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

	ip := team.Instance[0].PrivateIPAddress


	for _, instance := range team.Instance {
		if uint(instanceNumber) == instance.InstanceNumber {
			ip = instance.PrivateIPAddress
		}
	}
		
	
	// if team.Instance[0].PrivateIPAddress == "" {
	// 	return c.JSON(http.StatusBadRequest, Response{false, "インスタンスが存在しません"})
	// }

	// if id == "2" {
	// 	ip = team.Instance[1].PrivateIPAddress
	// } else if id == "3" {
	// 	ip = team.Instance[2].PrivateIPAddress
	// }
	if ip == "" {
		return c.JSON(http.StatusBadRequest, Response{false, "インスタンスが存在しません"})
	}

	task := &Task{}

	db.Where("team_id = ?", team.ID).Not("state = 'done'").First(task)
	if task.CmdStr != "" {
		return c.JSON(http.StatusNotAcceptable, Response{false, "すでに登録されています"})
	}

	cmdStr := fmt.Sprintf("/home/isucon/isucari/bin/benchmarker "+
		"-data-dir \"/home/isucon/isucari/initial-data\" "+
		"-payment-url \"http://160.251.13.26:5555\""+
		"-shipment-url \"http://160.251.13.26:7000\""+
		"-static-dir \"/home/isucon/isucari/webapp/public/static\" "+
		"-target-host \"%s\" "+
		"-target-url http://%s", ip, ip)
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

// activeになったらdbにipアドレスとかを含めて登録
func instanceInfo() {
	for {
		instance := <-checkInstance

		fmt.Println("receive instance")
		switch instance.Status {
		case BUILD:
			log.Println("wait building")
			waitBuilding(instance)
			// instance.Status = SHUTOFF
			// go func() { checkInstance <- instance }()
		case PRE_SHUTDOWN:
			log.Println("pre shutdown")
			instance.Status = SHUTDOWNING
			time.Sleep(5 * time.Second)
			conohaClient.ShutdownInstance(instance.Name)
			go func() { checkInstance <- instance }()
		case SHUTDOWNING:
			log.Println("shutdowning")
			waitShutdown(instance)
		case SHUTOFF:
			log.Println("shutoff")
			networkID := os.Getenv("CONOHA_NETWORK_ID")
			log.Printf("AttachPrivateNetwork name:%s networkID %s privateIP:%s\n", instance.Name, os.Getenv("CONOHA_NETWORK_ID"), instance.PrivateIPAddress)
			conohaClient.AttachPrivateNetwork(instance.Name, networkID, instance.PrivateIPAddress)
			conohaClient.StartInstance(instance.Name)
			instance.Status = STARTING
			go func() { checkInstance <- instance }()
		case STARTING:
			log.Println("wait starting")
			waitStarting(instance)
		case ACTIVE:
			log.Println("write to db")
			db.Model(&Instance{Name: instance.Name}).Update(instance)
		}
	}
}

func waitBuilding(instance *Instance) {
	_instance, err := conohaClient.GetInstanceInfo(instance.Name)
	if err != nil {
		fmt.Println(err)
	}

	if strings.ToUpper(_instance.Status) == ACTIVE {
		IPv4 := ""
		// instanceのipv4のアドレスを抜き出そうとしてるけどもっといいやり方がありそう
		for _, v := range _instance.Addresses {
			// tmp := (([]interface{})(v.([]interface{})))
			for _, vv := range ([]interface{})(v.([]interface{})) {
				if (vv.(map[string]interface{})["version"]).(float64) == 4 {
					IPv4 = (vv.(map[string]interface{})["addr"]).(string)
				}
				// if i == 0 {
				// 	for k, vvv := range vv.(map[string]interface{}) {
				// 		if k == "addr" {
				// 			IPv4 = vvv.(string)
				// 		}
				// 	}
				// }
			}
			// fmt.Println(tmp)
		}
		if IPv4 != "" {
			instance.GlobalIPAddress = IPv4
			instance.Status = PRE_SHUTDOWN
			// db.Model(&Instance{Name: instance.Name}).Update(instance)
		}
	}
	go func() {
		checkInstance <- instance
	}()

	time.Sleep(10 * time.Second)
}

func waitShutdown(instance *Instance) {
	_instance, err := conohaClient.GetInstanceInfo(instance.Name)
	if err != nil {
		fmt.Println(err)
	}
	if strings.ToUpper(_instance.Status) == SHUTOFF {
		instance.Status = SHUTOFF
	}
	go func() { checkInstance <- instance }()

	time.Sleep(10 * time.Second)
}

func waitStarting(instance *Instance) {
	_instance, err := conohaClient.GetInstanceInfo(instance.Name)
	if err != nil {
		fmt.Println(err)
	}
	if strings.ToUpper(_instance.Status) == ACTIVE {
		instance.Status = ACTIVE
	}
	go func() { checkInstance <- instance }()

	time.Sleep(10 * time.Second)
}
