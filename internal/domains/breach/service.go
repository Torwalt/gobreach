package breach

type BreachService struct {
	ByDomainRetriever ByDomainRetriever
	ByEmailRetriever  ByEmailRetriever
}

func NewService(byDomainRetriever ByDomainRetriever, byEmailRetriever ByEmailRetriever) *BreachService {
	return &BreachService{
		ByDomainRetriever: byDomainRetriever,
		ByEmailRetriever:  byEmailRetriever,
	}
}

func (s *BreachService) GetByEmail(email string) (Breach, Error) {
	return s.ByEmailRetriever.GetByEmail(email)
}

func (s *BreachService) GetByDomain(domain string) (Breach, Error) {
	return s.ByDomainRetriever.GetByDomain(domain)
}
