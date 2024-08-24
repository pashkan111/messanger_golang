package entities

type QueryParams struct {
	// QueryParams is a struct for query parameters
	// for pagination
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
