package pharmacyrepo

import (
	"context"
	"fmt"

	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/entities"
)

func (r *Repository) CreatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error {
	errBase := fmt.Sprintf("userrepo.CreatePharmacy(%v)", pharmacy)

	const query = `
		INSERT INTO pharmacies(address_id, name)
							VALUES ($1, $2)
`

	if _, err := r.master.ExecContext(
		ctx,
		query,
		pharmacy.Address.ID,
		pharmacy.Name,
	); err != nil {
		return fmt.Errorf("%s: failed to create user: %w", errBase, err)
	}

	return nil
}
