package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceDmsGroupsV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDmsGroupsV1Create,
		Read:   resourceDmsGroupsV1Read,
		Delete: resourceDmsGroupsV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Required: true,
				ForceNew: true,
			},
			"queue_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumed_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"produced_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"produced_deadletters": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_deadletters": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDmsGroupsV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV1Client, err := config.DmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dms group client: %s", err)
	}

	var getGroups []groups.GroupOps

	n := groups.GroupOps{
		Name: d.Get("name").(string),
	}
	getGroups = append(getGroups, n)

	createOpts := &groups.CreateOps{
		Groups: getGroups,
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	v, err := groups.Create(dmsV1Client, d.Get("queue_id").(string), createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud group: %s", err)
	}
	logp.Printf("[INFO] group Name: %s", v[0].Name)

	// Store the group ID now
	d.SetId(v[0].ID)
	d.Set("queue_id", d.Get("queue_id").(string))

	return resourceDmsGroupsV1Read(d, meta)
}

func resourceDmsGroupsV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	dmsV1Client, err := config.DmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dms group client: %s", err)
	}

	queueID := d.Get("queue_id").(string)
	page, err := groups.List(dmsV1Client, queueID, false).AllPages()
	if err != nil {
		return fmtp.Errorf("Error getting groups in queue %s: %s", queueID, err)
	}
	groupsList, err := groups.ExtractGroups(page)
	if len(groupsList) < 1 {
		return fmtp.Errorf("No matching resource found.")
	}

	if len(groupsList) > 1 {
		return fmtp.Errorf("Multiple resources matched;")
	}

	group := groupsList[0]
	logp.Printf("[DEBUG] Dms group %s: %+v", d.Id(), group)

	d.SetId(group.ID)
	d.Set("name", group.Name)
	d.Set("consumed_messages", group.ConsumedMessages)
	d.Set("available_messages", group.AvailableMessages)
	d.Set("produced_messages", group.ProducedMessages)
	d.Set("produced_deadletters", group.ProducedDeadletters)
	d.Set("available_deadletters", group.AvailableDeadletters)

	return nil
}

func resourceDmsGroupsV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV1Client, err := config.DmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dms group client: %s", err)
	}

	err = groups.Delete(dmsV1Client, d.Get("queue_id").(string), d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud group: %s", err)
	}

	logp.Printf("[DEBUG] Dms group %s deactivated.", d.Id())
	d.SetId("")
	return nil
}
