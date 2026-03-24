package db

import (
	"database/sql"
	"fmt"
	"o2c-graph/internal/core/domain"
	"strings"
)

// NodeRepository handles database queries for individual nodes
type NodeRepository struct {
	db *sql.DB
}

// NewNodeRepository creates a new node repository
func NewNodeRepository(db *sql.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

// stripPrefix removes the node ID prefix (e.g. "DEL_", "SO_", "BIL_") sent by the frontend
func stripPrefix(nodeID string) string {
	prefixes := []string{"SO_", "DEL_", "BIL_", "PAY_", "CUST_", "PRD_", "PLANT_"}
	for _, p := range prefixes {
		if strings.HasPrefix(nodeID, p) {
			return nodeID[len(p):]
		}
	}
	return nodeID
}

// GetNodeDetail retrieves a single node and its neighbors
func (nr *NodeRepository) GetNodeDetail(nodeType, nodeID string) (*domain.NodeDetailResponse, error) {
	nodeID = stripPrefix(nodeID)
	node, err := nr.getNodeByTypeAndID(nodeType, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	if node == nil {
		return nil, fmt.Errorf("node not found: %s_%s", nodeType, nodeID)
	}

	// Get neighboring nodes
	neighbors, err := nr.getNeighbors(nodeType, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get neighbors: %w", err)
	}

	// Get edges connecting this node to neighbors
	edges, err := nr.getConnectingEdges(nodeType, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get edges: %w", err)
	}

	return &domain.NodeDetailResponse{
		Node:      *node,
		Neighbors: neighbors,
		Edges:     edges,
	}, nil
}

// getNodeByTypeAndID retrieves a single node by type and ID
func (nr *NodeRepository) getNodeByTypeAndID(nodeType, nodeID string) (*domain.Node, error) {
	switch strings.ToLower(nodeType) {
	case "salesorder":
		return nr.getSalesOrderNode(nodeID)
	case "delivery":
		return nr.getDeliveryNode(nodeID)
	case "billingdocument":
		return nr.getBillingDocumentNode(nodeID)
	case "payment":
		return nr.getPaymentNode(nodeID)
	case "customer":
		return nr.getCustomerNode(nodeID)
	case "product":
		return nr.getProductNode(nodeID)
	case "plant":
		return nr.getPlantNode(nodeID)
	default:
		return nil, fmt.Errorf("unknown node type: %s", nodeType)
	}
}

// getSalesOrderNode retrieves a sales order node by ID
func (nr *NodeRepository) getSalesOrderNode(soID string) (*domain.Node, error) {
	query := `SELECT sales_order, total_net_amount, transaction_currency,
		overall_delivery_status, sold_to_party
	FROM sales_order_headers WHERE sales_order = $1`

	var id, currency, deliveryStatus, soldToParty sql.NullString
	var amount sql.NullFloat64

	err := nr.db.QueryRow(query, soID).Scan(&id, &amount, &currency, &deliveryStatus, &soldToParty)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Node{
		ID:    fmt.Sprintf("SO_%s", id.String),
		Type:  "SalesOrder",
		Label: id.String,
		Properties: map[string]interface{}{
			"sales_order":             id.String,
			"total_net_amount":        amount.Float64,
			"transaction_currency":    currency.String,
			"overall_delivery_status": deliveryStatus.String,
			"sold_to_party":           soldToParty.String,
		},
	}, nil
}

// getDeliveryNode retrieves a delivery node by ID
func (nr *NodeRepository) getDeliveryNode(docID string) (*domain.Node, error) {
	query := `SELECT delivery_document, creation_date, overall_goods_movement_status,
		overall_picking_status, shipping_point
	FROM outbound_delivery_headers WHERE delivery_document = $1`

	var id, movementStatus, pickingStatus, shippingPoint sql.NullString
	var creationDate sql.NullTime

	err := nr.db.QueryRow(query, docID).Scan(&id, &creationDate, &movementStatus, &pickingStatus, &shippingPoint)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Node{
		ID:    fmt.Sprintf("DEL_%s", id.String),
		Type:  "Delivery",
		Label: id.String,
		Properties: map[string]interface{}{
			"delivery_document":             id.String,
			"creation_date":                 creationDate.Time,
			"overall_goods_movement_status": movementStatus.String,
			"overall_picking_status":        pickingStatus.String,
			"shipping_point":                shippingPoint.String,
		},
	}, nil
}

// getBillingDocumentNode retrieves a billing document node by ID
func (nr *NodeRepository) getBillingDocumentNode(docID string) (*domain.Node, error) {
	query := `SELECT billing_document, total_net_amount, transaction_currency,
		billing_document_is_cancelled, company_code
	FROM billing_document_headers WHERE billing_document = $1`

	var id, currency, companyCode sql.NullString
	var amount sql.NullFloat64
	var isCancelled bool

	err := nr.db.QueryRow(query, docID).Scan(&id, &amount, &currency, &isCancelled, &companyCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Node{
		ID:    fmt.Sprintf("BIL_%s", id.String),
		Type:  "BillingDocument",
		Label: id.String,
		Properties: map[string]interface{}{
			"billing_document":              id.String,
			"total_net_amount":              amount.Float64,
			"transaction_currency":          currency.String,
			"billing_document_is_cancelled": isCancelled,
			"company_code":                  companyCode.String,
		},
	}, nil
}

// getPaymentNode retrieves a payment node by ID
func (nr *NodeRepository) getPaymentNode(docID string) (*domain.Node, error) {
	query := `SELECT accounting_document, amount_in_transaction_currency,
		transaction_currency, customer, posting_date
	FROM payments_ar WHERE accounting_document = $1`

	var id, currency, customer sql.NullString
	var amount sql.NullFloat64
	var postingDate sql.NullTime

	err := nr.db.QueryRow(query, docID).Scan(&id, &amount, &currency, &customer, &postingDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Node{
		ID:    fmt.Sprintf("PAY_%s", id.String),
		Type:  "Payment",
		Label: id.String,
		Properties: map[string]interface{}{
			"accounting_document":            id.String,
			"amount_in_transaction_currency": amount.Float64,
			"transaction_currency":           currency.String,
			"customer":                       customer.String,
			"posting_date":                   postingDate.Time,
		},
	}, nil
}

// getCustomerNode retrieves a customer node by ID
func (nr *NodeRepository) getCustomerNode(bpID string) (*domain.Node, error) {
	query := `SELECT business_partner, business_partner_name, business_partner_full_name,
		business_partner_is_blocked
	FROM business_partners WHERE business_partner = $1`

	var id, bpName, fullName sql.NullString
	var isBlocked bool

	err := nr.db.QueryRow(query, bpID).Scan(&id, &bpName, &fullName, &isBlocked)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	label := bpName.String
	if label == "" {
		label = fullName.String
	}

	return &domain.Node{
		ID:    fmt.Sprintf("CUST_%s", id.String),
		Type:  "Customer",
		Label: label,
		Properties: map[string]interface{}{
			"business_partner":            id.String,
			"business_partner_name":       bpName.String,
			"business_partner_full_name":  fullName.String,
			"business_partner_is_blocked": isBlocked,
		},
	}, nil
}

// getProductNode retrieves a product node by ID
func (nr *NodeRepository) getProductNode(productID string) (*domain.Node, error) {
	query := `SELECT p.product, p.product_type, p.base_unit,
		COALESCE(pd.product_description, '') as description
	FROM products p
	LEFT JOIN product_descriptions pd ON p.product = pd.product AND pd.language = 'EN'
	WHERE p.product = $1
	LIMIT 1`

	var id, productType, baseUnit, description sql.NullString

	err := nr.db.QueryRow(query, productID).Scan(&id, &productType, &baseUnit, &description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	label := id.String
	if description.String != "" {
		label = fmt.Sprintf("%s (%s)", id.String, description.String)
	}

	return &domain.Node{
		ID:    fmt.Sprintf("PRD_%s", id.String),
		Type:  "Product",
		Label: label,
		Properties: map[string]interface{}{
			"product":      id.String,
			"product_type": productType.String,
			"base_unit":    baseUnit.String,
			"description":  description.String,
		},
	}, nil
}

// getPlantNode retrieves a plant node by ID
func (nr *NodeRepository) getPlantNode(plantID string) (*domain.Node, error) {
	query := `SELECT plant, plant_name, valuation_area
	FROM plants WHERE plant = $1`

	var id, plantName, valuationArea sql.NullString

	err := nr.db.QueryRow(query, plantID).Scan(&id, &plantName, &valuationArea)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	label := plantName.String
	if label == "" {
		label = id.String
	}

	return &domain.Node{
		ID:    fmt.Sprintf("PLANT_%s", id.String),
		Type:  "Plant",
		Label: label,
		Properties: map[string]interface{}{
			"plant":          id.String,
			"plant_name":     plantName.String,
			"valuation_area": valuationArea.String,
		},
	}, nil
}

// getNeighbors finds all nodes connected to the given node
func (nr *NodeRepository) getNeighbors(nodeType, nodeID string) ([]domain.Node, error) {
	nodeID = stripPrefix(nodeID)
	var neighbors []domain.Node

	switch strings.ToLower(nodeType) {
	case "salesorder":
		// SalesOrder connected to: Delivery (via fulfillment), Customer (sold_to_party), Product (items)
		// Get deliveries
		deliveryNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT odi.delivery_document FROM outbound_delivery_items odi
			JOIN sales_order_items soi ON soi.sales_order = odi.reference_sd_document
			WHERE soi.sales_order = $1`,
			nodeID, "delivery")
		neighbors = append(neighbors, deliveryNodes...)

		// Get customers
		custNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT sold_to_party FROM sales_order_headers WHERE sales_order = $1`,
			nodeID, "customer")
		neighbors = append(neighbors, custNodes...)

		// Get products
		prodNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT material FROM sales_order_items WHERE sales_order = $1`,
			nodeID, "product")
		neighbors = append(neighbors, prodNodes...)

	case "delivery":
		// Delivery connected to: SalesOrder, BillingDocument, Plant
		soNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT reference_sd_document FROM outbound_delivery_items
			WHERE delivery_document = $1`,
			nodeID, "salesorder")
		neighbors = append(neighbors, soNodes...)

		billingNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT bdi.billing_document FROM billing_document_items bdi
			JOIN outbound_delivery_items odi ON odi.delivery_document = bdi.reference_sd_document
			WHERE odi.delivery_document = $1`,
			nodeID, "billingdocument")
		neighbors = append(neighbors, billingNodes...)

		plantNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT plant FROM outbound_delivery_items WHERE delivery_document = $1`,
			nodeID, "plant")
		neighbors = append(neighbors, plantNodes...)

	case "billingdocument":
		// BillingDocument connected to: Delivery, Payment, Customer
		deliveryNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT reference_sd_document FROM billing_document_items WHERE billing_document = $1`,
			nodeID, "delivery")
		neighbors = append(neighbors, deliveryNodes...)

		paymentNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT pa.accounting_document FROM payments_ar pa
			JOIN journal_entry_items_ar jea ON pa.clearing_accounting_document = jea.clearing_accounting_document
			WHERE jea.reference_document = $1`,
			nodeID, "payment")
		neighbors = append(neighbors, paymentNodes...)

	case "payment":
		// Payment connected to: BillingDocument
		billingNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT jea.reference_document FROM journal_entry_items_ar jea
			JOIN payments_ar pa ON pa.clearing_accounting_document = jea.clearing_accounting_document
			WHERE pa.accounting_document = $1`,
			nodeID, "billingdocument")
		neighbors = append(neighbors, billingNodes...)

	case "customer":
		// Customer connected to: SalesOrder
		soNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT sales_order FROM sales_order_headers WHERE sold_to_party = $1`,
			nodeID, "salesorder")
		neighbors = append(neighbors, soNodes...)

	case "product":
		// Product connected to: SalesOrder, Plant
		soNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT sales_order FROM sales_order_items WHERE material = $1`,
			nodeID, "salesorder")
		neighbors = append(neighbors, soNodes...)

		plantNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT plant FROM product_plants WHERE product = $1`,
			nodeID, "plant")
		neighbors = append(neighbors, plantNodes...)

	case "plant":
		// Plant connected to: Product, Delivery
		prodNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT product FROM product_plants WHERE plant = $1`,
			nodeID, "product")
		neighbors = append(neighbors, prodNodes...)

		deliveryNodes, _ := nr.getConnectedNodes(
			`SELECT DISTINCT delivery_document FROM outbound_delivery_items WHERE plant = $1`,
			nodeID, "delivery")
		neighbors = append(neighbors, deliveryNodes...)
	}

	return neighbors, nil
}

// getConnectedNodes executes a query that returns IDs and converts them to nodes
func (nr *NodeRepository) getConnectedNodes(query, param, nodeType string) ([]domain.Node, error) {
	rows, err := nr.db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var id sql.NullString
		if err := rows.Scan(&id); err != nil {
			continue
		}

		if id.Valid {
			node, _ := nr.getNodeByTypeAndID(nodeType, id.String)
			if node != nil {
				nodes = append(nodes, *node)
			}
		}
	}

	return nodes, rows.Err()
}

// getConnectingEdges retrieves all edges involving this node
func (nr *NodeRepository) getConnectingEdges(nodeType, nodeID string) ([]domain.Edge, error) {
	var edges []domain.Edge

	// Get outgoing edges
	outgoing, _ := nr.getEdgesFromNode(nodeType, nodeID, true)
	edges = append(edges, outgoing...)

	// Get incoming edges
	incoming, _ := nr.getEdgesFromNode(nodeType, nodeID, false)
	edges = append(edges, incoming...)

	// Deduplicate edges
	seen := make(map[string]bool)
	var unique []domain.Edge
	for _, edge := range edges {
		if !seen[edge.ID] {
			unique = append(unique, edge)
			seen[edge.ID] = true
		}
	}

	return unique, nil
}

// getEdgesFromNode gets edges where this node is source or target
func (nr *NodeRepository) getEdgesFromNode(nodeType, nodeID string, outgoing bool) ([]domain.Edge, error) {
	return nil, nil
}
