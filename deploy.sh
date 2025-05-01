#!/bin/bash
# AWS SSH TUI Portfolio Deployment Script

# Configuration - replace these with your actual values
AWS_KEY="path/to/your-key.pem"
AWS_HOST="ec2-user@your-ec2-instance-ip"
PROJECT_NAME="sh.kurttekin.com"
REMOTE_DIR="/home/ec2-user/${PROJECT_NAME}"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Building application locally...${NC}"
go build -o tuiserver ./cmd/tuiserver
if [ $? -ne 0 ]; then
    echo "Build failed! Aborting deployment."
    exit 1
fi

echo -e "${YELLOW}Copying files to server...${NC}"
# Create the directory if it doesn't exist
ssh -i $AWS_KEY $AWS_HOST "mkdir -p ${REMOTE_DIR}"

# Copy the application files
scp -i $AWS_KEY -r internal cmd go.mod go.sum LICENSE README.md tuiserver tuiserver.service $AWS_HOST:$REMOTE_DIR/

echo -e "${YELLOW}Setting up systemd service...${NC}"
ssh -i $AWS_KEY $AWS_HOST "sudo cp ${REMOTE_DIR}/tuiserver.service /etc/systemd/system/ && \
                         sudo systemctl daemon-reload && \
                         sudo systemctl enable tuiserver && \
                         sudo systemctl restart tuiserver && \
                         sudo systemctl status tuiserver"

echo -e "${GREEN}Deployment complete!${NC}"
echo -e "${GREEN}Your TUI SSH server is now running on port 2222${NC}"
echo -e "${GREEN}Users can connect with: ssh -t ${AWS_HOST#*@} -p 2222${NC}" 