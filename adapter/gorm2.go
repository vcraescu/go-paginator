package adapter

import (
	"gorm.io/gorm"
)

// GORM2Adapter gorm adapter to be passed to paginator constructor
type GORM2Adapter struct {
	db *gorm.DB
}

// NewGORM2Adapter gorm adapter constructor which receive the gorm db query.
func NewGORM2Adapter(db *gorm.DB) GORM2Adapter {
	return GORM2Adapter{db: db}
}

// Nums returns the number of records
func (a GORM2Adapter) Nums() int {
	var count int64
	a.db.Count(&count)

	return int(count)
}

// Slice stores into data argument a slice of the results.
// data must be a pointer to a slice of models.
func (a GORM2Adapter) Slice(offset, length int, data interface{}) error {
	// Work on a dedicated session to not offset the total count nums
	return a.db.Session(&gorm.Session{}).Limit(length).Offset(offset).Find(data).Error
}
