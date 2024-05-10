#!/bin/bash

# Clone your repository
cd /home/$AWS_USERNAME
if [ ! -d "$APP_NAME" ]; then
  git clone $GITHUB_REPOSITORY_URL
fi

# Pull the latest changes
cd $APP_NAME
git pull

# Build your Go application
go build -o $APP_EXECUTABLE

# Create a new systemd service
if [ ! -e "/etc/systemd/system/$APP_NAME.service" ]; then
  sudo bash -c 'cat > /etc/systemd/system/$APP_NAME.service <<EOF
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
  EOF'
fi

# Start the application service
sudo systemctl start $APP_NAME

# Enable the application service
sudo systemctl enable $APP_NAME

# Check the status of your service
sudo systemctl status $APP_NAME
