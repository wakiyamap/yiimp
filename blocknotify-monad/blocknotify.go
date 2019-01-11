// Copyright (c) 2015-2017 YiiMP

// Sample blocknotify wrapper tool compatible with decred notifications
// will call the standard bin/blocknotify yiimp tool on new block event.

// Note: this tool is connected directly to monad, not to the wallet!

package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/wakiyamap/monad/wire"

	"github.com/wakiyamap/monad/rpcclient"
	"github.com/wakiyamap/monautil"
)

const (
	processName = "blocknotify"    // set the full path if required
	stratumDest = "yaamp.com:3252" // stratum host:port
	coinId = "393"                 // decred database coin id

	monadUser = "yiimprpc"
	monadPass = "myMonadPassword"

	debug = false
)

func main() {
	// Only override the handlers for notifications you care about.
	// Also note most of these handlers will only be called if you register
	// for notifications.  See the documentation of the rpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := rpcclient.NotificationHandlers{

		OnFilteredBlockConnected: func(height int32, blockHeader *wire.BlockHeader,
			txs []*monautil.Tx) {
			// log.Printf("Block bytes: %v %v", blockHeader, transactions)
			var bhead *wire.BlockHeader = blockHeader
			str := bhead.BlockHash().String();
			args := []string{ stratumDest, coinId, str }
			out, err := exec.Command(processName, args...).Output()
			if err != nil {
				log.Printf("err %s", err)
			} else if debug {
				log.Printf("out %s", out)
			}
			if (debug) {
				log.Printf("Block connected: %s", str)
			}
		},

	}

	// Connect to local monad RPC server using websockets.
	// monadHomeDir := monautil.AppDataDir("monad", false)
	// folder := monadHomeDir
	folder := ""
	certs, err := ioutil.ReadFile(filepath.Join(folder, "rpc.cert"))
	if err != nil {
		certs = nil
		log.Printf("%s, trying without TLS...", err)
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         "127.0.0.1:9402",
		Endpoint:     "ws", // websocket

		User:         monadUser,
		Pass:         monadPass,

		DisableTLS: (certs == nil),
		Certificates: certs,
	}

	client, err := rpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatalln(err)
	}

	// Register for block connect and disconnect notifications.
	if err := client.NotifyBlocks(); err != nil {
		log.Fatalln(err)
	}
	log.Println("NotifyBlocks: Registration Complete")

	// Wait until the client either shuts down gracefully (or the user
	// terminates the process with Ctrl+C).
	client.WaitForShutdown()
}
