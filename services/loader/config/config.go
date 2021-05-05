package config

var ServerPort = 5000
var ProviderUrl = "wss://mainnet.infura.io/ws/v3/844de29fabee4fcebf315309262d0836"
var IpfsLink = []string{"https://ipfs.daonomic.com", "https://ipfs.io"}

var RabbitMQUrl = "amqp://rabbitmq:rabbitmq@rabbit1:5672/"
var RabbitMQQueueName = "imageSources"

var RedisUrl = "redis:6379"
var RedisJobQueue = "imageSources"

var DownloadImageBufferSize = 2
var DownloadImageMaxCount = -1 // -1 for ignoring
