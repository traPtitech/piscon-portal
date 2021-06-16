package main

import (
	"context"
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
	"github.com/traPtitech/piscon-portal/aws"
	"github.com/traPtitech/piscon-portal/conoha"
	"github.com/traPtitech/piscon-portal/model"
	"golang.org/x/crypto/acme/autocert"
)

type serverClient interface {
	CreateInstance(c context.Context, name string, subnetId string, privateIp string) error
	DeleteInstance(c context.Context, instanceId string) error
	StartInstance(c context.Context, instanceId string) error
	StopInstance(c context.Context, instanceId string) error
	GetInstanceInfo(c context.Context, instanceName string) (*model.Instance, error)
}

var (
	checkTask     chan struct{}
	sendWorker    chan *model.Task
	checkInstance chan *model.Instance
	db            *gorm.DB
	client        serverClient
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
	log.Println((fmt.Sprintf("%s:%s@(%s)/%s", user, pass, host, dbname) + "?charset=utf8mb4&parseTime=True&loc=Local"))
	_db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", user, pass, host, dbname)+"?charset=utf8mb4&parseTime=True&loc=Local")
	_db.BlockGlobalUpdate(true)
	db = _db
	return db, err
}

func getNewer(c echo.Context) error {
	teams := []model.Team{}
	db.Raw("SELECT * FROM results AS PI LEFT JOIN teams ON PI.team_id = teams.id WHERE PI.id =( SELECT po.id FROM results AS po LEFT JOIN teams ON po.team_id = teams.id WHERE pass = 1 AND PI.team_id = po.team_id AND score > 0 ORDER BY po.score DESC LIMIT 1 ) AND (PI.created_at > (CURRENT_TIME() - INTERVAL 1 day))").Scan(&teams)
	return c.JSON(http.StatusOK, teams)
}

func getTeam(c echo.Context) error {
	id := c.Param("id")
	team := model.Team{}
	db.Where("id = ?", id).Find(&team)

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されていません"})
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
			emptyInstance := &model.Instance{}
			emptyInstance.InstanceNumber = uint(i + 1)
			emptyInstance.Status = model.NOT_EXIST
			team.Instance = append(team.Instance, emptyInstance)
		}
	}

	return c.JSON(http.StatusOK, team)
}

func getUser(c echo.Context) error {
	name := c.Param("name")
	user := model.User{}
	db.Where("name = ?", name).Find(&user)

	if user.Name == "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されていません"})
	}

	return c.JSON(http.StatusOK, user)
}

func getAllTeam(c echo.Context) error {
	teams := []*model.Team{}
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
	questions := []*model.Question{}
	db.Find(&questions)
	// TODO: 改行対応

	return c.JSON(http.StatusOK, questions)
}

func postQuestions(c echo.Context) error {
	req := struct {
		Question string `json:"question"`
	}{}
	c.Bind(&req)
	question := &model.Question{
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
	question := &model.Question{
		Answer: req.Answer,
	}
	db.Model(question).Where("id = ?", id).Update(&question)

	return c.JSON(http.StatusOK, question)
}

func deleteQuestions(c echo.Context) error {
	id := c.Param("id")
	question := &model.Question{}
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
	user := &model.User{}
	c.Bind(user)

	u := &model.User{}
	db.Where("name = ?", user.Name).Find(u)

	if u.Name != "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されています"})
	}

	db.Create(user)
	return c.JSON(http.StatusCreated, user)
}

func createTeam(c echo.Context) error {
	requestBody := &struct {
		Name  string `json:"name"`
		Group string `json:"group"`
	}{}

	c.Bind(requestBody)

	if requestBody.Name == "" {
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: "リクエストボディの要素が足りません"})
	}

	t := &model.Team{}
	db.Where("name = ?", requestBody.Name).Find(t)

	if t.Name != "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されています"})
	}
	// pass := genPassword()

	team := &model.Team{
		Name:              requestBody.Name,
		MaxInstanceNumber: MAX_INSTANCE_NUMBER,
		Instance:          []*model.Instance{},
		Group:             requestBody.Group,
	}
	db.Create(team)
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
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: "instance number should be 1 or 2 or 3"})
	}

	name := fmt.Sprintf("%d-%d", teamId, instanceNumber)

	pass := os.Getenv("CONOHA_ISUCON_PASSWORD")
	i := &model.Instance{}
	db.Where("name = ?", name).Find(i)
	if i.Name != "" {
		return c.JSON(http.StatusConflict, "既に登録されています")
	}

	privateIP := fmt.Sprintf("172.16.0.%d", teamId*10+instanceNumber)

	log.Printf("Makeinstance name:%s pass %s privateIP:%s\n", name, pass, privateIP)
	err = client.MakeInstance(name, privateIP)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	instance := &model.Instance{
		Password:         pass,
		InstanceNumber:   uint(instanceNumber),
		TeamID:           uint(teamId),
		Name:             name,
		Status:           model.BUILDING,
		GlobalIPAddress:  "",
		PrivateIPAddress: privateIP,
	}
	go func() {
		fmt.Println("send chan")
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
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: "instance number should be 1 or 2 or 3"})
	}

	name := fmt.Sprintf("%d-%d", teamId, instanceNumber)
	i := &model.Instance{}
	if db.Where("name = ?", name).First(i).RecordNotFound() {
		return c.JSON(http.StatusNotFound, "指定したインスタンスが見つかりません")
	}

	err = client.DeleteInstance(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	i = &model.Instance{}
	db.Where("name = ?", name).Delete(i)

	return nil
}

func getAllResults(c echo.Context) error {
	teams := []*model.Team{}
	db.Find(&teams)
	for _, team := range teams {
		db.Where("team_id = ?", &team.ID).Preload("Messages").Find(&team.Results)
	}
	return c.JSON(http.StatusOK, teams)
}

func queBenchmark(c echo.Context) error {
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

	team := &model.Team{}
	db.Where("name = ?", name).Find(team)

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されていません"})
	}

	db.Model(team).Related(&team.Instance)

	ip := team.Instance[0].PrivateIPAddress

	for _, instance := range team.Instance {
		if uint(instanceNumber) == instance.InstanceNumber {
			ip = instance.PrivateIPAddress
		}
	}

	if ip == "" {
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: "インスタンスが存在しません"})
	}

	task := &model.Task{}

	db.Where("team_id = ?", team.ID).Not("state = 'done'").First(task)
	if task.CmdStr != "" {
		return c.JSON(http.StatusNotAcceptable, model.Response{
			Success: false,
			Message: "既に登録されています"})
	}

	cmdStr := fmt.Sprintf("/home/isucon/isucari/bin/benchmarker "+
		"-data-dir \"/home/isucon/isucari/initial-data\" "+
		"-payment-url \"http://172.16.0.1:5555\" "+
		"-shipment-url \"http://172.16.0.1:7000\" "+
		"-static-dir \"/home/isucon/isucari/webapp/public/static\" "+
		"-target-host \"%s\" "+
		"-target-url \"http://%s\"", ip, ip)
	t := &model.Task{
		CmdStr:    cmdStr,
		IP:        ip,
		State:     "waiting",
		TeamID:    team.ID,
		Betterize: req.Betterize,
	}
	fmt.Println(cmdStr)
	db.Create(t)

	go func() {
		sendWorker <- t
	}()

	return c.JSON(http.StatusCreated, model.Response{
		Success: false,
		Message: "キューに追加しました"})
}

func getBenchmarkQueue(c echo.Context) error {
	tasks := getTaskQueInfo()
	for _, task := range tasks {
		db.Model(task).Related(&task.Team)
	}
	return c.JSON(http.StatusOK, tasks)
}

func getTaskQueInfo() []*model.Task {
	tasks := []*model.Task{}
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
			client.ShutdownInstance(instance.Name)
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
