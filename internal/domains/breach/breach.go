package breach

import (
	"crypto/sha256"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// Breach is a domain model for account information leaked from a service/website.
// e.g. information about a user's account on XYZ.com
type Breach struct {
	email        string
	EmailHash    []byte
	Domain       string
	BreachedInfo []string
	BreachDate   time.Time
	BreachSource string
}

func New(email string, breachedInfo []string, breachDate time.Time, breachSource string) (*Breach, error) {

	err := IsEmail(email)
	if err != nil {
		return nil, NewError(BreachValidationErr, fmt.Sprintf("passed email is not valid: %v", err))
	}
	domain, err := GetDomainFromEmail(email)
	if err != nil {
		return nil, NewError(BreachValidationErr, fmt.Sprintf("could not parse domain from email: %v", err))
	}

	return &Breach{
		email:        email,
		EmailHash:    HashEmail(email),
		Domain:       domain,
		BreachedInfo: breachedInfo,
		BreachDate:   breachDate,
		BreachSource: breachSource,
	}, nil
}

// EmailHash returns the SHA256 hash of the email address.
func HashEmail(email string) []byte {
	h := sha256.New()
	h.Write([]byte(email))
	return h.Sum(nil)
}

// Returns the domain part of a valid RFC5322 email address
func GetDomainFromEmail(email string) (string, error) {
	spl := strings.Split(email, "@")

	if len(spl) != 2 {
		return "", NewError(
			BreachValidationErr, "no or more than one '@'")
	}

	return spl[1], nil
}

// Checks if email is RFC5322 compliant
func IsEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}
