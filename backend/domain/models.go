package domain

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// NullString is a wrapper for sql.NullString that handles JSON marshalling/unmarshalling
type NullString struct {
	sql.NullString
}

// MarshalJSON implements the json.Marshaler interface.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ns *NullString) UnmarshalJSON(b []byte) error {
	var s *string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

// CrawlResult holds the data extracted from the crawled URL
type CrawlResult struct {
	ID                  int               `json:"id"` // Added ID field
	HTMLVersion         string            `json:"html_version"`
	URL                 NullString        `json:"url"`
	PageTitle           string            `json:"page_title"`
	HeadingCounts       map[string]int    `json:"heading_counts"`
	InternalLinkCount   int               `json:"internal_link_count"`
	ExternalLinkCount   int               `json:"external_link_count"`
	InaccessibleLinkCount int             `json:"inaccessible_link_count"` // Only for the main URL in this implementation
	HasLoginForm        bool              `json:"has_login_form"`
	Error               string            `json:"error"`
	CreatedAt           time.Time         `json:"created_at"` // Added CreatedAt field
}

// LoginRequest defines the structure for the login POST request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CrawlRequest defines the structure for the POST request body
type CrawlRequest struct {
	URL string `json:"url" form:"url" binding:"required"`
}

// Claims defines the structure of the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// DeleteCrawlResultsRequest defines the structure for the batch delete POST request body
type DeleteCrawlResultsRequest struct {
	IDs []int `json:"ids" binding:"required"`
}
