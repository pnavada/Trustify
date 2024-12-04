package config

type Config struct {
	BlockchainSettings ConfigBlockchainSettings
	GenesisBlock       ConfigGenesisBlock
	Nodes              map[string]ConfigNode // This would be an actual node later
	peers              []string
}

type ConfigWallet struct {
	BitcoinAddress string
	PublicKey      string
	PrivateKey     string
}

type ConfigTransaction struct {
	Type            string
	Delay           int
	BuyerAddress    string
	SellerAddress   string
	ProductID       string
	Fee             int
	ReviewerAddress string
	Rating          int
	Amount          int
}

type ConfigNode struct {
	Wallet       ConfigWallet
	Transactions []ConfigTransaction
}

type ConfigBlockchainSettings struct {
	BlockSize              int
	TargetHash             string
	BlockConfirmationDepth int
	MiningReward           int
	ReviewReward           int
	RewardHalfTime         int
	Protocols              ConfigProtocol
}

type ConfigProtocol struct {
	GetBlocks ConfigGetBlocksProtocol
}

type ConfigGetBlocksProtocol struct {
	Timeout int
}

type ConfigGenesisBlock struct {
	BlockHash        string
	PreviousHash     string
	Nonce            int
	TargetHash       string
	Timestamp        int
	MerkleRoot       string
	TransactionCount int
	Transactions     []ConfigUTXOTransaction
}

type ConfigUTXOTransaction struct {
	ID      string
	Address string
	Amount  int
}

// Context - config.yml and this file

func LoadConfig(path string) *Config {
	// Load configuration from file based on the structures here and get everything ready
	// Identify the self hostname and load the node data accordingly and others go as list of peers
	// The final output would be accurately filled config
	return nil
}
