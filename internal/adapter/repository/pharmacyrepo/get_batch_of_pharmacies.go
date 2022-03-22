package pharmacyrepo

import (
	"context"
	"fmt"

	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/entities"
)

func (r *Repository) GetBatchOfPharmacies(
	ctx context.Context,
	lastPharmacyID int,
	limit int,
) ([]entities.Pharmacy, error) {
	errBase := fmt.Sprintf("userrepo.GetBatchOfPharmacies(%d, %d)", lastPharmacyID, limit)

	const query = `
		SELECT 
		       pharmacies.id AS pharmacy_id,
		       pharmacies.name,
		       pharmacies.is_blocked, 
		       addresses.id as address_id,
		       addresses.city,
		       addresses.house,
		       addresses.street
		FROM pharmacies LEFT JOIN addresses ON pharmacies.address_id = addresses.id
		WHERE pharmacies.id > $1
		LIMIT $2
`

	type row struct {
		PharmacyID int    `db:"pharmacy_id"`
		Name       string `db:"name"`
		IsBlocked  bool   `db:"is_blocked"`
		AddressID  int    `db:"address_id"`
		City       string `db:"city"`
		House      string `db:"house"`
		Street     string `db:"street"`
	}

	rows := make([]row, 0)
	if err := r.slave.SelectContext(ctx, &rows, query, lastPharmacyID, limit); err != nil {
		return nil, fmt.Errorf("%s: QueryError: %w", errBase, err)
	}

	result := make([]entities.Pharmacy, 0, len(rows))

	for _, row := range rows {
		result = append(result, entities.Pharmacy{
			ID:        row.PharmacyID,
			Name:      row.Name,
			IsBlocked: row.IsBlocked,
			Address: entities.Address{
				ID:     row.AddressID,
				City:   row.City,
				Street: row.Street,
				House:  row.House,
			},
		})
	}

	return result, nil
}
