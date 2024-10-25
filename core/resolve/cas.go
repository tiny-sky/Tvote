package resolve

import (
	"github.com/graphql-go/graphql"
	"github.com/tiny-sky/Tvote/core/dao"
)

func Cas(p graphql.ResolveParams) (interface{}, error) {
	currentTicket, err := dao.GetDB().GetTicket()
	if err != nil {
		return nil, err
	}
	return currentTicket.Ticket, nil
}
