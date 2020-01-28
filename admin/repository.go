package admin

import "github.com/Team-x-AAIT/Freelancing-website/entity"

//IRepository contains databse operation methods
type IRepository interface {
	Admin(email, password string) (*entity.Admin, []error)
	UpdateAdmin(user *entity.Admin) (*entity.Admin, []error)
	DeleteAdmin(id string) (*entity.Admin, []error)
	StoreAdmin(admin *entity.Admin) (*entity.Admin, []error)
}
