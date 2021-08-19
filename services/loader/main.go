package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/eth"
	"github.com/vladimir3322/stonent_go/events"
	"github.com/vladimir3322/stonent_go/ipfs"
	"github.com/vladimir3322/stonent_go/postgres"
	"github.com/vladimir3322/stonent_go/rabbitmq"
	"github.com/vladimir3322/stonent_go/server"
	"github.com/vladimir3322/stonent_go/stonent"
	"github.com/vladimir3322/stonent_go/tools/erc1155"
	"github.com/vladimir3322/stonent_go/tools/models"
	"github.com/vladimir3322/stonent_go/tools/utils"
	"sync"
)

func main() {
	configErr := config.InitConfig()

	if configErr != nil {
		panic(configErr)
	}

	go server.Run()
	rabbitmq.Init()
	ipfs.Init()
	postgres.Init()

	commonEthConnection := eth.GetEthClient(config.CommonProviderUrl)
	collectionsEthConnection := eth.GetEthClient(config.CollectionsProviderUrl)

	collections, collectionsErr := stonent.GetIndexedCollections(commonEthConnection)
	completedCollections := 0

	if collectionsErr != nil {
		panic(collectionsErr)
	}

	fmt.Println(fmt.Sprintf("Indexed collections set: %s", collections))

	startBlockNumber := uint64(0)
	latestBlockNumber := eth.GetLatestBlockNumber(collectionsEthConnection)

	for _, collection := range collections {
		go getEvents(collectionsEthConnection, collection, startBlockNumber, latestBlockNumber, func() {
			completedCollections += 1

			if len(collections) == completedCollections {
				fmt.Println("loader completed successfully")

				rabbitmq.SendNFTToRabbit(models.NFT{
					IsFinite: true,
				})
			}
		})

		go events.ListenEvents(collectionsEthConnection, collection, latestBlockNumber)
	}

	go stonent.ListenIndexedCollections(commonEthConnection, func(collection string, isActive bool) {
		if isActive {
			if !utils.Contains(collections, collection) {
				collections = append(collections, collection)

				go getEvents(collectionsEthConnection, collection, startBlockNumber, latestBlockNumber, func() {
					rabbitmq.SendNFTToRabbit(models.NFT{
						IsFinite: true,
					})
				})
				go events.ListenEvents(collectionsEthConnection, collection, latestBlockNumber)
			}
		} else {
			if utils.Contains(collections, collection) {
				collections = utils.RemoveByItem(collections, collection)
			}
		}

		fmt.Println(fmt.Sprintf("Indexed collections set has been changed %s", collections))
	})

	utils.WaitSignals()
}

func getEvents(ethConnection *ethclient.Client, address string, startBlock uint64, endBlock uint64, onFinished func()) {
	contract, err := erc1155.NewErc1155(common.HexToAddress(address), ethConnection)

	if err != nil {
		fmt.Println(fmt.Sprintf("whoops something went wrong: %s", err))
	}

	var waiter = &sync.WaitGroup{}

	waiter.Add(1)
	events.GetEvents(address, contract, startBlock, endBlock, waiter)
	go events.RunBuffer()
	waiter.Wait()
	onFinished()
}
