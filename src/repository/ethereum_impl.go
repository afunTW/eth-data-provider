package repository

type ethereumRpcImpl struct{}

func NewEthereumRpcImpl() EthereumRepository {
	return &ethereumRpcImpl{}
}

func (e *ethereumRpcImpl) GetBlockByNumber()      {}
func (e *ethereumRpcImpl) GetTransactionByHash()  {}
func (e *ethereumRpcImpl) GetTransactionReceipt() {}
