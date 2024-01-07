package repository

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ethereumIndexGormImpl struct {
	db *gorm.DB
}

func NewEthereumIndexGormImpl(dsn string) EthereumIndexRepository {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &ethereumIndexGormImpl{db: db}
}

func (r *ethereumIndexGormImpl) AddBlocks(records []*EthereumBlock) error {
	return r.addRecords(&records)
}

func (r *ethereumIndexGormImpl) AddTransactions(records []*EthereumTransaction) error {
	return r.addRecords(&records)
}

func (r *ethereumIndexGormImpl) AddLogs(records []*EthereumLog) error {
	return r.addRecords(records)
}

func (r *ethereumIndexGormImpl) GetLatestBlock(limit int) ([]*EthereumBlock, error) {
	var records []*EthereumBlock
	result := r.db.Order("block_num DESC").Limit(limit).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (r *ethereumIndexGormImpl) GetBlock(blockNum uint64) (*EthereumBlock, error) {
	var record *EthereumBlock
	result := r.db.Model(&EthereumBlock{}).
		Where("block_num = ?", blockNum).
		First(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func (r *ethereumIndexGormImpl) GetTransaction(txHash string) (*EthereumTransaction, error) {
	var record *EthereumTransaction
	result := r.db.Model(&EthereumTransaction{}).
		Where("tx_hash = ?", txHash).
		First(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func (r *ethereumIndexGormImpl) GetTransactions(blockNum uint64) ([]*EthereumTransaction, error) {
	var records []*EthereumTransaction
	result := r.db.Model(&EthereumTransaction{}).
		Where("block_num = ?", blockNum).
		Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (r *ethereumIndexGormImpl) GetLogs(txHash string) ([]*EthereumLog, error) {
	var records []*EthereumLog
	result := r.db.Model(&EthereumLog{}).
		Where("tx_hash = ?", txHash).
		Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (r *ethereumIndexGormImpl) addRecords(records interface{}) error {
	result := r.db.Create(records)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
