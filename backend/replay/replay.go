package replay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ReplayProgress represents progress update structure
type ReplayProgress struct {
	Rec       int   `json:"rec"`
	Total     int   `json:"total"`
	Timestamp int64 `json:"timestamp"`
}

// Read JSON records from file
func readRecords(filePath string) ([]map[string]interface{}, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var records []map[string]interface{}
	err = json.Unmarshal(file, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

// Send event to HEC
func sendToHEC(event map[string]interface{}, hecURL, hecToken string) error {
	jsonPayload, _ := json.Marshal(event)
	req, err := http.NewRequest("POST", hecURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Splunk "+hecToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Process and send records with updated timestamps
func ReplayRecords(filePath, hecURL, hecToken string, progressChan chan ReplayProgress) {
	records, err := readRecords(filePath)
	if err != nil {
		fmt.Println("Error reading scenario file:", err)
		close(progressChan)
		return
	}

	startTime := time.Now().UnixMilli()

	for index, record := range records {
		eventOffset, ok := record["@timestamp"].(float64)
		if !ok {
			eventOffset = 0
		}

		// eventTime := startTime + int64(eventOffset)
		eventTime := startTime - int64(eventOffset / 1000)

		if eventTime < startTime {
			fmt.Println("Warning: Adjusting eventTime, detected incorrect offset")
			fmt.Println("Base Time:", startTime)
			fmt.Println("Offset:", eventOffset)
			fmt.Println("Calculated Event Time:", eventTime, time.UnixMilli(eventTime))
			eventTime = startTime
		}

		record["@timestamp"] = eventTime

		err := sendToHEC(record, hecURL, hecToken)
		if err != nil {
			fmt.Println("Error sending record:", err)
			continue
		}

		fmt.Println("Timestamp Debug:", eventTime, time.UnixMilli(eventTime))

		progress := ReplayProgress{
			Rec:       index + 1,
			Total:     len(records),
			Timestamp: eventTime,
		}

		progressChan <- progress
	}

	close(progressChan)
}