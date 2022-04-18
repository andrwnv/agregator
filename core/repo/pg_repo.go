package repo

import (
	"github.com/andrwnv/event-aggregator/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgRepo struct {
	Database *gorm.DB
	url      string
}

func NewPgRepo() (pg *PgRepo) {
	pg = &PgRepo{}
	pg.url = "host=localhost user=postgres password=852456 dbname=event_aggregator port=5432 sslmode=disable"

	var err error
	pg.Database, err = gorm.Open(postgres.Open(pg.url), &gorm.Config{})
	if err != nil {
		utils.ReportCritical("Couldn't connect database")
	}

	return pg
}
