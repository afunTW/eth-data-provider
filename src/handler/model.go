package handler

// ==============================
// Request model
// ==============================

type ReqGetLatestBlocks struct {
	Limit uint16 `form:"limit"`
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

// ==============================
// Component model
// ==============================

type BlockInfo struct {
	BlockNum   int    `json:"block_num"`
	BlockHash  string `json:"block_hash"`
	BlockTime  uint64 `json:"block_time"`
	ParentHash string `json:"parent_hash"`
}
