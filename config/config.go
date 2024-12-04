package config

type Config struct {
	BlockchainSettings ConfigBlockchainSettings
	GenesisBlock       ConfigGenesisBlock
	Nodes              map[string]ConfigNode // This would be an actual node later
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

func LoadConfig(path string) *Config {
	// Load configuration from file
	return nil
}
