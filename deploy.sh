#!/bin/bash

# Go to the user's home directory
cd /home/$AWS_USERNAME/$REPOSITORY_NAME

# Tidy go environment
echo "Tidy go env"
go mod tidy

# Build the Go application
echo "Start build"
go build -o $APP_NAME .

# Check if systemd service exists, if not, create one
echo "Check if systemd exists"
if [ -e "/etc/systemd/system/$APP_NAME.service" ]; then
  echo "Service already exists"
  echo "Deleting old service"
  sudo rm -rf /etc/systemd/system/$APP_NAME.service
fi

# Write a new service file
echo "Creating new service."
sudo bash -c "cat > /etc/systemd/system/$APP_NAME.service" <<EOF
[Unit]
Description=$APP_NAME Service
After=network.target

[Service]
Environment="MONGODB_URI=$MONGODB_URI"
Environment="SECRET_KEY=$SECRET_KEY"
ExecStart=/home/$AWS_USERNAME/$REPOSITORY_NAME/$APP_NAME
WorkingDirectory=/home/$AWS_USERNAME/$REPOSITORY_NAME
User=$AWS_USERNAME
Group=$AWS_USERNAME
Restart=always

[Install]
WantedBy=multi-user.target
EOF

echo "Service created"

# Reload and start the service
echo "Restart daemon"
sudo systemctl daemon-reload

echo "Start service"
sudo systemctl start $APP_NAME

# Enable the application service
sudo systemctl enable $APP_NAME

# Check the status of the service
sudo systemctl status $APP_NAME