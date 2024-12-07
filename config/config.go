package config

import (
	"fmt"
	"log"
	"os"
	"trustify/logger"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BlockchainSettings ConfigBlockchainSettings `yaml:"blockchain_settings"`
	GenesisBlock       ConfigGenesisBlock       `yaml:"genesis_block"`
	Nodes              map[string]ConfigNode    `yaml:"nodes"`
}

type ConfigBlockchainSettings struct {
	BlockSize              int            `yaml:"block_size"`
	TargetHash             string         `yaml:"target_hash"`
	BlockConfirmationDepth int            `yaml:"block_confirmation_depth"`
	MiningReward           int            `yaml:"mining_reward"`
	ReviewReward           int            `yaml:"review_reward"`
	RewardHalfTime         int            `yaml:"reward_half_time"`
	MiningTimeout          int            `yaml:"mining_timeout"`
	Protocols              ConfigProtocol `yaml:"protocols"`
}

type ConfigProtocol struct {
	GetBlocks ConfigGetBlocksProtocol `yaml:"get_blocks"`
}

type ConfigGetBlocksProtocol struct {
	Timeout int `yaml:"timeout"`
}

type ConfigNode struct {
	Wallet       ConfigWallet        `yaml:"wallet"`
	Transactions []ConfigTransaction `yaml:"transactions"`
}

type ConfigWallet struct {
	PrivateKey string `yaml:"private_key"`
}

type ConfigTransaction struct {
	Type            string `yaml:"type"`
	Delay           int    `yaml:"delay"`
	BuyerAddress    string `yaml:"buyer_address,omitempty"`
	SellerAddress   string `yaml:"seller_address,omitempty"`
	ProductID       string `yaml:"product_id"`
	Fee             int    `yaml:"fee,omitempty"`
	ReviewerAddress string `yaml:"reviewer_address,omitempty"`
	Rating          int    `yaml:"rating,omitempty"`
	Amount          int    `yaml:"amount,omitempty"`
}

type ConfigGenesisBlock struct {
	BlockHash        string                    `yaml:"block_hash"`
	PreviousHash     string                    `yaml:"previous_hash"`
	Nonce            int                       `yaml:"nonce"`
	TargetHash       string                    `yaml:"target_hash"`
	Timestamp        int                       `yaml:"timestamp"`
	MerkleRoot       string                    `yaml:"merkle_root"`
	TransactionCount int                       `yaml:"transaction_count"`
	Transactions     ConfigGenesisTransactions `yaml:"transactions"`
}

type ConfigGenesisTransactions struct {
	ID      string                  `yaml:"id"`
	Data    string                  `yaml:"data"`
	Inputs  []ConfigUTXOTransaction `yaml:"inputs"`
	Outputs []ConfigUTXOTransaction `yaml:"outputs"`
}

type ConfigUTXOTransaction struct {
	ID      string `yaml:"id"`
	Address string `yaml:"address"`
	Amount  int    `yaml:"amount"`
}

func LoadConfig(path string) (*Config, error) {
	// Read the YAML file
	configData, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read config file: %v\n", err)
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal the YAML into the Config struct
	var cfg Config
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		log.Printf("Failed to parse config file: %v\n", err)
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Print the config object
	// logger.InfoLogger.Printf("Loaded config: %+v\n", cfg)
	logger.InfoLogger.Printf("Loaded config")

	// Return the populated Config struct
	return &cfg, nil
}
