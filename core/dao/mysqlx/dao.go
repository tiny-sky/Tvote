package mysqlx

import (
	"context"

	"github.com/tiny-sky/Tvote/core/dao/entity"
	"github.com/tiny-sky/Tvote/log"
	"gorm.io/gorm"
)

type MysqlDao interface {
	GetUserByName(name string) (*entity.User, error)
	GetTicket() (*entity.Ticket, error)
	SetUser(user *entity.User) error
	UpdateVotesByNames(users []*entity.User) error
	AddUsageByTicket(ticket *entity.Ticket) error
}

func (d *Mysql) GetUserByName(name string) (*entity.User, error) {
	user := d.query.User
	u, err := d.query.User.WithContext(context.Background()).Where(user.Name.Eq(name)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newUser := &entity.User{
				Name:  name,
				Votes: 0,
			}
			if createErr := d.query.User.WithContext(context.Background()).Create(newUser); createErr != nil {
				return nil, createErr
			}
			return newUser, nil
		}
		return nil, err
	}

	return u, nil
}

func (d *Mysql) GetTicket() (*entity.Ticket, error) {
	ticket, err := d.query.Ticket.WithContext(context.Background()).First()
	if err != nil {
		return &entity.Ticket{}, err
	}
	if ticket == nil {
		return &entity.Ticket{}, err
	}

	return ticket, nil
}

func (d *Mysql) SetUser(user *entity.User) error {
	err := d.query.User.WithContext(context.Background()).Create(user)
	if err != nil {
		log.Errorf("create user fail, err:%v\n", err)
		return err
	}
	return nil
}

func (d *Mysql) UpdateVotesByNames(users []*entity.User) error {
	var err error

	User := d.query.User
	for _, user := range users {
		if user == nil {
			continue
		}

		_, err = d.query.User.WithContext(context.Background()).Where(User.ID.Eq(user.ID)).
			UpdateSimple(User.Votes.Add(1))
		if err != nil {
			break
		}
	}
	return err
}

func (d *Mysql) AddUsageByTicket(ticket *entity.Ticket) error {
	Ticket := d.query.Ticket
	_, err := d.query.Ticket.WithContext(context.Background()).Where(Ticket.ID.Eq(ticket.ID)).
		UpdateSimple(Ticket.UsedCount.Add(1))
	return err
}
