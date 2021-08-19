package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const ServerPort = 8080

var StonentContractAddress = ""
var CommonProviderUrl = ""
var CollectionsProviderUrl = ""

const RabbitLogin = "guest"
const RabbitPassword = "guest"
const RabbitHost = "localhost"
const RabbitPort = "5672"
const RabbitQueueIndexing = "indexing"

const PostgresDbName = "postgres"
const PostgresSchema = "schema"
const PostgresLogin = "guest"
const PostgresPassword = "guest"
const PostgresHost = "localhost"
const PostgresPort = "5432"

var DownloadImagesBufferSize = 10
var DownloadImagesMaxCount = -1 // -1 for ignoring

func InitConfig() error {
	dotenvErr := godotenv.Load()

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
		parsedDownloadImagesMaxCount, parseDownloadImagesMaxCount := strconv.Atoi(unparsedDownloadImagesMaxCount)

		if parseDownloadImagesMaxCount != nil {
			return parseDownloadImagesMaxCount
		}

		DownloadImagesMaxCount = parsedDownloadImagesMaxCount
	}

	return nil
}
