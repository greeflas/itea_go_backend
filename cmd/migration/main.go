package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"log"
	"os"
)

const (
	cmdCreate  = "create"
	cmdMigrate = "migrate"
)

var Migrations = migrate.NewMigrations()

func main() {
	dsn := "postgres://postgres:pass@localhost:5432/itea_backend?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	migrator := migrate.NewMigrator(db, Migrations)

	ctx := context.Background()

	if err := migrator.Init(ctx); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("command name is missed. usage: go run cmd/migration/* [create|migrate]")
		os.Exit(1)
	}

	cmdName := os.Args[1]

	switch cmdName {
	case cmdCreate:
		name := os.Args[2]

		_, err := migrator.CreateGoMigration(ctx, name, migrate.WithPackageName("main"))
		if err != nil {
			log.Fatal(err)
		}
	case cmdMigrate:
		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Fatal(err)
		}

		if group.IsZero() {
			fmt.Printf("there are no new migrations to run (database is up to date)\n")
		} else {
			fmt.Printf("migrated to %s\n", group)
		}
	default:
		fmt.Printf("unknown command: %q", cmdName)
		os.Exit(1)
	}
}
