// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceAccount_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccount_basic("Test Account"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.salesforce_account.test", "id"),
					resource.TestCheckResourceAttr("data.salesforce_account.test", "name", "Test Account"),
				),
			},
		},
	})
}

func testAccDataSourceAccount_basic(name string) string {
	return fmt.Sprintf(`
resource "salesforce_account" "test" {
  name = "%s"
}

data "salesforce_account" "test" {
  name = salesforce_account.test.name
}
`, name)
} 