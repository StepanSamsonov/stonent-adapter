package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vladimir3322/stonent_go/erc1155"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//go pastEvents("0xd07dc4262bcdbf85190c01c996b4c06a461d2430")

	go watchEvents("0xd07dc4262bcdbf85190c01c996b4c06a461d2430")

	WaitSignals()

	//fmt.Println(api.GetLatestBlock(conn))
}

func pastEvents(address string) {
	conn, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/844de29fabee4fcebf315309262d0836")
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	contract, err := erc1155.NewErc1155(common.HexToAddress(address), conn)
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	opt := &bind.FilterOpts{} // todo нужно идти рекурсивно
	s := []*big.Int{}
	past, err := contract.FilterURI(opt, s)
	if err != nil {
		log.Fatalf("Failed FilterURI: %v", err)
	}
	notEmpty := true
	for notEmpty {
		notEmpty = past.Next()
		if notEmpty {
			fmt.Println("event log:", past.Event.Id, past.Event.Value)
		}
	}
}

func watchEvents(address string) {
	conn, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/844de29fabee4fcebf315309262d0836")
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	contract, err := erc1155.NewErc1155(common.HexToAddress(address), conn)
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	var blockNumber uint64 = 2900000 // todo нужно взять последний блок
	s := []*big.Int{}
	ch := make(chan *erc1155.Erc1155URI)
	opts := &bind.WatchOpts{}
	opts.Start = &blockNumber
	sub, err := contract.WatchURI(opts, ch, s)
	if err != nil {
		log.Fatalf("Failed WatchYearChanged: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-ch:
			fmt.Println("event log:", vLog.Id, vLog.Value) // pointer to event log
		}
	}
}
func WaitSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals
	fmt.Println("Got signal for exiting", sig)
}
