package breach

type ByDomainRetriever interface {
	GetByDomain(domain string) (Breach, Error)
}

type ByEmailRetriever interface {
	GetByEmail(email string) (Breach, Error)
}
