/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package waf

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/huaweicloud/golangsdk/openstack/waf_hw/v1/webtamperprotection_rules"
)

func TestAccWafRuleWebTamperProtection_basic(t *testing.T) {
	var rule rules.WebTamper
	randName := acctest.RandString(5)
	resourceName := "huaweicloud_waf_rule_web_tamper_protection.rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckWafWafRuleWebTamperProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafWafRuleWebTamperProtection_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleWebTamperProtectionExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "domain", "www.abc.com"),
					resource.TestCheckResourceAttr(resourceName, "path", "/a"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccWafRuleImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckWafWafRuleWebTamperProtectionDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_waf_rule_web_tamper_protection" {
			continue
		}

		policyID := rs.Primary.Attributes["policy_id"]
		_, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("WAF rule still exists")
		}
	}

	return nil
}

func testAccCheckWafRuleWebTamperProtectionExists(n string, rule *rules.WebTamper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
		}

		policyID := rs.Primary.Attributes["policy_id"]
		found, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("WAF web tamper protection rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccWafWafRuleWebTamperProtection_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "huaweicloud_waf_rule_web_tamper_protection" "rule_1" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  domain    = "www.abc.com"
  path      = "/a"
}
`, name)
}
