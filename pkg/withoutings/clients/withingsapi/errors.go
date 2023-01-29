package withingsapi

import "fmt"

type APIError int

func (apiErr APIError) Error() string {
	return fmt.Sprintf("Withings API returned error code %d", apiErr)
}
