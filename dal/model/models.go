package model

import (
	"gorm.io/gorm"
)

type BotResponse struct {
	gorm.Model
	IntentId   string
	IntentName string
}

type StatsRequest struct {
	N         int    `json:"n"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type StatsResponse struct {
	Cnt        int
	IntentId   string
	IntentName string
}

func (BotResponse) TableName() string {
	return "bot_response"
}
