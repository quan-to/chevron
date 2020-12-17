package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/slog"
)

const SMEncryptedDataOnly = false

var SyslogServer string
var SyslogFacility string
var PrivateKeyFolder string
var KeyPrefix string
var SKSServer string
var HttpPort int
var MaxKeyRingCache int
var EnableDatabase bool
var RethinkDBHost string
var RethinkDBPort int
var RethinkDBUsername string
var RethinkDBPassword string
var RethinkDBPoolSize int
var DatabaseName string
var MasterGPGKeyPath string
var MasterGPGKeyPasswordPath string
var MasterGPGKeyBase64Encoded bool
var KeysBase64Encoded bool
var IgnoreKubernetesCA bool
var VaultStorage bool
var VaultAddress string
var VaultRootToken string
var ReadonlyKeyPath bool
var VaultSkipVerify bool
var VaultUseUserpass bool
var VaultUsername string
var VaultPassword string
var VaultNamespace string
var VaultBackend string
var VaultSkipDataType bool
var VaultTokenTTL string
var AgentTargetURL string
var AgentTokenExpiration int
var AgentKeyFingerPrint string
var AgentBypassLogin bool
var OnDemandKeyLoad bool
var RequestIDHeader string

var DatabaseTokenManager bool
var DatabaseAuthManager bool
var Environment string

var AgentExternalURL string
var AgentAdminExternalURL string
var ShowLines bool
var SingleKeyMode bool
var SingleKeyPath string
var SingleKeyPassword string
var DatabaseDialect string
var ConnectionString string

var EnableRedis bool
var RedisHost string
var RedisUser string
var RedisPass string
var RedisDatabaseIndex int
var RedisTLSEnabled bool
var RedisMaxLocalObjects int
var RedisLocalObjectTTL time.Duration

// LogFormat allows to configure the output log format
var LogFormat slog.Format

func configDeprecationMessage(userConfig, newConfig string) {
	if newConfig != "" {
		slog.Warn("The configuration %q is currently deprecated. Please use %q instead.", userConfig, newConfig)
	} else {
		slog.Warn("The configuration %q is currently deprecated. Please refer to https://github.com/quan-to/chevron to more information")
	}
}

func Setup() {
	var err error
	// Pre init
	MaxKeyRingCache = -1
	HttpPort = -1
	RethinkDBPort = -1
	RethinkDBPoolSize = -1
	AgentTokenExpiration = -1
	ShowLines = false
	DatabaseDialect = ""
	ConnectionString = ""

	// Load envvars
	DatabaseDialect = strings.ToLower(os.Getenv("DATABASE_DIALECT"))
	ConnectionString = os.Getenv("CONNECTION_STRING")
	SyslogServer = os.Getenv("SYSLOG_IP")
	SyslogFacility = os.Getenv("SYSLOG_FACILITY")
	PrivateKeyFolder = os.Getenv("PRIVATE_KEY_FOLDER")
	SKSServer = os.Getenv("SKS_SERVER")
	KeyPrefix = os.Getenv("KEY_PREFIX")
	ShowLines = os.Getenv("SHOW_LINES") == "true"
	LogFormat = slog.ToFormat(os.Getenv("LOG_FORMAT"))
	slog.SetShowLines(ShowLines)

	var maxKeyRingCache = os.Getenv("MAX_KEYRING_CACHE_SIZE")
	if maxKeyRingCache != "" {
		i, err := strconv.ParseInt(maxKeyRingCache, 10, 32)
		if err != nil {
			slog.Error("Error parsing MAX_KEYRING_CACHE_SIZE: %s", err)
			panic(err)
		}
		MaxKeyRingCache = int(i)
	}

	var hp = os.Getenv("HTTP_PORT")
	if hp != "" {
		i, err := strconv.ParseInt(hp, 10, 32)
		if err != nil {
			slog.Error("Error parsing HTTP_PORT: %s", err)
			panic(err)
		}
		HttpPort = int(i)
	}

	EnableDatabase = strings.ToLower(os.Getenv("ENABLE_DATABASE")) == "true" || DatabaseDialect == "rethinkdb"

	if strings.ToLower(os.Getenv("ENABLE_RETHINKDB_SKS")) == "true" || DatabaseDialect == "rethinkdb" {
		slog.Warn("RethinkDB services are currently deprecated. Please check the project README for more information: https://github.com/quan-to/chevron")
		// Backwards compatibility
		DatabaseDialect = "rethinkdb"
		EnableDatabase = true
		RethinkDBHost = os.Getenv("RETHINKDB_HOST")
		RethinkDBUsername = os.Getenv("RETHINKDB_USERNAME")
		RethinkDBPassword = os.Getenv("RETHINKDB_PASSWORD")

		var rdbport = os.Getenv("RETHINKDB_PORT")
		if rdbport != "" {
			i, err := strconv.ParseInt(rdbport, 10, 32)
			if err != nil {
				slog.Error("Error parsing RETHINKDB_PORT: %s", err)
				panic(err)
			}
			RethinkDBPort = int(i)
		}

		var poolSize = os.Getenv("RETHINKDB_POOL_SIZE")
		if poolSize != "" {
			i, err := strconv.ParseInt(poolSize, 10, 32)
			if err != nil {
				slog.Error("Error parsing RETHINKDB_POOL_SIZE: %s", err)
				panic(err)
			}
			RethinkDBPoolSize = int(i)
		}
		DatabaseName = os.Getenv("DATABASE_NAME")
	}

	MasterGPGKeyPath = os.Getenv("MASTER_GPG_KEY_PATH")
	MasterGPGKeyPasswordPath = os.Getenv("MASTER_GPG_KEY_PASSWORD_PATH")
	MasterGPGKeyBase64Encoded = strings.ToLower(os.Getenv("MASTER_GPG_KEY_BASE64_ENCODED")) == "true"

	KeysBase64Encoded = strings.ToLower(os.Getenv("KEYS_BASE64_ENCODED")) == "true"
	IgnoreKubernetesCA = strings.ToLower(os.Getenv("IGNORE_KUBERNETES_CA")) == "true"

	VaultStorage = strings.ToLower(os.Getenv("VAULT_STORAGE")) == "true"
	VaultAddress = os.Getenv("VAULT_ADDRESS")
	VaultRootToken = os.Getenv("VAULT_ROOT_TOKEN")
	ReadonlyKeyPath = os.Getenv("READONLY_KEYPATH") == "true"
	VaultSkipVerify = os.Getenv("VAULT_SKIP_VERIFY") == "true"
	VaultUseUserpass = os.Getenv("VAULT_USE_USERPASS") == "true"
	VaultUsername = os.Getenv("VAULT_USERNAME")
	VaultPassword = os.Getenv("VAULT_PASSWORD")
	VaultNamespace = os.Getenv("VAULT_NAMESPACE")
	VaultBackend = os.Getenv("VAULT_BACKEND")
	VaultSkipDataType = os.Getenv("VAULT_SKIP_DATA_TYPE") == "true"
	VaultTokenTTL = os.Getenv("VAULT_TOKEN_TTL")
	AgentTargetURL = os.Getenv("AGENT_TARGET_URL")
	AgentKeyFingerPrint = os.Getenv("AGENT_KEY_FINGERPRINT")
	AgentBypassLogin = os.Getenv("AGENT_BYPASS_LOGIN") == "true"
	DatabaseTokenManager = os.Getenv("RETHINK_TOKEN_MANAGER") == "true" || os.Getenv("DATABASE_TOKEN_MANAGER") == "true"
	DatabaseAuthManager = os.Getenv("RETHINK_AUTH_MANAGER") == "true" || os.Getenv("AUTH_MANAGER") == "true"

	if (DatabaseAuthManager || DatabaseTokenManager) && !EnableDatabase {
		slog.Fatal("Database Auth / Token Manager requires a database configuration")
	}

	RequestIDHeader = os.Getenv("REQUESTID_HEADER")
	AgentExternalURL = os.Getenv("AGENT_EXTERNAL_URL")
	AgentAdminExternalURL = os.Getenv("AGENTADMIN_EXTERNAL_URL")

	Environment = os.Getenv("Environment")

	atke := os.Getenv("AGENT_TOKEN_EXPIRATION")

	if atke != "" {
		i, err := strconv.ParseInt(atke, 10, 32)
		if err != nil {
			slog.Fatal("Error parsing AGENT_TOKEN_EXPIRATION: %s", err)
		}
		AgentTokenExpiration = int(i)
	}

	OnDemandKeyLoad = os.Getenv("ON_DEMAND_KEY_LOAD") == "true"

	SingleKeyMode = os.Getenv("MODE") == "single_key"
	SingleKeyPath = os.Getenv("SINGLE_KEY_PATH")
	SingleKeyPassword = os.Getenv("SINGLE_KEY_PASSWORD")

	EnableRedis = os.Getenv("REDIS_ENABLE") == "true"
	RedisTLSEnabled = os.Getenv("REDIS_TLS_ENABLED") == "true"
	RedisHost = os.Getenv("REDIS_HOST")
	RedisUser = os.Getenv("REDIS_USER")
	RedisPass = os.Getenv("REDIS_PASS")

	redisDBIdx := os.Getenv("REDIS_DATABASE_INDEX")
	RedisDatabaseIndex = 0
	if redisDBIdx != "" {
		v, err := strconv.ParseInt(redisDBIdx, 10, 32)
		if err != nil {
			slog.Error("Invalid field REDIS_DATABASE_INDEX = %q - Invalid number", redisDBIdx)
			v = 0
		}
		RedisDatabaseIndex = int(v)
	}

	redisLocalObjectTTL := os.Getenv("REDIS_MAX_LOCAL_TTL")
	if redisLocalObjectTTL != "" {
		if RedisLocalObjectTTL, err = time.ParseDuration(redisLocalObjectTTL); err != nil {
			slog.Error("Invalid field REDIS_MAX_LOCAL_TTL = %q - Invalid Duration")
		}
	}

	redisMaxLocalObjects := os.Getenv("REDIS_MAX_LOCAL_OBJECTS")
	if redisMaxLocalObjects != "" {
		v, err := strconv.ParseInt(redisMaxLocalObjects, 10, 32)
		if err != nil {
			slog.Error("Invalid field REDIS_MAX_LOCAL_OBJECTS = %q - Invalid number", redisMaxLocalObjects)
			v = 0
		}
		RedisMaxLocalObjects = int(v)
	}

	// Set defaults if not defined
	if SyslogServer == "" {
		SyslogServer = "127.0.0.1"
	}

	if SyslogFacility == "" {
		SyslogFacility = "LOG_USER"
	}

	if PrivateKeyFolder == "" {
		PrivateKeyFolder = "./keys"
	}

	if MaxKeyRingCache == -1 {
		MaxKeyRingCache = 1000
	}

	if HttpPort == -1 {
		HttpPort = 5100
	}

	if RethinkDBHost == "" {
		RethinkDBHost = "127.0.0.1"
	}

	if VaultTokenTTL == "" {
		VaultTokenTTL = "768h"
	}

	if RethinkDBUsername == "" {
		RethinkDBUsername = "admin"
	}

	if RethinkDBPort == -1 {
		RethinkDBPort = 28015
	}

	if RethinkDBPoolSize == -1 {
		RethinkDBPoolSize = 10
	}

	if DatabaseName == "" {
		DatabaseName = "remote_signer"
	}

	if VaultAddress == "" {
		VaultAddress = "http://localhost:8200"
	}

	if VaultNamespace == "" {
		VaultNamespace = "remote-signer"
	}

	if VaultBackend == "" {
		VaultBackend = "secret"
	}

	if AgentTargetURL == "" {
		AgentTargetURL = "https://api.sandbox.contaquanto.com/all"
	}

	if AgentTokenExpiration == -1 {
		AgentTokenExpiration = 3600
	}

	if Environment == "" {
		Environment = "development"
	}

	if AgentExternalURL == "" {
		AgentExternalURL = "/agent"
	}

	if AgentAdminExternalURL == "" {
		AgentAdminExternalURL = "/agentAdmin"
	}

	if RedisMaxLocalObjects == 0 {
		RedisMaxLocalObjects = 100
	}

	if RedisLocalObjectTTL == 0 {
		RedisLocalObjectTTL = time.Minute * 5
	}

	if RedisHost == "" {
		RedisHost = "localhost:6379"
	}

	// Other stuff
	_ = os.Mkdir(PrivateKeyFolder, 0750)

	if Environment == "development" {
		slog.SetDebug(true)
		QuantoError.EnableStackTrace()
	} else {
		slog.SetDebug(false)
		QuantoError.DisableStackTrace()
	}

	// Deprecation
	if os.Getenv("RETHINK_TOKEN_MANAGER") != "" {
		configDeprecationMessage("RETHINK_TOKEN_MANAGER", "DATABASE_TOKEN_MANAGER")
	}
	if os.Getenv("RETHINK_AUTH_MANAGER") != "" {
		configDeprecationMessage("RETHINK_AUTH_MANAGER", "DATABASE_AUTH_MANAGER")
	}
	if os.Getenv("ENABLE_RETHINKDB_SKS") != "" {
		configDeprecationMessage("ENABLE_RETHINKDB_SKS", "DATABASE_DIALECT=rethinkdb")
	}
}

func init() {
	Setup()
}
