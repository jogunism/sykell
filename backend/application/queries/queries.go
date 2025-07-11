package queries

type GetTestMessageQuery struct {}

type GetCrawlResultsQuery struct {
	CurrPage int
	PageSize int
	Query    string
}