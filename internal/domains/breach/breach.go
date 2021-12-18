package breach

import (
	"crypto/sha256"
	"time"
)

// Breach is a domain model for account information leaked from a service.
// e.g. information about a user's account on XYZ.com
type Breach struct {
	email        string
	EmailHash    string
	Domain       string
	BreachedInfo []string
	BreachDate   time.Time
	BreachSource string
}

func New(email string, domain string, breachedInfo []string, breachDate time.Time, breachSource string) *Breach {
	return &Breach{
		email:        email,
		EmailHash:    HashEmail(email),
		Domain:       domain,
		BreachedInfo: breachedInfo,
		BreachDate:   breachDate,
		BreachSource: breachSource,
	}
}

// EmailHash returns the SHA256 hash of the email address.
func HashEmail(email string) string {
	h := sha256.New()
	h.Write([]byte(email))
	return string(h.Sum(nil))
}
