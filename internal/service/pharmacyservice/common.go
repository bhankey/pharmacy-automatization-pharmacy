package pharmacyservice

import (
	"context"
	"fmt"
	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/entities"
)

func (s *Service) CreatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error {
	errBase := fmt.Sprintf("pharmacyservice.CreatePharmacy(%v)", pharmacy)
	var err error

	pharmacy.Address.ID, err = s.addressRepo.CreateAddress(ctx, pharmacy.Address)
	if err != nil {
		return fmt.Errorf("%s: failed to create pharmacy address: %w", errBase, err)
	}

	if err := s.pharmacyRepo.CreatePharmacy(ctx, pharmacy); err != nil {
		return fmt.Errorf("%s: failed to create pharmacy: %w", errBase, err)
	}

	return nil
}

func (s *Service) GetBatchOfPharmacies(
	ctx context.Context,
	lastPharmacyID int,
	limit int,
) ([]entities.Pharmacy, error) {
	errBase := fmt.Sprintf("pharmacyservice.GetBatchOfPharmacies(%d, %d)", lastPharmacyID, limit)

	pharmacies, err := s.pharmacyRepo.GetBatchOfPharmacies(ctx, lastPharmacyID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get batch of pharmacies: %w", errBase, err)
	}

	return pharmacies, nil
}
