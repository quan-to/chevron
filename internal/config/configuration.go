package config

import (
	"github.com/quan-to/chevron/pkg/QuantoError"
	"github.com/quan-to/slog"
	"os"
	"strconv"
	"strings"
)

const SMEncryptedDataOnly = false

var SyslogServer string
var SyslogFacility string
var PrivateKeyFolder string
var KeyPrefix string
var SKSServer string
var HttpPort int
var MaxKeyRingCache int
var EnableRethinkSKS bool
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

var RethinkTokenManager bool
var RethinkAuthManager bool
var Environment string

var AgentExternalURL string
var AgentAdminExternalURL string
var ShowLines bool

// LogFormat allows to configure the output log format
var LogFormat slog.Format

func Setup() {
	// Pre init
	MaxKeyRingCache = -1
	HttpPort = -1
	RethinkDBPort = -1
	RethinkDBPoolSize = -1
	AgentTokenExpiration = -1
	ShowLines = false

	// Load envvars
	SyslogServer = os.Getenv("SYSLOG_IP")
	SyslogFacility = os.Getenv("SYSLOG_FACILITY")
	PrivateKeyFolder = os.Getenv("PRIVATE_KEY_FOLDER")
	SKSServer = os.Getenv("SKS_SERVER")
	KeyPrefix = os.Getenv("KEY_PREFIX")
	ShowLines = os.Getenv("SHOW_LINES") == "true"
	LogFormat = slog.ToFormat(os.Getenv("LOG_FORMAT"))

	if ShowLines {
		slog.SetShowLines(true)
	}

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

	EnableRethinkSKS = strings.ToLower(os.Getenv("ENABLE_RETHINKDB_SKS")) == "true"

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
	RethinkTokenManager = os.Getenv("RETHINK_TOKEN_MANAGER") == "true"
	RethinkAuthManager = os.Getenv("RETHINK_AUTH_MANAGER") == "true"

	if (RethinkAuthManager || RethinkTokenManager) && !EnableRethinkSKS {
		slog.Fatal("Rethink Auth / Token Manager requires Rethink SKS")
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

	// Other stuff
	_ = os.Mkdir(PrivateKeyFolder, 0770)

	if Environment == "development" {
		slog.SetDebug(true)
		QuantoError.EnableStackTrace()
	} else {
		slog.SetDebug(false)
		QuantoError.DisableStackTrace()
	}
}

func init() {
	Setup()
}
