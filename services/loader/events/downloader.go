package events

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/vladimir3322/stonent_go/ipfs"
	"github.com/vladimir3322/stonent_go/postgres"
	"github.com/vladimir3322/stonent_go/rabbitmq"
	"github.com/vladimir3322/stonent_go/tools/models"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

func getImageSource(ipfsPath string) (string, error) {
	imageMetadataRes, err := ipfs.Get(ipfsPath)

	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s: %s", ipfsPath, err))
	}

	re := regexp.MustCompile(`"ipfs://ipfs/.*?"`)
	imageUrl := re.Find(imageMetadataRes)

	if len(imageUrl) == 0 {
		return "", errors.New(fmt.Sprintf("error with %s: image source url not found", ipfsPath))
	}

	stringImageUrl := string(imageUrl)
	fixedStringImageUrl := strings.Replace(stringImageUrl, "\"", "", -1)
	parsedImageUrl, err := url.Parse(fixedStringImageUrl)

	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s: %s", ipfsPath, err))
	}

	imageSourceRes, err := ipfs.Get(parsedImageUrl.Path[1:])

	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s %s: %s", ipfsPath, parsedImageUrl.Path, err))
	}

	b64ImageSource := base64.StdEncoding.EncodeToString(imageSourceRes)

	return b64ImageSource, nil
}

func downloadImageWithWaiter(address string, nftId string, ipfsPath string, blockNumber uint64, waiter *sync.WaitGroup, cb func(isSucceed bool)) {
	defer waiter.Done()

	isSucceed := downloadImage(address, nftId, ipfsPath, blockNumber)

	if isSucceed {
		CountOfDownloaded += 1
	}

	cb(isSucceed)
}

func downloadImage(address string, nftId string, ipfsPath string, blockNumber uint64) bool {
	imageSource, err := getImageSource(ipfsPath)

	if err != nil {
		postgres.SaveRejectedImageByIPFS(address, nftId, ipfsPath, err)

		return false
	}

	rabbitmq.SendNFTToRabbit(models.NFT{
		NFTID:           nftId,
		ContractAddress: address,
		Data:            imageSource,
		BlockNumber:     blockNumber,
		IsFinite:        false,
	})

	return true
}
