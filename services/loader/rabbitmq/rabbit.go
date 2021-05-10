package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/tools/models"
	"log"
)

var rabbitConn *amqp.Connection

func initRabbit() *amqp.Connection {
	var err error
	rabbitAddr := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitLogin, config.RabbitPass, config.RabbitHost, config.RabbitPort)
	rabbitConn, err = amqp.Dial(rabbitAddr)
	if err != nil {
		handleError(err, "cant connect to rabbit")
		return nil
	}
	return rabbitConn
}

func getRabbitConn() *amqp.Connection {
	if rabbitConn == nil {
		return initRabbit()
	}
	return rabbitConn
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SendTestNFT() {
	// lets send test nft
	nft := models.NFT{
		NFTID:           "1",
		ContractAddress: "0x13123123123213",
		Data:            "",
	}
	SendNFTToRabbit(nft)
}
