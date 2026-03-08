# API Documentation

This document outlines the API routes available in the system, along with their expected request payloads and response shapes.

## Table of Contents
- [Authentication](#authentication)
- [Health Check Routes](#health-check-routes)
- [User Routes](#user-routes)
- [Webhook Routes](#webhook-routes)
- [Catalog Routes](#catalog-routes)
- [Inventory Routes](#inventory-routes)

---

## Authentication

Routes marked as **Protected: Yes** require authentication. The specific authentication mechanism is handled by `AuthMiddleware`, which typically relies on a valid session or token provided in the request headers or context to identify the `OrganizationID`.

---

## Health Check Routes

### Liveness Probe
- **Method:** `GET`
- **Path:** `/health/healthz`
- **Protected:** No
- **Response:** Basic health status indicating the server is running.

### Readiness Probe
- **Method:** `GET`
- **Path:** `/health/readyz`
- **Protected:** No
- **Response:** Detailed health status indicating if the server is ready to accept traffic (e.g., database connection is alive).

---

## User Routes

### Create User
- **Method:** `POST`
- **Path:** `/users/`
- **Protected:** No
- **Request Body:**
  ```json
  {
    "id": "string",
    "name": "string",
    "email": "string",
    "email_verified": true,
    "image": "string"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "id": "string",
    "name": "string",
    "email": "string",
    "email_verified": true,
    "image": "string (optional)",
    "created_at": "string (RFC3339 datetime)",
    "updated_at": "string (RFC3339 datetime)"
  }
  ```

---

## Webhook Routes

### Clerk Webhook
- **Method:** `POST`
- **Path:** `/webhooks/clerk`
- **Protected:** No
- **Description:** Endpoint to handle Clerk webhook events.

---

## Catalog Routes

*All Catalog routes are protected.*

### Create Product
- **Method:** `POST`
- **Path:** `/catalog/products`
- **Protected:** Yes
- **Request Body:**
  ```json
  {
    "name": "string",
    "description": "string (optional)",
    "base_unit": "string"
  }
  ```
- **Response (200 OK):** `Product` object.
  ```json
  {
    "id": "uuid",
    "organization_id": "string",
    "name": "string",
    "description": "string (optional)",
    "base_unit": "string",
    "status": "string (active | archived)",
    "created_at": "string (RFC3339 datetime)",
    "updated_at": "string (RFC3339 datetime)"
  }
  ```

### List Products
- **Method:** `GET`
- **Path:** `/catalog/products`
- **Protected:** Yes
- **Query Parameters:**
  - `limit` (integer, optional)
  - `offset` (integer, optional)
- **Response (200 OK):** Array of `Product` objects.

### Get Product
- **Method:** `GET`
- **Path:** `/catalog/products/{productID}`
- **Protected:** Yes
- **Path Parameters:**
  - `productID` (uuid)
- **Response (200 OK):** `Product` object.

### Update Product
- **Method:** `PATCH`
- **Path:** `/catalog/products/{productID}`
- **Protected:** Yes
- **Path Parameters:**
  - `productID` (uuid)
- **Request Body:**
  ```json
  {
    "name": "string (optional)",
    "description": "string (optional)",
    "status": "string (active | archived) (optional)"
  }
  ```
- **Response (200 OK):** Updated `Product` object.

### Archive Product
- **Method:** `POST`
- **Path:** `/catalog/products/{productID}/archive`
- **Protected:** Yes
- **Path Parameters:**
  - `productID` (uuid)
- **Response (200 OK):** Archived `Product` object.

### Create Variant
- **Method:** `POST`
- **Path:** `/catalog/products/{productID}/variants`
- **Protected:** Yes
- **Path Parameters:**
  - `productID` (uuid)
- **Request Body:**
  ```json
  {
    "sku": "string",
    "barcode": "string (optional)",
    "price": 0.0,
    "cost": 0.0, // optional
    "is_active": true
  }
  ```
- **Response (200 OK):** `ProductVariant` object.
  ```json
  {
    "id": "uuid",
    "organization_id": "string",
    "product_id": "uuid",
    "sku": "string",
    "barcode": "string (optional)",
    "price": 0.0,
    "cost": 0.0, // optional
    "is_active": true,
    "created_at": "string (RFC3339 datetime)",
    "updated_at": "string (RFC3339 datetime)"
  }
  ```

### List Variants by Product
- **Method:** `GET`
- **Path:** `/catalog/products/{productID}/variants`
- **Protected:** Yes
- **Path Parameters:**
  - `productID` (uuid)
- **Response (200 OK):** Array of `ProductVariant` objects.

### Get Variant
- **Method:** `GET`
- **Path:** `/catalog/variants/{variantID}`
- **Protected:** Yes
- **Path Parameters:**
  - `variantID` (uuid)
- **Response (200 OK):** `ProductVariant` object.

### Update Variant
- **Method:** `PATCH`
- **Path:** `/catalog/variants/{variantID}`
- **Protected:** Yes
- **Path Parameters:**
  - `variantID` (uuid)
- **Request Body:**
  ```json
  {
    "sku": "string (optional)",
    "barcode": "string (optional)",
    "price": 0.0, // optional
    "cost": 0.0, // optional
    "is_active": true // optional
  }
  ```
- **Response (200 OK):** Updated `ProductVariant` object.

---

## Inventory Routes

*All Inventory routes are protected.*

### Upsert Conversion
- **Method:** `POST`
- **Path:** `/inventory/products/{productID}/conversions`
- **Protected:** Yes
- **Path Parameters:**
  - `productID` (uuid)
- **Request Body:**
  ```json
  {
    "unit_from": "string",
    "unit_to": "string",
    "factor": 0.0,
    "precision": 0
  }
  ```
- **Response (200 OK):** `UnitConversion` object.
  ```json
  {
    "id": "uuid",
    "organization_id": "string",
    "product_id": "uuid",
    "unit_from": "string",
    "unit_to": "string",
    "factor": 0.0,
    "precision": 0,
    "created_at": "string (RFC3339 datetime)",
    "updated_at": "string (RFC3339 datetime)"
  }
  ```

### List Conversions by Product
- **Method:** `GET`
- **Path:** `/inventory/products/{productID}/conversions`
- **Protected:** Yes
- **Path Parameters:**
  - `productID` (uuid)
- **Response (200 OK):** Array of `UnitConversion` objects.

### Create Receipt
- **Method:** `POST`
- **Path:** `/inventory/variants/{variantID}/receipt`
- **Protected:** Yes
- **Path Parameters:**
  - `variantID` (uuid)
- **Request Body:**
  ```json
  {
    "quantity": 0,
    "unit": "string",
    "source_id": "uuid (optional)",
    "note": "string (optional)"
  }
  ```
- **Response (200 OK):** `InventoryEvent` object.
  ```json
  {
    "id": "uuid",
    "organization_id": "string",
    "product_variant_id": "uuid",
    "event_type": "string (e.g., purchase_receipt)",
    "quantity_change": 0,
    "source_type": "string (optional)",
    "source_id": "uuid (optional)",
    "note": "string (optional)",
    "created_at": "string (RFC3339 datetime)"
  }
  ```

### Create Adjustment
- **Method:** `POST`
- **Path:** `/inventory/variants/{variantID}/adjustment`
- **Protected:** Yes
- **Path Parameters:**
  - `variantID` (uuid)
- **Request Body:**
  ```json
  {
    "quantity_change": 0,
    "note": "string (optional)"
  }
  ```
- **Response (200 OK):** `InventoryEvent` object.

### Reserve Inventory
- **Method:** `POST`
- **Path:** `/inventory/variants/{variantID}/reserve`
- **Protected:** Yes
- **Path Parameters:**
  - `variantID` (uuid)
- **Request Body:**
  ```json
  {
    "quantity": 0,
    "expires_at": "string (RFC3339 datetime) (optional)",
    "order_id": "uuid (optional)"
  }
  ```
- **Response (200 OK):** `InventoryReservation` object.
  ```json
  {
    "id": "uuid",
    "organization_id": "string",
    "product_variant_id": "uuid",
    "order_id": "uuid (optional)",
    "quantity": 0,
    "status": "string (active | released | expired)",
    "expires_at": "string (RFC3339 datetime) (optional)",
    "created_at": "string (RFC3339 datetime)",
    "released_at": "string (RFC3339 datetime) (optional)"
  }
  ```
- **Response (409 Conflict):** If there is insufficient stock.

### Release Reservation
- **Method:** `POST`
- **Path:** `/inventory/reservations/{reservationID}/release`
- **Protected:** Yes
- **Path Parameters:**
  - `reservationID` (uuid)
- **Response (204 No Content):** Empty response on success.

### Get Variant Stock
- **Method:** `GET`
- **Path:** `/inventory/variants/{variantID}/stock`
- **Protected:** Yes
- **Path Parameters:**
  - `variantID` (uuid)
- **Response (200 OK):** `VariantStock` object.
  ```json
  {
    "total_stock": 0,
    "reserved_stock": 0,
    "available_stock": 0
  }
  ```
