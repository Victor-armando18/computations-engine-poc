#!/bin/bash
BASE_URL="http://localhost:8080/orders/123/patch"

echo "========== TEST 1: PATCH REPLACE =========="
curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: test-replace-1" \
  -d '[{ "op": "replace", "path": "/items/0/quantity", "value": 3 }]'
echo -e "\n"

echo "========== TEST 2: PATCH ADD =========="
curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: test-add-1" \
  -d '[{ "op": "add", "path": "/items/-", "value": { "sku": "ITEM-003", "quantity": 2, "price": 75 } }]'
echo -e "\n"

echo "========== TEST 3: PATCH REMOVE =========="
curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: test-remove-1" \
  -d '[{ "op": "remove", "path": "/items/1" }]'
echo -e "\n"

echo "========== TEST 4: IDEMPOTENCY (repeat TEST 1) =========="
curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: test-replace-1" \
  -d '[{ "op": "replace", "path": "/items/0/quantity", "value": 3 }]'
echo -e "\n"

echo "========== TEST 5: RULES DISPARADAS =========="
curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: test-rule-1" \
  -d '[{ "op": "replace", "path": "/items/0/quantity", "value": 11 }]'
echo -e "\n"
