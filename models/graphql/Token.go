package graphql

import "github.com/quan-to/graphql"

type Token struct {
	Value                 string
	UserName              string
	UserFullName          string
	Expiration            int64
	ExpirationDateTimeISO string
}

var GraphQLToken = graphql.NewObject(graphql.ObjectConfig{
	Name: "Token",
	Fields: graphql.Fields{
		"Value": &graphql.Field{
			Type:        graphql.String,
			Description: "Token Value. Use this for all authenticated calls",
		},
		"UserName": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the user this token belongs",
		},
		"Expiration": &graphql.Field{
			Type:        graphql.Float,
			Description: "Unix Epoch Timestamp when this token expires (in milisseconds)",
		},
		"ExpirationDateTimeISO": &graphql.Field{
			Type:        graphql.String,
			Description: "ISO DateTime when this token expires",
		},
		"UserFullName": &graphql.Field{
			Type:        graphql.String,
			Description: "Full name of the user",
		},
	},
})
