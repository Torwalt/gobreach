package breach

type BreachService struct {
	byDomainRetriever ByDomainRetriever
	byEmailRetriever  ByEmailRetriever
}

func NewService(byDomainRetriever ByDomainRetriever, byEmailRetriever ByEmailRetriever) *BreachService {
	return &BreachService{
		byDomainRetriever: byDomainRetriever,
		byEmailRetriever:  byEmailRetriever,
	}
}

// Retrieve all breaches found for an email
func (s *BreachService) GetByEmail(email string) ([]Breach, *Error) {
	empty := []Breach{}
	err := IsEmail(email)
	if err != nil {
		return empty, NewErrorf(BreachValidationErr, "invalid email address: %v", email)
	}

	bS, bErr := s.byEmailRetriever.GetByEmail(email)

	if bErr != nil {
		switch code := bErr.ErrCode; code {
		case BreachNotFoundErr:
			return empty, nil
		default:
			return empty, bErr
		}
	}

	return bS, nil
}

// Retrieve all breaches found for a domain
func (s *BreachService) GetByDomain(domain string) ([]Breach, *Error) {
	return s.byDomainRetriever.GetByDomain(domain)
}
