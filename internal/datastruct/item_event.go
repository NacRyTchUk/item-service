package datastruct

import "time"

type ItemEvent struct {
	Id          int64
	CampaignId  int64
	Name        string
	Description string
	Priority    int64
	Removed     bool
	EventTime   int64
}

func EventFrom(item *Item) (it ItemEvent) {
	it.Id = item.Id
	it.CampaignId = item.CampaignId
	it.Name = item.Name
	it.Description = ""
	it.Priority = item.Priority
	it.Removed = item.Removed
	it.EventTime = time.Now().Unix()

	if item.Description != nil {
		it.Description = *item.Description
	}
	return
}
