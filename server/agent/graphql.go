package agent

import (
	"github.com/quan-to/graphql"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/QuantoError"
	"github.com/quan-to/remote-signer/etc"
	mgql "github.com/quan-to/remote-signer/models/graphql"
	"net/http"
	"time"
)

const TokenManagerKey = "TokenManager"
const AuthManagerKey = "AuthManager"
const HTTPRequestKey = "HTTPRequest"

var RootManagementQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "ManagementQueries",
	Fields: graphql.Fields{
		"WhoAmI": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				req := p.Context.Value(HTTPRequestKey).(*http.Request)
				tm := p.Context.Value(TokenManagerKey).(etc.TokenManager)

				token := req.Header.Get("proxyToken")

				err := tm.Verify(token)
				if err != nil {
					e := QuantoError.New(QuantoError.InvalidFieldData, "proxyToken", "The specified proxyToken is either invalid or expired.", nil)
					return nil, e.ToFormattedError()
				}

				user := tm.GetUserData(token)

				return user.GetFullName(), nil
			},
		},
	},
})

var RootManagementMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "ManagementMutations",
	Fields: graphql.Fields{
		"Login": &graphql.Field{
			Type: mgql.GraphQLToken,
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Username to Login",
				},
				"password": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Password to Login",
				},
			},
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				tm := p.Context.Value(TokenManagerKey).(etc.TokenManager)
				am := p.Context.Value(AuthManagerKey).(etc.AuthManager)

				username := p.Args["username"].(string)
				password := p.Args["password"].(string)

				fingerPrint, fullname, err := am.LoginAuth(username, password)

				if err != nil {
					e := QuantoError.New(QuantoError.InvalidFieldData, "username/password", "Invalid username or password", nil)
					return nil, e.ToFormattedError()
				}

				createdAt := time.Now()
				exp := createdAt.Add(time.Second * time.Duration(remote_signer.AgentTokenExpiration))

				token := tm.AddUser(&etc.BasicUser{
					FingerPrint: fingerPrint,
					Username:    username,
					CreatedAt:   createdAt,
					FullName:    fullname,
				})

				return mgql.Token{
					Value:                 token,
					UserName:              username,
					Expiration:            exp.UnixNano() / 1e6, // ms
					ExpirationDateTimeISO: exp.Format(time.RFC3339),
				}, nil
			},
		},
		"ChangePassword": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"password": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "The new password",
				},
			},
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				req := p.Context.Value(HTTPRequestKey).(*http.Request)
				am := p.Context.Value(AuthManagerKey).(etc.AuthManager)
				tm := p.Context.Value(TokenManagerKey).(etc.TokenManager)

				token := req.Header.Get("proxyToken")

				err := tm.Verify(token)
				if err != nil {
					e := QuantoError.New(QuantoError.InvalidFieldData, "proxyToken", "The specified proxyToken is either invalid or expired.", nil)
					return "NOK", e.ToFormattedError()
				}

				password := p.Args["password"].(string)

				user := tm.GetUserData(token)
				err = am.ChangePassword(user.GetUsername(), password)

				if err != nil {
					e := QuantoError.New(QuantoError.InternalServerError, "server", "There was an error changing your password. Please try again.", nil)
					return "NOK", e.ToFormattedError()
				}

				return "OK", nil
			},
		},
	},
})
