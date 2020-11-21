package adapter

import (
	"github.com/vcraescu/go-paginator/v2"
	"gorm.io/gorm"
)

// GORMAdapter gorm adapter to be passed to paginator constructor
type GORMAdapter struct {
	db *gorm.DB
}

// NewGORMAdapter gorm adapter constructor which receive the gorm db query.
func NewGORMAdapter(db *gorm.DB) paginator.Adapter {
	return &GORMAdapter{db: db}
}

// Nums returns the number of records
func (a *GORMAdapter) Nums() (int64, error) {
	var count int64
	if err := a.db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Slice stores into data argument a slice of the results.
// data must be a pointer to a slice of models.
func (a *GORMAdapter) Slice(offset, length int, data interface{}) error {
	// Work on a dedicated session to not offset the total count nums
	return a.db.Session(&gorm.Session{}).Limit(length).Offset(offset).Find(data).Error
}
