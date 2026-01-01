package platform

import (
	"io"
	"net/http"
	"strings"
)

// CheckGrok checks if Grok AI is available by accessing its main page and looking for region restriction indicators.
// Returns true if no restriction message is found (available), false otherwise.
// Similar to CheckGemini, but checks for absence of block messages like in CheckOpenAI.
func CheckGrok(httpClient *http.Client) (bool, error) {
	req, err := http.NewRequest("GET", "https://grok.x.ai/", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil // Non-200 status might indicate block or redirect
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	lowerBody := strings.ToLower(string(body))
	// Check for common restriction phrases from Grok (based on user reports and searches)
	if strings.Contains(lowerBody, "not available in your region") ||
		strings.Contains(lowerBody, "not available in this region") ||
		strings.Contains(lowerBody, "regional restriction") ||
		strings.Contains(lowerBody, "unsupported country") {
		return false, nil
	}

	// Optional: If you find a positive indicator string (like in Gemini), add it here for confirmation.
	// For example: if !strings.Contains(lowerBody, "welcome | xai") { return false, nil }

	return true, nil
}
