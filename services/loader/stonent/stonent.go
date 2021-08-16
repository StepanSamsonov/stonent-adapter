package stonent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/tools/utils"
	"strconv"
	"strings"
)

type Log struct {
	TransactionHash string
}

type Transaction struct {
	Input string
}

type ParsedLog struct {
	Address  string
	IsActive bool
}

func parseCollectionChangedLog(ethClient *ethclient.Client, rawLog types.Log) (*ParsedLog, error) {
	jsonLog, err := rawLog.MarshalJSON()

	if err != nil {
		return nil, err
	}

	log := new(Log)
	parseLogErr := json.Unmarshal(jsonLog, &log)

	if parseLogErr != nil {
		return nil, parseLogErr
	}

	hash := common.HexToHash(log.TransactionHash)
	tx, _, err := ethClient.TransactionByHash(context.Background(), hash)

	if err != nil {
		fmt.Println(err)
	}

	transaction, err := tx.MarshalJSON()

	if err != nil {
		return nil, err
	}

	transactionData := new(Transaction)
	parseTransactionErr := json.Unmarshal(transaction, &transactionData)

	if parseTransactionErr != nil {
		return nil, parseTransactionErr
	}

	transactionValue := transactionData.Input
	protoAddress := transactionValue[len(transactionValue)-128 : len(transactionValue)-64]
	address := common.HexToAddress(protoAddress).String()
	value, err := strconv.Atoi(transactionValue[len(transactionValue)-64:])

	if err != nil {
		return nil, err
	}

	isActive := value == 1
	result := ParsedLog{
		Address:  address,
		IsActive: isActive,
	}

	return &result, nil
}

func GetIndexedCollections(ethClient *ethclient.Client) ([]string, error) {
	contractAddress := common.HexToAddress(config.StonentContractAddress)
	contractAbi, err := abi.JSON(strings.NewReader(ABI))

	if err != nil {
		return []string{}, err
	}

	filterEvents := [][]common.Hash{{
		contractAbi.Events["CollectionAdded"].ID,
		contractAbi.Events["CollectionRemoved"].ID,
	}}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			contractAddress,
		},
		Topics: filterEvents,
	}
	logs, err := ethClient.FilterLogs(context.Background(), query)

	if err != nil {
		return []string{}, err
	}

	res := []string{}

	for _, rawLog := range logs {
		parsedLog, err := parseCollectionChangedLog(ethClient, rawLog)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if parsedLog.IsActive {
			if !utils.Contains(res, parsedLog.Address) {
				res = append(res, parsedLog.Address)
			}
		} else {
			if utils.Contains(res, parsedLog.Address) {
				res = utils.RemoveByItem(res, parsedLog.Address)
			}
		}
	}

	return res, nil
}

func ListenIndexedCollections(ethClient *ethclient.Client, onChange func(address string, isActive bool)) {
	contractAddress := common.HexToAddress(config.StonentContractAddress)
	contractAbi, err := abi.JSON(strings.NewReader(ABI))

	if err != nil {
		fmt.Println(fmt.Sprintf("Error while listening collections: %s", err))
		return
	}

	filterEvents := [][]common.Hash{{
		contractAbi.Events["CollectionAdded"].ID,
		contractAbi.Events["CollectionRemoved"].ID,
	}}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			contractAddress,
		},
		Topics: filterEvents,
	}

	logsCh := make(chan types.Log)

	_, subscriptionError := ethClient.SubscribeFilterLogs(context.Background(), query, logsCh)

	if subscriptionError != nil {
		fmt.Println(fmt.Sprintf("Error while listening collections: %s", subscriptionError))
		return
	}

	for {
		select {
		case rawLog := <-logsCh:
			parsedLog, err := parseCollectionChangedLog(ethClient, rawLog)

			if err != nil {
				fmt.Println(err)
				continue
			}

			onChange(parsedLog.Address, parsedLog.IsActive)
		}
	}
}
