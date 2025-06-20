// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a new provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"salesforce": providerserver.NewProtocol6WithError(New()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func TestAccAccountResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAccountResourceConfig("Test Account"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("salesforce_account.test", "name", "Test Account"),
					resource.TestCheckResourceAttr("salesforce_account.test", "type", "Customer"),
					resource.TestCheckResourceAttr("salesforce_account.test", "industry", "Technology"),
					resource.TestCheckResourceAttrSet("salesforce_account.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "salesforce_account.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAccountResourceConfig("Updated Test Account"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("salesforce_account.test", "name", "Updated Test Account"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAccountResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "salesforce_account" "test" {
  name     = %[1]q
  type     = "Customer"
  industry = "Technology"
  phone    = "555-1234"
  website  = "https://example.com"
}
`, name)
} 