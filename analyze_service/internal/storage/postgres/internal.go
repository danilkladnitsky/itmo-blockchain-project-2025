package postgres

type (
	Pagination struct {
		Limit  *int
		Offset *int
	}
	Order struct {
		Column *string
		Value  *string
	}
)
