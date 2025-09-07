#!/bin/bash

set -e

echo "=== Deploying Hot Coffee Application ==="

# Build the application
echo "Building application..."
go build -o hot-coffee .

# Create systemd service file
sudo tee /etc/systemd/system/hot-coffee.service > /dev/null << EOF
[Unit]
Description=Hot Coffee Management System
After=network.target

[Service]
Type=simple
User=coffee
Group=coffee
WorkingDirectory=/opt/hot-coffee
ExecStart=/opt/hot-coffee/hot-coffee --port 8080 --dir /opt/hot-coffee/data
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Create application directory
sudo mkdir -p /opt/hot-coffee/data
sudo cp hot-coffee /opt/hot-coffee/
sudo chown -R coffee:coffee /opt/hot-coffee

# Initialize with sample data if needed
if [ ! -f /opt/hot-coffee/data/menu_items.json ]; then
    echo "Initializing sample data..."
    sudo -u coffee ./scripts/init_data.sh
fi

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable hot-coffee
sudo systemctl restart hot-coffee

# Check service status
sudo systemctl status hot-coffee

echo "=== Deployment Complete ==="
echo "Service is running on port 8080"
echo "Check logs with: sudo journalctl -u hot-coffee -f"
