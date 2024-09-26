package responses

type QueryResult[T any] struct {
	Data          []T `json:"data"`
	TotalElements int `json:"totalElements"`
}

func NewQueryResult[T any](totalElements int, data []T) *QueryResult[T] {
	return &QueryResult[T]{TotalElements: totalElements, Data: data}
}
