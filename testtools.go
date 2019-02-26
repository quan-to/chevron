package remote_signer

import (
	"fmt"
	"github.com/quan-to/remote-signer/SLog"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"time"
)

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

// RQLStart please don't judge me. That its what it takes to avoid RethinkDB Non-Atomic Database operations on tests :(
func RQLStart() (*exec.Cmd, error) {
	genPort := 28020 + rand.Intn(100)
	_ = os.RemoveAll("./rethinkdb_data")
	connString := fmt.Sprintf("127.0.0.1:%d", genPort)
	conn, err := net.DialTimeout("tcp", connString, time.Millisecond*500)
	if err == nil {
		// Retry another port
		_ = conn.Close()
		return RQLStart()
	}

	SLog.Info("Starting RethinkDB")
	cmd := exec.Command("rethinkdb", "--no-http-admin", "--bind", "127.0.0.1", "--cluster-port", fmt.Sprintf("%d", genPort+200), "--driver-port", fmt.Sprintf("%d", genPort))
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	err = cmd.Start()

	if err != nil {
		b, _ := cmd.Output()
		SLog.Error(string(b))
		panic(fmt.Errorf("cannot start rethinkdb: %s", err))
	}

	RethinkDBPort = genPort

	SLog.Info("Waiting for rethink to settle")
	time.Sleep(time.Second)
	retry := 0
	started := false

	for retry < 5 {
		conn, err := net.DialTimeout("tcp", connString, time.Millisecond*500)
		if err != nil {
			SLog.Error(err)
			retry++
			time.Sleep(time.Second)
			continue
		}
		_ = conn.Close()
		started = true
		break
	}

	if !started {
		return nil, fmt.Errorf("timeout waiting for rethinkdb to start")
	}

	return cmd, nil
}

func RQLStop(cmd *exec.Cmd) {
	SLog.Info("Stopping RethinkDB")
	err := cmd.Process.Signal(os.Kill)
	if err != nil {
		panic(fmt.Errorf("error killing rethinkdb: %s", err))
	}
	err = cmd.Process.Kill()
	if err != nil {
		panic(fmt.Errorf("error killing rethinkdb: %s", err))
	}

	err = os.RemoveAll("./rethinkdb_data")
	if err != nil {
		panic(fmt.Errorf("error erasing rethinkdb folder: %s", err))
	}

	_, err = cmd.Process.Wait()
	if err != nil {
		panic(fmt.Errorf("error waiting rethinkdb: %s", err))
	}

	time.Sleep(10 * time.Second)
}
