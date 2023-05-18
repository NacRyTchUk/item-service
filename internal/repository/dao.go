package repository

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type DAO interface {
	NewItemQuery() ItemQuery
	NewCampaignQuery() CampaignQuery
}

type dao struct {
	Cache *redis.Client
}

var db *sql.DB

func pgQb() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(db)
}

func NewDAO(cache *redis.Client) DAO {
	return &dao{
		Cache: cache,
	}
}

func NewDB() (*sql.DB, error) {
	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		viper.Get("database-user").(string),
		viper.Get("database-password").(string),
		viper.Get("database-host").(string),
		viper.Get("database-port").(int),
		viper.Get("database-dbname").(string))
	conn, err := sql.Open("pgx", cs)
	if err != nil {
		return nil, err
	}
	db = conn
	return conn, nil
}

func (d dao) NewItemQuery() ItemQuery {
	return itemQuery{
		db:  db,
		dao: d,
	}
}

func (d dao) NewCampaignQuery() CampaignQuery {
	return campaignQuery{db: db}
}
