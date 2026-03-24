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
	Model     string        `json:"model"`
	Messages  []GroqMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens,omitempty"`
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
	prompt := fmt.Sprintf(`You are a PostgreSQL query generator for an SAP Order-to-Cash database.
You ONLY answer questions about this dataset. If the question is unrelated
to SAP O2C business data, respond with exactly: OFFTOPIC

Database schema:
- products(product, product_type, product_old_id, gross_weight, net_weight, product_group, base_unit, division)
- product_descriptions(product, language, product_description)
- plants(plant, plant_name, sales_organization, distribution_channel)
- business_partners(business_partner, customer, business_partner_full_name, business_partner_name)
- business_partner_addresses(business_partner, city_name, country, postal_code, region, street_name)
- sales_order_headers(sales_order, sales_order_type, sold_to_party, creation_date, total_net_amount, transaction_currency, overall_delivery_status, requested_delivery_date)
- sales_order_items(sales_order, sales_order_item, material, requested_quantity, net_amount, production_plant, storage_location)
- outbound_delivery_headers(delivery_document, creation_date, overall_goods_movement_status, overall_picking_status, shipping_point)
- outbound_delivery_items(delivery_document, delivery_document_item, actual_delivery_quantity, plant, reference_sd_document, reference_sd_document_item)
- billing_document_headers(billing_document, billing_document_type, billing_document_date, billing_document_is_cancelled, total_net_amount, transaction_currency, company_code, accounting_document, sold_to_party)
- billing_document_items(billing_document, billing_document_item, material, billing_quantity, net_amount, reference_sd_document, reference_sd_document_item)
- journal_entry_items_ar(company_code, fiscal_year, accounting_document, gl_account, reference_document, customer, amount_in_transaction_currency, posting_date, clearing_date, clearing_accounting_document)
- payments_ar(company_code, fiscal_year, accounting_document, accounting_document_item, clearing_date, amount_in_transaction_currency, transaction_currency, customer, posting_date)

Key relationships:
- outbound_delivery_items.reference_sd_document = sales_order_headers.sales_order
- billing_document_items.reference_sd_document = outbound_delivery_items.delivery_document
- journal_entry_items_ar.reference_document = billing_document_headers.billing_document
- payments_ar linked via journal_entry_items_ar.clearing_accounting_document

Generate ONLY a valid PostgreSQL SELECT query. No explanation. No markdown. No code blocks.
Just the raw SQL query.

User query: %s`, userQuery)

	return callGroqAPI(apiKey, prompt)
}

// FormatAnswer calls Groq API to format the query results as a natural language answer
func FormatAnswer(apiKey, userQuery string, rows []map[string]interface{}) (string, error) {
	rowsJSON, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`Given this user question: %s
And this data result: %s
Write a clear, concise natural language answer in 2-3 sentences.
Only use the data provided. Do not add anything not in the data.`, userQuery, string(rowsJSON))

	return callGroqAPI(apiKey, prompt)
}

// callGroqAPI makes the actual HTTP request to Groq API
func callGroqAPI(apiKey, prompt string) (string, error) {
	logger := utils.GetLogger()

	if apiKey == "" {
		logger.Error(utils.CategoryGemini, "Groq API call failed - API key not set", fmt.Errorf("GROQ_API_KEY not set"))
		return "", fmt.Errorf("GROQ_API_KEY not set")
	}

	const groqURL = "https://api.groq.com/openai/v1/chat/completions"
	const model = "llama-3.1-8b-instant"
	const maxTokens = 1024

	reqBody := GroqRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages: []GroqMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.Error(utils.CategoryGemini, "Failed to marshal Groq request", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	logger.DebugWithData(utils.CategoryGemini, "Sending request to Groq API", map[string]interface{}{
		"model":         model,
		"endpoint":      groqURL,
		"prompt_length": len(prompt),
	})

	startTime := time.Now()

	// Create HTTP request with proper headers
	req, err := http.NewRequest("POST", groqURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Error(utils.CategoryGemini, "Failed to create Groq request", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(utils.CategoryGemini, "Failed to call Groq API", err)
		return "", fmt.Errorf("failed to call Groq API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(utils.CategoryGemini, "Failed to read Groq response", err)
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	duration := time.Since(startTime).Milliseconds()

	logger.InfoWithData(utils.CategoryGemini, "Groq API response received", map[string]interface{}{
		"model":       model,
		"status_code": resp.StatusCode,
		"duration_ms": duration,
	})

	if resp.StatusCode != http.StatusOK {
		logger.ErrorWithData(utils.CategoryGemini, "Groq API error response", fmt.Errorf("status code: %d", resp.StatusCode), map[string]interface{}{
			"status_code": resp.StatusCode,
			"response":    string(body),
			"model":       model,
			"duration_ms": duration,
		})
		return "", fmt.Errorf("Groq API error: status %d", resp.StatusCode)
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		logger.Error(utils.CategoryGemini, "Failed to parse Groq response", err)
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors in response
	if groqResp.Error != nil {
		errMsg := fmt.Errorf("Groq API error: %s", groqResp.Error.Message)
		logger.Error(utils.CategoryGemini, "Groq API returned error", errMsg)
		return "", errMsg
	}

	if len(groqResp.Choices) == 0 {
		logger.Warn(utils.CategoryGemini, "Groq API returned empty choices")
		return "", fmt.Errorf("empty response from Groq API")
	}

	response := groqResp.Choices[0].Message.Content

	logger.InfoWithData(utils.CategoryGemini, "Groq API call successful", map[string]interface{}{
		"model":           model,
		"duration_ms":     duration,
		"response_length": len(response),
	})

	return response, nil
}
