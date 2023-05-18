package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/blockloop/scan/v2"
	"item-service/internal/datastruct"
	"log"
	"time"
)

type ItemQuery interface {
	List() ([]datastruct.Item, error)
	Create(campaignId int64, name string) (*datastruct.Item, error)
	Update(id, campaignId int64, name string, description *string) (*datastruct.Item, error)
	Remove(id, campaignId int64) (*datastruct.Item, error)
}

type itemQuery struct {
	db  *sql.DB
	dao dao
}

func (i itemQuery) Create(campaignId int64, name string) (*datastruct.Item, error) {
	tx, err := i.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check if campaign exists
	_, err = i.dao.NewCampaignQuery().Get(campaignId)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Insert new line
	ib := pgQb().
		Insert(datastruct.ItemTableName).
		Columns("campaign_id", "name").
		Values(campaignId, name).
		Suffix("RETURNING *")
	var item datastruct.Item
	err = ib.QueryRow().Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description,
		&item.Priority, &item.Removed, &item.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error while item insertion: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commit transaction: %w", err)
	}

	i.dao.Cache.Del(context.Background(), "items")
	return &item, nil
}

func (i itemQuery) Update(id, campaignId int64, name string, description *string) (*datastruct.Item, error) {
	tx, err := i.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check if campaign exists
	_, err = i.dao.NewCampaignQuery().Get(campaignId)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	updateMap := make(map[string]interface{}, 2)
	updateMap["name"] = name
	if description != nil {
		updateMap["description"] = *description
	}
	// Update new line
	ub := pgQb().
		Update(datastruct.ItemTableName).
		SetMap(updateMap).
		Where(squirrel.And{squirrel.Eq{"id": id}, squirrel.Eq{"campaign_id": campaignId}, squirrel.Eq{"removed": false}}).
		Suffix("RETURNING *")
	var item datastruct.Item
	err = ub.QueryRow().Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description,
		&item.Priority, &item.Removed, &item.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error while item update: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commit transaction: %w", err)
	}

	i.dao.Cache.Del(context.Background(), "items")
	return &item, nil
}

func (i itemQuery) Remove(id, campaignId int64) (*datastruct.Item, error) {
	tx, err := i.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check if campaign exists
	_, err = i.dao.NewCampaignQuery().Get(campaignId)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	rb := pgQb().
		Update(datastruct.ItemTableName).
		Set("removed", true).
		Where(squirrel.And{squirrel.Eq{"id": id}, squirrel.Eq{"campaign_id": campaignId}, squirrel.Eq{"removed": false}}).
		Suffix("RETURNING *")
	var item datastruct.Item
	err = rb.QueryRow().Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description,
		&item.Priority, &item.Removed, &item.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error while item remove: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commit transaction: %w", err)
	}

	i.dao.Cache.Del(context.Background(), "items")
	return &item, nil
}

func (i itemQuery) List() ([]datastruct.Item, error) {
	var items []datastruct.Item

	cache, err := i.dao.Cache.Get(context.Background(), "items").Result()
	if err == nil {
		log.Println("cache hit")
		err = datastruct.ItemsUnmarshal([]byte(cache), &items)
		if err == nil {
			return items, nil
		}
	}
	log.Println("cache miss")

	tx, err := i.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get lines
	sb := pgQb().
		Select("*").
		From(datastruct.ItemTableName).
		Where(squirrel.Eq{"removed": false})
	rows, err := sb.Query()
	if err != nil {
		return nil, fmt.Errorf("error while getting items: %w", err)
	}
	err = scan.Rows(&items, rows)
	if err != nil {
		return nil, fmt.Errorf("error while rows scanning: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commit transaction: %w", err)
	}

	marshal, _ := datastruct.ItemsMarshal(items)
	err = i.dao.Cache.Set(context.Background(), "items", string(marshal), time.Minute).Err()
	if err != nil {
		log.Println(fmt.Errorf("cache err: %w", err))
	}

	return items, nil
}
