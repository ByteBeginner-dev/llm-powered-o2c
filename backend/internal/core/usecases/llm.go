package usecases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"o2c-graph/pkg/utils"
)

// GroqRequest is the request structure for Groq API (OpenAI-compatible format)
type GroqRequest struct {
	Model       string        `json:"model"`
	Messages    []GroqMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature"`
}

// GroqMessage represents a message in the request
type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GroqResponse is the response structure from Groq API
type GroqResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// GenerateSQL calls Groq API to generate SQL from natural language query
func GenerateSQL(apiKey, userQuery string) (string, error) {
	systemPrompt := `You are a PostgreSQL query generator for an SAP Order-to-Cash database.

CRITICAL: THIS IS POSTGRESQL - NOT SQL SERVER!
- Use NOW() for current timestamp (NOT GETDATE)
- Use CURRENT_DATE for current date
- Use INTERVAL for date arithmetic: NOW() - INTERVAL '1 month'
- Use EXTRACT(MONTH FROM column) for month functions

IMPORTANT:
- ANSWER ALL LEGITIMATE DATABASE QUESTIONS
- Respond with EXACTLY "OFFTOPIC" (no other text) for non-database questions like jokes, weather, general knowledge

VALID QUESTION TYPES (ALL answered from dataset):
- Count/Aggregate: "How many deliveries?", "How many orders?", "Total sales?"
- Analytics: "Top customers by revenue", "Products in most billing documents?"
- Status: "Which orders haven't been delivered?", "Payment status?"
- Details: "Show customer details for order X", "List all items in billing document Y"
- Tracing: "Trace order 12345", "Show fulfillment flow for delivery Z", "Trace full flow of billing document X"

INVALID - Respond with EXACTLY "OFFTOPIC" for:
- Unrelated topics: "Tell me a joke", "What is the weather?", "Give me Python code"
- Personal info not in data: "Who is Karthik?", "What is John's phone number?"

DATABASE SCHEMA - ALL 19 TABLES:

REFERENCE/MASTER DATA:
1. products(product, product_type, cross_plant_status, creation_date, created_by_user, last_change_date, is_marked_for_deletion, product_old_id, gross_weight, weight_unit, net_weight, product_group, base_unit, division, industry_sector)
2. product_descriptions(product, language, product_description)
3. plants(plant, plant_name, valuation_area, factory_calendar, sales_organization, address_id, plant_category, distribution_channel, division, language, is_marked_for_archiving)
4. product_plants(product, plant, country_of_origin, profit_center, mrp_type, availability_check_type)
5. product_storage_locations(product, plant, storage_location, physical_inventory_block_ind, date_of_last_posted_cnt)
6. business_partners(business_partner, customer, business_partner_category, business_partner_full_name, business_partner_name, business_partner_grouping, correspondence_language, created_by_user, creation_date, first_name, last_name, organization_bp_name1, industry, last_change_date, business_partner_is_blocked, is_marked_for_archiving)
7. business_partner_addresses(business_partner, address_id, validity_start_date, validity_end_date, address_uuid, address_time_zone, city_name, country, postal_code, region, street_name, transport_zone)
8. customer_company_assignments(customer, company_code, payment_terms, reconciliation_account, deletion_indicator, customer_account_group)
9. customer_sales_area_assignments(customer, sales_organization, distribution_channel, division, currency, customer_payment_terms, delivery_priority, incoterms_classification, incoterms_location1, shipping_condition)

SALES/TRANSACTIONAL DATA:
10. sales_order_headers(sales_order, sales_order_type, sales_organization, distribution_channel, sold_to_party, creation_date, created_by_user, last_change_datetime, total_net_amount, transaction_currency, overall_delivery_status, overall_ord_reltd_bilg_status, pricing_date, requested_delivery_date, header_billing_block_reason, delivery_block_reason, incoterms_classification, incoterms_location1, customer_payment_terms)
11. sales_order_items(sales_order, sales_order_item, sales_order_item_category, material, requested_quantity, requested_quantity_unit, transaction_currency, net_amount, material_group, production_plant, storage_location, sales_document_rjcn_reason, item_billing_block_reason)
12. sales_order_schedule_lines(sales_order, sales_order_item, schedule_line, confirmed_delivery_date, order_quantity_unit, confd_order_qty_by_matl_avail_check)

DELIVERY DATA:
13. outbound_delivery_headers(delivery_document, actual_goods_movement_date, creation_date, delivery_block_reason, hdr_general_incompletion_status, header_billing_block_reason, last_change_date, overall_goods_movement_status, overall_picking_status, overall_proof_of_delivery_status, shipping_point)
14. outbound_delivery_items(delivery_document, delivery_document_item, actual_delivery_quantity, delivery_quantity_unit, item_billing_block_reason, last_change_date, plant, reference_sd_document, reference_sd_document_item, storage_location, batch)

BILLING DATA:
15. billing_document_headers(billing_document, billing_document_type, creation_date, billing_document_date, last_change_datetime, billing_document_is_cancelled, cancelled_billing_document, total_net_amount, transaction_currency, company_code, fiscal_year, accounting_document, sold_to_party)
16. billing_document_cancellations(billing_document, billing_document_type, creation_date, billing_document_date, last_change_datetime, billing_document_is_cancelled, cancelled_billing_document, total_net_amount, transaction_currency, company_code, fiscal_year, accounting_document, sold_to_party)
17. billing_document_items(billing_document, billing_document_item, material, billing_quantity, billing_quantity_unit, net_amount, transaction_currency, reference_sd_document, reference_sd_document_item)

FINANCIAL/ACCOUNTING DATA:
18. journal_entry_items_ar(company_code, fiscal_year, accounting_document, accounting_document_item, gl_account, reference_document, cost_center, profit_center, transaction_currency, amount_in_transaction_currency, company_code_currency, amount_in_company_code_currency, posting_date, document_date, accounting_document_type, assignment_reference, last_change_datetime, customer, financial_account_type, clearing_date, clearing_accounting_document, clearing_doc_fiscal_year)
19. payments_ar(company_code, fiscal_year, accounting_document, accounting_document_item, clearing_date, clearing_accounting_document, clearing_doc_fiscal_year, amount_in_transaction_currency, transaction_currency, amount_in_company_code_currency, company_code_currency, customer, invoice_reference, sales_document, posting_date, document_date, gl_account, financial_account_type, profit_center, cost_center)

KEY RELATIONSHIPS & FOREIGN KEYS:
- sales_order_items.material → products.product
- sales_order_headers.sold_to_party → business_partners.business_partner
- outbound_delivery_items.reference_sd_document → sales_order_headers.sales_order
- billing_document_headers.sold_to_party → business_partners.business_partner
- billing_document_items.reference_sd_document → outbound_delivery_items.delivery_document
- journal_entry_items_ar.reference_document → billing_document_headers.billing_document
- payments_ar.accounting_document = journal_entry_items_ar.clearing_accounting_document

CRITICAL RULES:
1. NEVER USE UNION OR UNION ALL — use LEFT JOIN instead for multi-table queries
2. NEVER use SELECT * — always name specific columns with table aliases
3. For tracing flows across multiple tables, use LEFT JOINs in ONE single SELECT
4. Return ONLY a single SQL SELECT statement
5. No markdown, no backticks, no explanation, no code blocks — raw SQL only
6. Every column must have a table alias prefix (bdh.billing_document not just billing_document)
7. For multi-table queries always use table aliases (soh, bdi, odh etc)
8. Use COUNT(DISTINCT column) when counting unique documents
9. GROUP BY must include ALL non-aggregated columns
10. POSTGRESQL ONLY: use NOW(), CURRENT_DATE, INTERVAL — never GETDATE(), DATEDIFF()

WORKING EXAMPLES:

Example 1 - Simple count (no GROUP BY needed):
SELECT COUNT(DISTINCT sales_order) AS total_orders FROM sales_order_headers

Example 2 - Orders without deliveries:
SELECT soh.sales_order, soh.sold_to_party, soh.total_net_amount, soh.creation_date
FROM sales_order_headers soh
LEFT JOIN outbound_delivery_items odi ON odi.reference_sd_document = soh.sales_order
WHERE odi.delivery_document IS NULL

Example 3 - Products in most billing documents:
SELECT bdi.material, COUNT(DISTINCT bdi.billing_document) AS document_count
FROM billing_document_items bdi
GROUP BY bdi.material
ORDER BY document_count DESC
LIMIT 10

Example 4 - TRACE FULL FLOW (ALWAYS USE LEFT JOIN — NEVER UNION):
SELECT
  bdh.billing_document,
  bdh.billing_document_date,
  bdh.total_net_amount AS billed_amount,
  bdh.billing_document_is_cancelled,
  bdh.sold_to_party AS customer,
  bdi.billing_document_item,
  bdi.material,
  bdi.billing_quantity,
  bdi.reference_sd_document AS linked_delivery,
  odh.delivery_document,
  odh.overall_goods_movement_status AS delivery_status,
  odh.overall_picking_status,
  odi.reference_sd_document AS sales_order,
  soh.total_net_amount AS order_amount,
  soh.overall_delivery_status,
  soh.creation_date AS order_date,
  je.accounting_document AS journal_entry,
  je.amount_in_transaction_currency AS journal_amount,
  je.posting_date,
  je.clearing_date,
  p.accounting_document AS payment_doc,
  p.amount_in_transaction_currency AS payment_amount,
  p.clearing_date AS payment_date
FROM billing_document_headers bdh
LEFT JOIN billing_document_items bdi
  ON bdi.billing_document = bdh.billing_document
LEFT JOIN outbound_delivery_items odi
  ON odi.delivery_document = bdi.reference_sd_document
LEFT JOIN outbound_delivery_headers odh
  ON odh.delivery_document = odi.delivery_document
LEFT JOIN sales_order_headers soh
  ON soh.sales_order = odi.reference_sd_document
LEFT JOIN journal_entry_items_ar je
  ON je.reference_document = bdh.billing_document
LEFT JOIN payments_ar p
  ON p.accounting_document = je.clearing_accounting_document
WHERE bdh.billing_document = '90504248'

Example 5 - Top customers by revenue:
SELECT soh.sold_to_party, SUM(soh.total_net_amount) AS total_sales, COUNT(soh.sales_order) AS order_count
FROM sales_order_headers soh
GROUP BY soh.sold_to_party
ORDER BY total_sales DESC
LIMIT 10

Example 6 - Billing documents with no payment:
SELECT bdh.billing_document, bdh.total_net_amount, bdh.billing_document_date, bdh.sold_to_party
FROM billing_document_headers bdh
LEFT JOIN journal_entry_items_ar je ON je.reference_document = bdh.billing_document
LEFT JOIN payments_ar p ON p.accounting_document = je.clearing_accounting_document
WHERE p.accounting_document IS NULL
AND bdh.billing_document_is_cancelled = false

Example 7 - Deliveries with no billing document:
SELECT DISTINCT odi.delivery_document, odh.creation_date, odh.overall_goods_movement_status
FROM outbound_delivery_items odi
LEFT JOIN outbound_delivery_headers odh ON odh.delivery_document = odi.delivery_document
LEFT JOIN billing_document_items bdi ON bdi.reference_sd_document = odi.delivery_document
WHERE bdi.billing_document IS NULL

Example 8 - Top customers by sales with address:
SELECT bp.business_partner, bp.business_partner_name, COUNT(soh.sales_order) AS orders,
  SUM(soh.total_net_amount) AS total_spent, bpa.city_name, bpa.country
FROM business_partners bp
LEFT JOIN sales_order_headers soh ON soh.sold_to_party = bp.business_partner
LEFT JOIN business_partner_addresses bpa ON bpa.business_partner = bp.business_partner
GROUP BY bp.business_partner, bp.business_partner_name, bpa.city_name, bpa.country
ORDER BY total_spent DESC

Example 9 - Billing vs payment reconciliation:
SELECT bdh.billing_document, bdh.total_net_amount AS billed_amount,
  COALESCE(SUM(p.amount_in_transaction_currency), 0) AS paid_amount,
  (bdh.total_net_amount - COALESCE(SUM(p.amount_in_transaction_currency), 0)) AS outstanding
FROM billing_document_headers bdh
LEFT JOIN journal_entry_items_ar je ON je.reference_document = bdh.billing_document
LEFT JOIN payments_ar p ON p.accounting_document = je.clearing_accounting_document
GROUP BY bdh.billing_document, bdh.total_net_amount
ORDER BY outstanding DESC`

	return callGroqAPI(apiKey, systemPrompt, userQuery)
}

// FormatAnswer calls Groq API to format the query results as a natural language answer
func FormatAnswer(apiKey, userQuery string, rows []map[string]interface{}) (string, error) {
	rowsJSON, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return "", err
	}

	systemPrompt := `You are a helpful assistant that formats database query results into clear, natural language answers.`
	userMessage := fmt.Sprintf(`Given this user question: %s
And this data result: %s
Write a clear, concise natural language answer in 2-3 sentences.
Only use the data provided. Do not add anything not in the data.`, userQuery, string(rowsJSON))

	return callGroqAPI(apiKey, systemPrompt, userMessage)
}

// callGroqAPI makes the actual HTTP request to Groq API - PROPERLY SEPARATES SYSTEM AND USER MESSAGES
func callGroqAPI(apiKey, systemPrompt, userMessage string) (string, error) {
	logger := utils.GetLogger()

	if apiKey == "" {
		logger.Error(utils.CategoryGroq, "Groq API call failed - API key not set", fmt.Errorf("GROQ_API_KEY not set"))
		return "", fmt.Errorf("GROQ_API_KEY not set")
	}

	const groqURL = "https://api.groq.com/openai/v1/chat/completions"
	const model = "llama-3.3-70b-versatile"
	const maxTokens = 1024

	reqBody := GroqRequest{
		Model:       model,
		MaxTokens:   maxTokens,
		Temperature: 0.1,
		Messages: []GroqMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userMessage,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.Error(utils.CategoryGroq, "Failed to marshal Groq request", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	logger.DebugWithData(utils.CategoryGroq, "Sending request to Groq API", map[string]interface{}{
		"model":         model,
		"endpoint":      groqURL,
		"system_length": len(systemPrompt),
		"user_message":  userMessage,
	})

	startTime := time.Now()

	// Create HTTP request with proper headers
	req, err := http.NewRequest("POST", groqURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Error(utils.CategoryGroq, "Failed to create Groq request", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(utils.CategoryGroq, "Failed to call Groq API", err)
		return "", fmt.Errorf("failed to call Groq API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(utils.CategoryGroq, "Failed to read Groq response", err)
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	duration := time.Since(startTime).Milliseconds()

	logger.InfoWithData(utils.CategoryGroq, "Groq API response received", map[string]interface{}{
		"model":       model,
		"status_code": resp.StatusCode,
		"duration_ms": duration,
	})

	if resp.StatusCode != http.StatusOK {
		logger.ErrorWithData(utils.CategoryGroq, "Groq API error response", fmt.Errorf("status code: %d", resp.StatusCode), map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(body),
			"model":       model,
			"duration_ms": duration,
		})
		return "", fmt.Errorf("Groq API error: status %d", resp.StatusCode)
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		logger.Error(utils.CategoryGroq, "Failed to parse Groq response", err)
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors in response
	if groqResp.Error != nil {
		errMsg := fmt.Errorf("Groq API error: %s", groqResp.Error.Message)
		logger.Error(utils.CategoryGroq, "Groq API returned error", errMsg)
		return "", errMsg
	}

	if len(groqResp.Choices) == 0 {
		logger.Warn(utils.CategoryGroq, "Groq API returned empty choices")
		return "", fmt.Errorf("empty response from Groq API")
	}

	response := groqResp.Choices[0].Message.Content

	logger.InfoWithData(utils.CategoryGroq, "Groq API call successful", map[string]interface{}{
		"model":           model,
		"duration_ms":     duration,
		"response_length": len(response),
	})

	return response, nil
}
