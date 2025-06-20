# Salesforce Account Resource Example

This example demonstrates how to use the `salesforce_account` resource to create and manage Salesforce Account records.

## Prerequisites

1. A Salesforce org with API access
2. A Connected App configured with OAuth2
3. A user with System Administrator permissions

## Setup

1. Configure your Salesforce Connected App and obtain the required credentials
2. Set the required variables either through environment variables or a `.tfvars` file

## Usage

### Using Environment Variables

```bash
export SALESFORCE_CLIENT_ID="your-client-id"
export SALESFORCE_PRIVATE_KEY="your-private-key"
export SALESFORCE_USERNAME="your-username"
export SALESFORCE_API_VERSION="58.0"
export SALESFORCE_LOGIN_URL="https://login.salesforce.com"

terraform init
terraform plan
terraform apply
```

### Using a .tfvars file

Create a `terraform.tfvars` file:

```hcl
client_id   = "your-client-id"
private_key = "your-private-key"
username    = "your-username"
api_version = "58.0"
login_url   = "https://login.salesforce.com"
```

Then run:

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

This example creates two Salesforce Account records:

1. **Example Company** - A fully configured account with all fields populated
2. **Minimal Account** - An account with only the required `name` field

## Available Fields

The `salesforce_account` resource supports the following fields:

- `name` (Required) - The name of the account
- `account_number` (Optional) - Account number
- `type` (Optional) - Account type (e.g., Customer, Prospect, Partner)
- `industry` (Optional) - Industry (e.g., Technology, Healthcare, Finance)
- `phone` (Optional) - Phone number
- `website` (Optional) - Website URL

## Cleanup

To destroy the created resources:

```bash
terraform destroy
``` 