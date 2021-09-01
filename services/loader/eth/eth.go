package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

func GetEthClient(url string) *ethclient.Client {
	ethClient, err := ethclient.Dial(url)

	if err != nil {
		time.Sleep(time.Millisecond * 100)
		return GetEthClient(url)
	}

	return ethClient
}

func GetLatestBlockNumber(ethClient *ethclient.Client) uint64 {
	latestBlockNumber, err := ethClient.BlockNumber(context.Background())

	if err != nil {
		fmt.Println(fmt.Sprintf("failed to get latest block number with error: %s, retrying", err))
		time.Sleep(time.Millisecond * 100)
		return GetLatestBlockNumber(ethClient)
	}

	return latestBlockNumber
}
