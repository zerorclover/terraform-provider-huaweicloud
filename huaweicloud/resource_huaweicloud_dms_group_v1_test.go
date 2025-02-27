package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccDmsGroupsV1_basic(t *testing.T) {
	var group groups.Group
	var groupName = fmt.Sprintf("dms_group_%s", acctest.RandString(5))
	var queueName = fmt.Sprintf("dms_queue_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDms(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Group_basic(groupName, queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1GroupExists("huaweicloud_dms_group_v1.group_1", group),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_group_v1.group_1", "name", groupName),
				),
			},
		},
	})
}

func testAccCheckDmsV1GroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud group client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dms_group_v1" {
			continue
		}

		queueID := rs.Primary.Attributes["queue_id"]
		page, err := groups.List(dmsClient, queueID, false).AllPages()
		if err == nil {
			groupsList, err := groups.ExtractGroups(page)
			if err != nil {
				return fmtp.Errorf("Error getting groups in queue %s: %s", queueID, err)
			}
			if len(groupsList) > 0 {
				for _, group := range groupsList {
					if group.ID == rs.Primary.ID {
						return fmtp.Errorf("The Dms group still exists.")
					}
				}
			}
		}
	}
	return nil
}

func testAccCheckDmsV1GroupExists(n string, group groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud group client: %s", err)
		}

		queueID := rs.Primary.Attributes["queue_id"]
		page, err := groups.List(dmsClient, queueID, false).AllPages()
		if err != nil {
			return fmtp.Errorf("Error getting groups in queue %s: %s", queueID, err)
		}

		groupsList, err := groups.ExtractGroups(page)
		if len(groupsList) > 0 {
			for _, found := range groupsList {
				if found.ID == rs.Primary.ID {
					group = found
					return nil
				}
			}
		}
		return fmtp.Errorf("The Dms group not found.")
	}
}

func testAccDmsV1Group_basic(groupName string, queueName string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_dms_queue_v1" "queue_1" {
			name  = "%s"
		}
		resource "huaweicloud_dms_group_v1" "group_1" {
			name = "%s"
			queue_id = "${huaweicloud_dms_queue_v1.queue_1.id}"
		}
	`, queueName, groupName)
}
