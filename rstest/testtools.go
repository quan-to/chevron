package rstest

import (
	"fmt"
	"github.com/quan-to/slog"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"time"
)

// RQLStart please don't judge me. That its what it takes to avoid RethinkDB Non-Atomic Database operations on tests :(
func RQLStart() (*exec.Cmd, int, error) {
	genPort := 28020 + rand.Intn(100)
	_ = os.RemoveAll("./rethinkdb_data")
	connString := fmt.Sprintf("127.0.0.1:%d", genPort)
	conn, err := net.DialTimeout("tcp", connString, time.Millisecond*500)
	if err == nil {
		// Retry another port
		_ = conn.Close()
		return RQLStart()
	}

	slog.Info("Starting RethinkDB")
	cmd := exec.Command("rethinkdb", "--no-http-admin", "--bind", "127.0.0.1", "--cluster-port", fmt.Sprintf("%d", genPort+200), "--driver-port", fmt.Sprintf("%d", genPort))
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	err = cmd.Start()

	if err != nil {
		b, _ := cmd.Output()
		slog.Error(string(b))
		panic(fmt.Errorf("cannot start rethinkdb: %s", err))
	}

	//remote_signer.RethinkDBPort = genPort

	slog.Info("Waiting for rethink to settle")
	time.Sleep(time.Second)
	retry := 0
	started := false

	for retry < 5 {
		conn, err := net.DialTimeout("tcp", connString, time.Millisecond*500)
		if err != nil {
			slog.Error(err)
			retry++
			time.Sleep(time.Second)
			continue
		}
		_ = conn.Close()
		started = true
		break
	}

	if !started {
		return nil, genPort, fmt.Errorf("timeout waiting for rethinkdb to start")
	}

	return cmd, genPort, nil
}

// RQLStop stops a RethinkDB instance that has been started with RQLStart
func RQLStop(cmd *exec.Cmd) {
	slog.Info("Stopping RethinkDB")
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
