package breach_test

import (
	"encoding/hex"
	"gobreach/internal/domains/breach"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		Name  string
		Email string
	}{
		{Name: "ok", Email: "testUser@gmail.com"},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {
			b, err := breach.New(test.Email, []string{"Passwords"}, time.Now(), "h4cked_site.com")

			assert.NotNil(t, b)
			assert.Nil(t, err)
			assert.Equal(t, b.Domain, "gmail.com")
			assert.Equal(t, b.EmailHash, breach.HashEmail(test.Email))
		})
	}
}

func TestNewErrors(t *testing.T) {
	tests := []struct {
		Name  string
		Email string
	}{
		{Name: "test not email", Email: "not an email"},
		{Name: "test too many @", Email: "asd@ads.de@"},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {
			b, err := breach.New(test.Email, []string{"Passwords"}, time.Now(), "h4cked_site.com")

			assert.Nil(t, b)
			assert.NotNil(t, err)
		})
	}
}

func TestHashEmail(t *testing.T) {
	t.Run("TestHashEmail", func(t *testing.T) {
		actualHash := breach.HashEmail("testUser@gmail.com")
		expectedHash := "f7180f978fe1a7b5537c5257bd7cb737e942f336ac06cbeea0b77a45fffadea1"

		assert.NotNil(t, actualHash)
		assert.Equal(t, expectedHash, hex.EncodeToString(actualHash))
	})
}

func TestDomainFromEmail(t *testing.T) {
	tests := []struct {
		Name           string
		Email          string
		ExpectedDomain string
		Err            error
	}{
		{Name: "ok",
			Email:          "testUser@gmail.com",
			ExpectedDomain: "gmail.com",
			Err:            nil,
		},
		{
			Name:  "no @",
			Email: "testUsergmail.com",
			Err:   breach.NewError(breach.BreachValidationErr, "no or more than one '@'"),
		},
		{
			Name:  "too many @",
			Email: "testUser@gmail@.com",
			Err:   breach.NewError(breach.BreachValidationErr, "no or more than one '@'"),
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {
			domain, err := breach.GetDomainFromEmail(test.Email)

			assert.Equal(t, test.ExpectedDomain, domain)
			assert.Equal(t, test.Err, err)
		})
	}
}
