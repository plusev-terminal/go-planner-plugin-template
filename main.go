package main

import (
	"github.com/extism/go-pdk"
	m "github.com/plusev-terminal/go-plugin-common/meta"
)

// main function is required but not used in WASM plugins
func main() {}

// calendarImport is a host function that allows the plugin to import events into the calendar
// This function is provided by the PlusEV planner application and is called when you want to
// add events to the user's calendar. It takes the memory offset of the ImportData JSON and
// returns the memory offset of the ImportResult JSON.
//
//go:wasmimport extism:host/user calendar_import
func calendarImport(uint64) uint64

// meta is the required export function that provides plugin metadata to the host application
// This function is called by the PlusEV planner to get information about your plugin including:
// - Plugin identification and versioning
// - Required permissions and resource access
// - Display information for the plugin marketplace
//
//go:wasmexport meta
func meta() int32 {
	pdk.OutputJSON(m.Meta{
		// Unique identifier for your plugin - should be descriptive and unique
		PluginID: "example-planner-plugin",

		// Display name shown to users in the plugin marketplace
		Name: "Example Planner Plugin",

		// Must be "plusev_planner" for planner plugins
		AppID: "plusev_planner",

		// Category helps users find your plugin. Common categories:
		// "Import" - for importing events from external sources
		// "Export" - for exporting calendar data
		// "Utility" - for calendar utilities and tools
		Category: "Import",

		// Brief description of what your plugin does
		Description: "An example plugin that demonstrates how to import events into the PlusEV planner",

		// Your name or organization
		Author: "Your Name",

		// Semantic version (major.minor.patch)
		Version: "1.0.0",

		// Optional: URL to your plugin's repository
		Repository: "https://github.com/your-username/your-plugin-repo",

		// Optional: Tags for better discoverability
		Tags: []string{"example", "demo", "template"},

		// Optional: Contact information
		Contacts: []m.AuthorContact{
			{
				Kind:  "email",
				Value: "your-email@example.com",
			},
		},

		// Resource access permissions - be specific about what your plugin needs
		Resources: m.ResourceAccess{
			// Network access - specify exact URLs/patterns your plugin will access
			// Only request access to the specific domains and endpoints you need
			AllowedNetworkTargets: []m.NetworkTargetRule{
				{Pattern: "https://jsonplaceholder.typicode.com/*"}, // Example API for demo
				// Add more patterns as needed:
				// {Pattern: "https://api.yourservice.com/*"},
			},

			// File system write access - specify directories if needed
			// nil means no file write access (recommended for most plugins)
			FsWriteAccess: nil,

			// Standard output access - usually true for logging
			StdoutAccess: true,

			// Standard error access - usually true for error logging
			StderrAccess: true,
		},
	})

	return 0
}
