package events

import (
	"errors"
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

func GetById(contract *erc1155.Erc1155, id *big.Int) (string, error) {
	opt := &bind.FilterOpts{}
	s := []*big.Int{id}
	event, err := contract.FilterURI(opt, s)

	if err != nil {
		return "", err
	}

	isExist := event.Next()

	if !isExist {
		return "", errors.New("event not found")
	}

	imageSource, err := getImageSource(config.IpfsLink[0], event.Event.Value)

	if err != nil {
		return "", err
	}

	return imageSource, nil
}

func GetEvents(address string, contract *erc1155.Erc1155, startBlock uint64, endBlock uint64, waiter *sync.WaitGroup) {
	defer waiter.Done()

	if IsExceededImagesLimitCount() {
		waiter.Done()
		return
	}

	if startBlock <= endBlock {
		var s []*big.Int

		opt := &bind.FilterOpts{Start: startBlock, End: &endBlock}
		past, err := contract.FilterURI(opt, s)

		if IsExceededImagesLimitCount() {
			waiter.Done()
			return
		}

		if err != nil {
			var middleBlock = (startBlock + endBlock) / 2

			waiter.Add(1)
			go GetEvents(address, contract, startBlock, middleBlock, waiter)
			waiter.Add(1)
			go GetEvents(address, contract, middleBlock+1, endBlock, waiter)
			return
		}

		notEmpty := true
		ipfsNodeIndex := 0

		for notEmpty {
			if IsExceededImagesLimitCount() {
				waiter.Done()
				return
			}

			notEmpty = past.Next()

			if notEmpty {
				waiter.Add(1)

				go pushToBuffer(BufferItem{
					address:  address,
					nftId:    past.Event.Id.String(),
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

func ListenEvents(address string, contract *erc1155.Erc1155, startBlock uint64) {
	var s []*big.Int

	opts := &bind.WatchOpts{Start: &startBlock}
	ch := make(chan *erc1155.Erc1155URI)
	watcher, err := contract.WatchURI(opts, ch, s)

	if err != nil {
		fmt.Println("failed listening events:", err)
	}

	ipfsNodeIndex := 0

	for {
		select {
		case err := <-watcher.Err():
			fmt.Println("failed listening events:", err)
		case Event := <-ch:
			fmt.Println(fmt.Sprintf("received event from listening: %s", Event.Id))

			go downloadImage(address, Event.Id.String(), config.IpfsLink[ipfsNodeIndex], Event.Value)

			ipfsNodeIndex += 1
			ipfsNodeIndex %= len(config.IpfsLink)
		}
	}
}
