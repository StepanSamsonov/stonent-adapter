package config

const ServerPort = 8080
const ProviderUrl = "wss://rinkeby.infura.io/ws/v3/ebf385aedfcc4f2b9d34a97ee6a86f93"
const StonentContractAddress = "0xFa9aF655Ef79445ECBb73389914e2ab16A31F62D"

const RabbitLogin = "guest"
const RabbitPass = "guest"
const RabbitHost = "localhost"
const RabbitPort = "5672"
const RabbitQueueIndexing = "indexing"

const PostgresDbName = "postgres"
const PostgresSchema = "schema"
const PostgresLogin = "guest"
const PostgresPassword = "guest"
const PostgresHost = "localhost"
const PostgresPort = "5432"

const DownloadImageBufferSize = 10
const DownloadImageMaxCount = 1000 // -1 for ignoring
