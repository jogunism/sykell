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
	GetAll(page, pageSize int) ([]domain.CrawlResult, error)
	Delete(id int) error // Added Delete method
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

// GetAll retrieves a paginated list of CrawlResults from the database
func (r *mysqlCrawlResultRepository) GetAll(page, pageSize int) ([]domain.CrawlResult, error) {
	offset := (page - 1) * pageSize

	rows, err := r.db.Query(`
		SELECT id, html_version, page_title, heading_counts,
			internal_link_count, external_link_count, inaccessible_link_count,
			has_login_form, error, created_at
		FROM crawl_results
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query crawl results: %w", err)
	}
	defer rows.Close()

	var results []domain.CrawlResult
	for rows.Next() {
		var result domain.CrawlResult
		var headingCountsJSON []byte
		// var created_at []byte // To scan TIMESTAMP - no need to declare here, scan directly into time.Time

		err := rows.Scan(
			&result.ID,
			&result.HTMLVersion,
			&result.PageTitle,
			&headingCountsJSON,
			&result.InternalLinkCount,
			&result.ExternalLinkCount,
			&result.InaccessibleLinkCount,
			&result.HasLoginForm,
			&result.Error,
			&result.CreatedAt, // Scan directly into time.Time
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan crawl result row: %w", err)
		}

		// Unmarshal HeadingCounts JSON
		if len(headingCountsJSON) > 0 {
			err = json.Unmarshal(headingCountsJSON, &result.HeadingCounts)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal heading counts JSON: %w", err)
			}
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return results, nil
}

// Delete deletes a CrawlResult from the database by ID
func (r *mysqlCrawlResultRepository) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM crawl_results WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no record found with id %d", id)
	}

	return nil
}
