package main

import (
	"testing"
	"time"
)

func getDb() *Database {
	db := OpenDatabase(MEMORY_DB)
	return db
}

func TestNewDatabase(t *testing.T) {
	db := getDb()
	defer db.Close()

	if db.Filename != MEMORY_DB {
		t.Fatal("Database not in memory")
	}

	db.Close()
}

func TestRecordsVisits(t *testing.T) {
	db := getDb()
	defer db.Close()

	v := Visit{"hello", "a", time.Now()}
	if err := db.recordVisit(v); err != nil {
		t.Error(err)
	}

	cnt := db.GetVisitCount(v.Post)
	if cnt != 1 {
		t.Errorf("count should be 1, was %d", cnt)
	}
}

func TestMultipleVisitsSameHash(t *testing.T) {
	db := getDb()
	defer db.Close()

	v := Visit{"post", "a", time.Now()}

	db.recordVisit(v)
	db.recordVisit(v)
	db.recordVisit(v)

	cnt := db.GetVisitCount(v.Post)
	if cnt != 1 {
		t.Errorf("count should be 1, was %d", cnt)
	}

	v.Hash = "b"
	db.recordVisit(v)
	db.recordVisit(v)
	db.recordVisit(v)

	cnt = db.GetVisitCount(v.Post)
	if cnt != 2 {
		t.Errorf("count should be 2, was %d", cnt)
	}
}
