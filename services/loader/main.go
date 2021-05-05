package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/events"
	"github.com/vladimir3322/stonent_go/server"
	"github.com/vladimir3322/stonent_go/tools/erc1155"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	//go getEvents("0xd07dc4262bcdbf85190c01c996b4c06a461d2430", 0, 12291943) // Почти все
	//go getEvents("0xd07dc4262bcdbf85190c01c996b4c06a461d2430", 12211843, 12291943) // Много картин
	go getEvents("0xd07dc4262bcdbf85190c01c996b4c06a461d2430", 12291940, 12291943) // 4 картины

	//go listenEvents("0xd07dc4262bcdbf85190c01c996b4c06a461d2430", 12291943)

	//go redis.ConsumeEvents()
    go server.Run()
	WaitSignals()

	//fmt.Println(api.GetLatestBlock(conn))
}

func getEvents(address string, startBlock uint64, endBlock uint64) {
	conn, err := ethclient.Dial(config.ProviderUrl)

	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	contract, err := erc1155.NewErc1155(common.HexToAddress(address), conn)
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	var waiter = &sync.WaitGroup{}

	waiter.Add(1)
	events.GetEvents(contract, startBlock, endBlock, waiter)
	go events.RunBuffer()
	waiter.Wait()

	fmt.Println("Finished")
}

func listenEvents(address string, startBlock uint64) {
	conn, err := ethclient.Dial(config.ProviderUrl)

	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	contract, err := erc1155.NewErc1155(common.HexToAddress(address), conn)

	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	events.ListenEvents(contract, startBlock)
}

func WaitSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals
	fmt.Println("Got signal for exiting", sig)
}
