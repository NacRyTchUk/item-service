package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"item-service/internal/datastruct"
)

type CampaignQuery interface {
	Get(campaignId int64) (*datastruct.Campaign, error)
}

type campaignQuery struct {
	db *sql.DB
}

func (c campaignQuery) Get(campaignId int64) (*datastruct.Campaign, error) {
	qb := pgQb().
		Select("name").
		From(datastruct.CampaignTableName).
		Where(squirrel.Eq{"id": campaignId})
	var name string
	err := qb.QueryRow().Scan(&name)
	if err != nil {
		return nil, err
	}
	return &datastruct.Campaign{
		Id:   campaignId,
		Name: name,
	}, nil
}
