#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== Testing Coffee Shop Management API ==="
echo

echo "1. Testing Inventory API"
echo "Getting all inventory items:"
curl -s "$BASE_URL/inventory" | jq .
echo
echo

echo "2. Testing Menu API"
echo "Getting all menu items:"
curl -s "$BASE_URL/menu" | jq .
echo
echo

echo "3. Testing Order Creation"
echo "Creating a new order:"
ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Alice Smith",
    "items": [
      {
        "product_id": "latte",
        "quantity": 2
      },
      {
        "product_id": "blueberry_muffin",
        "quantity": 1
      }
    ]
  }')

echo "$ORDER_RESPONSE" | jq .
ORDER_ID=$(echo "$ORDER_RESPONSE" | jq -r '.order_id')
echo
echo

echo "4. Getting the created order:"
curl -s "$BASE_URL/orders/$ORDER_ID" | jq .
echo
echo

echo "5. Getting all orders:"
curl -s "$BASE_URL/orders" | jq .
echo
echo

echo "6. Closing the order:"
curl -s -X POST "$BASE_URL/orders/$ORDER_ID/close"
echo "Order closed"
echo
echo

echo "7. Testing Reports"
echo "Getting total sales:"
curl -s "$BASE_URL/reports/total-sales" | jq .
echo
echo

echo "Getting popular items:"
curl -s "$BASE_URL/reports/popular-items" | jq .
echo
echo

echo "8. Checking updated inventory:"
curl -s "$BASE_URL/inventory" | jq .
echo

echo "=== API Test Complete ==="