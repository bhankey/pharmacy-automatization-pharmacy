package container

import (
	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/adapter/repository/addressrepo"
	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/adapter/repository/pharmacyrepo"
	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/delivery/grpc/v1/pharmacy"
	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/service/pharmacyservice"
)

func (c *Container) GetPharmacyGRPCHandler() *pharmacy.GRPCHandler {
	const key = "PharmacyGRPCHandler"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*pharmacy.GRPCHandler)
		if ok {
			return typedDependency
		}
	}

	typedDependency := pharmacy.NewPharmacyGRPCHandler(c.getPharmacySrv())

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getPharmacySrv() *pharmacyservice.Service {
	const key = "PharmacySrv"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*pharmacyservice.Service)
		if ok {
			return typedDependency
		}
	}

	typedDependency := pharmacyservice.NewPharmacyService(
		c.getPharmacyStorage(),
		c.getAddressStorage(),
	)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getAddressStorage() *addressrepo.Repository {
	const key = "AddressStorage"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*addressrepo.Repository)
		if ok {
			return typedDependency
		}
	}

	typedDependency := addressrepo.NewAddressRepo(c.masterPostgresDB, c.slavePostgresDB)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getPharmacyStorage() *pharmacyrepo.Repository {
	const key = "PharmacyStorage"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*pharmacyrepo.Repository)
		if ok {
			return typedDependency
		}
	}

	typedDependency := pharmacyrepo.NewPharmacyRepo(c.masterPostgresDB, c.slavePostgresDB)

	c.dependencies[key] = typedDependency

	return typedDependency
}
