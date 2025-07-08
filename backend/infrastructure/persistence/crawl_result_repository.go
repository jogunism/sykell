package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/domain"
)

// CrawlResultRepository defines the interface for storing CrawlResult
type CrawlResultRepository interface {
	Save(result domain.CrawlResult) error
}

// mysqlCrawlResultRepository implements CrawlResultRepository for MySQL
type mysqlCrawlResultRepository struct {
	db *sql.DB
}

// NewMySQLCrawlResultRepository creates a new MySQLCrawlResultRepository
func NewMySQLCrawlResultRepository(db *sql.DB) CrawlResultRepository {
	return &mysqlCrawlResultRepository{db: db}
}

// Save saves a CrawlResult to the database
func (r *mysqlCrawlResultRepository) Save(result domain.CrawlResult) error {
	// Convert HeadingCounts map to JSON string
	headingCountsJSON, err := json.Marshal(result.HeadingCounts)
	if err != nil {
		return fmt.Errorf("failed to marshal heading counts: %w", err)
	}

	stmt, err := r.db.Prepare(`
		INSERT INTO crawl_results (
			html_version, page_title, heading_counts,
			internal_link_count, external_link_count, inaccessible_link_count,
			has_login_form, error
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		result.HTMLVersion,
		result.PageTitle,
		headingCountsJSON,
		result.InternalLinkCount,
		result.ExternalLinkCount,
		result.InaccessibleLinkCount,
		result.HasLoginForm,
		result.Error,
	)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
