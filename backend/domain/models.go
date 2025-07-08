package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

// CrawlResult holds the data extracted from the crawled URL
type CrawlResult struct {
	HTMLVersion         string            `json:"html_version"`
	PageTitle           string            `json:"page_title"`
	HeadingCounts       map[string]int    `json:"heading_counts"`
	InternalLinkCount   int               `json:"internal_link_count"`
	ExternalLinkCount   int               `json:"external_link_count"`
	InaccessibleLinkCount int             `json:"inaccessible_link_count"` // Only for the main URL in this implementation
	HasLoginForm        bool              `json:"has_login_form"`
	Error               string            `json:"error,omitempty"`
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
