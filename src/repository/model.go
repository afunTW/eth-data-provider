package repository

/*
	CREATE TABLE blocks (
	   id INT AUTO_INCREMENT PRIMARY KEY,
	   block_num BIGINT UNSIGNED NOT NULL,
	   block_hash VARCHAR(66) NOT NULL,  -- 66 chars to accommodate 0x prefix and 64 hex characters
	   block_time BIGINT UNSIGNED NOT NULL,
	   parent_hash VARCHAR(66) NOT NULL,
	   UNIQUE (block_num),
	   UNIQUE (block_hash)

);
*/
type EthereumBlock struct {
	Id         int    `gorm:"column:id;primaryKey;autoIncrement"`
	BlockNum   uint64 `gorm:"column:block_num;uniqueIndex;"`
	BlockHash  string `gorm:"column:block_hash;uniqueIndex;"`
	BlockTime  uint64 `gorm:"column:block_time;"`
	ParentHash string `gorm:"column:parent_hash;"`
}

func (e *EthereumBlock) TableName() string { return "blocks" }

/*
	CREATE TABLE transactions (
	   id INT AUTO_INCREMENT PRIMARY KEY,
	   tx_hash VARCHAR(66) NOT NULL,  -- 66 chars for 0x prefix and 64 hex characters
	   block_num BIGINT UNSIGNED NOT NULL,
	   from_address VARCHAR(42) NOT NULL,  -- 42 chars for Ethereum addresses (0x + 40 hex)
	   to_address VARCHAR(42),
	   nonce BIGINT UNSIGNED,
	   data BINARY,
	   value BIGINT UNSIGNED NOT NULL,
	   FOREIGN KEY (block_num) REFERENCES blocks(block_num),
	   UNIQUE (tx_hash)

);
*/
type EthereumTransaction struct {
	Id          int    `gorm:"column:id;primaryKey;autoIncrement"`
	TxHash      string `gorm:"column:tx_hash;uniqueIndex"`
	BlockNum    uint64 `gorm:"column:block_num"`
	FromAddress string `gorm:"column:from_address"`
	ToAddress   string `gorm:"column:to_address"`
	Nonce       uint64 `gorm:"column:nonce"`
	Data        []byte `gorm:"column:data"`
	Value       uint64 `gorm:"column:value"`
}

func (e *EthereumTransaction) TableName() string { return "transactions" }
