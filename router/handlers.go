package router

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal/model"
	"github.com/traPtitech/piscon-portal/oauth"
	"gorm.io/gorm"
)

const (
	MAX_INSTANCE_NUMBER = 3
)

type Handlers struct {
	client        model.ServerClient
	db            *gorm.DB
	checkInstance chan *model.Instance
	sendWorker    chan *model.Task
	authConf      *oauth.OauthClient
}

func NewHandlers(c model.ServerClient, db *gorm.DB, ci chan *model.Instance, sw chan *model.Task) *Handlers {
	return &Handlers{
		client:        c,
		db:            db,
		checkInstance: ci,
		sendWorker:    sw,
		authConf:      oauth.New(),
	}
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
		"-payment-url \"http://10.0.145.247:5555\" "+
		"-shipment-url \"http://10.0.145.247:7000\" "+
		"-static-dir \"/home/isucon/isucari/webapp/public/static\" "+
		"-target-host \"%s\" "+
		"-target-url \"http://%s\"", ip, ip)
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
	h.db.Where("id = ?", id).Preload("Instance").Preload("Results").Find(&team)

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されていません"})
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
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
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
		return c.JSON(http.StatusConflict, model.Response{
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
	h.db.Preload("Instance").Where("name = ?", requestBody.Name).Find(t)

	if t.Name == "" {
		ins := initializeInstances()
		team := &model.Team{
			Name:              requestBody.Name,
			MaxInstanceNumber: MAX_INSTANCE_NUMBER,
			Instance:          *ins,
			Group:             requestBody.Group,
		}
		h.db.Create(team)
		return c.JSON(http.StatusCreated, team)
	}
	return c.JSON(http.StatusCreated, t)
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
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: err.Error()})
	}

	if instanceNumber != 1 && instanceNumber != 2 && instanceNumber != 3 {
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: "instance number should be 1 or 2 or 3"})
	}

	name := fmt.Sprintf("%d-%d", teamId, instanceNumber)

	pass := genPassword()
	i := &model.Instance{}
	h.db.Where("name = ?", name).Find(i)
	if i.Status == model.ACTIVE {
		return c.JSON(http.StatusConflict, model.Response{
			Success: false,
			Message: "起動中です"})
	}

	privateIP := fmt.Sprintf("10.0.0.%d", teamId*10+instanceNumber)
	id, err := h.client.CreateInstance(name, privateIP, pass)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Success: false,
			Message: err.Error()})
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
	h.db.Model(instance).Where("id = ?", instanceNumber).Where("team_id = ?", instance.TeamID).Updates(instance)
	return c.JSON(http.StatusCreated, instance)
}

func (h *Handlers) DeleteInstance(c echo.Context) error {
	log.Println("delete command received")
	instanceNumber, err := strconv.Atoi(c.Param("instance_number"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: err.Error()})
	}
	teamId, err := strconv.Atoi(c.Param("team_id"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: err.Error()})
	}

	if instanceNumber != 1 && instanceNumber != 2 && instanceNumber != 3 {
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: "instance number should be 1 or 2 or 3"})
	}

	name := fmt.Sprintf("%d-%d", teamId, instanceNumber)
	i := &model.Instance{}
	if err := h.db.Where("name = ?", name).First(i).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "指定したインスタンスが見つかりません"})
	}

	err = h.client.DeleteInstance(i.InstanceId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Success: false,
			Message: "Internal server error"})
	}
	ins := emptyInstance(int(i.InstanceNumber))

	h.db.Where("instance_id = ?", i.InstanceId).Updates(ins)

	return c.JSON(http.StatusNoContent, nil)
}

func (h *Handlers) GetAllResults(c echo.Context) error {
	teams := []*model.Team{}
	h.db.Find(&teams)
	for _, team := range teams {
		h.db.Where("team_id = ?", &team.ID).Preload("Messages").Find(&team.Results)
		if team.Results == nil {
			team.Results = []*model.Result{}
		}
	}
	return c.JSON(http.StatusOK, teams)
}

func (h *Handlers) QueBenchmark(c echo.Context) error {
	instanceNumber, err := strconv.Atoi(c.Param("instance_number"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: err.Error()})
	}
	name := c.Param("name")

	req := struct {
		Betterize string `json:"betterize"`
	}{}

	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: err.Error()})
	}

	team := &model.Team{}
	if err = h.db.Where("name = ?", name).Preload("Instance").Find(team).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Success: false,
			Message: err.Error()})
	}

	if team.Name == "" {
		return c.JSON(http.StatusNotFound, model.Response{
			Success: false,
			Message: "登録されていません"})
	}

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

	if err = h.db.Where("team_id = ?", team.ID).Not("state = ? ", "done").First(task).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, model.Response{
				Success: false,
				Message: err.Error()})
		}
	}
	if task.CmdStr != "" {
		return c.JSON(http.StatusConflict, model.Response{
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
	if err = h.db.Create(t).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Success: false,
			Message: err.Error()})
	}

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
		h.db.Preload("Team").Find(&task)
	}
	return c.JSON(http.StatusOK, tasks)
}

func (h *Handlers) GetTeamMember(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("team_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Message: err.Error()})
	}
	var member []model.User
	if err = h.db.Where("team_id = ?", id).Find(&member).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, model.Response{
				Success: false,
				Message: err.Error()})
		}
	}
	return c.JSON(http.StatusOK, member)
}

func (h *Handlers) getTaskQueInfo() []*model.Task {
	tasks := []*model.Task{}
	h.db.Table("tasks").Joins("LEFT JOIN teams ON `teams`.id = `tasks`.team_id").Not("state = 'done'").Find(&tasks)
	return tasks
}
func initializeInstances() *[]*model.Instance {
	res := []*model.Instance{}
	for i := 0; i < MAX_INSTANCE_NUMBER; i++ {
		ins := emptyInstance(i + 1)
		res = append(res, ins)
	}
	return &res
}

func emptyInstance(n int) *model.Instance {
	emptyInstance := &model.Instance{}
	emptyInstance.InstanceNumber = uint(n)
	emptyInstance.Status = model.NOT_EXIST
	return emptyInstance
}
