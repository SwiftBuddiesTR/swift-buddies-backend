#!/bin/bash

# Go to the user's home directory
cd /home/$AWS_USERNAME

# If the repository directory doesn't exist, clone it
if [ ! -d "$APP_NAME" ]; then
  git clone $REPOSITORY_URL
fi

# Navigate to the repository directory (which is now guaranteed to exist)
cd $APP_NAME

# Pull the latest changes
git pull

# Build the Go application
go build -o $APP_EXECUTABLE .

# Check if systemd service exists, if not, create one
if [ ! -e "/etc/systemd/system/$APP_NAME.service" ]; then
  sudo bash -c 'cat > /etc/systemd/system/$APP_NAME.service' <<EOF
[Unit]
Description=$APP_NAME Service
After=network.target

[Service]
ExecStart=/home/$AWS_USERNAME/$APP_NAME/$APP_EXECUTABLE
WorkingDirectory=/home/$AWS_USERNAME/$APP_NAME
User=$AWS_USERNAME
Group=$AWS_USERNAME
Restart=always

[Install]
WantedBy=multi-user.target
EOF
fi

# Start the application service
sudo systemctl start $APP_NAME

# Enable the application service
sudo systemctl enable $APP_NAME

# Check the status of the service
sudo systemctl status $APP_NAME