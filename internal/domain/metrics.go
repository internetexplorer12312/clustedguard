package domain

// MetricsSample is a point-in-time resource snapshot for charts.
type MetricsSample struct {
	ServerID      string  `json:"serverId"`
	Timestamp     int64   `json:"timestamp"`
	CPUPercent    float64 `json:"cpuPercent"`
	MemPercent    float64 `json:"memPercent"`
	DiskPercent   float64 `json:"diskPercent"`
	MemAvailBytes uint64  `json:"memAvailBytes"`
	DiskFreeBytes uint64  `json:"diskFreeBytes"`
}

// AlertKind identifies which resource triggered the alert.
type AlertKind string

const (
	AlertCPU   AlertKind = "cpu"
	AlertMemory AlertKind = "memory"
	AlertDisk  AlertKind = "disk"
)

// Alert is a threshold violation notification.
type Alert struct {
	ID        string    `json:"id"`
	ServerID  string    `json:"serverId"`
	ServerName string   `json:"serverName"`
	Kind      AlertKind `json:"kind"`
	Value     float64   `json:"value"`
	Threshold float64   `json:"threshold"`
	Message   string    `json:"message"`
	CreatedAt int64     `json:"createdAt"`
	Read      bool      `json:"read"`
}
