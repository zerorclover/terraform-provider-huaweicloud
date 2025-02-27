package huaweicloud

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDDSFlavorV3DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSFlavorV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSFlavorV3DataSourceID("data.huaweicloud_dds_flavors_v3.flavor"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dds_flavors_v3.flavor", "engine_name", "DDS-Community"),
				),
			},
		},
	})
}

func testAccCheckDDSFlavorV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find DDS Flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("DDS Flavor data source ID not set")
		}

		return nil
	}
}

var testAccDDSFlavorV3DataSource_basic = `
data "huaweicloud_dds_flavors_v3" "flavor" {
    engine_name = "DDS-Community"
    vcpus = 8
}
`
