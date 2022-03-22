package pharmacyservice

import (
	"context"

	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/entities"
)

type Service struct {
	pharmacyRepo PharmacyRepo
	addressRepo  AddressRepo
}

type PharmacyRepo interface {
	CreatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error
	GetBatchOfPharmacies(ctx context.Context, lastPharmacyID int, limit int) ([]entities.Pharmacy, error)
}

type AddressRepo interface {
	CreateAddress(ctx context.Context, address entities.Address) (int, error)
}

func NewPharmacyService(
	pharmacyRepo PharmacyRepo,
	addressRepo AddressRepo,
) *Service {
	return &Service{
		pharmacyRepo: pharmacyRepo,
		addressRepo:  addressRepo,
	}
}
