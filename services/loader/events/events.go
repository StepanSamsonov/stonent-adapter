package events

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/vladimir3322/stonent_go/erc1155"
	"github.com/vladimir3322/stonent_go/services/loader/config"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"sync"
)

type iImageMetadata struct {
	Image string
}

func downloadImageWithWaiter(ipfsHost string, ipfsPath string, waiter *sync.WaitGroup) {
	defer waiter.Done()

	downloadImage(ipfsHost, ipfsPath)
}

func downloadImage(ipfsHost string, ipfsPath string) {
	ipfsMetadataUrl := ipfsHost + ipfsPath
	imageMetadataRes, err := http.Get(ipfsMetadataUrl)

	if err != nil {
		fmt.Println("Error with:", ipfsMetadataUrl, err)
		return
	}
	if imageMetadataRes.StatusCode != http.StatusOK {
		fmt.Println("Error with:", ipfsMetadataUrl, "invalid response code:", imageMetadataRes.StatusCode)
		return
	}

	defer imageMetadataRes.Body.Close()

	var jsonBody iImageMetadata
	imageMetadataParserErr := json.NewDecoder(imageMetadataRes.Body).Decode(&jsonBody)

	if imageMetadataParserErr != nil {
		fmt.Println("Error with:", ipfsMetadataUrl, imageMetadataParserErr)
		return
	}

	parsedImageUrl, err := url.Parse(jsonBody.Image)

	if err != nil {
		fmt.Println("Error with:", ipfsMetadataUrl, err)
		return
	}

	imageSourceUrl := ipfsHost + "/ipfs" + parsedImageUrl.Path
	imageSourceRes, err := http.Get(imageSourceUrl)

	if err != nil {
		fmt.Println("Error with:", imageSourceUrl, err)
		return
	}
	if imageSourceRes.StatusCode != http.StatusOK {
		fmt.Println("Error with:", imageSourceUrl, "invalid response code:", imageSourceRes.StatusCode)
		return
	}

	defer imageSourceRes.Body.Close()

	imageSource, err := ioutil.ReadAll(imageSourceRes.Body)

	if err != nil {
		fmt.Println("Error with:", imageSourceUrl, err)
		return
	}

	fmt.Println("Image size:", len(imageSource))
}

func GetEvents(contract *erc1155.Erc1155, startBlock uint64, endBlock uint64, waiter *sync.WaitGroup) {
	defer waiter.Done()

	if startBlock <= endBlock {
		opt := &bind.FilterOpts{Start: startBlock, End: &endBlock}
		s := []*big.Int{}
		past, err := contract.FilterURI(opt, s)

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
			notEmpty = past.Next()
			if notEmpty {
				waiter.Add(1)

				go downloadImageWithWaiter(config.IpfsLink[ipfsNodeIndex], past.Event.Value, waiter)

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
