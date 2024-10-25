package dao

import (
	"github.com/tiny-sky/Tvote/core/dao/entity"
	"github.com/tiny-sky/Tvote/core/dao/mysqlx"
	"github.com/tiny-sky/Tvote/core/dao/redisx"
)

type db struct {
	redisDao redisx.RedisDao
	mysqlDao mysqlx.MysqlDao
}

type DB interface {
	GetVotesByName(name string) (int, error)
	GetTicket() (*entity.Ticket, error)
	GetUser(name string) (*entity.User, error)
	CreateTicket(ticket *entity.Ticket) error
	CreateUserByName(name string) error
	UpdateVotes(users []*entity.User) error
	AddUsageByTicket(ticket *entity.Ticket) error
}

var dao db

func NewDB() {
	dao = db{
		redisDao: &redisx.Rdb,
		mysqlDao: &mysqlx.Mydb,
	}
}

func GetDB() DB {
	return &dao
}
