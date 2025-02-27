package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/queues"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccDmsQueuesV1_basic(t *testing.T) {
	var queue queues.Queue
	var queueName = fmt.Sprintf("dms_queue_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDms(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsV1QueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Queue_basic(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1QueueExists("huaweicloud_dms_queue_v1.queue_1", queue),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_queue_v1.queue_1", "name", queueName),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_queue_v1.queue_1", "queue_mode", "NORMAL"),
				),
			},
		},
	})
}

func TestAccDmsQueuesV1_FIFOmode(t *testing.T) {
	var queue queues.Queue
	var queueName = fmt.Sprintf("dms_queue_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDms(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsV1QueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Queue_FIFOmode(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1QueueExists("huaweicloud_dms_queue_v1.queue_1", queue),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_queue_v1.queue_1", "name", queueName),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_queue_v1.queue_1", "description", "test create dms queue"),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_queue_v1.queue_1", "queue_mode", "FIFO"),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_queue_v1.queue_1", "redrive_policy", "enable"),
					resource.TestCheckResourceAttr(
						"huaweicloud_dms_queue_v1.queue_1", "max_consume_count", "80"),
				),
			},
		},
	})
}

func testAccCheckDmsV1QueueDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud queue client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dms_queue_v1" {
			continue
		}

		_, err := queues.Get(dmsClient, rs.Primary.ID, false).Extract()
		if err == nil {
			return fmtp.Errorf("The Dms queue still exists.")
		}
	}
	return nil
}

func testAccCheckDmsV1QueueExists(n string, queue queues.Queue) resource.TestCheckFunc {
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
			return fmtp.Errorf("Error creating HuaweiCloud queue client: %s", err)
		}

		v, err := queues.Get(dmsClient, rs.Primary.ID, false).Extract()
		if err != nil {
			return fmtp.Errorf("Error getting HuaweiCloud queue: %s, err: %s", rs.Primary.ID, err)
		}
		if v.ID != rs.Primary.ID {
			return fmtp.Errorf("The Dms queue not found.")
		}
		queue = *v
		return nil
	}
}

func testAccDmsV1Queue_basic(queueName string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_dms_queue_v1" "queue_1" {
			name  = "%s"
		}
	`, queueName)
}

func testAccDmsV1Queue_FIFOmode(queueName string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_dms_queue_v1" "queue_1" {
			name  = "%s"
			description  = "test create dms queue"
			queue_mode  = "FIFO"
			redrive_policy  = "enable"
          max_consume_count = 80
		}
	`, queueName)
}
