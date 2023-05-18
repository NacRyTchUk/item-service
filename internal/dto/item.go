package dto

import "time"

type Item struct {
	Id          int64
	CampaignId  int64
	Name        string
	Description *string
	Priority    int64
	Removed     bool
	CreatedAt   time.Time
}
