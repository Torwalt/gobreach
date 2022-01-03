package persistence

import (
	"database/sql"
	"fmt"
	"gobreach/cmd/server/config"
	"gobreach/internal/domains/breach"
	"testing"
	"time"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	_ "github.com/lib/pq"
)

func TestBreachDBCreate(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	defer mockDb.Close()

	assert.Nil(t, err)

	bDB := NewBreachDB(mockDb)
	breach, err := breach.New(
		"test@test.de",
		[]string{"test"},
		time.Now(),
		"test hack",
	)
	mock.ExpectBegin()
	expectedSQL := fmt.Sprintf(
		`INSERT INTO Breach (email_hash, domain, breached_info, breach_date, breach_source) 
		 VALUES (%v, %v, %v, %v, %v)`,
		breach.EmailHash, breach.Domain, breach.BreachedInfo, breach.BreachDate, breach.BreachSource)

	mock.ExpectExec(expectedSQL)
	bDB.Create(breach)
}

func TestBreachDBGetByEmail(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	defer mockDb.Close()

	assert.Nil(t, err)

	bDB := NewBreachDB(mockDb)
	breach, err := breach.New(
		"test@test.de",
		[]string{"test"},
		time.Now(),
		"test hack",
	)

	mock.ExpectBegin()
	expectedSQL := fmt.Sprintf(`
		SELECT
			email_hash,
			domain,
			breached_info,
			breach_date,
			breach_source
		FROM
			Breach
		WHERE
			email_hash = %v
	`, breach.EmailHash)
	expectedRows := sqlmock.NewRows(
		[]string{"email_hash", "domain", "breachd_info", "breached_date", "breach_source"}).
		AddRow(breach.EmailHash, breach.Domain, strings.Join(breach.BreachedInfo, ","), breach.BreachDate, breach.BreachSource)

	mock.ExpectPrepare(expectedSQL).ExpectQuery().WithArgs(breach.EmailHash).WillReturnRows(expectedRows)
	bDB.Create(breach)
}

/*   integration tests   */

func TestBreachDBCreateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	conf := config.FromEnv()
	connS := fmt.Sprintf(
		PgDbConnS, conf.PostgresHost, conf.PostgresPort, conf.PostgresUser, conf.PostgresSecret, conf.PostgresName)

	db, err := sql.Open("postgres", connS)
	assert.Nil(t, err)

	bDB := NewBreachDB(db)
	breach, err := breach.New(
		"test@test.de",
		[]string{"test"},
		time.Now(),
		"test hack",
	)
	err = bDB.Create(breach)
	assert.Nil(t, err)
}
