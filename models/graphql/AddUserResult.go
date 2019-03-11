package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/quan-to/remote-signer/etc"
)

type AddUserResult struct {
	etc.BasicUser
	Password string
}

var GraphQLAddUserResult = graphql.NewObject(graphql.ObjectConfig{
	Name: "AddUserResult",
	Fields: graphql.Fields{
		"UserName": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the user this token belongs",
		},
		"FullName": &graphql.Field{
			Type:        graphql.String,
			Description: "Full name of the user",
		},
		"FingerPrint": &graphql.Field{
			Type:        graphql.String,
			Description: "FingerPrint of the key user has access",
		},
		"Password": &graphql.Field{
			Type:        graphql.String,
			Description: "Auto-generated password",
		},
	},
})
