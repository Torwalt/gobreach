package persistence

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSQL(t *testing.T) {
	sql, err := LoadSQL(Create)

	assert.Nil(t, err)
	assert.NotNil(t, sql)

	eSQL := `
	INSERT INTO Breach (email_hash, domain, breached_info, breach_date, breach_source)
VALUES (:email_hash, :domain, :breached_info, :breach_date, :breach_source)
	`
	assert.Equal(t, strings.TrimSpace(eSQL), strings.TrimSpace(sql))
}

