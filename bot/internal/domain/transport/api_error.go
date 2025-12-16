package transport

import (
	"fmt"
)

type APIError struct {
	StatusCode int
	URL        string
	Body       []byte
}

func (e APIError) Error() string {
	return fmt.Sprintf("API request to %s failed with status %d: %s", e.URL, e.StatusCode, string(e.Body))
}
