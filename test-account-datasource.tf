# Test configuration for the account data source
# This file can be used to manually test the account data source

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

# Create a test account first
resource "salesforce_account" "test" {
  name          = "Test Account for Data Source"
  account_number = "TEST-001"
  type          = "Customer"
  industry      = "Technology"
  phone         = "555-123-4567"
  website       = "https://test.example.com"
}

# Test the account data source
data "salesforce_account" "test_lookup" {
  name = salesforce_account.test.name
}

# Outputs to verify the data source works
output "created_account_id" {
  description = "ID of the created account"
  value       = salesforce_account.test.id
}

output "lookup_account_id" {
  description = "ID from the data source lookup"
  value       = data.salesforce_account.test_lookup.id
}

output "lookup_account_name" {
  description = "Name from the data source lookup"
  value       = data.salesforce_account.test_lookup.name
}

output "lookup_account_type" {
  description = "Type from the data source lookup"
  value       = data.salesforce_account.test_lookup.type
}

output "lookup_account_industry" {
  description = "Industry from the data source lookup"
  value       = data.salesforce_account.test_lookup.industry
}

output "lookup_account_phone" {
  description = "Phone from the data source lookup"
  value       = data.salesforce_account.test_lookup.phone
}

output "lookup_account_website" {
  description = "Website from the data source lookup"
  value       = data.salesforce_account.test_lookup.website
}

# Test with a non-existent account (should fail)
data "salesforce_account" "non_existent" {
  name = "Non-Existent Account Name"
}

output "non_existent_account_id" {
  description = "This should fail"
  value       = data.salesforce_account.non_existent.id
} 