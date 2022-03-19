package pharmacyrepo

import "github.com/jmoiron/sqlx"

type Repository struct {
	master *sqlx.DB
	slave  *sqlx.DB
}

func NewPharmacyRepo(master *sqlx.DB, slave *sqlx.DB) *Repository {
	return &Repository{
		master: master,
		slave:  slave,
	}
}
