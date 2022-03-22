package pharmacy

import (
	"context"

	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/entities"
	"github.com/bhankey/pharmacy-automatization-pharmacy/pkg/api/pharmacyproto"
)

type GRPCHandler struct {
	pharmacyproto.UnimplementedPharmacyServiceServer // Must be

	pharmacySrv Srv
}

type Srv interface {
	GetBatchOfPharmacies(
		ctx context.Context,
		lastPharmacyID int,
		limit int,
	) ([]entities.Pharmacy, error)
	CreatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error
}

func NewPharmacyGRPCHandler(pharmacySrv Srv) *GRPCHandler {
	return &GRPCHandler{
		UnimplementedPharmacyServiceServer: pharmacyproto.UnimplementedPharmacyServiceServer{},
		pharmacySrv:                        pharmacySrv,
	}
}
