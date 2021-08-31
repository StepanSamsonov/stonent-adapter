package models

type NFT struct {
	NFTID           string `json:"nftID"`
	ContractAddress string `json:"contractAddress"`
	Data            string `json:"data"`
	BlockNumber     uint64 `json:"blockNumber"`
	IsFinite        bool   `json:"isFinite"`
}
