package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDmsMaintainWindowV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDms(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsMaintainWindowV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsMaintainWindowV1DataSourceID("data.huaweicloud_dms_maintainwindow_v1.maintainwindow1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_maintainwindow_v1.maintainwindow1", "seq", "1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dms_maintainwindow_v1.maintainwindow1", "begin", "22"),
				),
			},
		},
	})
}

func testAccCheckDmsMaintainWindowV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find Dms maintainwindow data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Dms maintainwindow data source ID not set")
		}

		return nil
	}
}

var testAccDmsMaintainWindowV1DataSource_basic = fmt.Sprintf(`
data "huaweicloud_dms_maintainwindow_v1" "maintainwindow1" {
seq = 1
}
`)
