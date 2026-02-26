package queries

type GetNotCompletedOrdersQuery struct{}

func NewGetNotCompletedOrdersQuery() (GetNotCompletedOrdersQuery, error) {
	return GetNotCompletedOrdersQuery{}, nil
}
