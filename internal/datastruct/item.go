package datastruct

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"time"
)

const ItemTableName = "items"

type Item struct {
	Id          int64     `db:"id"`
	CampaignId  int64     `db:"campaign_id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	Priority    int64     `db:"priority"`
	Removed     bool      `db:"removed"`
	CreatedAt   time.Time `db:"created_at"`
}

type wrapItem struct {
	Id          int64
	CampaignId  int64
	Name        string
	Description *string
	Priority    int64
	Removed     bool
	CreatedAt   time.Time
}

func (it Item) MarshalBinary() ([]byte, error) {
	w := wrapItem{
		Id:          it.Id,
		CampaignId:  it.CampaignId,
		Name:        it.Name,
		Description: it.Description,
		Priority:    it.Priority,
		Removed:     it.Removed,
		CreatedAt:   it.CreatedAt,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(w); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (it *Item) UnmarshalBinary(data []byte) error {
	w := wrapItem{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&w); err != nil {
		return err
	}
	it.Id = w.Id
	it.CampaignId = w.CampaignId
	it.Name = w.Name
	it.Description = w.Description
	it.Priority = w.Priority
	it.Removed = w.Removed
	it.CreatedAt = w.CreatedAt
	return nil
}

func ItemsMarshal(items []Item) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	for _, v := range items {
		w := wrapItem{
			Id:          v.Id,
			CampaignId:  v.CampaignId,
			Name:        v.Name,
			Description: v.Description,
			Priority:    v.Priority,
			Removed:     v.Removed,
			CreatedAt:   v.CreatedAt,
		}
		if err := enc.Encode(w); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func ItemsUnmarshal(data []byte, items *[]Item) error {
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	for {
		w := wrapItem{}
		if err := dec.Decode(&w); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		*items = append(*items, Item{
			Id:          w.Id,
			CampaignId:  w.CampaignId,
			Name:        w.Name,
			Description: w.Description,
			Priority:    w.Priority,
			Removed:     w.Removed,
			CreatedAt:   w.CreatedAt,
		})
	}
	return nil
}
