package responses

type QueryResult[T any] struct {
	Data          []T `json:"data"`
	TotalElements int `json:"totalElements"`
}

func NewQueryResult[T any](totalElements int, data []T) *QueryResult[T] {
	return &QueryResult[T]{TotalElements: totalElements, Data: data}
}

// PaginatedResult holds the paginated data along with metadata.
type PaginatedResult[T any] struct {
	Data          []T   `json:"data"`
	TotalElements int64 `json:"totalElements"`
	Page          int   `json:"page"`
	Limit         int   `json:"limit"`
}
