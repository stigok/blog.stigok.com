package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Database *sql.DB
	Filename string
}

func NewDatabase(filename string) *Database {
	db := Database{
		Filename: filename,
	}

	file, err := os.Create(db.Filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Printf("%s created", db.Filename)

	sqliteDatabase, _ := sql.Open("sqlite3", db.Filename)

	db.Database = sqliteDatabase
	return &db
}

func (db *Database) Close() {
	db.Database.Close()
}

func (db *Database) CreateTables() {
	var statements []string = []string{
		`CREATE TABLE IF NOT EXISTS visits (
			"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			"post" TEXT,
			"hash" TEXT,
			"date" DATETIME DEFAULT (DATETIME('now'))
		);
		CREATE INDEX IF NOT EXISTS idx_visits_post
		ON visits (post);
		`,
	}

	for _, stmt := range statements {
		s, err := db.Database.Prepare(stmt)
		if err != nil {
			log.Fatal(err.Error())
		}
		s.Exec()
	}
}

func (db *Database) RecordVisit(post, hash string) error {
	stmt, err := db.Database.Prepare(
		`INSERT INTO visits (post, hash) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("Failed to record visit: %w", err)
	}
	stmt.Exec(post, hash)
	return nil
}

func (db *Database) GetVisitCount(post string) int {
	stmt, err := db.Database.Prepare(
		`SELECT COUNT(DISTINCT hash)
		 FROM visits
		 WHERE post=? GROUP BY post;
	`)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	err = stmt.QueryRow(post).Scan(&count)

	if err == sql.ErrNoRows {
		log.Printf("post not found: %s", post)
	} else if err != nil {
		log.Printf("failed to get post count for %s: %s", post, err)
	}

	return count
}
