package resolve

import (
	"log"

	"github.com/graphql-go/graphql"
	"github.com/tiny-sky/Tvote/core/dao"
	"github.com/tiny-sky/Tvote/core/dao/entity"
)

func Query(p graphql.ResolveParams) (interface{}, error) {
	name := p.Args["name"].(string)
	log.Println("query name is ", name)
	votes, err := dao.GetDB().GetVotesByName(name)
	if err != nil {
		return 0, err
	}
	return entity.User{Name: name, Votes: votes}, nil
}
