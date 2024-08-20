package dbconnect

import (
	"fmt"
	"log"
	"myapp/config"
	"strconv"
	"time"

	tarantool "github.com/tarantool/go-tarantool"
)

const tarantoolMaxRetries = 10
const tarantoolRetryDelay = 5 * time.Second

func ConnectToTarantool() (*tarantool.Connection, error) {
	addr := string(config.CFG.TarantoolCFG.Host) + ":" + strconv.Itoa(int(config.CFG.TarantoolCFG.Port))

	var tConn *tarantool.Connection
	var err error

	opts := tarantool.Opts{
		User: "root",
		Pass: "root",
	}

	for i := 0; i < tarantoolMaxRetries; i++ {
		tConn, err = tarantool.Connect(addr, opts)
		if err != nil {
			log.Printf("Failed to connect to Tarantool (attempt %d/%d): %v", i+1, tarantoolMaxRetries, err)
			time.Sleep(tarantoolRetryDelay)
			continue
		}

		// Проверка соединения с Tarantool
		resp, err := tConn.Ping()
		if err != nil || resp.Code != tarantool.OkCode {
			log.Printf("Failed to ping Tarantool (attempt %d/%d): %v", i+1, tarantoolMaxRetries, err)
			tConn.Close()
			time.Sleep(tarantoolRetryDelay)
			continue
		}

		log.Printf("Successfully connected to Tarantool. Address: %s", addr)
		return tConn, nil
	}

	return nil, fmt.Errorf("failed to connect to Tarantool after %d attempts", tarantoolMaxRetries)
}
