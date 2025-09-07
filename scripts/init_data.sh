#!/bin/bash
set -euo pipefail

DATA_DIR="${1:-data}" # можно передать путь (по умолчанию ./data)

mkdir -p "$DATA_DIR"

# Inventory
cat > "$DATA_DIR/inventory.json" << 'EOF'
[
  {
    "ingredient_id": "espresso_shot",
    "name": "Espresso Shot",
    "quantity": 500,
    "unit": "shots"
  },
  {
    "ingredient_id": "milk",
    "name": "Milk",
    "quantity": 5000,
    "unit": "ml"
  },
  {
    "ingredient_id": "flour",
    "name": "Flour",
    "quantity": 10000,
    "unit": "g"
  },
  {
    "ingredient_id": "blueberries",
    "name": "Blueberries",
    "quantity": 2000,
    "unit": "g"
  },
  {
    "ingredient_id": "sugar",
    "name": "Sugar",
    "quantity": 5000,
    "unit": "g"
  },
  {
    "ingredient_id": "cocoa_powder",
    "name": "Cocoa Powder",
    "quantity": 1000,
    "unit": "g"
  },
  {
    "ingredient_id": "vanilla_syrup",
    "name": "Vanilla Syrup",
    "quantity": 2000,
    "unit": "ml"
  },
  {
    "ingredient_id": "whipped_cream",
    "name": "Whipped Cream",
    "quantity": 1500,
    "unit": "ml"
  }
]
EOF

# Menu
cat > "$DATA_DIR/menu_items.json" << 'EOF'
[
  {
    "product_id": "latte",
    "name": "Caffe Latte",
    "description": "Espresso with steamed milk",
    "price": 3.50,
    "ingredients": [
      {
        "ingredient_id": "espresso_shot",
        "quantity": 1
      },
      {
        "ingredient_id": "milk",
        "quantity": 200
      }
    ]
  },
  {
    "product_id": "cappuccino",
    "name": "Cappuccino",
    "description": "Espresso with steamed milk and foam",
    "price": 3.25,
    "ingredients": [
      {
        "ingredient_id": "espresso_shot",
        "quantity": 1
      },
      {
        "ingredient_id": "milk",
        "quantity": 150
      }
    ]
  },
  {
    "product_id": "espresso",
    "name": "Espresso",
    "description": "Strong and bold coffee",
    "price": 2.50,
    "ingredients": [
      {
        "ingredient_id": "espresso_shot",
        "quantity": 1
      }
    ]
  },
  {
    "product_id": "mocha",
    "name": "Caffe Mocha",
    "description": "Espresso with chocolate and steamed milk",
    "price": 4.00,
    "ingredients": [
      {
        "ingredient_id": "espresso_shot",
        "quantity": 1
      },
      {
        "ingredient_id": "milk",
        "quantity": 180
      },
      {
        "ingredient_id": "cocoa_powder",
        "quantity": 15
      },
      {
        "ingredient_id": "whipped_cream",
        "quantity": 30
      }
    ]
  },
  {
    "product_id": "vanilla_latte",
    "name": "Vanilla Latte",
    "description": "Latte with vanilla syrup",
    "price": 3.75,
    "ingredients": [
      {
        "ingredient_id": "espresso_shot",
        "quantity": 1
      },
      {
        "ingredient_id": "milk",
        "quantity": 200
      },
      {
        "ingredient_id": "vanilla_syrup",
        "quantity": 20
      }
    ]
  },
  {
    "product_id": "blueberry_muffin",
    "name": "Blueberry Muffin",
    "description": "Freshly baked muffin with blueberries",
    "price": 2.50,
    "ingredients": [
      {
        "ingredient_id": "flour",
        "quantity": 100
      },
      {
        "ingredient_id": "blueberries",
        "quantity": 20
      },
      {
        "ingredient_id": "sugar",
        "quantity": 30
      }
    ]
  }
]
EOF

# Orders
echo '[]' > "$DATA_DIR/orders.json"

echo "Sample data initialized successfully!"
echo "Files created in $DATA_DIR:"
ls -1 "$DATA_DIR"
