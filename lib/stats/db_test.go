package main

import (
	"testing"
)

func getDb() *Database {
	db := NewDatabase("::memory::")
	db.CreateTables()
	return db
}

func TestRecordsVisits(t *testing.T) {
	db := getDb()
	defer db.Close()

	if err := db.RecordVisit("hello", "a"); err != nil {
		t.Error(err)
	}

	cnt := db.GetVisitCount("hello")
	if cnt != 1 {
		t.Errorf("count should be 1, was %d", cnt)
	}
}

func TestMultipleVisitsSameHash(t *testing.T) {
	db := getDb()
	defer db.Close()

	p := "post"

	db.RecordVisit(p, "a")
	db.RecordVisit(p, "a")
	db.RecordVisit(p, "a")

	cnt := db.GetVisitCount(p)
	if cnt != 1 {
		t.Errorf("count should be 1, was %d", cnt)
	}

	db.RecordVisit(p, "b")
	db.RecordVisit(p, "b")
	db.RecordVisit(p, "b")

	cnt = db.GetVisitCount(p)
	if cnt != 2 {
		t.Errorf("count should be 2, was %d", cnt)
	}
}
