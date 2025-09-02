# Getting Started with PlusEV Planner Plugins

This guide will walk you through creating your first PlusEV planner plugin using this template.

## Prerequisites

- Go 1.24+ installed
- TinyGo installed for WASM compilation
- Task runner (optional but recommended)

## Step 1: Customize Plugin Metadata

Edit `main.go` and update the plugin metadata:

```go
pdk.OutputJSON(m.Meta{
    PluginID:    "my-awesome-plugin",           // Change this!
    Name:        "My Awesome Plugin",           // Change this!
    Description: "Import events from my API",   // Change this!
    Author:      "Your Name",                   // Change this!
    // ... update other fields as needed
})
```

## Step 2: Update Network Permissions

Update the `AllowedNetworkTargets` to match your API:

```go
AllowedNetworkTargets: []m.NetworkTargetRule{
    {Pattern: "https://your-api.com/*"},
    {Pattern: "https://auth.your-api.com/*"},
},
```

## Step 3: Implement Your Import Logic

In `importer.go`, replace the demo API call with your actual data source:

```go
func fetchDemoEvents(logger *logging.Logger, from, to time.Time) ([]pi.ImportEvent, error) {
    // Replace this with your API call
    req := requester.Request{
        Method: "GET",
        URL:    "https://your-api.com/events",
        Headers: map[string]string{
            "Authorization": "Bearer YOUR_TOKEN",
        },
    }
    
    resp, err := requester.Send(&req, nil)
    // ... parse your response format
}
```

## Step 4: Parse Your Data Format

Update the parsing logic to match your data format:

```go
// Example for different formats:

// JSON API response
type YourEvent struct {
    Name      string `json:"name"`
    StartTime string `json:"start_time"`
    EndTime   string `json:"end_time"`
}

// Parse and convert to ImportEvent
event := pi.ImportEvent{
    Title:     yourEvent.Name,
    StartDate: parseTime(yourEvent.StartTime),
    EndDate:   parseTime(yourEvent.EndTime),
    Timezone:  "UTC",
}
```

## Step 5: Build the wasm file

```bash
# Build the plugin
task build

# Or manually
tinygo build -o plugin.wasm -target wasip1 main.go importer.go
```

## Step 6: Test with PlusEV

1. Load your plugin in PlusEV
2. Try importing events for a small date range
3. Check the logs for any errors
4. Verify events appear in the calendar

## Common Issues and Solutions

### Network Permission Denied
- Check that your API URLs match the `AllowedNetworkTargets` patterns
- Use exact patterns, wildcards only work at the end

### No Events Imported
- Check your date parsing logic
- Verify the API returns data for your date range
- Use logging to debug the parsing process

### Build Errors
- Make sure you're using TinyGo, not regular Go
- Check that all imports are available
- Verify you're targeting `wasip1`

Happy plugin development! 🚀
