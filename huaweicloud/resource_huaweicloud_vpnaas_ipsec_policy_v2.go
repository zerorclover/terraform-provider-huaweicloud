package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceVpnIPSecPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpnIPSecPolicyV2Create,
		Read:   resourceVpnIPSecPolicyV2Read,
		Update: resourceVpnIPSecPolicyV2Update,
		Delete: resourceVpnIPSecPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encapsulation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pfs": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transform_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
			},
			"lifetime": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"units": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVpnIPSecPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	encapsulationMode := resourceIPSecPolicyV2EncapsulationMode(d.Get("encapsulation_mode").(string))
	authAlgorithm := resourceIPSecPolicyV2AuthAlgorithm(d.Get("auth_algorithm").(string))
	encryptionAlgorithm := resourceIPSecPolicyV2EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
	pfs := resourceIPSecPolicyV2PFS(d.Get("pfs").(string))
	transformProtocol := resourceIPSecPolicyV2TransformProtocol(d.Get("transform_protocol").(string))
	lifetime := resourceIPSecPolicyV2LifetimeCreateOpts(d.Get("lifetime").(*schema.Set))

	opts := VpnIPSecPolicyCreateOpts{
		ipsecpolicies.CreateOpts{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			TenantID:            d.Get("tenant_id").(string),
			EncapsulationMode:   encapsulationMode,
			AuthAlgorithm:       authAlgorithm,
			EncryptionAlgorithm: encryptionAlgorithm,
			PFS:                 pfs,
			TransformProtocol:   transformProtocol,
			Lifetime:            &lifetime,
		},
		MapValueSpecs(d),
	}

	logp.Printf("[DEBUG] Create IPSec policy: %#v", opts)

	policy, err := ipsecpolicies.Create(networkingClient, opts).Extract()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForIPSecPolicyCreation(networkingClient, policy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()

	logp.Printf("[DEBUG] IPSec policy created: %#v", policy)

	d.SetId(policy.ID)

	return resourceVpnIPSecPolicyV2Read(d, meta)
}

func resourceVpnIPSecPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	logp.Printf("[DEBUG] Retrieve information about IPSec policy: %s", d.Id())

	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	policy, err := ipsecpolicies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "IPSec policy")
	}

	logp.Printf("[DEBUG] Read HuaweiCloud IPSec policy %s: %#v", d.Id(), policy)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("tenant_id", policy.TenantID)
	d.Set("encapsulation_mode", policy.EncapsulationMode)
	d.Set("encryption_algorithm", policy.EncryptionAlgorithm)
	d.Set("transform_protocol", policy.TransformProtocol)
	d.Set("pfs", policy.PFS)
	d.Set("auth_algorithm", policy.AuthAlgorithm)
	d.Set("region", GetRegion(d, config))

	// Set the lifetime
	var lifetimeMap map[string]interface{}
	lifetimeMap = make(map[string]interface{})
	lifetimeMap["units"] = policy.Lifetime.Units
	lifetimeMap["value"] = policy.Lifetime.Value
	var lifetime []map[string]interface{}
	lifetime = append(lifetime, lifetimeMap)
	if err := d.Set("lifetime", &lifetime); err != nil {
		logp.Printf("[WARN] unable to set IPSec policy lifetime")
	}

	return nil
}

func resourceVpnIPSecPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var hasChange bool
	opts := ipsecpolicies.UpdateOpts{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	if d.HasChange("auth_algorithm") {
		opts.AuthAlgorithm = resourceIPSecPolicyV2AuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = resourceIPSecPolicyV2EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("transform_protocol") {
		opts.TransformProtocol = resourceIPSecPolicyV2TransformProtocol(d.Get("transform_protocol").(string))
		hasChange = true
	}

	if d.HasChange("pfs") {
		opts.PFS = resourceIPSecPolicyV2PFS(d.Get("pfs").(string))
		hasChange = true
	}

	if d.HasChange("encapsulation_mode") {
		opts.EncapsulationMode = resourceIPSecPolicyV2EncapsulationMode(d.Get("encapsulation_mode").(string))
		hasChange = true
	}

	if d.HasChange("lifetime") {
		lifetime := resourceIPSecPolicyV2LifetimeUpdateOpts(d.Get("lifetime").(*schema.Set))
		opts.Lifetime = &lifetime
		hasChange = true
	}

	logp.Printf("[DEBUG] Updating IPSec policy with id %s: %#v", d.Id(), opts)

	if hasChange {
		_, err = ipsecpolicies.Update(networkingClient, d.Id(), opts).Extract()
		if err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"ACTIVE"},
			Refresh:    waitForIPSecPolicyUpdate(networkingClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return err
		}
	}
	return resourceVpnIPSecPolicyV2Read(d, meta)
}

func resourceVpnIPSecPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	logp.Printf("[DEBUG] Destroy IPSec policy: %s", d.Id())

	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForIPSecPolicyDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return err
	}

	return nil
}

func waitForIPSecPolicyDeletion(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := ipsecpolicies.Delete(networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
			if errCode.Actual == 409 {
				return nil, "ACTIVE", nil
			}
		}

		return nil, "ACTIVE", err
	}
}

func waitForIPSecPolicyCreation(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ipsecpolicies.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func waitForIPSecPolicyUpdate(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ipsecpolicies.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func resourceIPSecPolicyV2TransformProtocol(trp string) ipsecpolicies.TransformProtocol {
	var protocol ipsecpolicies.TransformProtocol
	switch trp {
	case "esp":
		protocol = ipsecpolicies.TransformProtocolESP
	case "ah":
		protocol = ipsecpolicies.TransformProtocolAH
	case "ah-esp":
		protocol = ipsecpolicies.TransformProtocolAHESP
	}
	return protocol

}
func resourceIPSecPolicyV2PFS(pfsString string) ipsecpolicies.PFS {
	var pfs ipsecpolicies.PFS
	switch pfsString {
	case "group2":
		pfs = ipsecpolicies.PFSGroup2
	case "group5":
		pfs = ipsecpolicies.PFSGroup5
	case "group14":
		pfs = ipsecpolicies.PFSGroup14
	}
	return pfs

}
func resourceIPSecPolicyV2EncryptionAlgorithm(encryptionAlgo string) ipsecpolicies.EncryptionAlgorithm {
	var alg ipsecpolicies.EncryptionAlgorithm
	switch encryptionAlgo {
	case "3des":
		alg = ipsecpolicies.EncryptionAlgorithm3DES
	case "aes-128":
		alg = ipsecpolicies.EncryptionAlgorithmAES128
	case "aes-256":
		alg = ipsecpolicies.EncryptionAlgorithmAES256
	case "aes-192":
		alg = ipsecpolicies.EncryptionAlgorithmAES192
	}
	return alg
}
func resourceIPSecPolicyV2AuthAlgorithm(authAlgo string) ipsecpolicies.AuthAlgorithm {
	var alg ipsecpolicies.AuthAlgorithm
	switch authAlgo {
	case "md5":
		alg = ipsecpolicies.AuthAlgorithmMD5
	case "sha1":
		alg = ipsecpolicies.AuthAlgorithmSHA1
	case "sha2-256":
		alg = ipsecpolicies.AuthAlgorithmSHA256
	case "sha2-384":
		alg = ipsecpolicies.AuthAlgorithmSHA384
	case "sha2-512":
		alg = ipsecpolicies.AuthAlgorithmSHA512
	}
	return alg
}
func resourceIPSecPolicyV2EncapsulationMode(encMode string) ipsecpolicies.EncapsulationMode {
	var mode ipsecpolicies.EncapsulationMode
	switch encMode {
	case "tunnel":
		mode = ipsecpolicies.EncapsulationModeTunnel
	case "transport":
		mode = ipsecpolicies.EncapsulationModeTransport
	}
	return mode
}

func resourceIPSecPolicyV2LifetimeCreateOpts(d *schema.Set) ipsecpolicies.LifetimeCreateOpts {
	lifetime := ipsecpolicies.LifetimeCreateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		lifetime.Units = resourceIPSecPolicyV2Unit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetime.Value = value
	}
	return lifetime
}

func resourceIPSecPolicyV2Unit(units string) ipsecpolicies.Unit {
	var unit ipsecpolicies.Unit
	switch units {
	case "seconds":
		unit = ipsecpolicies.UnitSeconds
	case "kilobytes":
		unit = ipsecpolicies.UnitKilobytes
	}
	return unit
}

func resourceIPSecPolicyV2LifetimeUpdateOpts(d *schema.Set) ipsecpolicies.LifetimeUpdateOpts {
	lifetimeUpdateOpts := ipsecpolicies.LifetimeUpdateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		lifetimeUpdateOpts.Units = resourceIPSecPolicyV2Unit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetimeUpdateOpts.Value = value
	}
	return lifetimeUpdateOpts

}
