package agentapi

// Metrics is the JSON payload returned by the agent /metrics endpoint.
type Metrics struct {
	Timestamp            int64   `json:"timestamp"`
	CPUPercent           float64 `json:"cpuPercent"`
	MemoryUsedPercent    float64 `json:"memoryUsedPercent"`
	MemoryAvailableBytes uint64  `json:"memoryAvailableBytes"`
	MemoryTotalBytes     uint64  `json:"memoryTotalBytes"`
	DiskUsedPercent      float64 `json:"diskUsedPercent"`
	DiskFreeBytes        uint64  `json:"diskFreeBytes"`
	DiskTotalBytes       uint64  `json:"diskTotalBytes"`
}
