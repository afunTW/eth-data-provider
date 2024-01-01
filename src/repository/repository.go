package repository

type EthereumRepository interface {
	GetBlockByNumber()
	GetTransactionByHash()
	GetTransactionReceipt()
}
