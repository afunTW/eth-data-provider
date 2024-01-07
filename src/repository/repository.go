package repository

type Tabler interface {
	TableName() string
}

type EthereumIndexRepository interface {
	AddBlocks(records []*EthereumBlock) error
}
