package dao

import (
	"github.com/tiny-sky/Tvote/core/dao/entity"
	"github.com/tiny-sky/Tvote/log"
)

func (d *db) GetVotesByName(name string) (int, error) {
	// 首先尝试从 Redis 获取数据
	user, err := d.redisDao.GetUserByName(name)
	if err == nil && user != nil {
		return user.Votes, nil
	}

	// 如果 Redis 中没有，从 MySQL 获取
	user, err = d.mysqlDao.GetUserByName(name)
	if err != nil {
		return -1, err
	}

	// 将结果缓存回 Redis
	if err := d.redisDao.SetUser(user); err != nil {
		log.Warnf("Failed to cache votes for name %s: %v", name, err)
	}

	return user.Votes, nil
}

func (d *db) GetTicket() (*entity.Ticket, error) {
	ticket, err := d.redisDao.GetTicket()
	if err == nil && ticket != nil {
		return ticket, nil
	}

	return nil, err
}

func (d *db) GetUser(name string) (*entity.User, error) {
	user, err := d.redisDao.GetUserByName(name)
	if err == nil && user != nil {
		return user, nil
	}

	user, err = d.mysqlDao.GetUserByName(name)
	if err != nil {
		return nil, err
	}

	if err := d.redisDao.SetUser(user); err != nil {
		log.Warnf("Failed to cache user %s: %v", name, err)
	}

	return user, nil
}

func (d *db) CreateTicket(ticket *entity.Ticket) error {
	// 保存到 Redis
	if err := d.redisDao.SetTicket(ticket); err != nil {
		log.Warnf("Failed to cache ticket: %v", err)
	}

	return nil
}

func (d *db) CreateUserByName(name string) error {
	user := &entity.User{
		Name:  name,
		Votes: 0,
	}

	if err := d.mysqlDao.SetUser(user); err != nil {
		return err
	}

	if err := d.redisDao.SetUser(user); err != nil {
		log.Warnf("Failed to cache user %s: %v", name, err)
	}

	return nil
}

func (d *db) UpdateVotes(users []*entity.User) error {
	// 更新 MySQL
	err := d.mysqlDao.UpdateVotesByNames(users)
	if err != nil {
		return err
	}

	// 更新 Redis 缓存
	for _, user := range users {
		if err := d.redisDao.SetUser(user); err != nil {
			log.Warnf("Failed to cache user %s: %v", user.Name, err)
		}
	}

	return nil
}

func (d *db) AddUsageByTicket(ticket *entity.Ticket) error {
	ticket.UsedCount++
	if err := d.redisDao.UpdateTicket(ticket); err != nil {
		log.Warnf("Failed to update cache for ticket: %v", err)
	}

	return nil
}
