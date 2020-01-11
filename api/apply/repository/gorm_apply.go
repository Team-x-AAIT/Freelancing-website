package repository

import (
	"fmt"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)
// creating connection 
type ApplyGormRepo struct {
	conn *gorm.DB
}
// initializing connection
func NewApplyGormRepo(dbconn *gorm.DB) *ApplyGormRepo {
	return &ApplyGormRepo{conn: dbconn}
}
// store applies with gorm  
func (ar *ApplyGormRepo) StoreApply(apply *entity.Apply) (*entity.Apply, []error) {
	aly := apply
	errs := ar.conn.Create(aly).GetErrors()

	for _, err := range errs {
		pqerr := err.(*pq.Error)
		fmt.Println(pqerr)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return aly, nil
}
