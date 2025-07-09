package commands

type CrawlCommand struct {
	URL string
}

type DeleteCrawlResultsCommand struct {
	IDs []int
}
