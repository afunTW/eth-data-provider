package repository

type EthereumRepository interface {
	GetBlockByNumber(blockNumber string) (*EthereumBlock, error)
	GetTransactionByHash()
	GetTransactionReceipt()
}
