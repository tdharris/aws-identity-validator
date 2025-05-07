#!/usr/bin/env bash

set -euo pipefail

# Configuration from environment variables
ECR_REGISTRY="${ECR_REGISTRY:-}"  # Example: 123456789012.dkr.ecr.us-east-2.amazonaws.com
ECR_REPOSITORY="${ECR_REPOSITORY:-aws-identity-validator}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
AWS_REGION="${AWS_REGION:-us-east-2}"
AWS_ACCOUNT_ID="${AWS_ACCOUNT_ID:-}"
BINARY_NAME="aws-identity-validator"
SKIP_ECR="${SKIP_ECR:-false}"

# Default local image name
LOCAL_IMAGE_NAME="aws-identity-validator:${IMAGE_TAG}"

# Build the binary first
echo "🔨 Building Go binary..."
go build -o "$BINARY_NAME"
if [ -f "$BINARY_NAME" ]; then
    echo "✅ Binary build successful: $BINARY_NAME"
else
    echo "❌ Binary build failed"
    exit 1
fi

# Build Docker image
echo "🔨 Building Docker image..."
docker build -t "${LOCAL_IMAGE_NAME}" -f Dockerfile .

# Skip ECR push if requested
if [ "$SKIP_ECR" = "true" ]; then
    echo "🚫 Skipping ECR push as requested"
    exit 0
fi

# Check for required environment variables for ECR push
if [[ -z "$AWS_ACCOUNT_ID" && -z "$ECR_REGISTRY" ]]; then
    echo "Error: Either AWS_ACCOUNT_ID or ECR_REGISTRY must be set for ECR push"
    echo "Example usage: AWS_ACCOUNT_ID=123456789012 AWS_REGION=us-east-2 $0"
    echo "To skip ECR push: SKIP_ECR=true $0"
    exit 1
fi

# Construct ECR registry URL if not provided directly
if [[ -z "$ECR_REGISTRY" ]]; then
    ECR_REGISTRY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
fi

# Full ECR image reference
ECR_IMAGE="${ECR_REGISTRY}/${ECR_REPOSITORY}:${IMAGE_TAG}"

echo "🚀 Preparing to push to ECR: ${ECR_IMAGE}"

# Create the ECR repository if it doesn't exist
echo "📦 Checking if ECR repository exists..."
if ! aws ecr describe-repositories --repository-names "${ECR_REPOSITORY}" --region "${AWS_REGION}" >/dev/null 2>&1; then
    # prompt for confirmation
    read -p "ECR repository ${ECR_REPOSITORY} does not exist. Do you want to create it? (y/n): " confirm
    if [[ "$confirm" != "y" ]]; then
        echo "🚫 Exiting without creating ECR repository."
        exit 1
    fi
    echo "📦 Creating ECR repository ${ECR_REPOSITORY}..."
    aws ecr create-repository --repository-name "${ECR_REPOSITORY}" --region "${AWS_REGION}"
    echo "📦 ECR repository created."
fi

# Get ECR login password for Skopeo
echo "🔑 Getting ECR login credentials..."
ECR_PASSWORD=$(aws ecr get-login-password --region "${AWS_REGION}")

# Using Skopeo to sync the image to ECR
echo "🚀 Copying image to ECR using Skopeo..."
skopeo copy \
    --dest-creds "AWS:${ECR_PASSWORD}" \
    "docker-daemon:${LOCAL_IMAGE_NAME}" \
    "docker://${ECR_IMAGE}"

echo "✅ Successfully pushed ${LOCAL_IMAGE_NAME} to ${ECR_IMAGE}"

# Print Helm deployment instructions
echo ""
echo "🚢 To deploy with Helm, update the values.yaml file and run:"
echo "helm install aws-identity-validator ./chart/ \\"
echo "  --set image.repository=${ECR_REGISTRY}/${ECR_REPOSITORY} \\"
echo "  --set image.tag=${IMAGE_TAG} \\"
echo "  --set serviceAccount.annotations.\"eks\\.amazonaws\\.com/role-arn\"=arn:aws:iam::${AWS_ACCOUNT_ID}:role/YOUR_ROLE_NAME \\"
echo "  --set aws.region=${AWS_REGION}"
