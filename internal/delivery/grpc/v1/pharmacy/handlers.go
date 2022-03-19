package pharmacy

import (
	"context"
	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/entities"
	"github.com/bhankey/pharmacy-automatization-pharmacy/pkg/api/pharmacyproto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *GRPCHandler) CreatePharmacy(ctx context.Context, req *pharmacyproto.NewPharmacy) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	err := h.pharmacySrv.CreatePharmacy(ctx, entities.Pharmacy{
		Name:      req.GetName(),
		IsBlocked: false,
		Address: entities.Address{
			City:   req.Address.GetCity(),
			Street: req.Address.GetStreet(),
			House:  req.Address.GetHouse(),
		},
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) GetPharmacies(ctx context.Context, req *pharmacyproto.PaginationRequest) (*pharmacyproto.Pharmacies, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	pharmacies, err := h.pharmacySrv.GetBatchOfPharmacies(ctx, int(req.GetLastId()), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}

	resp := make([]*pharmacyproto.Pharmacy, 0, len(pharmacies))
	for _, pharmacy := range pharmacies {
		resp = append(resp, &pharmacyproto.Pharmacy{
			Id:   int64(pharmacy.ID),
			Name: pharmacy.Name,
			Address: &pharmacyproto.Address{
				City:   pharmacy.Address.City,
				Street: pharmacy.Address.Street,
				House:  pharmacy.Address.House,
			},
		})
	}

	return &pharmacyproto.Pharmacies{
		Pharmacies: resp,
	}, nil
}
