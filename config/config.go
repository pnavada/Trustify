package config

type Config struct {
	BlockchainSettings ConfigBlockchainSettings `yaml:"blockchain_settings"`
	GenesisBlock       ConfigGenesisBlock       `yaml:"genesis_block"`
	Nodes              map[string]ConfigNode    // This would be an actual node later
	Transactions       []ConfigTransaction      `yaml:"transactions"`
}

type ConfigWallet struct {
	BitcoinAddress string `yaml:"bitcoin_address"`
	PublicKey      string `yaml:"public_key"`
	PrivateKey     string `yaml:"private_key"`
}

type ConfigTransaction struct {
	Type            string `yaml:"type"`
	Delay           int    `yaml:"delay"`
	BuyerAddress    string `yaml:"buyer_address,omitempty"`
	SellerAddress   string `yaml:"seller_address,omitempty"`
	ProductID       string `yaml:"product_id,omitempty"`
	Fee             int    `yaml:"fee,omitempty"`
	ReviewerAddress string `yaml:"reviewer_address,omitempty"`
	Rating          int    `yaml:"rating,omitempty"`
}

type ConfigNode struct {
	Wallet ConfigWallet `yaml:"wallet"`
}

type ConfigBlockchainSettings struct {
	BlockSize              int            `yaml:"block_size"`
	TargetHash             string         `yaml:"target_hash"`
	BlockConfirmationDepth int            `yaml:"block_confirmation_depth"`
	MiningReward           int            `yaml:"mining_reward"`
	ReviewReward           int            `yaml:"review_reward"`
	RewardHalfTime         int            `yaml:"reward_half_time"`
	Protocols              ConfigProtocol `yaml:"protocols"`
}

type ConfigProtocol struct {
	GetBlocks ConfigGetBlocksProtocol `yaml:"get_blocks"`
}

type ConfigGetBlocksProtocol struct {
	Timeout int `yaml:"timeout"`
}

type ConfigGenesisBlock struct {
	BlockHash        string                `yaml:"block_hash"`
	PreviousHash     string                `yaml:"previous_hash"`
	Nonce            int                   `yaml:"nonce"`
	TargetHash       string                `yaml:"target_hash"`
	Timestamp        int                   `yaml:"timestamp"`
	MerkleRoot       string                `yaml:"merkle_root"`
	TransactionCount int                   `yaml:"transaction_count"`
	Transactions     []ConfigGenesisOutput `yaml:"transactions"`
}

type ConfigGenesisOutput struct {
	Type    string             `yaml:"type"`
	Outputs []ConfigUTXOOutput `yaml:"outputs"`
}

type ConfigUTXOOutput struct {
	ID      string `yaml:"id"`
	Address string `yaml:"address"`
	Amount  int    `yaml:"amount"`
}

func LoadConfig(path string) *Config {
	// Load configuration from file
	return nil
}
