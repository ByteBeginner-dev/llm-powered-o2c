package domain

// Node represents a vertex in the O2C graph
type Node struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties"`
}

// Edge represents a connection between two nodes
type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"`
}

// GraphResponse is the API response structure for /api/graph
type GraphResponse struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// NodeDetailResponse returns a single node with its neighbors and connecting edges
type NodeDetailResponse struct {
	Node      Node   `json:"node"`
	Neighbors []Node `json:"neighbors"`
	Edges     []Edge `json:"edges"`
}

// ChatRequest is the request body for /api/chat
type ChatRequest struct {
	Query string `json:"query"`
}

// ChatResponse is the response for /api/chat
type ChatResponse struct {
	Answer string                   `json:"answer"`
	SQL    string                   `json:"sql,omitempty"`
	Rows   []map[string]interface{} `json:"rows,omitempty"`
}
