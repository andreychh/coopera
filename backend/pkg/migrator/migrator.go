package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

type Migration struct {
	Version  string
	Name     string
	UpPath   string
	DownPath string
}

// Migrate — запускает все миграции "up"
func Migrate(path, dsn, schema string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Устанавливаем search_path на схему проекта
	if _, err := db.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s; SET search_path TO %s;`, schema, schema)); err != nil {
		return fmt.Errorf("failed to set search_path to schema %s: %w", schema, err)
	}

	if err := ensureTable(db, schema); err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}

	migs, err := loadMigrations(path)
	if err != nil {
		return err
	}

	applied, err := loadAppliedVersions(db, schema)
	if err != nil {
		return err
	}

	for _, m := range migs {
		if applied[m.Version] {
			continue
		}

		fmt.Printf("Applying migration:\n  Version: %s\n  Up: %s\n", m.Version, filepath.Base(m.UpPath))

		sqlBytes, err := ioutil.ReadFile(m.UpPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", m.UpPath, err)
		}

		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("error executing migration %s: %w", m.Version, err)
		}

		_, err = db.Exec(fmt.Sprintf(`INSERT INTO %s.schema_migrations (version) VALUES ($1)`, schema), m.Version)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", m.Version, err)
		}

		fmt.Printf("✅ Applied migration: %s (%s)\n", m.Version, m.Name)
	}

	return nil
}

// ensureTable — создаёт таблицу версий в нужной схеме
func ensureTable(db *sql.DB, schema string) error {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.schema_migrations (
			version VARCHAR(20) PRIMARY KEY
		)
	`, schema))
	return err
}

// loadAppliedVersions — возвращает уже применённые версии
func loadAppliedVersions(db *sql.DB, schema string) (map[string]bool, error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT version FROM %s.schema_migrations`, schema))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, nil
}

// loadMigrations — читает файлы миграций из каталога
func loadMigrations(path string) ([]Migration, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	m := map[string]*Migration{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		if !strings.HasSuffix(name, ".up.sql") && !strings.HasSuffix(name, ".down.sql") {
			continue
		}

		base := strings.TrimSuffix(strings.TrimSuffix(name, ".up.sql"), ".down.sql")
		parts := strings.SplitN(base, "_", 2)
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		desc := parts[1]
		fullPath := filepath.Join(path, name)

		if _, ok := m[version]; !ok {
			m[version] = &Migration{
				Version: version,
				Name:    desc,
			}
		}

		if strings.HasSuffix(name, ".up.sql") {
			m[version].UpPath = fullPath
		} else if strings.HasSuffix(name, ".down.sql") {
			m[version].DownPath = fullPath
		}
	}

	var out []Migration
	for _, v := range m {
		if v.UpPath == "" {
			return nil, errors.New("missing .up.sql for version " + v.Version)
		}
		out = append(out, *v)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Version < out[j].Version
	})

	return out, nil
}
