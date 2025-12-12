package constants

import "time"

const (
	// APITimeout is the default timeout for API calls
	APITimeout = 30 * time.Second
	// APIModel is the Gemini model to use
	APIModel = "gemini-2.5-flash-lite"
)
