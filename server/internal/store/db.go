package store

import (
	"fmt"
	"net/url"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "modernc.org/sqlite" // pure-Go driver registered as "sqlite"
)

// DriverName is the database/sql driver used under GORM (modernc.org/sqlite).
const DriverName = "sqlite"

// Create opens a new writable database at file and creates the schema.
func Create(file string) (*gorm.DB, error) {
	db, err := open(file, false)
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Article{}, &SidebarDoc{}, &Asset{}); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	if err := db.Exec(ftsSchema).Error; err != nil {
		return nil, fmt.Errorf("create fts table: %w", err)
	}
	return db, nil
}

// OpenReadOnly opens an existing build database without write access.
func OpenReadOnly(file string) (*gorm.DB, error) { return open(file, true) }

func open(file string, readOnly bool) (*gorm.DB, error) {
	// A relative path would be parsed as the authority of the file: URI.
	abs, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	file = filepath.ToSlash(abs)
	q := url.Values{}
	q.Add("_pragma", "foreign_keys(1)")
	if readOnly {
		q.Set("mode", "ro")
		q.Add("_pragma", "query_only(1)")
	} else {
		q.Add("_pragma", "journal_mode(WAL)")
	}
	dsn := (&url.URL{Scheme: "file", Path: file, RawQuery: q.Encode()}).String()
	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: DriverName, DSN: dsn}), &gorm.Config{
		Logger:                 logger.Discard,
		SkipDefaultTransaction: readOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", file, err)
	}
	return db, nil
}

// Close releases the database connections.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// StripForSearch reduces a copied build database to what search needs: the FTS
// index and the article metadata. HTML, tags, assets, and the sidebar are removed.
func StripForSearch(file string) error {
	db, err := open(file, false)
	if err != nil {
		return err
	}
	defer func() { _ = Close(db) }()
	for _, stmt := range []string{
		`UPDATE articles SET html = '', tags = '[]'`,
		`DELETE FROM assets`,
		`DELETE FROM sidebar_docs`,
		`PRAGMA journal_mode = DELETE`,
		`VACUUM`,
	} {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}
	return nil
}
