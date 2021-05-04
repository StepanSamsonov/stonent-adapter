package events

import (
	"encoding/json"
	"fmt"
	"github.com/vladimir3322/stonent_go/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

func downloadImageWithWaiter(ipfsHost string, ipfsPath string, waiter *sync.WaitGroup, cb func(isSucceed bool)) {
	defer waiter.Done()

	isSucceed := downloadImage(ipfsHost, ipfsPath)
	cb(isSucceed)
}

func downloadImage(ipfsHost string, ipfsPath string) bool {
	ipfsMetadataUrl := ipfsHost + ipfsPath
	imageMetadataRes, err := http.Get(ipfsMetadataUrl)

	if countOfDownloaded >= config.DownloadImageMaxCount {
		return false
	}
	if err != nil {
		fmt.Println("Error with:", ipfsMetadataUrl, err)
		return false
	}
	if imageMetadataRes.StatusCode != http.StatusOK {
		fmt.Println("Error with:", ipfsMetadataUrl, "invalid response code:", imageMetadataRes.StatusCode)
		return false
	}

	defer imageMetadataRes.Body.Close()

	var jsonBody iImageMetadata
	imageMetadataParserErr := json.NewDecoder(imageMetadataRes.Body).Decode(&jsonBody)

	if countOfDownloaded >= config.DownloadImageMaxCount {
		return false
	}
	if imageMetadataParserErr != nil {
		fmt.Println("Error with:", ipfsMetadataUrl, imageMetadataParserErr)
		return false
	}

	parsedImageUrl, err := url.Parse(jsonBody.Image)

	if countOfDownloaded >= config.DownloadImageMaxCount {
		return false
	}
	if err != nil {
		fmt.Println("Error with:", ipfsMetadataUrl, err)
		return false
	}

	imageSourceUrl := ipfsHost + "/ipfs" + parsedImageUrl.Path
	imageSourceRes, err := http.Get(imageSourceUrl)

	if countOfDownloaded >= config.DownloadImageMaxCount {
		return false
	}
	if err != nil {
		fmt.Println("Error with:", ipfsMetadataUrl, imageSourceUrl, err)
		return false
	}
	if imageSourceRes.StatusCode != http.StatusOK {
		fmt.Println("Error with:", ipfsMetadataUrl, imageSourceUrl, "invalid response code:", imageSourceRes.StatusCode)
		return false
	}

	defer imageSourceRes.Body.Close()

	imageSource, err := ioutil.ReadAll(imageSourceRes.Body)

	if countOfDownloaded >= config.DownloadImageMaxCount {
		return false
	}
	if err != nil {
		fmt.Println("Error with:", imageSourceUrl, err)
		return false
	}

	fmt.Println("Downloaded image size:", len(imageSource))
	//rabbitmq.PushEvent(imageSource)

	return true
}
