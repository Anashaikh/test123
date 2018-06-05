# go-cag

go-cag is a client library for accessing the CAG API.

For documentation of the CAG API, see https://webhookv3.alerting.monitoring.bskyb.com/api/doc.

## Usage

```go
import "github.com/sky-uk/go-cag"
```

Create a new CAG client.

```go
cagClient := cag.NewClient(
	*cagURL,
	*cagAPIKey,
	*sslSkipVerify,
)
```

Create an alert.

```go
cagAlert := cagClient.NewAlert(&cag.AlertData{
        MonitoredItem: "test.hostname",
        AlertSummary: "Memory Spike",
        DetailedDescription: "The memory on node: test.hostname has spiked to 150% usage, please check this out",
        AssignTo: "Your Spark Team",
        Severity: "warn",
        HelpURL: "https://your-helpful-docs",
        MonitoringGroup: "My System",
        MonitoringSystem: "My Amazing Client",
})

cagAlertResponse, statusCode, err := cagAlert.Create()
if err != nil {
    fmt.Printf("Failed to create CAG alert, recieved status code \"%i\" and error: %e\n", statusCode, err)
    exit 1
} else {
    fmt.Printf("Created CAG alert: %v\n", cagAlertResponse)
}
```

## Supported Resources

Support for resources is being added as needed.

- [x] Status
- [x] Heartbeat
- [x] History
- [x] Permissions
- [x] Submit
