package main

import (
	"context"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`
ALTER TABLE "public"."users" ADD COLUMN age int NULL
`)
		return err
	}, nil)
}
