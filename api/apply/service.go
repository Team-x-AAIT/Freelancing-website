package apply

import "github.com/Team-x-AAIT/Freelancing-website/api/entity"

// Service interface for apply
type ApplyService interface {
	StoreApply(apply *entity.Apply) (*entity.Apply, []error)
}
