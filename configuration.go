package remote_signer

import (
	"github.com/quan-to/remote-signer/SLog"
	"os"
	"strconv"
	"strings"
)

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
var VaultPathPrefix string
var ReadonlyKeyPath bool
var VaultSkipVerify bool

var varStack []map[string]interface{}

func PushVariables() {
	if varStack == nil {
		varStack = make([]map[string]interface{}, 0)
	}

	insMap := map[string]interface{}{
		"SyslogServer":              SyslogServer,
		"SyslogFacility":            SyslogFacility,
		"PrivateKeyFolder":          PrivateKeyFolder,
		"KeyPrefix":                 KeyPrefix,
		"SKSServer":                 SKSServer,
		"HttpPort":                  HttpPort,
		"MaxKeyRingCache":           MaxKeyRingCache,
		"EnableRethinkSKS":          EnableRethinkSKS,
		"RethinkDBHost":             RethinkDBHost,
		"RethinkDBPort":             RethinkDBPort,
		"RethinkDBUsername":         RethinkDBUsername,
		"RethinkDBPassword":         RethinkDBPassword,
		"RethinkDBPoolSize":         RethinkDBPoolSize,
		"DatabaseName":              DatabaseName,
		"MasterGPGKeyPath":          MasterGPGKeyPath,
		"MasterGPGKeyPasswordPath":  MasterGPGKeyPasswordPath,
		"MasterGPGKeyBase64Encoded": MasterGPGKeyBase64Encoded,
		"KeysBase64Encoded":         KeysBase64Encoded,
		"IgnoreKubernetesCA":        IgnoreKubernetesCA,
		"VaultAddress":              VaultAddress,
		"VaultRootToken":            VaultRootToken,
		"VaultStorage":              VaultStorage,
		"VaultPathPrefix":           VaultPathPrefix,
		"ReadonlyKeyPath":           ReadonlyKeyPath,
		"VaultSkipVerify":           VaultSkipVerify,
	}

	varStack = append(varStack, insMap)
}

func PopVariables() {
	if len(varStack) == 0 {
		return
	}

	insMap := varStack[len(varStack)-1]
	varStack = varStack[:len(varStack)-1]

	SyslogServer = insMap["SyslogServer"].(string)
	SyslogFacility = insMap["SyslogFacility"].(string)
	PrivateKeyFolder = insMap["PrivateKeyFolder"].(string)
	KeyPrefix = insMap["KeyPrefix"].(string)
	SKSServer = insMap["SKSServer"].(string)
	HttpPort = insMap["HttpPort"].(int)
	MaxKeyRingCache = insMap["MaxKeyRingCache"].(int)
	EnableRethinkSKS = insMap["EnableRethinkSKS"].(bool)
	RethinkDBHost = insMap["RethinkDBHost"].(string)
	RethinkDBPort = insMap["RethinkDBPort"].(int)
	RethinkDBUsername = insMap["RethinkDBUsername"].(string)
	RethinkDBPassword = insMap["RethinkDBPassword"].(string)
	RethinkDBPoolSize = insMap["RethinkDBPoolSize"].(int)
	DatabaseName = insMap["DatabaseName"].(string)
	MasterGPGKeyPath = insMap["MasterGPGKeyPath"].(string)
	MasterGPGKeyPasswordPath = insMap["MasterGPGKeyPasswordPath"].(string)
	MasterGPGKeyBase64Encoded = insMap["MasterGPGKeyBase64Encoded"].(bool)
	KeysBase64Encoded = insMap["KeysBase64Encoded"].(bool)
	IgnoreKubernetesCA = insMap["IgnoreKubernetesCA"].(bool)
	VaultAddress = insMap["VaultAddress"].(string)
	VaultRootToken = insMap["VaultRootToken"].(string)
	VaultStorage = insMap["VaultStorage"].(bool)
	VaultPathPrefix = insMap["VaultPathPrefix"].(string)
	ReadonlyKeyPath = insMap["ReadonlyKeyPath"].(bool)
	VaultSkipVerify = insMap["VaultSkipVerify"].(bool)
}

func Setup() {
	// Pre init
	MaxKeyRingCache = -1
	HttpPort = -1
	RethinkDBPort = -1
	RethinkDBPoolSize = -1

	// Load envvars
	SyslogServer = os.Getenv("SYSLOG_IP")
	SyslogFacility = os.Getenv("SYSLOG_FACILITY")
	PrivateKeyFolder = os.Getenv("PRIVATE_KEY_FOLDER")
	SKSServer = os.Getenv("SKS_SERVER")
	KeyPrefix = os.Getenv("KEY_PREFIX")

	var maxKeyRingCache = os.Getenv("MAX_KEYRING_CACHE_SIZE")
	if maxKeyRingCache != "" {
		i, err := strconv.ParseInt(maxKeyRingCache, 10, 32)
		if err != nil {
			SLog.Error("Error parsing MAX_KEYRING_CACHE_SIZE: %s", err)
			panic(err)
		}
		MaxKeyRingCache = int(i)
	}

	var hp = os.Getenv("HTTP_PORT")
	if hp != "" {
		i, err := strconv.ParseInt(hp, 10, 32)
		if err != nil {
			SLog.Error("Error parsing HTTP_PORT: %s", err)
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
			SLog.Error("Error parsing RETHINKDB_PORT: %s", err)
			panic(err)
		}
		RethinkDBPort = int(i)
	}

	var poolSize = os.Getenv("RETHINKDB_POOL_SIZE")
	if poolSize != "" {
		i, err := strconv.ParseInt(poolSize, 10, 32)
		if err != nil {
			SLog.Error("Error parsing RETHINKDB_POOL_SIZE: %s", err)
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
	VaultPathPrefix = os.Getenv("VAULT_PATH_PREFIX")
	ReadonlyKeyPath = os.Getenv("READONLY_KEYPATH") == "true"
	VaultSkipVerify = os.Getenv("VAULT_SKIP_VERIFY") == "true"

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

	// Other stuff
	_ = os.Mkdir(PrivateKeyFolder, 0770)
}

func init() {
	Setup()
}
