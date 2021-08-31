package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const Mode = "PROD"
const ServerPort = 8080

var StonentContractAddress = ""
var CommonProviderUrl = ""
var CollectionsProviderUrl = ""

var RabbitHost = ""

const RabbitLogin = "guest"
const RabbitPassword = "guest"
const RabbitPort = "5672"
const RabbitQueueIndexing = "indexing"

var PostgresHost = ""

const PostgresDbName = "postgres"
const PostgresSchema = "schema"
const PostgresLogin = "guest"
const PostgresPassword = "guest"
const PostgresPort = "5432"

var DownloadImagesBufferSize = 10
var DownloadImagesMaxCount = -1 // -1 for ignoring

func InitConfig() error {
	var dotenvErr error

	if Mode == "DEV" {
		dotenvErr = godotenv.Load("../../loader.env")
		RabbitHost = "localhost"
		PostgresHost = "localhost"
	} else {
		dotenvErr = godotenv.Load("./loader.env")
		RabbitHost = "rabbitmq"
		PostgresHost = "postgres"
	}

	if dotenvErr != nil {
		return dotenvErr
	}

	StonentContractAddress = os.Getenv("STONENT_CONTRACT_ADDRESS")

	if len(StonentContractAddress) == 0 {
		return errors.New("STONENT_CONTRACT_ADDRESS must be specified")
	}

	CommonProviderUrl = os.Getenv("COMMON_PROVIDER_URL")

	if len(CommonProviderUrl) == 0 {
		return errors.New("COMMON_PROVIDER_URL must be specified")
	}

	CollectionsProviderUrl = os.Getenv("COLLECTIONS_PROVIDER_URL")

	if len(CollectionsProviderUrl) == 0 {
		return errors.New("COLLECTIONS_PROVIDER_URL must be specified")
	}

	unparsedDownloadImagesBufferSize := os.Getenv("DOWNLOAD_IMAGES_BUFFER_SIZE")

	if len(unparsedDownloadImagesBufferSize) != 0 {
		parsedDownloadImagesBufferSize, parseDownloadImagesBufferSizeErr := strconv.Atoi(unparsedDownloadImagesBufferSize)

		if parseDownloadImagesBufferSizeErr != nil {
			return parseDownloadImagesBufferSizeErr
		}

		DownloadImagesBufferSize = parsedDownloadImagesBufferSize
	}

	unparsedDownloadImagesMaxCount := os.Getenv("DOWNLOAD_IMAGES_MAX_COUNT")

	if len(unparsedDownloadImagesMaxCount) != 0 {
		parsedDownloadImagesMaxCount, parseDownloadImagesMaxCountErr := strconv.Atoi(unparsedDownloadImagesMaxCount)

		if parseDownloadImagesMaxCountErr != nil {
			DownloadImagesMaxCount = -1
		}

		DownloadImagesMaxCount = parsedDownloadImagesMaxCount
	}

	return nil
}
