package main

import (
	"database/sql"
	"flag"
	"log"
	"path"

	"github.com/flohansen/dasher/internal/api"
	"github.com/flohansen/dasher/internal/datastore"
	"github.com/flohansen/dasher/internal/notification"
	"github.com/flohansen/dasher/internal/routes"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	_ "github.com/mattn/go-sqlite3"
)

func run() error {
	dataPath := flag.String("data", "/data", "The full path to the sqlite3 database file")
	flag.Parse()

	db, err := sql.Open("sqlite3", path.Join(*dataPath, "dasher.db"))
	if err != nil {
		return errors.Wrap(err, "sql open")
	}

	rpc := grpc.NewServer()
	store := datastore.NewSQLite(db)
	notifier := notification.NewFeatureNotifier(rpc, store)
	migrator := datastore.NewSQLMigrator(db)
	routes := routes.New(store, notifier)

	return api.New(
		api.WithLogging(),
		api.WithMigrator(migrator),
		api.WithHttpHandler(":3000", routes),
		api.WithNetListenerServer(":50051", rpc),
	).Start()
}

func main() {
	log.Fatal(run())
}
