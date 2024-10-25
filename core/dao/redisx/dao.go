package redisx

import (
	"encoding/json"
	"errors"

	"time"

	"github.com/go-redis/redis"
	"github.com/tiny-sky/Tvote/core/dao/entity"
)

type RedisDao interface {
	GetUserByName(name string) (*entity.User, error)
	SetUser(votes *entity.User) error
	GetTicket() (*entity.Ticket, error)
	SetTicket(ticket *entity.Ticket) error
	UpdateTicket(ticket *entity.Ticket) error
}

func (r *Redis) GetUserByName(name string) (*entity.User, error) {
	val, err := r.Client.Get("user:" + name).Result()
	if err == redis.Nil {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	var user entity.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Redis) SetUser(user *entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = r.Client.Set("user:"+user.Name, data, 0).Result()
	return err
}

func (r *Redis) GetTicket() (*entity.Ticket, error) {
	val, err := r.Client.Get("ticket").Result()
	if err == redis.Nil {
		return nil, errors.New("ticket not found")
	} else if err != nil {
		return nil, err
	}

	var ticket entity.Ticket
	err = json.Unmarshal([]byte(val), &ticket)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *Redis) SetTicket(ticket *entity.Ticket) error {
	data, err := json.Marshal(ticket)
	if err != nil {
		return err
	}

	_, err = r.Client.Set("ticket", data, time.Duration(ticket.ExpiresAt-time.Now().Unix())*time.Second).Result()
	return err
}

func (r *Redis) UpdateTicket(ticket *entity.Ticket) error {
	ttl, err := r.Client.TTL("ticket").Result()
	if err != nil {
		return err
	}
	if ttl < 0 && ttl != -1 {
		return err
	}

	data, err := json.Marshal(ticket)
	if err != nil {
		return err
	}
	_, err = r.Client.Set("ticket", data, ttl).Result()
	if err != nil {
		return err
	}

	return nil
}
