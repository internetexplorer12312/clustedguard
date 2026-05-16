package domain

// ServerStatus represents the health state of a node.
type ServerStatus string

const (
	StatusUnknown  ServerStatus = "unknown"
	StatusOnline   ServerStatus = "online"
	StatusOffline  ServerStatus = "offline"
	StatusDegraded ServerStatus = "degraded"
)

// ServerRole describes the node's role inside a cluster.
type ServerRole string

const (
	RoleMaster ServerRole = "master"
	RoleWorker ServerRole = "worker"
	RoleAny    ServerRole = "any"
)

// Server is a managed node in the cluster.
type Server struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Host      string       `json:"host"`
	Port      int          `json:"port"`
	Role      ServerRole   `json:"role"`
	Status    ServerStatus `json:"status"`
	Tags      []string     `json:"tags"`
	CheckType string       `json:"checkType"`
	CheckPath string       `json:"checkPath"`
	LastCheck int64        `json:"lastCheck"`
	LatencyMs int64        `json:"latencyMs"`
	ClusterID string       `json:"clusterId"`
	Notes     string       `json:"notes"`

	// Agent-based monitoring
	UseAgent    bool    `json:"useAgent"`
	AgentPort   int     `json:"agentPort"`
	AgentToken  string  `json:"agentToken"`
	CpuThreshold  float64 `json:"cpuThreshold"`
	MemThreshold  float64 `json:"memThreshold"`
	DiskThreshold float64 `json:"diskThreshold"`

	// Latest metrics from agent (cached on server record)
	CpuPercent   float64 `json:"cpuPercent"`
	MemPercent   float64 `json:"memPercent"`
	DiskPercent  float64 `json:"diskPercent"`
	MemAvailBytes uint64 `json:"memAvailBytes"`
	DiskFreeBytes uint64 `json:"diskFreeBytes"`
}
