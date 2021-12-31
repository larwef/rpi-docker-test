package storage

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed schema/*
var content embed.FS

func migrateUp(db *sql.DB) error {
	n, err := migrate.Exec(db, "postgres", migrate.AssetMigrationSource{
		Asset:    getAsset,
		AssetDir: getAssetDir,
		Dir:      "schema",
	}, migrate.Up)
	if err != nil {
		return err
	}
	log.Printf("Database migration: %d migrations applied up\n", n)
	return nil
}

func getAsset(path string) ([]byte, error) {
	return content.ReadFile(path)
}

func getAssetDir(path string) ([]string, error) {
	var result []string
	if err := fs.WalkDir(content, path, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() {
			result = append(result, d.Name())
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}
