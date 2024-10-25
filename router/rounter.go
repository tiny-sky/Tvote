package router

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/tiny-sky/Tvote/core/resolve"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"votes": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

func getRootQueryFields() graphql.Fields {
	return graphql.Fields{
		"query": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: resolve.Query,
		},
		"cas": &graphql.Field{
			Name:    "Query",
			Type:    graphql.String,
			Resolve: resolve.Cas,
		},
	}
}

func getMutationFields() graphql.Fields {
	return graphql.Fields{
		"vote": &graphql.Field{
			Name: "Vote",
			Type: graphql.NewList(userType),
			Args: graphql.FieldConfigArgument{
				"names": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
				},
				"ticket": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: resolve.Vote,
		},
	}
}

func CreateSchema() (graphql.Schema, error) {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: getRootQueryFields(),
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: getMutationFields(),
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: mutation,
	})
}
