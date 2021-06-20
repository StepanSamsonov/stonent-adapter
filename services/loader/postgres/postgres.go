package postgres

import (
    "context"
    "fmt"
    "github.com/go-pg/pg/v10"
    "github.com/vladimir3322/stonent_go/config"
)

var connection *pg.DB = nil
var ctx = context.Background()


func createNewConnection() *pg.DB {
    address := fmt.Sprintf("%s:%s", config.PostgresHost, config.PostgresPort)
    connection = pg.Connect(&pg.Options{
        User: config.PostgresLogin,
        Password: config.PostgresPassword,
        Addr: address,
        Database: config.PostgresDbName,
        OnConnect: func(ctx context.Context, conn *pg.Conn) error {
            _, err := conn.Exec("set search_path=?", config.PostgresSchema)

            return err
        },
    })

    return connection
}

func getConnection() *pg.DB {
    if connection == nil {
        return createNewConnection()
    }

    err := connection.Ping(ctx)

    if err == nil {
        return connection
    }

    return createNewConnection()
}

func Init() {
    getConnection()
}

func SaveRejectedImageByIPFS(contractAddress string, nftId string, ipfsPath string, description error) {
    type RejectedImage struct {
        tableName struct{} `pg:"rejected_images_by_ipfs"`
        Id string `json:"id" pg:"id,pk"`
        ContractAddress string `json:"contract_address" pg:"contract_address"`
        NftId string `json:"nft_id" pg:"nft_id"`
        IpfsPath string `json:"ipfs_path" pg:"ipfs_path"`
        Description string `json:"description" pg:"description"`
    }

    rejectedImage := RejectedImage{
        ContractAddress: contractAddress,
        NftId: nftId,
        IpfsPath: ipfsPath,
        Description: fmt.Sprintf("%s", description),
    }
    connection := getConnection()
    _, err := connection.Model(&rejectedImage).Insert()

    if err != nil {
        fmt.Println(fmt.Sprintf("Error with saving rejected image by IPFS: %s", err))
    }
}
