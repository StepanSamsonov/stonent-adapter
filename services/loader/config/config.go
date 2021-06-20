package config

const ServerPort = 8080
const ProviderUrl = "wss://mainnet.infura.io/ws/v3/844de29fabee4fcebf315309262d0836"

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

const MlUrl = "http://localhost:9090"

const DownloadImageBufferSize = 10
const DownloadImageMaxCount = 1000 // -1 for ignoring
