package apply

import "github.com/Team-x-AAIT/Freelancing-website/api/entity"

//ApplyRepository is repository interface for apply
type ApplyRepository interface {
	StoreApply(apply *entity.Apply) (*entity.Apply, []error)
}
