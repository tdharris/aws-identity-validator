# AWS Identity Validator

A diagnostic utility for validating AWS IAM Roles for Service Accounts (IRSA) in Kubernetes environments.

## Overview

AWS Identity Validator is a troubleshooting tool that helps validate whether your Kubernetes pods can properly authenticate to AWS services using IAM Roles for Service Accounts (IRSA). The tool verifies:

- Proper IRSA environment configuration within the pod
- Successful AWS authentication using both AWS SDK v1 and v2
- Identity information returned by STS GetCallerIdentity

This is particularly useful when troubleshooting EKS pods that need to access AWS services or when verifying your IRSA setup.

## How It Works

When deployed, AWS Identity Validator:

1. Checks for the presence of IRSA environment variables and token file
2. Attempts to authenticate with AWS using SDK v1
3. Attempts to authenticate with AWS using SDK v2
4. Reports detailed identity information from both SDK versions

### Example Output

```console
AWS Identity Validator
===========================

[IRSA Environment Check]
✓ AWS_WEB_IDENTITY_TOKEN_FILE env var is set: /var/run/secrets/eks.amazonaws.com/serviceaccount/token
✓ Token file exists
✓ Token file size: 1246 bytes
✓ AWS_ROLE_ARN env var is set: arn:aws:iam::<account-id>:role/<role-name>

[AWS SDK v1]
✓ Successfully authenticated with AWS SDK v1
Account ID: <account-id>
User ID: <user-id>:<session-id>
ARN: arn:aws:sts::<account-id>:assumed-role/<role-name>/<session-id>
Provider: WebIdentityCredentials
Access Key ID: <access-key-id>
Using temporary credentials (has session token)

[AWS SDK v2]
✓ Successfully authenticated with AWS SDK v2
Account ID: <account-id>
User ID: <user-id>:<session-id>
ARN: arn:aws:sts::<account-id>:assumed-role/<role-name>/<session-id>
Provider: WebIdentityCredentials
Access Key ID: <access-key-id>
Using temporary credentials (has session token)
```

## Deployment

### Prerequisites

- Kubernetes cluster with EKS IAM Roles for Service Accounts (IRSA) configured
- kubectl configured to access your cluster
- Helm v3

### Installation

Deploy using Helm:

```bash
helm install aws-identity-validator ./chart \
  --set serviceAccount.annotations."eks\.amazonaws\.com/role-arn"="arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME" \
  --set aws.region=us-east-2
```

Replace `ACCOUNT_ID` and `ROLE_NAME` with your AWS account ID and the IAM role name.

Alternatively, create a values file:

```yaml
serviceAccount:
  annotations:
    eks.amazonaws.com/role-arn: "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME"
aws:
  region: us-east-2
```

And deploy with:

```bash
helm install aws-identity-validator ./chart -f your-values.yaml
```

### Viewing Results

Check the pod logs to see authentication results:

```bash
kubectl logs job/aws-identity-validator
```

## Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `serviceAccount.name` | Name of the service account | `aws-identity-validator` |
| `serviceAccount.annotations` | Annotations for the service account | `{}` |
| `aws.region` | AWS region to use | `us-east-2` |
| `image.repository` | Image repository | `aws-identity-validator` |
| `image.tag` | Image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `env` | Additional environment variables | `[]` |
| `job.backoffLimit` | Number of retries before considering job as failed | `0` |

## Building Locally

### Prerequisites

- Go 1.20+
- Docker (for building container image)

### Build Steps

1. Build the binary:
   ```bash
   cd src
   go build -o aws-identity-validator
   ```

2. Build the container:
   ```bash
   docker build -t aws-identity-validator:latest .
   ```

## Troubleshooting

If the validator fails to authenticate:

1. Verify your EKS cluster has IRSA configured properly
2. Check that the IAM role exists and has proper trust relationships
3. Ensure the serviceAccount has the correct annotation with the proper IAM role ARN
4. Verify the IAM role has the necessary permissions

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
