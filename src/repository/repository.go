package repository

type Tabler interface {
	TableName() string
}

type EthereumIndexRepository interface {
	AddBlocks(records []*EthereumBlock) error
	AddTransactions(records []*EthereumTransaction) error
	AddLogs(records []*EthereumLog) error
}
