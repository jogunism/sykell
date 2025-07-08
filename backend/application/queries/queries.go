package queries

type GetTestMessageQuery struct {}

type GetCrawlResultsQuery struct {
	Page     int
	PageSize int
}