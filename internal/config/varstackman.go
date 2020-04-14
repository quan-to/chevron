package config

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
		"ReadonlyKeyPath":           ReadonlyKeyPath,
		"VaultSkipVerify":           VaultSkipVerify,
		"VaultUseUserpass":          VaultUseUserpass,
		"VaultUsername":             VaultUsername,
		"VaultPassword":             VaultPassword,
		"VaultNamespace":            VaultNamespace,
		"VaultBackend":              VaultBackend,
		"VaultSkipDataType":         VaultSkipDataType,
		"AgentTargetURL":            AgentTargetURL,
		"AgentTokenExpiration":      AgentTokenExpiration,
		"AgentKeyFingerPrint":       AgentKeyFingerPrint,
		"AgentBypassLogin":          AgentBypassLogin,
		"RethinkTokenManager":       RethinkTokenManager,
		"RethinkAuthManager":        RethinkAuthManager,
		"Environment":               Environment,
		"AgentExternalURL":          AgentExternalURL,
		"AgentAdminExternalURL":     AgentAdminExternalURL,
		"OnDemandKeyLoad":           OnDemandKeyLoad,
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
	ReadonlyKeyPath = insMap["ReadonlyKeyPath"].(bool)
	VaultSkipVerify = insMap["VaultSkipVerify"].(bool)
	VaultUseUserpass = insMap["VaultUseUserpass"].(bool)
	VaultUsername = insMap["VaultUsername"].(string)
	VaultPassword = insMap["VaultPassword"].(string)
	VaultNamespace = insMap["VaultNamespace"].(string)
	VaultBackend = insMap["VaultBackend"].(string)
	VaultSkipDataType = insMap["VaultSkipDataType"].(bool)
	AgentTargetURL = insMap["AgentTargetURL"].(string)
	AgentTokenExpiration = insMap["AgentTokenExpiration"].(int)
	AgentKeyFingerPrint = insMap["AgentKeyFingerPrint"].(string)
	AgentBypassLogin = insMap["AgentBypassLogin"].(bool)
	RethinkTokenManager = insMap["RethinkTokenManager"].(bool)
	RethinkAuthManager = insMap["RethinkAuthManager"].(bool)
	Environment = insMap["Environment"].(string)
	AgentExternalURL = insMap["AgentExternalURL"].(string)
	AgentAdminExternalURL = insMap["AgentAdminExternalURL"].(string)
	OnDemandKeyLoad = insMap["OnDemandKeyLoad"].(bool)
}
