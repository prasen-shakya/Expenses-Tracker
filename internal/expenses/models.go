package expenses

type LogExpenseParams struct {
	Vendor      string  `json:"vendor"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type Expense struct {
	UserId      int64   `json:"user_id"`
	Vendor      string  `json:"vendor"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}
