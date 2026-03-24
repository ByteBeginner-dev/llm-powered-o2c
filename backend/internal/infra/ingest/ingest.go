package ingest

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Run orchestrates data ingestion from JSONL files into PostgreSQL.
// It processes all 19 entity types in FK dependency order.
// Uses ON CONFLICT DO NOTHING to ensure idempotency.
func Run(db *sql.DB, dataDir string) error {
	// Ingestion order respects FK dependencies
	ingestFuncs := []struct {
		name string
		fn   func(*sql.DB, string) error
	}{
		{"products", ingestProducts},
		{"product_descriptions", ingestProductDescriptions},
		{"plants", ingestPlants},
		{"product_plants", ingestProductPlants},
		{"product_storage_locations", ingestProductStorageLocations},
		{"business_partners", ingestBusinessPartners},
		{"business_partner_addresses", ingestBusinessPartnerAddresses},
		{"customer_company_assignments", ingestCustomerCompanyAssignments},
		{"customer_sales_area_assignments", ingestCustomerSalesAreaAssignments},
		{"sales_order_headers", ingestSalesOrderHeaders},
		{"sales_order_items", ingestSalesOrderItems},
		{"sales_order_schedule_lines", ingestSalesOrderScheduleLines},
		{"outbound_delivery_headers", ingestOutboundDeliveryHeaders},
		{"outbound_delivery_items", ingestOutboundDeliveryItems},
		{"billing_document_headers", ingestBillingDocumentHeaders},
		{"billing_document_cancellations", ingestBillingDocumentCancellations},
		{"billing_document_items", ingestBillingDocumentItems},
		{"journal_entry_items_ar", ingestJournalEntryItemsAR},
		{"payments_ar", ingestPaymentsAR},
	}

	for _, item := range ingestFuncs {
		if err := item.fn(db, dataDir); err != nil {
			return fmt.Errorf("ingest %s failed: %w", item.name, err)
		}
	}

	return nil
}

// readFolder reads all *.jsonl files in a folder and returns parsed JSON objects
func readFolder(folderPath string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read folder %s: %w", folderPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process .jsonl files
		if !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		filePath := filepath.Join(folderPath, entry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", filePath, err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(line), &obj); err != nil {
				file.Close()
				return nil, fmt.Errorf("failed to parse JSON line in %s: %w", filePath, err)
			}

			result = append(result, obj)
		}

		if err := scanner.Err(); err != nil {
			file.Close()
			return nil, fmt.Errorf("scanner error in %s: %w", filePath, err)
		}

		file.Close()
	}

	return result, nil
}

// toString safely converts an interface{} to sql.NullString
func toString(v interface{}) sql.NullString {
	if v == nil {
		return sql.NullString{Valid: false}
	}
	if str, ok := v.(string); ok {
		if str == "" {
			return sql.NullString{Valid: false}
		}
		return sql.NullString{String: str, Valid: true}
	}
	return sql.NullString{Valid: false}
}

// toFloat safely converts an interface{} to sql.NullFloat64
func toFloat(v interface{}) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{Valid: false}
	}

	switch val := v.(type) {
	case float64:
		return sql.NullFloat64{Float64: val, Valid: true}
	case string:
		if val == "" {
			return sql.NullFloat64{Valid: false}
		}
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return sql.NullFloat64{Valid: false}
		}
		return sql.NullFloat64{Float64: f, Valid: true}
	}
	return sql.NullFloat64{Valid: false}
}

// toBool safely converts an interface{} to bool
func toBool(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.ToLower(val) == "true"
	case float64:
		return val != 0
	}
	return false
}

// toTime safely converts an interface{} to sql.NullTime (RFC3339 format)
func toTime(v interface{}) sql.NullTime {
	if v == nil {
		return sql.NullTime{Valid: false}
	}

	str, ok := v.(string)
	if !ok {
		return sql.NullTime{Valid: false}
	}

	if str == "" {
		return sql.NullTime{Valid: false}
	}

	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: t, Valid: true}
}

// ingestProducts ingests product master data
func ingestProducts(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "products")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO products (
		product, product_type, cross_plant_status, creation_date, created_by_user,
		last_change_date, is_marked_for_deletion, product_old_id, gross_weight,
		weight_unit, net_weight, product_group, base_unit, division, industry_sector
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["product"]),
			toString(rec["productType"]),
			toString(rec["crossPlantStatus"]),
			toTime(rec["creationDate"]),
			toString(rec["createdByUser"]),
			toTime(rec["lastChangeDate"]),
			toBool(rec["isMarkedForDeletion"]),
			toString(rec["productOldId"]),
			toFloat(rec["grossWeight"]),
			toString(rec["weightUnit"]),
			toFloat(rec["netWeight"]),
			toString(rec["productGroup"]),
			toString(rec["baseUnit"]),
			toString(rec["division"]),
			toString(rec["industrySector"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting product: %w", err)
		}
	}

	return nil
}

// ingestProductDescriptions ingests product descriptions
func ingestProductDescriptions(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "product_descriptions")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO product_descriptions (product, language, product_description)
	VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["product"]),
			toString(rec["language"]),
			toString(rec["productDescription"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting product_description: %w", err)
		}
	}

	return nil
}

// ingestPlants ingests plant master data
func ingestPlants(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "plants")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO plants (
		plant, plant_name, valuation_area, factory_calendar, sales_organization,
		address_id, plant_category, distribution_channel, division, language,
		is_marked_for_archiving
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["plant"]),
			toString(rec["plantName"]),
			toString(rec["valuationArea"]),
			toString(rec["factoryCalendar"]),
			toString(rec["salesOrganization"]),
			toString(rec["addressId"]),
			toString(rec["plantCategory"]),
			toString(rec["distributionChannel"]),
			toString(rec["division"]),
			toString(rec["language"]),
			toBool(rec["isMarkedForArchiving"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting plant: %w", err)
		}
	}

	return nil
}

// ingestProductPlants ingests product-plant relationships
func ingestProductPlants(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "product_plants")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO product_plants (
		product, plant, country_of_origin, profit_center, mrp_type,
		availability_check_type
	) VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["product"]),
			toString(rec["plant"]),
			toString(rec["countryOfOrigin"]),
			toString(rec["profitCenter"]),
			toString(rec["mrpType"]),
			toString(rec["availabilityCheckType"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting product_plant: %w", err)
		}
	}

	return nil
}

// ingestProductStorageLocations ingests product storage locations
func ingestProductStorageLocations(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "product_storage_locations")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO product_storage_locations (
		product, plant, storage_location, physical_inventory_block_ind,
		date_of_last_posted_cnt
	) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["product"]),
			toString(rec["plant"]),
			toString(rec["storageLocation"]),
			toString(rec["physicalInventoryBlockInd"]),
			toTime(rec["dateOfLastPostedCntUnRstrcdStk"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting product_storage_location: %w", err)
		}
	}

	return nil
}

// ingestBusinessPartners ingests business partner master data
func ingestBusinessPartners(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "business_partners")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO business_partners (
		business_partner, customer, business_partner_category,
		business_partner_full_name, business_partner_name, business_partner_grouping,
		correspondence_language, created_by_user, creation_date, first_name,
		last_name, organization_bp_name1, industry, last_change_date,
		business_partner_is_blocked, is_marked_for_archiving
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["businessPartner"]),
			toString(rec["customer"]),
			toString(rec["businessPartnerCategory"]),
			toString(rec["businessPartnerFullName"]),
			toString(rec["businessPartnerName"]),
			toString(rec["businessPartnerGrouping"]),
			toString(rec["correspondenceLanguage"]),
			toString(rec["createdByUser"]),
			toTime(rec["creationDate"]),
			toString(rec["firstName"]),
			toString(rec["lastName"]),
			toString(rec["organizationBpName1"]),
			toString(rec["industry"]),
			toTime(rec["lastChangeDate"]),
			toBool(rec["businessPartnerIsBlocked"]),
			toBool(rec["isMarkedForArchiving"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting business_partner: %w", err)
		}
	}

	return nil
}

// ingestBusinessPartnerAddresses ingests business partner addresses
func ingestBusinessPartnerAddresses(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "business_partner_addresses")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO business_partner_addresses (
		business_partner, address_id, validity_start_date, validity_end_date,
		address_uuid, address_time_zone, city_name, country, postal_code,
		region, street_name, transport_zone
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["businessPartner"]),
			toString(rec["addressId"]),
			toTime(rec["validityStartDate"]),
			toTime(rec["validityEndDate"]),
			toString(rec["addressUuid"]),
			toString(rec["addressTimeZone"]),
			toString(rec["cityName"]),
			toString(rec["country"]),
			toString(rec["postalCode"]),
			toString(rec["region"]),
			toString(rec["streetName"]),
			toString(rec["transportZone"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting business_partner_address: %w", err)
		}
	}

	return nil
}

// ingestCustomerCompanyAssignments ingests customer-company assignments
func ingestCustomerCompanyAssignments(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "customer_company_assignments")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO customer_company_assignments (
		customer, company_code, payment_terms, reconciliation_account,
		deletion_indicator, customer_account_group
	) VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["customer"]),
			toString(rec["companyCode"]),
			toString(rec["paymentTerms"]),
			toString(rec["reconciliationAccount"]),
			toBool(rec["deletionIndicator"]),
			toString(rec["customerAccountGroup"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting customer_company_assignment: %w", err)
		}
	}

	return nil
}

// ingestCustomerSalesAreaAssignments ingests customer sales area assignments
func ingestCustomerSalesAreaAssignments(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "customer_sales_area_assignments")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO customer_sales_area_assignments (
		customer, sales_organization, distribution_channel, division,
		currency, customer_payment_terms, delivery_priority,
		incoterms_classification, incoterms_location1, shipping_condition
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["customer"]),
			toString(rec["salesOrganization"]),
			toString(rec["distributionChannel"]),
			toString(rec["division"]),
			toString(rec["currency"]),
			toString(rec["customerPaymentTerms"]),
			toString(rec["deliveryPriority"]),
			toString(rec["incotermsClassification"]),
			toString(rec["incotermsLocation1"]),
			toString(rec["shippingCondition"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting customer_sales_area_assignment: %w", err)
		}
	}

	return nil
}

// ingestSalesOrderHeaders ingests sales order headers
func ingestSalesOrderHeaders(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "sales_order_headers")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO sales_order_headers (
		sales_order, sales_order_type, sales_organization, distribution_channel,
		sold_to_party, creation_date, created_by_user, last_change_datetime,
		total_net_amount, transaction_currency, overall_delivery_status,
		overall_ord_reltd_bilg_status, pricing_date, requested_delivery_date,
		header_billing_block_reason, delivery_block_reason,
		incoterms_classification, incoterms_location1, customer_payment_terms
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["salesOrder"]),
			toString(rec["salesOrderType"]),
			toString(rec["salesOrganization"]),
			toString(rec["distributionChannel"]),
			toString(rec["soldToParty"]),
			toTime(rec["creationDate"]),
			toString(rec["createdByUser"]),
			toTime(rec["lastChangeDateTime"]),
			toFloat(rec["totalNetAmount"]),
			toString(rec["transactionCurrency"]),
			toString(rec["overallDeliveryStatus"]),
			toString(rec["overallOrdReltdBillgStatus"]),
			toTime(rec["pricingDate"]),
			toTime(rec["requestedDeliveryDate"]),
			toString(rec["headerBillingBlockReason"]),
			toString(rec["deliveryBlockReason"]),
			toString(rec["incotermsClassification"]),
			toString(rec["incotermsLocation1"]),
			toString(rec["customerPaymentTerms"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting sales_order_header: %w", err)
		}
	}

	return nil
}

// ingestSalesOrderItems ingests sales order items
func ingestSalesOrderItems(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "sales_order_items")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO sales_order_items (
		sales_order, sales_order_item, sales_order_item_category, material,
		requested_quantity, requested_quantity_unit, transaction_currency,
		net_amount, material_group, production_plant, storage_location,
		sales_document_rjcn_reason, item_billing_block_reason
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["salesOrder"]),
			toString(rec["salesOrderItem"]),
			toString(rec["salesOrderItemCategory"]),
			toString(rec["material"]),
			toFloat(rec["requestedQuantity"]),
			toString(rec["requestedQuantityUnit"]),
			toString(rec["transactionCurrency"]),
			toFloat(rec["netAmount"]),
			toString(rec["materialGroup"]),
			toString(rec["productionPlant"]),
			toString(rec["storageLocation"]),
			toString(rec["salesDocumentRjcnReason"]),
			toString(rec["itemBillingBlockReason"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting sales_order_item: %w", err)
		}
	}

	return nil
}

// ingestSalesOrderScheduleLines ingests sales order schedule lines
func ingestSalesOrderScheduleLines(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "sales_order_schedule_lines")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO sales_order_schedule_lines (
		sales_order, sales_order_item, schedule_line, confirmed_delivery_date,
		order_quantity_unit, confd_order_qty_by_matl_avail_check
	) VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["salesOrder"]),
			toString(rec["salesOrderItem"]),
			toString(rec["scheduleLine"]),
			toTime(rec["confirmedDeliveryDate"]),
			toString(rec["orderQuantityUnit"]),
			toFloat(rec["confdOrderQtyByMatlAvailCheck"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting sales_order_schedule_line: %w", err)
		}
	}

	return nil
}

// ingestOutboundDeliveryHeaders ingests outbound delivery headers
func ingestOutboundDeliveryHeaders(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "outbound_delivery_headers")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO outbound_delivery_headers (
		delivery_document, actual_goods_movement_date, creation_date,
		delivery_block_reason, hdr_general_incompletion_status,
		header_billing_block_reason, last_change_date,
		overall_goods_movement_status, overall_picking_status,
		overall_proof_of_delivery_status, shipping_point
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["deliveryDocument"]),
			toTime(rec["actualGoodsMovementDate"]),
			toTime(rec["creationDate"]),
			toString(rec["deliveryBlockReason"]),
			toString(rec["hdrGeneralIncompletionStatus"]),
			toString(rec["headerBillingBlockReason"]),
			toTime(rec["lastChangeDate"]),
			toString(rec["overallGoodsMovementStatus"]),
			toString(rec["overallPickingStatus"]),
			toString(rec["overallProofOfDeliveryStatus"]),
			toString(rec["shippingPoint"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting outbound_delivery_header: %w", err)
		}
	}

	return nil
}

// ingestOutboundDeliveryItems ingests outbound delivery items
func ingestOutboundDeliveryItems(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "outbound_delivery_items")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO outbound_delivery_items (
		delivery_document, delivery_document_item, actual_delivery_quantity,
		delivery_quantity_unit, item_billing_block_reason, last_change_date,
		plant, reference_sd_document, reference_sd_document_item,
		storage_location, batch
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["deliveryDocument"]),
			toString(rec["deliveryDocumentItem"]),
			toFloat(rec["actualDeliveryQuantity"]),
			toString(rec["deliveryQuantityUnit"]),
			toString(rec["itemBillingBlockReason"]),
			toTime(rec["lastChangeDate"]),
			toString(rec["plant"]),
			toString(rec["referenceSdDocument"]),
			toString(rec["referenceSdDocumentItem"]),
			toString(rec["storageLocation"]),
			toString(rec["batch"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting outbound_delivery_item: %w", err)
		}
	}

	return nil
}

// ingestBillingDocumentHeaders ingests billing document headers
func ingestBillingDocumentHeaders(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "billing_document_headers")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO billing_document_headers (
		billing_document, billing_document_type, creation_date,
		billing_document_date, last_change_datetime, billing_document_is_cancelled,
		cancelled_billing_document, total_net_amount, transaction_currency,
		company_code, fiscal_year, accounting_document, sold_to_party
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["billingDocument"]),
			toString(rec["billingDocumentType"]),
			toTime(rec["creationDate"]),
			toTime(rec["billingDocumentDate"]),
			toTime(rec["lastChangeDateTime"]),
			toBool(rec["billingDocumentIsCancelled"]),
			toString(rec["cancelledBillingDocument"]),
			toFloat(rec["totalNetAmount"]),
			toString(rec["transactionCurrency"]),
			toString(rec["companyCode"]),
			toString(rec["fiscalYear"]),
			toString(rec["accountingDocument"]),
			toString(rec["soldToParty"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting billing_document_header: %w", err)
		}
	}

	return nil
}

// ingestBillingDocumentCancellations ingests billing document cancellations
func ingestBillingDocumentCancellations(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "billing_document_cancellations")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO billing_document_cancellations (
		billing_document, billing_document_type, creation_date,
		billing_document_date, last_change_datetime, billing_document_is_cancelled,
		cancelled_billing_document, total_net_amount, transaction_currency,
		company_code, fiscal_year, accounting_document, sold_to_party
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["billingDocument"]),
			toString(rec["billingDocumentType"]),
			toTime(rec["creationDate"]),
			toTime(rec["billingDocumentDate"]),
			toTime(rec["lastChangeDateTime"]),
			toBool(rec["billingDocumentIsCancelled"]),
			toString(rec["cancelledBillingDocument"]),
			toFloat(rec["totalNetAmount"]),
			toString(rec["transactionCurrency"]),
			toString(rec["companyCode"]),
			toString(rec["fiscalYear"]),
			toString(rec["accountingDocument"]),
			toString(rec["soldToParty"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting billing_document_cancellation: %w", err)
		}
	}

	return nil
}

// ingestBillingDocumentItems ingests billing document items
func ingestBillingDocumentItems(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "billing_document_items")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO billing_document_items (
		billing_document, billing_document_item, material,
		billing_quantity, billing_quantity_unit, net_amount, transaction_currency,
		reference_sd_document, reference_sd_document_item
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["billingDocument"]),
			toString(rec["billingDocumentItem"]),
			toString(rec["material"]),
			toFloat(rec["billingQuantity"]),
			toString(rec["billingQuantityUnit"]),
			toFloat(rec["netAmount"]),
			toString(rec["transactionCurrency"]),
			toString(rec["referenceSdDocument"]),
			toString(rec["referenceSdDocumentItem"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting billing_document_item: %w", err)
		}
	}

	return nil
}

// ingestJournalEntryItemsAR ingests journal entry items for accounts receivable
func ingestJournalEntryItemsAR(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "journal_entry_items_accounts_receivable")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO journal_entry_items_ar (
		company_code, fiscal_year, accounting_document, accounting_document_item,
		gl_account, reference_document, cost_center, profit_center,
		transaction_currency, amount_in_transaction_currency, company_code_currency,
		amount_in_company_code_currency, posting_date, document_date,
		accounting_document_type, assignment_reference, last_change_datetime,
		customer, financial_account_type, clearing_date, clearing_accounting_document,
		clearing_doc_fiscal_year
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["companyCode"]),
			toString(rec["fiscalYear"]),
			toString(rec["accountingDocument"]),
			toString(rec["accountingDocumentItem"]),
			toString(rec["glAccount"]),
			toString(rec["referenceDocument"]),
			toString(rec["costCenter"]),
			toString(rec["profitCenter"]),
			toString(rec["transactionCurrency"]),
			toFloat(rec["amountInTransactionCurrency"]),
			toString(rec["companyCodeCurrency"]),
			toFloat(rec["amountInCompanyCodeCurrency"]),
			toTime(rec["postingDate"]),
			toTime(rec["documentDate"]),
			toString(rec["accountingDocumentType"]),
			toString(rec["assignmentReference"]),
			toTime(rec["lastChangeDateTime"]),
			toString(rec["customer"]),
			toString(rec["financialAccountType"]),
			toTime(rec["clearingDate"]),
			toString(rec["clearingAccountingDocument"]),
			toString(rec["clearingDocFiscalYear"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting journal_entry_item_ar: %w", err)
		}
	}

	return nil
}

// ingestPaymentsAR ingests payment records for accounts receivable
func ingestPaymentsAR(db *sql.DB, dataDir string) error {
	folderPath := filepath.Join(dataDir, "payments_accounts_receivable")
	records, err := readFolder(folderPath)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO payments_ar (
		company_code, fiscal_year, accounting_document, accounting_document_item,
		clearing_date, clearing_accounting_document, clearing_doc_fiscal_year,
		amount_in_transaction_currency, transaction_currency, amount_in_company_code_currency,
		company_code_currency, customer, invoice_reference, sales_document,
		posting_date, document_date, gl_account, financial_account_type,
		profit_center, cost_center
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	ON CONFLICT DO NOTHING`

	for _, rec := range records {
		_, err := db.Exec(stmt,
			toString(rec["companyCode"]),
			toString(rec["fiscalYear"]),
			toString(rec["accountingDocument"]),
			toString(rec["accountingDocumentItem"]),
			toTime(rec["clearingDate"]),
			toString(rec["clearingAccountingDocument"]),
			toString(rec["clearingDocFiscalYear"]),
			toFloat(rec["amountInTransactionCurrency"]),
			toString(rec["transactionCurrency"]),
			toFloat(rec["amountInCompanyCodeCurrency"]),
			toString(rec["companyCodeCurrency"]),
			toString(rec["customer"]),
			toString(rec["invoiceReference"]),
			toString(rec["salesDocument"]),
			toTime(rec["postingDate"]),
			toTime(rec["documentDate"]),
			toString(rec["glAccount"]),
			toString(rec["financialAccountType"]),
			toString(rec["profitCenter"]),
			toString(rec["costCenter"]),
		)
		if err != nil {
			return fmt.Errorf("error inserting payment_ar: %w", err)
		}
	}

	return nil
}
