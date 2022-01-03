package persistence

import (
	"fmt"
	"os"
	"path/filepath"
)

// const sqlPath = "gobreach/internal/ports/persistence/sql"
const sqlPath = "./sql"

type Statement string

var (
	Create     Statement = "create.sql"
	GetByEmail Statement = "get_by_email.sql"
)

func LoadSQL(s Statement) (string, error) {
	p := fmt.Sprintf("%v/%v", sqlPath, s)
	fp, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

	d, err := os.ReadFile(fp)
	if err != nil {
		return "", err
	}
	sql := string(d)

	return sql, nil
}
