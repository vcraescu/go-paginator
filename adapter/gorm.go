package adapter

import (
	"github.com/jinzhu/gorm"
)

// GORMAdapter gorm adapter to be passed to paginator constructor
type GORMAdapter struct {
	db *gorm.DB
}

// NewGORMAdapter gorm adapter constructor which receive the gorm db query.
func NewGORMAdapter(db *gorm.DB) GORMAdapter {
	return GORMAdapter{db: db}
}

// Nums returns the number of records
func (a GORMAdapter) Nums() int {
	var count int
	a.db.Count(&count)

	return count
}

// Slice stores into data argument a slice of the results.
// data must be a pointer to a slice of models.
func (a GORMAdapter) Slice(offset, length int, data interface{}) error {
	return a.db.Limit(length).Offset(offset).Find(data).Error
}
