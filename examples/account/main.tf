terraform {
  required_providers {
    salesforce = {
      source = "hashicorp/salesforce"
    }
  }
}

provider "salesforce" {
  client_id   = var.client_id
  private_key = var.private_key
  username    = var.username
  api_version = var.api_version
  login_url   = var.login_url
}

# Create a new Salesforce Account
resource "salesforce_account" "example" {
  name          = "Example Company"
  account_number = "ACC-001"
  type          = "Customer"
  industry      = "Technology"
  phone         = "555-123-4567"
  website       = "https://example.com"
}

# Create another account with minimal fields
resource "salesforce_account" "minimal" {
  name = "Minimal Account"
}

# Output the created account IDs
output "example_account_id" {
  description = "ID of the example account"
  value       = salesforce_account.example.id
}

output "minimal_account_id" {
  description = "ID of the minimal account"
  value       = salesforce_account.minimal.id
} 