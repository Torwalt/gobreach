package persistence

import (
	"database/sql"
	"fmt"
	"gobreach/internal/domains/breach"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

const PgDbConnS = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

type BreachDB struct {
	Conn *sqlx.DB
}

type dbBreach struct {
	EmailHash    []byte    `db:"email_hash"`
	Domain       string    `db:"domain"`
	BreachedInfo string    `db:"breached_info"`
	BreachDate   time.Time `db:"breach_date"`
	BreachSource string    `db:"breach_source"`
}

func fromBreach(b *breach.Breach) *dbBreach {
	bI := strings.Join(b.BreachedInfo[:], ",")
	return &dbBreach{
		EmailHash:    b.EmailHash,
		Domain:       b.Domain,
		BreachedInfo: bI,
		BreachDate:   b.BreachDate,
		BreachSource: b.BreachSource,
	}
}

func toBreach(dbB *dbBreach) (*breach.Breach, *breach.Error) {
	bI := strings.Split(dbB.BreachedInfo, ",")
	b, err := breach.New(
		"",
		bI,
		dbB.BreachDate,
		dbB.BreachSource,
	)
	if err != nil {
		return nil, breach.NewErrorf(breach.BreachRepositoryErr, "could not cast to Breach: %v", err)
	}
	return b, nil

}

func NewBreachDB(db *sql.DB) *BreachDB {
	xdb := sqlx.NewDb(db, "postgres")
	return &BreachDB{Conn: xdb}
}

func (db *BreachDB) Create(b *breach.Breach) *breach.Error {
	tx, err := db.Conn.Beginx()
	if err != nil {
		return breach.NewErrorf(breach.BreachRepositoryErr, "could not connect to db: %v", err)
	}

	dbb := fromBreach(b)
	sql, err := LoadSQL(Create)
	if err != nil {
		return breach.NewErrorf(breach.BreachRepositoryErr, "could not load create sql: %v", err)
	}

	fmt.Print(dbb)
	fmt.Print(sql)
	_, err = tx.NamedExec(sql, dbb)
	if err != nil {
		return breach.NewErrorf(breach.BreachRepositoryErr, "error when creating breach: %v", err)
	}

	return nil
}

func (db *BreachDB) Update(b *breach.Breach) error {
	return nil
}

func (db *BreachDB) GetByEmailHash(h []byte) (*breach.Breach, error) {
	tx, err := db.Conn.Beginx()
	if err != nil {
		return nil, breach.NewErrorf(breach.BreachRepositoryErr, "could not connect to db: %v", err)
	}

	sqlS, err := LoadSQL(GetByEmail)
	if err != nil {
		return nil, breach.NewErrorf(breach.BreachRepositoryErr, "could not load GetByEmail sql: %v", err)
	}

	row := tx.QueryRow(sqlS, h)
	dbB := &dbBreach{}
	err = row.Scan(*&dbB.BreachDate)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, breach.NewErrorf(breach.BreachRepositoryErr, "could not scan row : %v", err)
	}

	b, err := toBreach(dbB)
	if err != nil {
		return nil, breach.NewErrorf(breach.BreachRepositoryErr, "error when creating Breach: %v", err)
	}
	return b, nil
}

func (db *BreachDB) GetByDomain(d string) (*[]breach.Breach, error) {
	return &[]breach.Breach{}, nil
}
