package usecases

import (
	"o2c-graph/internal/adapter/db"
	"o2c-graph/internal/core/domain"
)

// NodeUsecase handles business logic for node details
type NodeUsecase struct {
	repo *db.NodeRepository
}

// NewNodeUsecase creates a new node usecase
func NewNodeUsecase(repo *db.NodeRepository) *NodeUsecase {
	return &NodeUsecase{repo: repo}
}

// GetNodeDetail retrieves a node with its neighbors and connecting edges
func (nu *NodeUsecase) GetNodeDetail(nodeType, nodeID string) (*domain.NodeDetailResponse, error) {
	return nu.repo.GetNodeDetail(nodeType, nodeID)
}
