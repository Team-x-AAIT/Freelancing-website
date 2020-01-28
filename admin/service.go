package admin

import "github.com/Team-x-AAIT/Freelancing-website/entity"

//IService contains databse operation methods
type IService interface {
	Admin(email, password string) (*entity.Admin, []error)
	UpdateAdmin(user *entity.Admin) (*entity.Admin, []error)
	DeleteAdmin(id string) (*entity.Admin, []error)
	StoreAdmin(admin *entity.Admin) (*entity.Admin, []error)
}
