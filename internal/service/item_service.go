package service

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"item-service/internal/datastruct"
	"item-service/internal/dto"
	"item-service/internal/repository"
	"log"
)

type ItemService interface {
	List() ([]dto.Item, error)
	Create(campaignId int64, name string) (dto.Item, error)
	Update(id, campaignId int64, name string, description *string) (dto.Item, error)
	Remove(id, campaignId int64) error
}

type itemService struct {
	dao    repository.DAO
	Broker *nats.Conn
}

func NewItemService(
	dao repository.DAO,
	broker *nats.Conn,
) ItemService {
	return &itemService{
		dao:    dao,
		Broker: broker,
	}
}

func (i itemService) List() ([]dto.Item, error) {
	list, err := i.dao.NewItemQuery().List()
	if err != nil {
		return nil, err
	}

	items := make([]dto.Item, len(list))
	for i, v := range list {
		items[i] = dto.Item{
			Id:          v.Id,
			CampaignId:  v.CampaignId,
			Name:        v.Name,
			Description: v.Description,
			Priority:    v.Priority,
			Removed:     v.Removed,
			CreatedAt:   v.CreatedAt,
		}
	}
	return items, nil
}

func (i itemService) Create(campaignId int64, name string) (dto.Item, error) {
	create, err := i.dao.NewItemQuery().Create(campaignId, name)
	if err != nil {
		return dto.Item{}, fmt.Errorf("error while creating a item: %w", err)
	}

	event := datastruct.EventFrom(create)
	marshal, _ := json.Marshal(event)
	err = i.Broker.Publish("logs", marshal)
	if err != nil {
		log.Println(fmt.Errorf("error while logging event: %w", err))
	}

	return dto.Item{
		Id:          create.Id,
		CampaignId:  create.CampaignId,
		Name:        create.Name,
		Description: create.Description,
		Priority:    create.Priority,
		Removed:     create.Removed,
		CreatedAt:   create.CreatedAt,
	}, nil
}

func (i itemService) Update(id, campaignId int64, name string, description *string) (dto.Item, error) {
	update, err := i.dao.NewItemQuery().Update(id, campaignId, name, description)
	if err != nil {
		return dto.Item{}, err
	}

	event := datastruct.EventFrom(update)
	marshal, _ := json.Marshal(event)
	err = i.Broker.Publish("logs", marshal)
	if err != nil {
		log.Println(fmt.Errorf("error while logging event: %w", err))
	}

	return dto.Item{
		Id:          update.Id,
		CampaignId:  update.CampaignId,
		Name:        update.Name,
		Description: update.Description,
		Priority:    update.Priority,
		Removed:     update.Removed,
		CreatedAt:   update.CreatedAt,
	}, err
}

func (i itemService) Remove(id, campaignId int64) error {
	remove, err := i.dao.NewItemQuery().Remove(id, campaignId)
	if err != nil {
		return err
	}

	event := datastruct.EventFrom(remove)
	marshal, _ := json.Marshal(event)
	err = i.Broker.Publish("logs", marshal)
	if err != nil {
		log.Println(fmt.Errorf("error while logging event: %w", err))
	}

	return nil
}
