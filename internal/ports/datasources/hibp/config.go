package hibp

type hibpConfig struct {
	host       string
	apiKey     string
	maxRetries int
}

func NewhibpConfig(host, apiKey string, maxRetries int) hibpConfig {
	return hibpConfig{
		host:       host,
		apiKey:     apiKey,
		maxRetries: maxRetries,
	}
}
