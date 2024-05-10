#!/bin/bash

# Go to the user's home directory
cd /home/$AWS_USERNAME/$REPOSITORY_NAME

# If the repository directory doesn't exist, clone it
# if [ ! -d "$APP_NAME" ]; then
#   git clone $REPOSITORY_URL
# fi

# Navigate to the repository directory (which is now guaranteed to exist)
# cd $REPOSITORY_NAME

echo "Tidy go env"
go mod tidy
# Build the Go application
echo "Start build"
go build -o $APP_NAME .

# Check if systemd service exists, if not, create one
echo "Check if systemd exists"
if [ ! -e "/etc/systemd/system/$APP_NAME.service" ]; then
  echo "Create service."
  sudo bash -c "cat > /etc/systemd/system/$APP_NAME.service" <<EOF
[Unit]
Description=$APP_NAME Service
After=network.target

[Service]
ExecStart=/home/$AWS_USERNAME/$REPOSITORY_NAME/$APP_NAME
WorkingDirectory=/home/$AWS_USERNAME/$REPOSITORY_NAME
User=$AWS_USERNAME
Group=$AWS_USERNAME
Restart=always

[Install]
WantedBy=multi-user.target
EOF
else
  echo "Service already exists"
fi

echo "Restart daemon"
sudo systemctl daemon-reload

echo "Start service"
# Start the application service
sudo systemctl start $APP_NAME

# Enable the application service
sudo systemctl enable $APP_NAME

# Check the status of the service
sudo systemctl status $APP_NAME