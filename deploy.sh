#!/bin/bash
echo "APP_NAME is ${{ vars.APP_NAME }}"
echo "REPOSITORY_URL is ${{ vars.REPOSITORY_URL }}"
echo "Username is ${{ secrets.AWS_USERNAME }}"
echo "REPOSITORY_NAME is ${{ vars.REPOSITORY_NAME }}"
            
cd /home/${{ secrets.AWS_USERNAME }}
if [ ! -d "${{ vars.APP_NAME }}" ]; then
  echo "Pulling changes from repository"
cd ${{ vars.REPOSITORY_NAME }}
  git restore .
  git pull
  echo "Pull finished"
else
  echo "Cloning repository"
  git clone ${{ env.REPOSITORY_URL }}
  echo "Clone finished"
fi
# Go to the user's home directory
cd /home/$AWS_USERNAME

# If the repository directory doesn't exist, clone it
# if [ ! -d "$APP_NAME" ]; then
#   git clone $REPOSITORY_URL
# fi

# Navigate to the repository directory (which is now guaranteed to exist)
cd $REPOSITORY_NAME

go mod tidy
# Build the Go application
go build -o $APP_EXECUTABLE .

# Check if systemd service exists, if not, create one
if [ ! -e "/etc/systemd/system/$APP_NAME.service" ]; then
  sudo bash -c 'cat > /etc/systemd/system/$APP_NAME.service' <<EOF
[Unit]
Description=$APP_NAME Service
After=network.target

[Service]
ExecStart=/home/$AWS_USERNAME/$REPOSITORY_NAME/$APP_EXECUTABLE
WorkingDirectory=/home/$AWS_USERNAME/$REPOSITORY_NAME
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