package agent

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"
	mgql "github.com/quan-to/chevron/pkg/models/graphql"
	"github.com/quan-to/slog"
	"time"
)

const TokenManagerKey = "TokenManager"
const AuthManagerKey = "AuthManager"
const HTTPRequestKey = "HTTPRequest"
const LoggedUserKey = "LoggerUser"

var amGqlLog = slog.Scope("Agent-GQL")

var RootManagementQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "ManagementQueries",
	Fields: graphql.Fields{
		"WhoAmI": &graphql.Field{
			Type:    graphql.String,
			Resolve: resolveWhoAmI,
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
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Username to Login",
				},
				"password": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Password to Login",
				},
				"expiresAfter": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Number of seconds since creation when the generated token will expire. If 0, defaults to server default.",
				},
			},
			Resolve: resolveLogin,
		},
		"AddUser": &graphql.Field{
			Type: mgql.GraphQLAddUserResult,
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Login of the new user",
				},
				"fullname": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The Full Name of the new user",
				},
				"fingerPrint": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "The fingerPrint that this user will use. Defaults to server Default",
				},
			},
			Resolve: resolveAddUser,
		},
		"ChangePassword": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"password": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The new password",
				},
			},
			Resolve: resolveChangePassword,
		},
		"GenerateToken": &graphql.Field{
			Type: mgql.GraphQLToken,
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Optional username to inform the user. It doesn't create anything and it isn't validated. Its just for info purpose.",
				},
				"fullname": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Optional fullname to inform the user. It doesn't create anything and it isn't validated. Its just for info purpose.",
				},
				"fingerPrint": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Fingerprint of the key to give access to. Defaults to Agent Default",
				},
				"expiresAfter": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Number of seconds since creation when the generated token will expire. If 0, defaults to server default.",
				},
			},
			Resolve: resolveGenerateToken,
		},
		"InvalidateToken": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"token": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The token to be invalidated",
				},
			},
			Resolve: resolveInvalidateToken,
		},
	},
})

func resolveWhoAmI(p graphql.ResolveParams) (i interface{}, e error) {
	lu := p.Context.Value(LoggedUserKey).(interfaces.UserData)

	if lu == nil {
		e := QuantoError.New(QuantoError.PermissionDenied, "proxyToken", "You need to be logged in to use this query", nil)
		return nil, e.ToFormattedError()
	}

	return lu.GetFullName(), nil
}

func resolveLogin(p graphql.ResolveParams) (i interface{}, e error) {
	tm := p.Context.Value(TokenManagerKey).(interfaces.TokenManager)
	am := p.Context.Value(AuthManagerKey).(interfaces.AuthManager)

	username := p.Args["username"].(string)
	password := p.Args["password"].(string)

	fingerPrint, fullname, err := am.LoginAuth(username, password)

	if err != nil {
		e := QuantoError.New(QuantoError.InvalidFieldData, "username/password", "Invalid username or password", nil)
		return nil, e.ToFormattedError()
	}

	createdAt := time.Now()

	expTime := config.AgentTokenExpiration
	exp := createdAt.Add(time.Second * time.Duration(expTime))

	if p.Args["expiresAfter"] != nil {
		expTime = p.Args["expiresAfter"].(int)
		exp = createdAt.Add(time.Second * time.Duration(expTime))
	}

	token := tm.AddUserWithExpiration(&models.BasicUser{
		FingerPrint: fingerPrint,
		Username:    username,
		CreatedAt:   createdAt,
		FullName:    fullname,
	}, expTime)

	return mgql.Token{
		Value:                 token,
		UserName:              username,
		Expiration:            exp.UnixNano() / 1e6, // ms
		ExpirationDateTimeISO: exp.Format(time.RFC3339),
	}, nil
}

func resolveAddUser(p graphql.ResolveParams) (i interface{}, e error) {
	var username, fullname, fingerPrint, password string

	lu := p.Context.Value(LoggedUserKey).(interfaces.UserData)
	if lu == nil {
		e := QuantoError.New(QuantoError.PermissionDenied, "proxyToken", "You need to be logged in to use this query", nil)
		return nil, e.ToFormattedError()
	}

	if lu.GetUsername() != "admin" {
		e := QuantoError.New(QuantoError.PermissionDenied, "username", "Only the administrator can add users", nil)
		return nil, e.ToFormattedError()
	}

	am := p.Context.Value(AuthManagerKey).(interfaces.AuthManager)

	username = p.Args["username"].(string)
	fullname = p.Args["fullname"].(string)

	if p.Args["fingerPrint"] != nil {
		fingerPrint = p.Args["fingerPrint"].(string)
	} else {
		fingerPrint = config.AgentKeyFingerPrint
	}

	password = tools.GeneratePassword()

	err := am.LoginAdd(username, password, fullname, fingerPrint)
	if err != nil {
		e := QuantoError.New(QuantoError.InternalServerError, "server", "There was an error adding the user. Please try again.", err.Error())
		return nil, e.ToFormattedError()
	}

	amGqlLog.Info("Added new user %s (%s)", fullname, username)

	return mgql.AddUserResult{
		BasicUser: models.BasicUser{
			Username:    username,
			FullName:    fullname,
			FingerPrint: fingerPrint,
		},
		Password: password,
	}, nil
}

func resolveChangePassword(p graphql.ResolveParams) (i interface{}, e error) {
	lu := p.Context.Value(LoggedUserKey).(interfaces.UserData)
	if lu == nil {
		e := QuantoError.New(QuantoError.PermissionDenied, "proxyToken", "You need to be logged in to use this query", nil)
		return nil, e.ToFormattedError()
	}

	am := p.Context.Value(AuthManagerKey).(interfaces.AuthManager)
	password := p.Args["password"].(string)

	err := am.ChangePassword(lu.GetUsername(), password)

	if err != nil {
		amGqlLog.Error("Error changing user %s password: %s", lu.GetUsername(), err)
		e := QuantoError.New(QuantoError.InternalServerError, "server", "There was an error changing your password. Please try again.", err.Error())
		return "NOK", e.ToFormattedError()
	}

	amGqlLog.Info("Changed Password for (%s)", lu.GetFullName(), lu.GetUsername())

	return "OK", nil
}

func resolveGenerateToken(p graphql.ResolveParams) (i interface{}, e error) {
	var username, fullname, fingerPrint string
	lu := p.Context.Value(LoggedUserKey).(interfaces.UserData)
	if lu == nil {
		e := QuantoError.New(QuantoError.PermissionDenied, "proxyToken", "You need to be logged in to use this query", nil)
		return nil, e.ToFormattedError()
	}

	if lu.GetUsername() != "admin" {
		e := QuantoError.New(QuantoError.PermissionDenied, "username", "Only the administrator can add users", nil)
		return nil, e.ToFormattedError()
	}

	tm := p.Context.Value(TokenManagerKey).(interfaces.TokenManager)

	if p.Args["username"] != nil {
		username = p.Args["username"].(string)
	} else {
		u, _ := uuid.NewRandom()
		username = u.String()
	}

	if p.Args["fullname"] != nil {
		fullname = p.Args["fullname"].(string)
	} else {
		fullname = username
	}

	if p.Args["fingerPrint"] != nil {
		fingerPrint = p.Args["fingerPrint"].(string)
	} else {
		fingerPrint = config.AgentKeyFingerPrint
	}

	expiration := 0

	if p.Args["expiresAfter"] != nil {
		expiration = p.Args["expiresAfter"].(int)
	}

	if expiration == 0 {
		expiration = config.AgentTokenExpiration
	}

	bu := models.BasicUser{
		FingerPrint: fingerPrint,
		Username:    username,
		FullName:    fullname,
		CreatedAt:   time.Now(),
	}

	exp := bu.GetCreatedAt().Add(time.Duration(expiration) * time.Second)

	amGqlLog.Await("Creating Token for key %s with expiration at %s", fingerPrint, exp.Format(time.RFC3339))

	token := tm.AddUserWithExpiration(&bu, expiration)

	amGqlLog.Done("Generated Token for %s (%s)", fullname, username)

	return mgql.Token{
		Value:                 token,
		UserName:              username,
		UserFullName:          fullname,
		Expiration:            exp.UnixNano() / 1e6, // ms
		ExpirationDateTimeISO: exp.Format(time.RFC3339),
	}, nil
}

func resolveInvalidateToken(p graphql.ResolveParams) (i interface{}, e error) {
	tm := p.Context.Value(TokenManagerKey).(interfaces.TokenManager)
	token := p.Args["token"].(string)

	err := tm.InvalidateToken(token)
	if err != nil {
		return "NOK", err
	}

	return "OK", nil
}
