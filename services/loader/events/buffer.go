package events

import (
	"github.com/vladimir3322/stonent_go/config"
	"sync"
)

type BufferItem struct {
	address     string
	nftId       string
	ipfsPath    string
	blockNumber uint64
	waiter      *sync.WaitGroup
}

var Buffer = make(chan BufferItem)
var bufferSize = 0
var CountOfFound = 0
var CountOfDownloaded = 0
var mutex = sync.Mutex{}

func IsExceededImagesLimitCount() bool {
	return config.DownloadImagesMaxCount != -1 && CountOfDownloaded >= config.DownloadImagesMaxCount
}

func pushToBuffer(item BufferItem) {
	mutex.Lock()
	defer mutex.Unlock()

	CountOfFound += 1

	for bufferSize >= config.DownloadImagesBufferSize {

	}

	if IsExceededImagesLimitCount() {
		item.waiter.Done()
		return
	}

	bufferSize += 1
	Buffer <- item
}

func RunBuffer() {
	for {
		select {
		case item := <-Buffer:
			if IsExceededImagesLimitCount() {
				item.waiter.Done()
				return
			}

			go downloadImageWithWaiter(item.address, item.nftId, item.ipfsPath, item.blockNumber, item.waiter, func(isSucceed bool) {
				bufferSize -= 1
			})
		}
	}
}
