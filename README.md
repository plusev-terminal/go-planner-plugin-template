# PlusEV Planner Plugin Template

This template provides a comprehensive starting point for creating planner plugins for PlusEV. It demonstrates how to import events from external sources into a PlusEV planner calendar.

## 🚀 Quick Start

1. **Clone or copy this template**
2. **Update the plugin metadata** in `main.go`
3. **Implement your event import logic** in `importer.go`
4. **Build and test** your plugin

## 📁 Project Structure

```
├── main.go          # Plugin metadata and entry point
├── importer.go      # Event import implementation
├── go.mod           # Go module dependencies
├── Taskfile.yml     # Build tasks
└── README.md        # This file
```

## 🔧 Core Components

### Plugin Metadata (`main.go`)

The `meta()` function defines your plugin's identity and permissions:

```go
//go:wasmexport meta
func meta() int32 {
    pdk.OutputJSON(m.Meta{
        PluginID:    "your-plugin-id",           // Unique identifier
        Name:        "Your Plugin Name",          // Display name
        AppID:       "plusev_planner",           // Must be "plusev_planner"
        Category:    "Import",                   // Plugin category
        Description: "What your plugin does",    // Brief description
        Author:      "Your Name",               // Your name/org
        Version:     "1.0.0",                  // Semantic version
        
        // Resource permissions
        Resources: m.ResourceAccess{
            AllowedNetworkTargets: []m.NetworkTargetRule{
                {Pattern: "https://your-api.com/*"},
            },
            StdoutAccess: true,
            StderrAccess: true,
        },
    })
    return 0
}
```

### Event Import (`importer.go`)

The `import_events()` function is called when users want to import events:

```go
//go:wasmexport import_events
func import_events() int32 {
    // 1. Parse input (date range)
    job := pi.ImportJob{}
    pdk.InputJSON(&job)
    
    // 2. Fetch data from your source
    events := fetchEventsFromAPI(job.From, job.To)
    
    // 3. Import to calendar
    importEventsToCalendar(events)
    
    return 0 // Success
}
```

## 🌐 Making HTTP Requests

Use the `requester` package from [go-plugin-common](http://github.com/plusev-terminal/go-plugin-common) for HTTP calls:

```go
req := requester.Request{
    Method: "GET",
    URL:    "https://api.example.com/events",
    Headers: map[string]string{
        "Authorization": "Bearer your-token",
        "Accept":       "application/json",
    },
}

resp, err := requester.Send(&req, nil)
if err != nil {
    return err
}

// Parse response
var data YourDataType
json.Unmarshal(resp.Body, &data)
```

## 📊 Structured Logging

Use the `logging` package for debugging. The log records show up in the plugin details within PlusEV:

```go
logger := logging.NewLogger("your-plugin")

// Simple messages
logger.Info("Starting import")
logger.Error("Import failed")

// With structured data
logger.InfoWithData("Import completed", map[string]any{
    "events_imported": count,
    "duration_ms":     elapsed,
})
```

## 📅 Creating Events

Convert your data to `ImportEvent` format:

```go
event := pi.ImportEvent{
    Title:     "Meeting with Team",
    StartDate: time.Now(),
    EndDate:   time.Now().Add(time.Hour),
    Timezone:  "UTC",
    Notes:     "Quarterly planning meeting",
    Tags:      []string{"work", "planning"},
}
```

## 🔒 Security & Permissions

### Network Access
Only request access to specific URL paths that you really need:

```go
AllowedNetworkTargets: []m.NetworkTargetRule{
    {Pattern: "https://api.calendar-service.com/*"},
    {Pattern: "https://auth.calendar-service.com/token"},
}
```

### File System Access
Most plugins don't need file access:

```go
FsWriteAccess: nil, // No file write access
```

If you need file access, specify exact paths:

```go
FsWriteAccess: map[string]string{"/path/on/your/system": "/path/within/the/plugin"},
```

> These are path mappings. The key represents the path in your machines file system. The value on the other hand is the path to which to map to. For example the mapping `map[string]string{"/home/user": "/"}` would mean that in order to read or write the file `/home/user/my-file.txt` you would need to access it using the path `/my-file.txt` within the plugins code. Because `/home/user` maps to `/` within the plugins environment.

## 🛠️ Development Workflow

### 1. Setup Dependencies

Make sure you have the `go-plugin-common` module available:

```bash
go get github.com/plusev-terminal/go-plugin-common
```

### 2. Build Plugin

```bash
# Using Task (recommended)
task build

# Or manually with TinyGo
tinygo build -o plugin.wasm -target wasip1 .
```

### 3. Test Plugin

```bash
# Test with the PlusEV application
# (Specific testing instructions depend on your PlusEV setup)
```

## 🔄 Common Patterns

### Error Handling
```go
if err != nil {
    logger.ErrorWithData("Operation failed", map[string]any{
        "error": err.Error(),
        "context": "additional context",
    })
    pdk.SetError(fmt.Errorf("descriptive error: %v", err))
    return 1
}
```

**Happy coding!** 🚀

For more information about PlusEV plugins, visit the [PlusEV documentation](https://github.com/plusev-terminal/).
