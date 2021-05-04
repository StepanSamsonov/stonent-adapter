package events

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/tools/erc1155"
	"math/big"
	"sync"
)

type iImageMetadata struct {
	Image string
}

func GetEvents(contract *erc1155.Erc1155, startBlock uint64, endBlock uint64, waiter *sync.WaitGroup) {
	defer waiter.Done()

	if countOfDownloaded >= config.DownloadImageMaxCount {
		waiter.Done()
		return
	}

	if startBlock <= endBlock {
		opt := &bind.FilterOpts{Start: startBlock, End: &endBlock}
		s := []*big.Int{}
		past, err := contract.FilterURI(opt, s)

		if countOfDownloaded >= config.DownloadImageMaxCount {
			waiter.Done()
			return
		}

		if err != nil {
			var middleBlock = (startBlock + endBlock) / 2

			waiter.Add(1)
			go GetEvents(contract, startBlock, middleBlock, waiter)
			waiter.Add(1)
			go GetEvents(contract, middleBlock+1, endBlock, waiter)
			return
		}

		notEmpty := true
		ipfsNodeIndex := 0

		for notEmpty {
			if countOfDownloaded >= config.DownloadImageMaxCount {
				waiter.Done()
				return
			}

			notEmpty = past.Next()
			if notEmpty {
				waiter.Add(1)

				go pushToBuffer(BufferItem{
					ipfsHost: config.IpfsLink[ipfsNodeIndex],
					ipfsPath: past.Event.Value,
					waiter:   waiter,
				})

				ipfsNodeIndex += 1
				ipfsNodeIndex %= len(config.IpfsLink)
			}
		}
	} else {
		return
	}
}

func ListenEvents(contract *erc1155.Erc1155, startBlock uint64) {
	s := []*big.Int{}
	ch := make(chan *erc1155.Erc1155URI)
	opts := &bind.WatchOpts{Start: &startBlock}
	watcher, err := contract.WatchURI(opts, ch, s)

	if err != nil {
		fmt.Println("Failed listening events:", err)
	}

	ipfsNodeIndex := 0

	for {
		select {
		case err := <-watcher.Err():
			fmt.Println("Failed listening events:", err)
		case Event := <-ch:
			fmt.Println(Event.Value)

			go downloadImage(config.IpfsLink[ipfsNodeIndex], Event.Value)

			ipfsNodeIndex += 1
			ipfsNodeIndex %= len(config.IpfsLink)
		}
	}
}
