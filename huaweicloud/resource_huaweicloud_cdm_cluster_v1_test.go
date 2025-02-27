// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCdmClusterV1_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdmClusterV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdmClusterV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdmClusterV1Exists(),
				),
			},
		},
	})
}

func testAccCdmClusterV1_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_cdm_flavors" "test" {}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "%s"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_cdm_cluster" "cluster" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  flavor_id         = data.huaweicloud_cdm_flavors.test.flavors[0].id
  name              = "%s"
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  version           = "2.8.2"
}`, rName, rName)
}

func testAccCheckCdmClusterV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.CdmV11Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cdm_cluster" {
			continue
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if err == nil {
			return fmtp.Errorf("huaweicloud_cdm_cluster still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckCdmClusterV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CdmV11Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_cdm_cluster.cluster"]
		if !ok {
			return fmtp.Errorf("Error checking huaweicloud_cdm_cluster.cluster exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return fmtp.Errorf("Error checking huaweicloud_cdm_cluster.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmtp.Errorf("huaweicloud_cdm_cluster.cluster is not exist")
			}
			return fmtp.Errorf("Error checking huaweicloud_cdm_cluster.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}
