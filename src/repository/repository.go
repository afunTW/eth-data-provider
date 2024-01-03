package repository

type EthereumRepository interface {
	GetBlockByNumber(blockNumber string) (*EthereumBlock, error)
	GetTransactionByHash(txHash string) (*EthereumTransaction, error)
	GetTransactionReceipt(txHash string) (*EthereumTransactionReceipt, error)
}
