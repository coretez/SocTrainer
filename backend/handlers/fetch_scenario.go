package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	// "strconv" //Convert string
	//"backend/config"   //Not used currently
)

type SiteConfig struct {
	Site  string `json:"site"`
	Token string `json:"token"`
}
type SanitizeRule struct {
	Find    string `json:"find"`
	Replace string `json:"replace"`
}

type SanitizationRules struct {
	Replace []SanitizeRule `json:"replace"`
	Drop    []string       `json:"drop"`
	Set     []struct {
		Field   string `json:"field"`
		Replace string `json:"replace"`
	} `json:"set"`
	IPObfuscate []string `json:"ipObfuscate"`
}

type SearchConfig struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	SearchStr string `json:"search_str"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

func extractHits(response map[string]interface{}) []interface{} {
	// Check if "response" exists and is a map
	respData, ok := response["response"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing or invalid 'response' field in JSON")
		return nil
	}

	// Check if "hits" exists and is a map
	hitsData, ok := respData["hits"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing or invalid 'hits' field in JSON")
		return nil
	}

	// Check if "hits.hits" exists and is an array
	hitsInterface, ok := hitsData["hits"]
	if !ok {
		fmt.Println("Missing 'hits.hits' in response")
		return nil
	}

	hits, ok := hitsInterface.([]interface{})
	if !ok {
		fmt.Println("Invalid format for 'hits.hits', expected an array")
		return nil
	}

	return hits
}

func fetchData(site string, token string, search SearchConfig, gridAccount string) ([]interface{}, error) {
	// Validate input parameters
	if site == "" || token == "" {
		return nil, fmt.Errorf("missing site or token in siteConfig")
	}

	method := "POST"

	baseURL := site + "api/ds/get_index_zoom_histogram_lv3"
	url := baseURL
	if gridAccount != "" {
		url += "?gridaccount=" + gridAccount
	}

	// Set headers
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
		"fluencytoken": token,
	}

	// Construct the JSON payload
	payload := map[string]interface{}{
		"kargs": map[string]interface{}{
			"partition": "default",
			"dataType":  "event",
			"options": map[string]interface{}{
				"dateFacetField": "@timestamp",
				"facets": map[string]interface{}{
					"facets":         []map[string]interface{}{},
					"mustFilters":    []map[string]interface{}{},
					"mustNotFilters": []map[string]interface{}{},
					"dateFacets": []map[string]interface{}{
						{
							"name": "dateHistogram",
							"key":  "@timestamp",
						},
					},
				},
				"dataType":    "event",
				"searchStr":   search.SearchStr,
				"sortField":   "@timestamp",
				"sortOrder":   "asc",
				"range_from":  search.StartTime,
				"range_to":    search.EndTime,
				"fetchOffset": float64(0),
				"fetchLimit":  float64(1000),
			},
		},
	}

	// Convert struct to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("❌ Error marshaling JSON:", err)
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// ✅ Pretty-print JSON before sending
	formattedPayload, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println("🚀 Sending JSON Payload:\n", string(formattedPayload))

	// Create request body from JSON
	reqBody := bytes.NewReader(payloadBytes)

	// Create HTTP request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Print response status for debugging
	fmt.Printf("HTTP Status Code: %d\n", resp.StatusCode)

	// Handle HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: HTTP %d - %s", resp.StatusCode, resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if response body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body from API")
	}

	// Print raw response for debugging
	fmt.Println("Raw API Response:", string(body))

	// Parse JSON response
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract response.hits.hits
	hits := extractHits(jsonResponse)

	return hits, nil
}

func replaceFields(obj interface{}, find string, replace string) interface{} {
	// Replace occurrences of a pattern in the object.
	switch v := obj.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = replaceFields(value, find, replace)
		}
		return v
	case []interface{}:
		for i, item := range v {
			v[i] = replaceFields(item, find, replace)
		}
		return v
	case string:
		re := regexp.MustCompile(regexp.QuoteMeta(find))
		return re.ReplaceAllString(v, replace)
	default:
		return obj
	}
}

func obfuscateIP(ip string) string {
	// Obfuscate IP addresses to 10.x.x.x format for internal ranges.
	privateRanges := []string{"10.", "172.", "192.168."}
	for _, prefix := range privateRanges {
		if strings.HasPrefix(ip, prefix) {
			parts := strings.Split(ip, ".")
			if len(parts) == 4 {
				parts[1] = "x"
				parts[2] = "x"
				parts[3] = "x"
				return strings.Join(parts, ".")
			}
		}
	}
	return ip
}

func processData(data map[string]interface{}, sanitizeRules SanitizationRules) []map[string]interface{} {
	// Process data: adjust timestamps, obfuscate fields, and apply sanitization rules.
	hitsInterface, ok := data["response"].(map[string]interface{})["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		fmt.Println("No Hits")
		return nil
	}
	hits := InterfaceSlice(hitsInterface)

	processedEvents := make([]map[string]interface{}, 0)
	var firstTimestamp int64
	for i, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			fmt.Printf("Skipping bad record %d", i)
			continue
		}
		timestampFloat, ok := source["@timestamp"].(float64)
		timestamp := int64(timestampFloat)
		if !ok {
			fmt.Println("Type assertion failed")
			continue
		}
		if firstTimestamp == 0 {
			firstTimestamp = timestamp
		}
		source["@timestamp"] = firstTimestamp - timestamp
		processedEvents = append(processedEvents, source)
	}
	return processedEvents
}

func saveToFile(data []map[string]interface{}, filename string) error {
	// Save processed data to a file.
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode data to JSON: %w", err)
	}
	return nil
}

func InterfaceSlice(slice []interface{}) []interface{} {
	//This assumes everything works.
	s := make([]interface{}, len(slice))
	for i, v := range slice {
		s[i] = v
	}
	return s
}

// GetDataHandler takes searchStr, range_from, and range_to from the request body
// and uses them to make a request to the remote API.
func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 1. Decode the request body
	var req struct {
		SearchStr   string `json:"searchString"`
		RangeFrom   int64  `json:"startTime"`
		RangeTo     int64  `json:"endTime"`
		GridAccount string `json:"gridAccount"` //Added query parameter
		Site        string `json:"site"`
		Token       string `json:"token"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	data, err := fetchData(req.Site, req.Token, SearchConfig{ //All code in functions
		SearchStr: req.SearchStr,
		StartTime: req.RangeFrom,
		EndTime:   req.RangeTo,
	}, req.GridAccount)

	if err != nil {
		http.Error(w, "Failed to fetch remote", http.StatusInternalServerError)
		log.Printf("Failed to get data from remote API: %v", err)
		return
	}

	js, err := json.Marshal(data) //Now is valid json
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
