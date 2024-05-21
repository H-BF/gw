package casbinprovider

import (
	"os"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2/persist"
)

func newPGAdapter() (persist.Adapter, error) {
	dsn := os.Getenv("PG_DSN")
	return pgadapter.NewAdapter(dsn)
}
