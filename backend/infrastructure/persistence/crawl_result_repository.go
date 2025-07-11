package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"backend/domain"
)

// CrawlResultRepository defines the interface for storing CrawlResult
type CrawlResultRepository interface {
	Save(result domain.CrawlResult) (int, error)
	GetAll(page, pageSize int, query, sortingJson string) ([]domain.CrawlResult, int, error)
	GetTotalCount(query, sortingJson string) (int, error)
	DeleteMany(ids []int) error
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
func (r *mysqlCrawlResultRepository) Save(result domain.CrawlResult) (int, error) {
	// Convert HeadingCounts map to JSON string
	headingCountsJSON, err := json.Marshal(result.HeadingCounts)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal heading counts: %w", err)
	}

	stmt, err := r.db.Prepare(`
		INSERT INTO crawl_results (
			html_version, url, page_title, heading_counts,
			internal_link_count, external_link_count, inaccessible_link_count,
			has_login_form, error
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		result.HTMLVersion,
		result.URL.String,
		result.PageTitle,
		headingCountsJSON,
		result.InternalLinkCount,
		result.ExternalLinkCount,
		result.InaccessibleLinkCount,
		result.HasLoginForm,
		result.Error,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(id), nil
}

// GetAll retrieves a paginated list of CrawlResults from the database
func (r *mysqlCrawlResultRepository) GetAll(page, pageSize int, query, sortingJson string) ([]domain.CrawlResult, int, error) {
	offset := (page - 1) * pageSize

	baseQuery := `
		SELECT id, html_version, url, page_title, heading_counts,
			internal_link_count, external_link_count, inaccessible_link_count,
			has_login_form, error, created_at
		FROM crawl_results
	`

	var args []interface{}
	whereClause := ""

	if query != "" {
		searchQuery := "%" + query + "%"
		whereClause = " WHERE page_title LIKE ? OR url LIKE ?"
		args = append(args, searchQuery, searchQuery)
	}

	// Get total count first
	totalCount, err := r.GetTotalCount(query, sortingJson)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	orderByClause, orderArgs, err := buildOrderByClause(sortingJson)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build order by clause: %w", err)
	}

	// Append ORDER BY, LIMIT, OFFSET for the main query
	baseQuery += whereClause + orderByClause + `
		LIMIT ? OFFSET ?
	`
	args = append(args, orderArgs...)
	args = append(args, pageSize, offset)

	// log.Printf("query: %s with args: %v", baseQuery, args)
	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query crawl results: %w", err)
	}
	defer rows.Close()

	var results []domain.CrawlResult
	for rows.Next() {
		var result domain.CrawlResult
		var headingCountsJSON []byte

		err := rows.Scan(
			&result.ID,
			&result.HTMLVersion,
			&result.URL,
			&result.PageTitle,
			&headingCountsJSON,
			&result.InternalLinkCount,
			&result.ExternalLinkCount,
			&result.InaccessibleLinkCount,
			&result.HasLoginForm,
			&result.Error,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan crawl result row: %w", err)
		}

		// Unmarshal HeadingCounts JSON
		if len(headingCountsJSON) > 0 {
			err = json.Unmarshal(headingCountsJSON, &result.HeadingCounts)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal heading counts JSON: %w", err)
			}
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error after scanning rows: %w", err)
	}

	return results, totalCount, nil
}

// GetTotalCount retrieves the total number of crawl results from the database
func (r *mysqlCrawlResultRepository) GetTotalCount(query, sortingJson string) (int, error) {
	var count int
	baseQuery := "SELECT COUNT(*) FROM crawl_results"
	var args []interface{}

	if query != "" {
		searchQuery := "%" + query + "%"
		baseQuery += " WHERE page_title LIKE ? OR url LIKE ?"
		args = append(args, searchQuery, searchQuery)
	}

	orderByClause, orderArgs, err := buildOrderByClause(sortingJson)
	if err != nil {
		return 0, fmt.Errorf("failed to build order by clause for total count: %w", err)
	}
	baseQuery += orderByClause // Add order by to count query if needed (though usually not for count)
	args = append(args, orderArgs...)

	err = r.db.QueryRow(baseQuery, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total count of crawl results: %w", err)
	}
	return count, nil
}

// buildOrderByClause parses the sorting JSON and builds the ORDER BY clause
func buildOrderByClause(sortingJson string) (string, []interface{}, error) {
	orderByClause := " ORDER BY id DESC" // Default order
	var orderArgs []interface{}

	if sortingJson == "" {
		return orderByClause, orderArgs, nil
	}

	var sortingMap map[string]bool
	if err := json.Unmarshal([]byte(sortingJson), &sortingMap); err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal sorting JSON: %w", err)
	}

	var orderByParts []string
	for col, isAsc := range sortingMap {
		if isAsc {
			orderByParts = append(orderByParts, fmt.Sprintf("%s ASC", col))
		} else {
			orderByParts = append(orderByParts, fmt.Sprintf("%s DESC", col))
		}
	}

	if len(orderByParts) > 0 {
		orderByClause = " ORDER BY " + strings.Join(orderByParts, ", ")
	}

	return orderByClause, orderArgs, nil
}

// DeleteMany deletes multiple CrawlResults from the database by IDs
func (r *mysqlCrawlResultRepository) DeleteMany(ids []int) error {
	if len(ids) == 0 {
		return nil // Nothing to delete
	}

	// Build the IN clause dynamically
	placeholders := strings.Repeat("?, ", len(ids)-1) + "?"
	query := fmt.Sprintf("DELETE FROM crawl_results WHERE id IN (%s)", placeholders)

	// Convert []int to []interface{} for Exec
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare batch delete statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute batch delete statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for batch delete: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no records found for batch delete with provided IDs")
	}

	return nil
}
