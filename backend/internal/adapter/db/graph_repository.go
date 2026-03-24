package db

import (
	"database/sql"
	"fmt"
	"o2c-graph/internal/core/domain"
)

// GraphRepository handles all database queries for graph operations
type GraphRepository struct {
	db *sql.DB
}

// NewGraphRepository creates a new graph repository
func NewGraphRepository(db *sql.DB) *GraphRepository {
	return &GraphRepository{db: db}
}

// GetSalesOrderNodes retrieves all sales order nodes
func (gr *GraphRepository) GetSalesOrderNodes() ([]domain.Node, error) {
	query := `SELECT sales_order, total_net_amount, transaction_currency, 
		overall_delivery_status, sold_to_party
	FROM sales_order_headers`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var soID, currency, deliveryStatus, soldToParty sql.NullString
		var amount sql.NullFloat64

		if err := rows.Scan(&soID, &amount, &currency, &deliveryStatus, &soldToParty); err != nil {
			return nil, err
		}

		node := domain.Node{
			ID:    fmt.Sprintf("SO_%s", soID.String),
			Type:  "SalesOrder",
			Label: soID.String,
			Properties: map[string]interface{}{
				"sales_order":             soID.String,
				"total_net_amount":        amount.Float64,
				"transaction_currency":    currency.String,
				"overall_delivery_status": deliveryStatus.String,
				"sold_to_party":           soldToParty.String,
			},
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// GetDeliveryNodes retrieves all delivery nodes
func (gr *GraphRepository) GetDeliveryNodes() ([]domain.Node, error) {
	query := `SELECT delivery_document, creation_date, overall_goods_movement_status, 
		overall_picking_status, shipping_point
	FROM outbound_delivery_headers`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var docID, movementStatus, pickingStatus, shippingPoint sql.NullString
		var creationDate sql.NullTime

		if err := rows.Scan(&docID, &creationDate, &movementStatus, &pickingStatus, &shippingPoint); err != nil {
			return nil, err
		}

		node := domain.Node{
			ID:    fmt.Sprintf("DEL_%s", docID.String),
			Type:  "Delivery",
			Label: docID.String,
			Properties: map[string]interface{}{
				"delivery_document":             docID.String,
				"creation_date":                 creationDate.Time,
				"overall_goods_movement_status": movementStatus.String,
				"overall_picking_status":        pickingStatus.String,
				"shipping_point":                shippingPoint.String,
			},
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// GetBillingDocumentNodes retrieves all billing document nodes
func (gr *GraphRepository) GetBillingDocumentNodes() ([]domain.Node, error) {
	query := `SELECT billing_document, total_net_amount, transaction_currency,
		billing_document_is_cancelled, company_code
	FROM billing_document_headers`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var docID, currency, companyCode sql.NullString
		var amount sql.NullFloat64
		var isCancelled bool

		if err := rows.Scan(&docID, &amount, &currency, &isCancelled, &companyCode); err != nil {
			return nil, err
		}

		node := domain.Node{
			ID:    fmt.Sprintf("BIL_%s", docID.String),
			Type:  "BillingDocument",
			Label: docID.String,
			Properties: map[string]interface{}{
				"billing_document":              docID.String,
				"total_net_amount":              amount.Float64,
				"transaction_currency":          currency.String,
				"billing_document_is_cancelled": isCancelled,
				"company_code":                  companyCode.String,
			},
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// GetPaymentNodes retrieves all payment nodes
func (gr *GraphRepository) GetPaymentNodes() ([]domain.Node, error) {
	query := `SELECT accounting_document, amount_in_transaction_currency,
		transaction_currency, customer, posting_date
	FROM payments_ar`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var docID, currency, customer sql.NullString
		var amount sql.NullFloat64
		var postingDate sql.NullTime

		if err := rows.Scan(&docID, &amount, &currency, &customer, &postingDate); err != nil {
			return nil, err
		}

		node := domain.Node{
			ID:    fmt.Sprintf("PAY_%s", docID.String),
			Type:  "Payment",
			Label: docID.String,
			Properties: map[string]interface{}{
				"accounting_document":            docID.String,
				"amount_in_transaction_currency": amount.Float64,
				"transaction_currency":           currency.String,
				"customer":                       customer.String,
				"posting_date":                   postingDate.Time,
			},
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// GetCustomerNodes retrieves all customer nodes
func (gr *GraphRepository) GetCustomerNodes() ([]domain.Node, error) {
	query := `SELECT business_partner, business_partner_name, business_partner_full_name,
		business_partner_is_blocked
	FROM business_partners`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var bpID, bpName, fullName sql.NullString
		var isBlocked bool

		if err := rows.Scan(&bpID, &bpName, &fullName, &isBlocked); err != nil {
			return nil, err
		}

		label := bpName.String
		if label == "" {
			label = fullName.String
		}

		node := domain.Node{
			ID:    fmt.Sprintf("CUST_%s", bpID.String),
			Type:  "Customer",
			Label: label,
			Properties: map[string]interface{}{
				"business_partner":            bpID.String,
				"business_partner_name":       bpName.String,
				"business_partner_full_name":  fullName.String,
				"business_partner_is_blocked": isBlocked,
			},
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// GetProductNodes retrieves all product nodes
func (gr *GraphRepository) GetProductNodes() ([]domain.Node, error) {
	query := `SELECT p.product, p.product_type, p.base_unit,
		COALESCE(pd.product_description, '') as description
	FROM products p
	LEFT JOIN product_descriptions pd ON p.product = pd.product AND pd.language = 'EN'
	GROUP BY p.product, p.product_type, p.base_unit, pd.product_description`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var productID, productType, baseUnit, description sql.NullString

		if err := rows.Scan(&productID, &productType, &baseUnit, &description); err != nil {
			return nil, err
		}

		label := productID.String
		if description.String != "" {
			label = fmt.Sprintf("%s (%s)", productID.String, description.String)
		}

		node := domain.Node{
			ID:    fmt.Sprintf("PRD_%s", productID.String),
			Type:  "Product",
			Label: label,
			Properties: map[string]interface{}{
				"product":      productID.String,
				"product_type": productType.String,
				"base_unit":    baseUnit.String,
				"description":  description.String,
			},
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// GetPlantNodes retrieves all plant nodes
func (gr *GraphRepository) GetPlantNodes() ([]domain.Node, error) {
	query := `SELECT plant, plant_name, valuation_area FROM plants`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var plantID, plantName, valuationArea sql.NullString

		if err := rows.Scan(&plantID, &plantName, &valuationArea); err != nil {
			return nil, err
		}

		label := plantName.String
		if label == "" {
			label = plantID.String
		}

		node := domain.Node{
			ID:    fmt.Sprintf("PLANT_%s", plantID.String),
			Type:  "Plant",
			Label: label,
			Properties: map[string]interface{}{
				"plant":          plantID.String,
				"plant_name":     plantName.String,
				"valuation_area": valuationArea.String,
			},
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// GetSalesOrderDeliveryEdges creates edges from SO to Delivery
func (gr *GraphRepository) GetSalesOrderDeliveryEdges() ([]domain.Edge, error) {
	query := `
	SELECT DISTINCT soi.sales_order, odi.delivery_document
	FROM sales_order_items soi
	JOIN outbound_delivery_items odi ON odi.reference_sd_document = soi.sales_order
	WHERE soi.sales_order IS NOT NULL AND odi.delivery_document IS NOT NULL
	`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edges []domain.Edge
	for rows.Next() {
		var soID, deliveryID sql.NullString

		if err := rows.Scan(&soID, &deliveryID); err != nil {
			return nil, err
		}

		sourceID := fmt.Sprintf("SO_%s", soID.String)
		targetID := fmt.Sprintf("DEL_%s", deliveryID.String)
		edgeID := fmt.Sprintf("%s-%s", sourceID, targetID)

		edge := domain.Edge{
			ID:     edgeID,
			Source: sourceID,
			Target: targetID,
			Label:  "fulfilled_by",
		}
		edges = append(edges, edge)
	}

	return edges, rows.Err()
}

// GetDeliveryBillingEdges creates edges from Delivery to BillingDocument
func (gr *GraphRepository) GetDeliveryBillingEdges() ([]domain.Edge, error) {
	query := `
	SELECT DISTINCT odi.delivery_document, bdi.billing_document
	FROM outbound_delivery_items odi
	JOIN billing_document_items bdi ON bdi.reference_sd_document = odi.delivery_document
	WHERE odi.delivery_document IS NOT NULL AND bdi.billing_document IS NOT NULL
	`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edges []domain.Edge
	for rows.Next() {
		var deliveryID, billingID sql.NullString

		if err := rows.Scan(&deliveryID, &billingID); err != nil {
			return nil, err
		}

		sourceID := fmt.Sprintf("DEL_%s", deliveryID.String)
		targetID := fmt.Sprintf("BIL_%s", billingID.String)
		edgeID := fmt.Sprintf("%s-%s", sourceID, targetID)

		edge := domain.Edge{
			ID:     edgeID,
			Source: sourceID,
			Target: targetID,
			Label:  "invoiced_as",
		}
		edges = append(edges, edge)
	}

	return edges, rows.Err()
}

// GetBillingPaymentEdges creates edges from BillingDocument to Payment
func (gr *GraphRepository) GetBillingPaymentEdges() ([]domain.Edge, error) {
	query := `
	SELECT DISTINCT jea.reference_document, pa.accounting_document
	FROM journal_entry_items_ar jea
	JOIN payments_ar pa ON pa.clearing_accounting_document = jea.clearing_accounting_document
	WHERE jea.reference_document IS NOT NULL AND pa.accounting_document IS NOT NULL
	`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edges []domain.Edge
	for rows.Next() {
		var billingID, paymentID sql.NullString

		if err := rows.Scan(&billingID, &paymentID); err != nil {
			return nil, err
		}

		sourceID := fmt.Sprintf("BIL_%s", billingID.String)
		targetID := fmt.Sprintf("PAY_%s", paymentID.String)
		edgeID := fmt.Sprintf("%s-%s", sourceID, targetID)

		edge := domain.Edge{
			ID:     edgeID,
			Source: sourceID,
			Target: targetID,
			Label:  "paid_by",
		}
		edges = append(edges, edge)
	}

	return edges, rows.Err()
}

// GetCustomerSOEdges creates edges from Customer to SalesOrder
func (gr *GraphRepository) GetCustomerSOEdges() ([]domain.Edge, error) {
	query := `
	SELECT DISTINCT soh.sold_to_party, soh.sales_order
	FROM sales_order_headers soh
	WHERE soh.sold_to_party IS NOT NULL AND soh.sales_order IS NOT NULL
	`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edges []domain.Edge
	for rows.Next() {
		var customerID, soID sql.NullString

		if err := rows.Scan(&customerID, &soID); err != nil {
			return nil, err
		}

		sourceID := fmt.Sprintf("CUST_%s", customerID.String)
		targetID := fmt.Sprintf("SO_%s", soID.String)
		edgeID := fmt.Sprintf("%s-%s", sourceID, targetID)

		edge := domain.Edge{
			ID:     edgeID,
			Source: sourceID,
			Target: targetID,
			Label:  "placed_order",
		}
		edges = append(edges, edge)
	}

	return edges, rows.Err()
}

// GetSOProductEdges creates edges from SalesOrder to Product
func (gr *GraphRepository) GetSOProductEdges() ([]domain.Edge, error) {
	query := `
	SELECT DISTINCT soi.sales_order, soi.material
	FROM sales_order_items soi
	WHERE soi.sales_order IS NOT NULL AND soi.material IS NOT NULL
	`

	rows, err := gr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edges []domain.Edge
	for rows.Next() {
		var soID, productID sql.NullString

		if err := rows.Scan(&soID, &productID); err != nil {
			return nil, err
		}

		sourceID := fmt.Sprintf("SO_%s", soID.String)
		targetID := fmt.Sprintf("PRD_%s", productID.String)
		edgeID := fmt.Sprintf("%s-%s", sourceID, targetID)

		edge := domain.Edge{
			ID:     edgeID,
			Source: sourceID,
			Target: targetID,
			Label:  "contains_item",
		}
		edges = append(edges, edge)
	}

	return edges, rows.Err()
}
