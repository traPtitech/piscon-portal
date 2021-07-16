package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	ACTIVE  = "ACTIVE"
	SHUTOFF = "SHUTOFF"

	BUILDING     = "BUILDING"
	SHUTDOWNING  = "SHUTDOWNING"
	NOT_EXIST    = "NOT_EXIST"
	STARTING     = "STARTING"
	PRE_SHUTDOWN = "PRE_SHUTDOWN"
)

type ServerClient interface {
	CreateInstance(name string, privateIp string, pwd string) (*string, error) //return InstanceID (TODO)
	DeleteInstance(instanceId string) error
	StartInstance(instanceId string) error
	StopInstance(instanceId string) error
	GetInstanceInfo(id string) (*Instance, error) //TODO IDにする
}

type Response struct {
	Success bool   `json:"suceess"`
	Message string `json:"message"`
}

type Output struct {
	Pass     bool     `json:"pass"`
	Score    int64    `json:"score"`
	Campaign int64    `json:"campaign"`
	Language string   `json:"language"`
	Messages []string `json:"messages"`
}

type Message struct {
	gorm.Model
	ResultId uint   `json:"result_id"`
	Text     string `json:"text"`
}

type Team struct {
	gorm.Model
	Name              string      `gorm:"unique size:50" json:"name"`
	Instance          []*Instance `json:"instance"`
	Results           []*Result   `json:"results"`
	MaxInstanceNumber int         `json:"max_instance_number"`
	Group             string      `json:"group"`
}

type User struct {
	gorm.Model
	Name       string `gorm:"unique size:50" json:"name"`
	ScreenName string `json:"screen_name"`
	IconFileId string `json:"iconFileId"`
	TeamID     uint   `json:"team_id"`
}

type Instance struct {
	gorm.Model
	TeamID           uint   `json:"team_id"`
	GlobalIPAddress  string `json:"global_ip_address"`
	PrivateIPAddress string `json:"private_ip_address"`
	Password         string `json:"password"`
	InstanceNumber   uint   `json:"instance_number"`
	InstanceId       string `json:"instance_id"`
	Status           string `json:"status"`
	Name             string `json:"name"`
}

type Result struct {
	ID        int        `gorm:"AUTO_INCREMENT" json:"id"`
	TeamID    uint       `json:"team_id"`
	TaskID    uint       `json:"task_id"`
	Pass      bool       `json:"pass"`
	Score     int64      `json:"score"`
	Campaign  int64      `json:"campaign"`
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
