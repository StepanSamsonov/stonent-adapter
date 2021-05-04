package events

import (
	"github.com/vladimir3322/stonent_go/config"
	"sync"
)

type BufferItem struct {
	ipfsHost string
	ipfsPath string
	waiter   *sync.WaitGroup
}

var Buffer = make(chan BufferItem)
var bufferSize = 0
var countOfDownloaded = 0
var mutex = sync.Mutex{}

func pushToBuffer(item BufferItem) {
	mutex.Lock()

	for bufferSize >= config.DownloadImageBufferSize {

	}

	if countOfDownloaded >= config.DownloadImageMaxCount {
		item.waiter.Done()
		return
	}

	bufferSize += 1
	Buffer <- item

	mutex.Unlock()
}

func RunBuffer() {
	for {
		select {
		case item := <-Buffer:
			if countOfDownloaded >= config.DownloadImageMaxCount {
				item.waiter.Done()
				return
			}

			go downloadImageWithWaiter(item.ipfsHost, item.ipfsPath, item.waiter, func(isSucceed bool) {
				bufferSize -= 1

				if isSucceed {
					countOfDownloaded += 1
				}
			})
		}
	}
}
