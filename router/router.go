package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MAX_INSTANCE_NUMBER = 2
)

type Handlers struct {
	client        model.ServerClient
	db            *gorm.DB
	checkInstance chan *model.Instance
	sendWorker    chan *model.Task
}

func genPassword() string {
	pass := ""
	gen := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	for i := 0; i < 12; i++ {
		pass += string(gen[rand.Intn(len(gen))])
	}
	return pass
}

func formatCommand(ip string) string {
	return fmt.Sprintf("/home/isucon/isucari/bin/benchmarker "+
		"-data-dir \"/home/isucon/isucari/initial-data\" "+
		"-payment-url \"http://172.16.0.1:5555\" "+
		"-shipment-url \"http://172.16.0.1:7000\" "+
		"-static-dir \"/home/isucon/isucari/webapp/public/static\" "+
		"-target-host \"%s\" "+
		"-target-url \"http://%s\"", ip, ip)
}

func EstablishConnection() (*gorm.DB, error) {
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

// traPかどうかの認証
// TODO: Fix ユーザーネーム認証
func middlewareAuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.NoContent(http.StatusForbidden)
		}
		req, _ := http.NewRequest("GET", "https://q.trap.jp/api/v3/users/me", nil)
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		res, _ := client.Do(req)
		if res.StatusCode != 200 {
			return c.NoContent(http.StatusForbidden)
		}
		return next(c)
	}
}

func (h *Handlers) GetNewer(c echo.Context) error {
	teams := []model.Team{}
	// チームIDのうち結果が存在するものをとってきて、かつ一回以上パスしており正の点数を取っていて、、かつ一日以内の者でもっとも得点が高いを一つ選択する
	h.db.Raw("SELECT * FROM results AS PI LEFT JOIN teams ON PI.team_id = teams.id WHERE PI.id =( SELECT po.id FROM results AS po LEFT JOIN teams ON po.team_id = teams.id WHERE pass = 1 AND PI.team_id = po.team_id AND score > 0 ORDER BY po.score DESC LIMIT 1 ) AND (PI.created_at > (CURRENT_TIME() - INTERVAL 1 day))").Scan(&teams)
	return c.JSON(http.StatusOK, teams)
}

func (h *Handlers) GetTeam(c echo.Context) error {
	id := c.Param("id")
	team := model.Team{}
	h.db.Where("id = ?", id).Find(&team)

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されていません"})
	}
	// Resultの中にあるMessageをPreloadする
	// belongs to で紐づいていいるやつのデータもとりだす。
	// Related
	h.db.Model(&team).Where("team_id = ?", &team.ID).Preload("Instance").Preload("Results.Messages").Find(&team)
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

func (h *Handlers) GetUser(c echo.Context) error {
	name := c.Param("name")
	user := model.User{}
	h.db.Where("name = ?", name).Find(&user)

	if user.Name == "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されていません"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handlers) GetAllTeam(c echo.Context) error {
	teams := []*model.Team{}
	h.db.Preload("Results").Preload("Instance").Preload("Results.Messages").Find(&teams)
	return c.JSON(http.StatusOK, teams)
}

func (h *Handlers) GetQuestions(c echo.Context) error {
	questions := []*model.Question{}
	h.db.Find(&questions)
	// TODO: 改行対応

	return c.JSON(http.StatusOK, questions)
}

func (h *Handlers) PostQuestions(c echo.Context) error {
	req := struct {
		Question string `json:"question"`
	}{}
	c.Bind(&req)
	question := &model.Question{
		Question: req.Question,
	}
	h.db.Create(question)
	return c.NoContent(http.StatusCreated)
}

func (h *Handlers) PutQuestions(c echo.Context) error {
	id := c.Param("id")
	req := struct {
		Answer string `json:"answer"`
	}{}
	c.Bind(&req)
	question := &model.Question{
		Answer: req.Answer,
	}
	h.db.Model(question).Where("id = ?", id).Updates(question)

	return c.JSON(http.StatusOK, question)
}

func (h *Handlers) DeleteQuestions(c echo.Context) error {
	id := c.Param("id")
	question := &model.Question{}
	h.db.Model(question).Where("id = ?", id).Delete(question)

	return c.JSON(http.StatusOK, question)
}

func (h *Handlers) CreateUser(c echo.Context) error {
	user := &model.User{}
	c.Bind(user)

	u := &model.User{}
	h.db.Where("name = ?", user.Name).Find(u)

	if u.Name != "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されています"})
	}

	h.db.Create(user)
	return c.JSON(http.StatusCreated, user)
}

func (h *Handlers) CreateTeam(c echo.Context) error {
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
	h.db.Where("name = ?", requestBody.Name).Find(t)

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
	h.db.Create(team)
	return c.JSON(http.StatusCreated, team)
}

func (h *Handlers) CreateInstance(c echo.Context) error {

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

	pass := os.Getenv("ISUCON_PASSWORD") //???情報はあるけど使われてなさそう
	i := &model.Instance{}
	h.db.Where("name = ?", name).Find(i)
	if i.Name != "" {
		return c.JSON(http.StatusConflict, "既に登録されています")
	}

	privateIP := fmt.Sprintf("172.16.0.%d", teamId*10+instanceNumber)

	log.Printf("CreateInstance name:%s pass %s privateIP:%s\n", name, pass, privateIP)
	id, err := h.client.CreateInstance(context.TODO(), name, privateIP)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	instance := &model.Instance{
		Password:         pass,
		InstanceNumber:   uint(instanceNumber),
		InstanceId:       *id,
		TeamID:           uint(teamId),
		Name:             name,
		Status:           model.BUILDING,
		GlobalIPAddress:  "",
		PrivateIPAddress: privateIP,
	}
	go func() {
		fmt.Println("send chan")
		h.checkInstance <- instance
	}()
	h.db.Create(instance)

	return nil
}

func (h *Handlers) DeleteInstance(c echo.Context) error {
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
	if err := h.db.Where("name = ?", name).First(i).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, "指定したインスタンスが見つかりません")
	}

	err = h.client.DeleteInstance(context.TODO(), i.InstanceId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	i = &model.Instance{}
	h.db.Where("name = ?", name).Delete(i)

	return nil
}

func (h *Handlers) GetAllResults(c echo.Context) error {
	teams := []*model.Team{}
	h.db.Find(&teams)
	for _, team := range teams {
		h.db.Where("team_id = ?", &team.ID).Preload("Messages").Find(&team.Results)
	}
	return c.JSON(http.StatusOK, teams)
}

func (h *Handlers) QueBenchmark(c echo.Context) error {
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
	h.db.Where("name = ?", name).Find(team)

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されていません"})
	}

	h.db.Model(team).Preload("Instance").Find(team.Instance)

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

	h.db.Where("team_id = ?", team.ID).Not("state = 'done'").First(task)
	if task.CmdStr != "" {
		return c.JSON(http.StatusNotAcceptable, model.Response{
			Success: false,
			Message: "既に登録されています"})
	}

	cmdStr := formatCommand(ip)
	t := &model.Task{
		CmdStr:    cmdStr,
		IP:        ip,
		State:     "waiting",
		TeamID:    team.ID,
		Betterize: req.Betterize,
	}
	fmt.Println(cmdStr)
	h.db.Create(t)

	go func() {
		h.sendWorker <- t
	}()

	return c.JSON(http.StatusCreated, model.Response{
		Success: false,
		Message: "キューに追加しました"})
}

func (h *Handlers) GetBenchmarkQueue(c echo.Context) error {
	tasks := h.getTaskQueInfo()
	for _, task := range tasks {
		h.db.Model(task).Preload("Team").Find(&task.Team)
	}
	return c.JSON(http.StatusOK, tasks)
}

func (h *Handlers) getTaskQueInfo() []*model.Task {
	tasks := []*model.Task{}
	h.db.Table("tasks").Joins("LEFT JOIN teams ON `teams`.id = `tasks`.team_id").Not("state = 'done'").Find(&tasks)
	return tasks
}
