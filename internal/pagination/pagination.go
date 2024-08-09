package pagination

type Pagination struct {
	TotalData   int  `json:"total_data"`
	CurrentPage int  `json:"current_page"`
	TotalPage   int  `json:"total_page"`
	PerPage     int  `json:"per_page"`
	NextPage    *int `json:"next_page"`
	PrevPage    *int `json:"prev_page"`
}

type GetResourcePaginatedResponse struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func NewPagination(totalData int, currentPage, perPage int) Pagination {
	totalPage := totalData / perPage
	if totalData%perPage != 0 {
		totalPage++
	}

	var nextPage *int
	if currentPage < totalPage {
		nextPage = new(int)
		*nextPage = currentPage + 1
	}

	var prevPage *int
	if currentPage > 1 {
		prevPage = new(int)
		*prevPage = currentPage - 1
	}

	return Pagination{
		TotalData:   totalData,
		CurrentPage: currentPage,
		TotalPage:   totalPage,
		PerPage:     perPage,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}
