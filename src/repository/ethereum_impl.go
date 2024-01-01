package repository

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type jsonRpcRequest struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

type jsonRpcResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result"`
}

func (r *jsonRpcResponse) ToBlock() (*EthereumBlock, error) {
	var block EthereumBlock
	blockBytes, err := json.Marshal(r.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(blockBytes, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *jsonRpcResponse) ToTransaction() (*EthereumTransaction, error) {
	var tx EthereumTransaction
	txBytes, err := json.Marshal(r.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(txBytes, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

type ethereumImpl struct {
	host   string
	client *http.Client
}

func NewEthereumJsonRpcImpl(host string) EthereumRepository {
	return &ethereumImpl{host: host, client: &http.Client{}}
}

func (e *ethereumImpl) rpcCall(body io.Reader) (*jsonRpcResponse, error) {
	req, err := http.NewRequest("POST", e.host, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse response
	var respBody jsonRpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	return &respBody, nil
}

func (e *ethereumImpl) GetBlockByNumber(blockNum string) (*EthereumBlock, error) {
	// prepare
	reqBody := jsonRpcRequest{
		JsonRpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{blockNum, true},
		Id:      1,
	}
	log.Debugf("GetBlockByNumber: prepare rpc request %+v\n", reqBody)
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	// rpc call
	respBody, err := e.rpcCall(bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	// transform
	block, err := respBody.ToBlock()
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethereumImpl) GetTransactionByHash(txHash string) (*EthereumTransaction, error) {
	// prepare
	reqBody := jsonRpcRequest{
		JsonRpc: "2.0",
		Method:  "eth_getTransactionByHash",
		Params:  []interface{}{txHash},
		Id:      1,
	}
	log.Debugf("GetTransactionByHash: prepare rpc request %+v\n", reqBody)
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	// rpc call
	respBody, err := e.rpcCall(bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	// transform
	tx, err := respBody.ToTransaction()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (e *ethereumImpl) GetTransactionReceipt() {}
