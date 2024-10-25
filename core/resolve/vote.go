package resolve

import (
	"errors"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/tiny-sky/Tvote/core/dao"
	"github.com/tiny-sky/Tvote/core/dao/entity"
)

// 定义User类型
type User struct {
	Name  string
	Votes int
}

func Vote(p graphql.ResolveParams) (interface{}, error) {
	namesParam := p.Args["names"].([]interface{})
	ticket := p.Args["ticket"].(string)

	currentTicket, err := dao.GetDB().GetTicket()
	if err != nil {
		return nil, errors.New("failed to retrieve ticket")
	}

	if err = checksafe(currentTicket, ticket); err != nil {
		return nil, errors.New("illegal ticket")
	}

	if err = dao.GetDB().AddUsageByTicket(currentTicket); err != nil {
		return nil, errors.New("add ticket's usage failed")
	}

	var result []*entity.User

	for _, nameparam := range namesParam {
		name, _ := nameparam.(string)

		user, err := dao.GetDB().GetUser(name)
		if err != nil {
			return nil, errors.New("failed to get user")
		}
		user.Votes += 1
		result = append(result, user)
	}

	err = dao.GetDB().UpdateVotes(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func checksafe(curticket *entity.Ticket, ticket string) error {
	if ticket != curticket.Ticket {
		return errors.New("invalid ticket")
	}

	currentTime := time.Now().Unix()
	if currentTime > curticket.ExpiresAt {
		return errors.New("ticket has expired")
	}

	if curticket.UsedCount >= curticket.MaxUsage {
		return errors.New("ticket usage limit reached")
	}

	return nil
}
