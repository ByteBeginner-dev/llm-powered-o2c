-- ============================================================
-- SAP Order-to-Cash (O2C) PostgreSQL Schema
-- ============================================================

-- -----------------------------------------------
-- 1. PRODUCTS (master data)
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS products (
    product             TEXT PRIMARY KEY,
    product_type        TEXT,
    cross_plant_status  TEXT,
    creation_date       TIMESTAMPTZ,
    created_by_user     TEXT,
    last_change_date    TIMESTAMPTZ,
    is_marked_for_deletion BOOLEAN DEFAULT FALSE,
    product_old_id      TEXT,
    gross_weight        NUMERIC,
    weight_unit         TEXT,
    net_weight          NUMERIC,
    product_group       TEXT,
    base_unit           TEXT,
    division            TEXT,
    industry_sector     TEXT
);

-- -----------------------------------------------
-- 2. PRODUCT DESCRIPTIONS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS product_descriptions (
    product             TEXT REFERENCES products(product),
    language            TEXT,
    product_description TEXT,
    PRIMARY KEY (product, language)
);

-- -----------------------------------------------
-- 3. PLANTS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS plants (
    plant                           TEXT PRIMARY KEY,
    plant_name                      TEXT,
    valuation_area                  TEXT,
    factory_calendar                TEXT,
    sales_organization              TEXT,
    address_id                      TEXT,
    plant_category                  TEXT,
    distribution_channel            TEXT,
    division                        TEXT,
    language                        TEXT,
    is_marked_for_archiving         BOOLEAN DEFAULT FALSE
);

-- -----------------------------------------------
-- 4. PRODUCT PLANTS (product ↔ plant link)
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS product_plants (
    product             TEXT REFERENCES products(product),
    plant               TEXT REFERENCES plants(plant),
    country_of_origin   TEXT,
    profit_center       TEXT,
    mrp_type            TEXT,
    availability_check_type TEXT,
    PRIMARY KEY (product, plant)
);

-- -----------------------------------------------
-- 5. PRODUCT STORAGE LOCATIONS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS product_storage_locations (
    product             TEXT REFERENCES products(product),
    plant               TEXT REFERENCES plants(plant),
    storage_location    TEXT,
    physical_inventory_block_ind TEXT,
    date_of_last_posted_cnt TIMESTAMPTZ,
    PRIMARY KEY (product, plant, storage_location)
);

-- -----------------------------------------------
-- 6. BUSINESS PARTNERS (Customers)
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS business_partners (
    business_partner            TEXT PRIMARY KEY,
    customer                    TEXT,
    business_partner_category   TEXT,
    business_partner_full_name  TEXT,
    business_partner_name       TEXT,
    business_partner_grouping   TEXT,
    correspondence_language     TEXT,
    created_by_user             TEXT,
    creation_date               TIMESTAMPTZ,
    first_name                  TEXT,
    last_name                   TEXT,
    organization_bp_name1       TEXT,
    industry                    TEXT,
    last_change_date            TIMESTAMPTZ,
    business_partner_is_blocked BOOLEAN DEFAULT FALSE,
    is_marked_for_archiving     BOOLEAN DEFAULT FALSE
);

-- -----------------------------------------------
-- 7. BUSINESS PARTNER ADDRESSES
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS business_partner_addresses (
    business_partner    TEXT REFERENCES business_partners(business_partner),
    address_id          TEXT,
    validity_start_date TIMESTAMPTZ,
    validity_end_date   TIMESTAMPTZ,
    address_uuid        TEXT,
    address_time_zone   TEXT,
    city_name           TEXT,
    country             TEXT,
    postal_code         TEXT,
    region              TEXT,
    street_name         TEXT,
    transport_zone      TEXT,
    PRIMARY KEY (business_partner, address_id)
);

-- -----------------------------------------------
-- 8. CUSTOMER COMPANY ASSIGNMENTS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS customer_company_assignments (
    customer                TEXT REFERENCES business_partners(business_partner),
    company_code            TEXT,
    payment_terms           TEXT,
    reconciliation_account  TEXT,
    deletion_indicator      BOOLEAN DEFAULT FALSE,
    customer_account_group  TEXT,
    PRIMARY KEY (customer, company_code)
);

-- -----------------------------------------------
-- 9. CUSTOMER SALES AREA ASSIGNMENTS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS customer_sales_area_assignments (
    customer                    TEXT REFERENCES business_partners(business_partner),
    sales_organization          TEXT,
    distribution_channel        TEXT,
    division                    TEXT,
    currency                    TEXT,
    customer_payment_terms      TEXT,
    delivery_priority           TEXT,
    incoterms_classification    TEXT,
    incoterms_location1         TEXT,
    shipping_condition          TEXT,
    PRIMARY KEY (customer, sales_organization, distribution_channel, division)
);

-- -----------------------------------------------
-- 10. SALES ORDER HEADERS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS sales_order_headers (
    sales_order                     TEXT PRIMARY KEY,
    sales_order_type                TEXT,
    sales_organization              TEXT,
    distribution_channel            TEXT,
    sold_to_party                   TEXT REFERENCES business_partners(business_partner),
    creation_date                   TIMESTAMPTZ,
    created_by_user                 TEXT,
    last_change_datetime            TIMESTAMPTZ,
    total_net_amount                NUMERIC,
    transaction_currency            TEXT,
    overall_delivery_status         TEXT,
    overall_ord_reltd_bilg_status   TEXT,
    pricing_date                    TIMESTAMPTZ,
    requested_delivery_date         TIMESTAMPTZ,
    header_billing_block_reason     TEXT,
    delivery_block_reason           TEXT,
    incoterms_classification        TEXT,
    incoterms_location1             TEXT,
    customer_payment_terms          TEXT
);

-- -----------------------------------------------
-- 11. SALES ORDER ITEMS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS sales_order_items (
    sales_order                     TEXT REFERENCES sales_order_headers(sales_order),
    sales_order_item                TEXT,
    sales_order_item_category       TEXT,
    material                        TEXT REFERENCES products(product),
    requested_quantity              NUMERIC,
    requested_quantity_unit         TEXT,
    transaction_currency            TEXT,
    net_amount                      NUMERIC,
    material_group                  TEXT,
    production_plant                TEXT,
    storage_location                TEXT,
    sales_document_rjcn_reason      TEXT,
    item_billing_block_reason       TEXT,
    PRIMARY KEY (sales_order, sales_order_item)
);

-- -----------------------------------------------
-- 12. SALES ORDER SCHEDULE LINES
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS sales_order_schedule_lines (
    sales_order                         TEXT REFERENCES sales_order_headers(sales_order),
    sales_order_item                    TEXT,
    schedule_line                       TEXT,
    confirmed_delivery_date             TIMESTAMPTZ,
    order_quantity_unit                 TEXT,
    confd_order_qty_by_matl_avail_check NUMERIC,
    PRIMARY KEY (sales_order, sales_order_item, schedule_line)
);

-- -----------------------------------------------
-- 13. OUTBOUND DELIVERY HEADERS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS outbound_delivery_headers (
    delivery_document               TEXT PRIMARY KEY,
    actual_goods_movement_date      TIMESTAMPTZ,
    creation_date                   TIMESTAMPTZ,
    delivery_block_reason           TEXT,
    hdr_general_incompletion_status TEXT,
    header_billing_block_reason     TEXT,
    last_change_date                TIMESTAMPTZ,
    overall_goods_movement_status   TEXT,
    overall_picking_status          TEXT,
    overall_proof_of_delivery_status TEXT,
    shipping_point                  TEXT
);

-- -----------------------------------------------
-- 14. OUTBOUND DELIVERY ITEMS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS outbound_delivery_items (
    delivery_document           TEXT REFERENCES outbound_delivery_headers(delivery_document),
    delivery_document_item      TEXT,
    actual_delivery_quantity    NUMERIC,
    delivery_quantity_unit      TEXT,
    item_billing_block_reason   TEXT,
    last_change_date            TIMESTAMPTZ,
    plant                       TEXT,
    reference_sd_document       TEXT,   -- links back to sales_order
    reference_sd_document_item  TEXT,
    storage_location            TEXT,
    batch                       TEXT,
    PRIMARY KEY (delivery_document, delivery_document_item)
);

-- -----------------------------------------------
-- 15. BILLING DOCUMENT HEADERS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS billing_document_headers (
    billing_document                TEXT PRIMARY KEY,
    billing_document_type           TEXT,
    creation_date                   TIMESTAMPTZ,
    billing_document_date           TIMESTAMPTZ,
    last_change_datetime            TIMESTAMPTZ,
    billing_document_is_cancelled   BOOLEAN DEFAULT FALSE,
    cancelled_billing_document      TEXT,
    total_net_amount                NUMERIC,
    transaction_currency            TEXT,
    company_code                    TEXT,
    fiscal_year                     TEXT,
    accounting_document             TEXT,
    sold_to_party                   TEXT REFERENCES business_partners(business_partner)
);

-- -----------------------------------------------
-- 16. BILLING DOCUMENT CANCELLATIONS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS billing_document_cancellations (
    billing_document                TEXT PRIMARY KEY,
    billing_document_type           TEXT,
    creation_date                   TIMESTAMPTZ,
    billing_document_date           TIMESTAMPTZ,
    last_change_datetime            TIMESTAMPTZ,
    billing_document_is_cancelled   BOOLEAN DEFAULT TRUE,
    cancelled_billing_document      TEXT,
    total_net_amount                NUMERIC,
    transaction_currency            TEXT,
    company_code                    TEXT,
    fiscal_year                     TEXT,
    accounting_document             TEXT,
    sold_to_party                   TEXT REFERENCES business_partners(business_partner)
);

-- -----------------------------------------------
-- 17. BILLING DOCUMENT ITEMS
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS billing_document_items (
    billing_document            TEXT REFERENCES billing_document_headers(billing_document),
    billing_document_item       TEXT,
    material                    TEXT REFERENCES products(product),
    billing_quantity            NUMERIC,
    billing_quantity_unit       TEXT,
    net_amount                  NUMERIC,
    transaction_currency        TEXT,
    reference_sd_document       TEXT,   -- links to delivery_document
    reference_sd_document_item  TEXT,
    PRIMARY KEY (billing_document, billing_document_item)
);

-- -----------------------------------------------
-- 18. JOURNAL ENTRY ITEMS (Accounts Receivable)
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS journal_entry_items_ar (
    company_code                    TEXT,
    fiscal_year                     TEXT,
    accounting_document             TEXT,
    accounting_document_item        TEXT,
    gl_account                      TEXT,
    reference_document              TEXT,   -- links to billing_document
    cost_center                     TEXT,
    profit_center                   TEXT,
    transaction_currency            TEXT,
    amount_in_transaction_currency  NUMERIC,
    company_code_currency           TEXT,
    amount_in_company_code_currency NUMERIC,
    posting_date                    TIMESTAMPTZ,
    document_date                   TIMESTAMPTZ,
    accounting_document_type        TEXT,
    assignment_reference            TEXT,
    last_change_datetime            TIMESTAMPTZ,
    customer                        TEXT REFERENCES business_partners(business_partner),
    financial_account_type          TEXT,
    clearing_date                   TIMESTAMPTZ,
    clearing_accounting_document    TEXT,
    clearing_doc_fiscal_year        TEXT,
    PRIMARY KEY (company_code, fiscal_year, accounting_document, accounting_document_item)
);

-- -----------------------------------------------
-- 19. PAYMENTS (Accounts Receivable)
-- -----------------------------------------------
CREATE TABLE IF NOT EXISTS payments_ar (
    company_code                    TEXT,
    fiscal_year                     TEXT,
    accounting_document             TEXT,
    accounting_document_item        TEXT,
    clearing_date                   TIMESTAMPTZ,
    clearing_accounting_document    TEXT,
    clearing_doc_fiscal_year        TEXT,
    amount_in_transaction_currency  NUMERIC,
    transaction_currency            TEXT,
    amount_in_company_code_currency NUMERIC,
    company_code_currency           TEXT,
    customer                        TEXT REFERENCES business_partners(business_partner),
    invoice_reference               TEXT,
    sales_document                  TEXT,
    posting_date                    TIMESTAMPTZ,
    document_date                   TIMESTAMPTZ,
    gl_account                      TEXT,
    financial_account_type          TEXT,
    profit_center                   TEXT,
    cost_center                     TEXT,
    PRIMARY KEY (company_code, fiscal_year, accounting_document, accounting_document_item)
);

-- ============================================================
-- INDEXES for fast JOIN performance
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_soi_sales_order ON sales_order_items(sales_order);
CREATE INDEX IF NOT EXISTS idx_soi_material ON sales_order_items(material);
CREATE INDEX IF NOT EXISTS idx_odi_delivery ON outbound_delivery_items(delivery_document);
CREATE INDEX IF NOT EXISTS idx_odi_ref_doc ON outbound_delivery_items(reference_sd_document);
CREATE INDEX IF NOT EXISTS idx_bdi_billing ON billing_document_items(billing_document);
CREATE INDEX IF NOT EXISTS idx_bdi_ref_doc ON billing_document_items(reference_sd_document);
CREATE INDEX IF NOT EXISTS idx_bdi_material ON billing_document_items(material);
CREATE INDEX IF NOT EXISTS idx_bdh_sold_to ON billing_document_headers(sold_to_party);
CREATE INDEX IF NOT EXISTS idx_je_ref_doc ON journal_entry_items_ar(reference_document);
CREATE INDEX IF NOT EXISTS idx_je_customer ON journal_entry_items_ar(customer);
CREATE INDEX IF NOT EXISTS idx_pay_customer ON payments_ar(customer);
CREATE INDEX IF NOT EXISTS idx_soh_sold_to ON sales_order_headers(sold_to_party);