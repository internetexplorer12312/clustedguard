package domain

// Cluster groups servers for unified monitoring and control.
type Cluster struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ServerIDs   []string `json:"serverIds"`
	CreatedAt   int64    `json:"createdAt"`
	UpdatedAt   int64    `json:"updatedAt"`
}
