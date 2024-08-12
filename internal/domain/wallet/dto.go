package wallet

type TopUpRequest struct {
	Amount int64 `json:"amount" binding:"required,min=10000"`
}

type TopUpResponse struct {
	RedirectURL string `json:"redirect_url"`
}

type GetBalanceResponse struct {
	Balance int64 `json:"balance"`
}

// GetMidtransTransactionsRequest paginated
type GetMidtransTransactionsRequest struct {
	Page  int `form:"page" binding:"required,min=1"`
	Limit int `form:"limit" binding:"required,min=1,max=30"`
}
