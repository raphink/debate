#!/bin/bash
# Setup script for configuring gcloud with multiple accounts

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== GCP Multi-Account Configuration ===${NC}"
echo ""

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo -e "${YELLOW}gcloud CLI not found. Install from: https://cloud.google.com/sdk/docs/install${NC}"
    exit 1
fi

# List current gcloud configurations
echo -e "${GREEN}Current gcloud configurations:${NC}"
gcloud config configurations list
echo ""

# Ask user how they want to configure
echo "How would you like to configure gcloud for this project?"
echo "1. Use an existing gcloud configuration"
echo "2. Create a new gcloud configuration for this project"
echo "3. Use environment variables only (.gcloudrc file)"
read -p "Choose (1/2/3): " choice

case $choice in
    1)
        # Use existing config
        echo ""
        echo -e "${GREEN}Available configurations:${NC}"
        gcloud config configurations list
        echo ""
        read -p "Enter the configuration name to use: " config_name
        
        # Create .gcloudrc with the config name
        cat > .gcloudrc << EOF
# Use existing gcloud configuration
CLOUDSDK_ACTIVE_CONFIG_NAME=$config_name
EOF
        
        echo -e "${GREEN}Created .gcloudrc to use configuration: $config_name${NC}"
        echo "The deploy.sh script will automatically use this configuration"
        ;;
        
    2)
        # Create new config
        echo ""
        read -p "Enter a name for the new configuration (e.g., debate-personal): " config_name
        read -p "Enter your GCP account email: " account_email
        
        # Create the configuration
        gcloud config configurations create "$config_name" 2>/dev/null || echo "Configuration already exists, updating..."
        gcloud config configurations activate "$config_name"
        gcloud config set account "$account_email"
        gcloud config set project debate-480911
        
        # Authenticate if not already
        if ! gcloud auth list --filter="account:$account_email" --format="value(account)" | grep -q "$account_email"; then
            echo ""
            echo -e "${YELLOW}Authenticating with $account_email...${NC}"
            gcloud auth login "$account_email"
        fi
        
        # Create .gcloudrc
        cat > .gcloudrc << EOF
# Use dedicated gcloud configuration for this project
CLOUDSDK_ACTIVE_CONFIG_NAME=$config_name
EOF
        
        echo ""
        echo -e "${GREEN}✓ Created configuration: $config_name${NC}"
        echo -e "${GREEN}✓ Created .gcloudrc${NC}"
        echo "The deploy.sh script will automatically use this configuration"
        ;;
        
    3)
        # Use environment variables
        echo ""
        read -p "Enter your GCP account email: " account_email
        
        cat > .gcloudrc << EOF
# GCP configuration using environment variables
CLOUDSDK_CORE_ACCOUNT=$account_email
CLOUDSDK_CORE_PROJECT=debate-480911
EOF
        
        echo -e "${GREEN}Created .gcloudrc with account: $account_email${NC}"
        echo ""
        echo -e "${YELLOW}Note: You may need to run 'gcloud auth login $account_email' if not already authenticated${NC}"
        ;;
        
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}=== Setup Complete! ===${NC}"
echo ""
echo "Your .gcloudrc file:"
cat .gcloudrc
echo ""
echo "Next steps:"
echo "1. Run './deploy.sh --backend-only' to test the configuration"
echo "2. The .gcloudrc file is in .gitignore and won't be committed"
echo ""
echo "To switch accounts later, either:"
echo "  - Edit .gcloudrc manually"
echo "  - Run this script again"
echo "  - Use: gcloud config configurations activate <config-name>"
