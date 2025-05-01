#!/bin/bash
# Run this script as root 

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Please run this script as root or with sudo${NC}"
  exit 1
fi

echo -e "${YELLOW}Setting up TUI portfolio on default SSH port (22)...${NC}"

# Step 1: Create the TUI user
echo -e "${YELLOW}Creating dedicated user for TUI access...${NC}"
if id "tui-portfolio" &>/dev/null; then
    echo "User tui-portfolio already exists"
else
    useradd -m tui-portfolio
    
    # Generate a random password
    TUI_PASSWORD=$(openssl rand -base64 12)
    echo "tui-portfolio:$TUI_PASSWORD" | chpasswd
    
    echo -e "${GREEN}Created user tui-portfolio with password: $TUI_PASSWORD${NC}"
    echo -e "${YELLOW}NOTE: This password is only needed if not using SSH key authentication${NC}"
fi

# Step 2: Check for the tuiserver binary
echo -e "${YELLOW}Checking for tuiserver binary...${NC}"
if [ -f "/usr/local/bin/tuiserver" ]; then
    echo "tuiserver binary already exists"
else
    # Ask for the path to the tuiserver binary
    echo -e "${YELLOW}Enter the path to your compiled tuiserver binary:${NC}"
    read -p "Path: " TUISERVER_PATH
    
    if [ ! -f "$TUISERVER_PATH" ]; then
        echo -e "${RED}Error: File not found at $TUISERVER_PATH${NC}"
        exit 1
    fi
    
    # Copy to /usr/local/bin and set permissions
    cp "$TUISERVER_PATH" /usr/local/bin/tuiserver
    chmod +x /usr/local/bin/tuiserver
    chown root:root /usr/local/bin/tuiserver
    
    echo -e "${GREEN}Installed tuiserver to /usr/local/bin/${NC}"
fi

# Step 3: Configure SSH
echo -e "${YELLOW}Configuring SSH server...${NC}"
SSHD_CONFIG="/etc/ssh/sshd_config"

# Check if the configuration already exists
if grep -q "Match User tui-portfolio" $SSHD_CONFIG; then
    echo "SSH configuration already exists for tui-portfolio"
else
    # Add configuration to sshd_config
    cat << EOF >> $SSHD_CONFIG

# Configuration for TUI Portfolio
Match User tui-portfolio
    ForceCommand /usr/local/bin/tuiserver
    PermitTTY yes
    AllowTcpForwarding no
    X11Forwarding no
EOF

    echo -e "${GREEN}Added SSH configuration for TUI portfolio${NC}"
fi

# Step 4: Configure for public key authentication only? (optional)
echo -e "${YELLOW}Would you like to require SSH key authentication for the TUI user?${NC}"
echo -e "${YELLOW}This is more secure but requires users to set up SSH keys.${NC}"
read -p "Require SSH keys? (y/n): " REQUIRE_KEYS

if [[ $REQUIRE_KEYS == "y" || $REQUIRE_KEYS == "Y" ]]; then
    # Check if already configured
    if grep -q "PasswordAuthentication no" $SSHD_CONFIG; then
        echo "Key authentication already configured"
    else
        # Update the Match block to add key authentication requirement
        sed -i '/Match User tui-portfolio/,/X11Forwarding no/ s/X11Forwarding no/X11Forwarding no\n    PasswordAuthentication no\n    AuthenticationMethods publickey/' $SSHD_CONFIG
        
        echo -e "${GREEN}Configured SSH to require key authentication${NC}"
        echo -e "${YELLOW}NOTE: Users will need to add their public SSH keys to:${NC}"
        echo -e "${YELLOW}/home/tui-portfolio/.ssh/authorized_keys${NC}"
        
        # Create .ssh directory for the user
        mkdir -p /home/tui-portfolio/.ssh
        chmod 700 /home/tui-portfolio/.ssh
        touch /home/tui-portfolio/.ssh/authorized_keys
        chmod 600 /home/tui-portfolio/.ssh/authorized_keys
        chown -R tui-portfolio:tui-portfolio /home/tui-portfolio/.ssh
    fi
fi

# Step 5: Restart SSH service
echo -e "${YELLOW}Restarting SSH service...${NC}"
systemctl restart sshd

# Check if SSH restart was successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN}SSH service restarted successfully${NC}"
else
    echo -e "${RED}Failed to restart SSH service. Please check configuration.${NC}"
    exit 1
fi

# Step 6: Success message
echo -e "\n${GREEN}==================================================${NC}"
echo -e "${GREEN}TUI Portfolio setup complete!${NC}"
echo -e "${GREEN}==================================================${NC}"
echo -e "${YELLOW}Users can now connect with:${NC}"
echo -e "${GREEN}ssh tui-portfolio@$(curl -s http://checkip.amazonaws.com)${NC}"
echo -e "${YELLOW}If you set up a domain name, users can connect with:${NC}"
echo -e "${GREEN}ssh tui-portfolio@yourdomain.com${NC}"
echo -e "${YELLOW}For a custom subdomain like tui.yourdomain.com, they can simply use:${NC}"
echo -e "${GREEN}ssh tui.yourdomain.com${NC}"

# Remind about firewall
echo -e "\n${YELLOW}Don't forget to ensure your security group allows SSH (port 22)${NC}"

exit 0 
