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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &ethereumIndexGormImpl{db: db}
}

func (r *ethereumIndexGormImpl) AddBlocks(records []*EthereumBlock) error {
	result := r.db.Create(&records)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
