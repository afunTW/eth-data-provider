package router

// ==============================
// Request model
// ==============================

type ReqGetLatestBlocks struct {
	Limit uint16 `form:"limit,default=1" binding="required,min=1,max=1000"`
}

// ==============================
// Response model
// ==============================

type RespStatus struct {
	Message   string `json:"message"`
	Timestamp uint64 `json:"timestamp"`
}

type RespGetLatestBlocks struct {
	Blocks []BlockInfo `json:"blocks"`
	RespStatus
}

type RespGetBlockDetail struct {
	BlockInfo
	Transactions []string `json:"transactions"`
	RespStatus
}

type RespGetTransactionDetail struct {
	TransactionInfo
	Logs []TransactionEventLogInfo `json:"logs"`
	RespStatus
}

// ==============================
// Component model
// ==============================

type BlockInfo struct {
	BlockNum   int    `json:"block_num"`
	BlockHash  string `json:"block_hash"`
	BlockTime  uint64 `json:"block_time"`
	ParentHash string `json:"parent_hash"`
}

type TransactionInfo struct {
	TxHash string `json:"tx_hash"`
	From   string `json:"from"`
	To     string `json:"to"`
	Nonce  int    `json:"nonce"`
	Data   string `json:"data"`
	Value  string `json:"value"`
}

type TransactionEventLogInfo struct {
	Index int    `json:"index"`
	Data  string `json:"data"`
}
