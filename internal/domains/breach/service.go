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
	err := IsEmail(email)
	if err != nil {
		return []Breach{}, NewErrorf(BreachValidationErr, "invalid email address: %v", email)
	}

	bS, errr := s.byEmailRetriever.GetByEmail(email)

	if errr != nil {
		switch code := errr.ErrCode; code {
		case BreachNotFoundErr:
			return []Breach{}, nil
		default:
			return []Breach{}, errr
		}
	}

	return bS, nil
}

// Retrieve all breaches found for a domain
func (s *BreachService) GetByDomain(domain string) ([]Breach, *Error) {
	return s.byDomainRetriever.GetByDomain(domain)
}
