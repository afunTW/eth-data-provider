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
