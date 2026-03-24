package usecases

import (
	"fmt"
	"o2c-graph/internal/adapter/db"
	"o2c-graph/internal/core/domain"
)

// GraphUsecase handles all graph-related business logic
type GraphUsecase struct {
	repo *db.GraphRepository
}

// NewGraphUsecase creates a new instance of GraphUsecase
func NewGraphUsecase(repo *db.GraphRepository) *GraphUsecase {
	return &GraphUsecase{repo: repo}
}

// GetGraph retrieves all nodes and edges from the database
func (gu *GraphUsecase) GetGraph() (*domain.GraphResponse, error) {
	nodes := []domain.Node{}
	edges := []domain.Edge{}

	// 1. Get Sales Order nodes
	soNodes, err := gu.repo.GetSalesOrderNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sales order nodes: %w", err)
	}
	nodes = append(nodes, soNodes...)

	// 2. Get Delivery nodes
	deliveryNodes, err := gu.repo.GetDeliveryNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delivery nodes: %w", err)
	}
	nodes = append(nodes, deliveryNodes...)

	// 3. Get Billing Document nodes
	billingNodes, err := gu.repo.GetBillingDocumentNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch billing document nodes: %w", err)
	}
	nodes = append(nodes, billingNodes...)

	// 4. Get Payment nodes
	paymentNodes, err := gu.repo.GetPaymentNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment nodes: %w", err)
	}
	nodes = append(nodes, paymentNodes...)

	// 5. Get Customer nodes
	customerNodes, err := gu.repo.GetCustomerNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customer nodes: %w", err)
	}
	nodes = append(nodes, customerNodes...)

	// 6. Get Product nodes
	productNodes, err := gu.repo.GetProductNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product nodes: %w", err)
	}
	nodes = append(nodes, productNodes...)

	// 7. Get Plant nodes
	plantNodes, err := gu.repo.GetPlantNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plant nodes: %w", err)
	}
	nodes = append(nodes, plantNodes...)

	// 8. Build edges
	soDeliveryEdges, err := gu.repo.GetSalesOrderDeliveryEdges()
	if err == nil {
		edges = append(edges, soDeliveryEdges...)
	}

	deliveryBillingEdges, err := gu.repo.GetDeliveryBillingEdges()
	if err == nil {
		edges = append(edges, deliveryBillingEdges...)
	}

	billingPaymentEdges, err := gu.repo.GetBillingPaymentEdges()
	if err == nil {
		edges = append(edges, billingPaymentEdges...)
	}

	customerSOEdges, err := gu.repo.GetCustomerSOEdges()
	if err == nil {
		edges = append(edges, customerSOEdges...)
	}

	soProductEdges, err := gu.repo.GetSOProductEdges()
	if err == nil {
		edges = append(edges, soProductEdges...)
	}

	return &domain.GraphResponse{
		Nodes: nodes,
		Edges: edges,
	}, nil
}
