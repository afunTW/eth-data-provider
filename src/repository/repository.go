package repository

type Tabler interface {
	TableName() string
}

type EthereumIndexRepository interface {
	AddBlocks(records []*EthereumBlock) error
	AddTransactions(records []*EthereumTransaction) error
	AddLogs(records []*EthereumLog) error
	GetLatestBlock(limit int) ([]*EthereumBlock, error)
	GetBlock(blockNum uint64) (*EthereumBlock, error)
	GetTransactions(blockNum uint64) ([]*EthereumTransaction, error)
}
