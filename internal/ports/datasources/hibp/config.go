package hibp

const BaseURL = "https://haveibeenpwned.com/api/v3/"

type hibpConfig struct {
	BaseURL string
	APIKey  string
}

func NewhibpConfig(baseURL string, apiKey string) hibpConfig {
	if baseURL == "" {
		baseURL = BaseURL
	}
	return hibpConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}
