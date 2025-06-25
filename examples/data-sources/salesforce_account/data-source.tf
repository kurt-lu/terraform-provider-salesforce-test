# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "salesforce_account" "example" {
  name = "Example Company"
}

output "account_id" {
  value = data.salesforce_account.example.id
}

output "account_type" {
  value = data.salesforce_account.example.type
}

output "account_industry" {
  value = data.salesforce_account.example.industry
} 