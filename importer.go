package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/extism/go-pdk"
	"github.com/plusev-terminal/go-plugin-common/logging"
	pi "github.com/plusev-terminal/go-plugin-common/planner/import"
	"github.com/plusev-terminal/go-plugin-common/requester"
)

// DemoPost represents a post from JSONPlaceholder API (used for demonstration)
// In a real plugin, replace this with the structure that matches your data source
type DemoPost struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

// import_events is the main export function that handles importing events into the calendar
// This function is called by the PlusEV planner when the user wants to import events.
// The host application will pass an ImportJob JSON containing the date range to import.
//
//go:wasmexport import_events
func import_events() int32 {
	// Initialize logger with your plugin name
	// The logger helps you debug issues and provides structured logging
	logger := logging.NewLogger("example-planner-plugin")

	// Parse the input JSON from the host application
	// The host sends an ImportJob with From and To dates
	job := pi.ImportJob{}
	err := pdk.InputJSON(&job)
	if err != nil {
		logger.ErrorWithData("Failed to parse input JSON", map[string]any{
			"error": err.Error(),
		})
		return 1 // Return 1 to indicate error
	}

	// Log the import job details for debugging
	logger.InfoWithData("Starting event import", map[string]any{
		"from": job.From.Format("2006-01-02"),
		"to":   job.To.Format("2006-01-02"),
	})

	// Example: Fetch data from a demo API (JSONPlaceholder)
	// In a real plugin, replace this with your actual data source
	events, err := fetchDemoEvents(logger, job.From, job.To)
	if err != nil {
		logger.ErrorWithData("Failed to fetch events", map[string]any{
			"error": err.Error(),
		})
		pdk.SetError(fmt.Errorf("failed to fetch events: %v", err))
		return 1
	}

	// Import the events using the host calendar_import function
	success, err := importEventsToCalendar(logger, events)
	if err != nil {
		logger.ErrorWithData("Failed to import events to calendar", map[string]any{
			"error": err.Error(),
		})
		return 1
	}

	if !success {
		logger.Error("Calendar import was not successful")
		return 1
	}

	logger.InfoWithData("Successfully imported events", map[string]any{
		"count": len(events),
	})

	return 0 // Return 0 to indicate success
}

// fetchDemoEvents demonstrates how to fetch data from an external API
// Replace this function with your own logic to fetch events from your data source
func fetchDemoEvents(logger *logging.Logger, from, to time.Time) ([]pi.ImportEvent, error) {
	// Create HTTP request to fetch demo data
	// In this example, we're using JSONPlaceholder which provides fake data for testing
	req := requester.Request{
		Method: "GET",
		URL:    "https://jsonplaceholder.typicode.com/posts",
		Headers: map[string]string{
			"Accept":     "application/json",
			"User-Agent": "PlusEV-Plugin/1.0",
		},
	}

	// Send the HTTP request
	resp, err := requester.Send(&req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}

	// Parse the response
	var posts []DemoPost
	err = json.Unmarshal(resp.Body, &posts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Convert the external data to ImportEvent format
	events := []pi.ImportEvent{}

	// For demo purposes, we'll create events from the posts
	// In a real plugin, you would parse your actual event data here
	for i, post := range posts {
		// Limit to first 5 posts for demo (you don't want 100 events!)
		if i >= 5 {
			break
		}

		// Create a demo event based on the post
		// In a real plugin, you would extract actual event information
		eventDate := from.AddDate(0, 0, i) // Spread events across the date range

		event := pi.ImportEvent{
			Title:     fmt.Sprintf("Demo Event: %s", truncateString(post.Title, 50)),
			StartDate: eventDate,
			EndDate:   eventDate.Add(time.Hour), // 1-hour events
			Timezone:  "UTC",
			Notes:     fmt.Sprintf("Demo event created from post ID %d. Content: %s", post.ID, truncateString(post.Body, 100)),
			Tags:      []string{"demo", "example"}, // Optional tags for categorization
		}

		events = append(events, event)
	}

	return events, nil
}

// importEventsToCalendar sends the events to the host application for import
// This function handles the communication with the PlusEV planner's calendar system
func importEventsToCalendar(logger *logging.Logger, events []pi.ImportEvent) (bool, error) {
	// Create the import data structure
	importData := pi.ImportData{
		Events: events,
	}

	// Allocate memory and marshal the events to JSON
	// This is required for WASM memory management
	mem, err := pdk.AllocateJSON(importData)
	if err != nil {
		return false, fmt.Errorf("failed to allocate memory for import data: %v", err)
	}

	// Call the host function to import events
	// The calendarImport function is provided by the PlusEV planner
	ptr := calendarImport(mem.Offset())

	// Read the response from the host
	resultMem := pdk.FindMemory(ptr)
	respData := resultMem.ReadBytes()

	// Parse the import result
	var result pi.ImportResult
	if err := json.Unmarshal(respData, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal import result: %v", err)
	}

	// Check if the import was successful
	if !result.Success {
		return false, fmt.Errorf("calendar import failed: %s", result.Error)
	}

	return true, nil
}

// Helper Functions
// ================

// truncateString limits a string to a maximum length and adds "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// mapValue is a type constraint for values that can be extracted from JSON maps
type mapValue interface {
	float64 | string | bool
}

// ifThen returns trueValue if condition is true, otherwise falseValue
// This is a generic ternary operator function
func ifThen[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// getValue safely extracts a value from a map with optional default value
// This is useful when parsing JSON data that might have missing fields
func getValue[T mapValue](key string, data map[string]any, defaultValue ...T) T {
	value, ok := data[key].(T)
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		var zero T
		return zero
	}

	var zero T
	if len(defaultValue) > 0 && value == zero {
		return defaultValue[0]
	}

	return value
}

// anyMatches checks if any value in the slice matches the predicate function
func anyMatches[T comparable](predicate func(T) bool, values ...T) bool {
	for _, v := range values {
		if predicate(v) {
			return true
		}
	}
	return false
}
