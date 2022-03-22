package addressrepo

import (
	"context"
	"fmt"

	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/entities"
)

func (r *Repository) CreateAddress(ctx context.Context, address entities.Address) (int, error) {
	errBase := fmt.Sprintf("addresses.CreateAddress(%v)", address)

	const query string = `
		INSERT INTO addresses(city, street, house)
					VALUES ($1, $2, $3)
		RETURNING id
`

	var ID int
	err := r.master.QueryRowxContext(
		ctx,
		query,
		address.City,
		address.Street,
		address.House,
	).Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("%s: QueryError: %w", errBase, err)
	}

	return ID, nil
}
