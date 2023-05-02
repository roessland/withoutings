package withings

import "fmt"

type APIError int

const (
	ErrInvalidToken APIError = 401
)

func (apiErr APIError) Error() string {
	return fmt.Sprintf("Withings API returned error code %d", apiErr)
}
