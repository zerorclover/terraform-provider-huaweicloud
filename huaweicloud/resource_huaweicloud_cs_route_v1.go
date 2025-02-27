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
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceCsRouteV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCsRouteV1Create,
		Read:   resourceCsRouteV1Read,
		Delete: resourceCsRouteV1Delete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"peering_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCsRouteV1UserInputParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"terraform_resource_data": d,
		"cluster_id":              d.Get("cluster_id"),
		"destination":             d.Get("destination"),
		"peering_id":              d.Get("peering_id"),
	}
}

func resourceCsRouteV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CloudStreamV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	opts := resourceCsRouteV1UserInputParams(d)

	params, err := buildCsRouteV1CreateParameters(opts, nil)
	if err != nil {
		return fmtp.Errorf("Error building the request body of api(create), err=%s", err)
	}
	r, err := sendCsRouteV1CreateRequest(d, params, client)
	if err != nil {
		return fmtp.Errorf("Error creating CsRouteV1, err=%s", err)
	}

	id, err := navigateValue(r, []string{"route", "id"}, nil)
	if err != nil {
		return fmtp.Errorf("Error constructing id, err=%s", err)
	}
	d.SetId(convertToStr(id))

	return resourceCsRouteV1Read(d, meta)
}

func resourceCsRouteV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CloudStreamV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	res := make(map[string]interface{})

	v, err := fetchCsRouteV1ByList(d, client)
	if err != nil {
		return err
	}
	res["list"] = fillCsRouteV1ListRespBody(v)

	states, err := flattenCsRouteV1Options(res)
	if err != nil {
		return err
	}

	return setCsRouteV1States(d, states)
}

func resourceCsRouteV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CloudStreamV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	url, err := replaceVars(d, "reserved_cluster/{cluster_id}/peering/{peering_id}/route/{id}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	logp.Printf("[DEBUG] Deleting Route %q", d.Id())
	r := golangsdk.Result{}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:      successHTTPCodes,
		JSONBody:     nil,
		JSONResponse: nil,
		MoreHeaders:  map[string]string{"Content-Type": "application/json"},
	})
	if r.Err != nil {
		return fmtp.Errorf("Error deleting Route %q, err=%s", d.Id(), r.Err)
	}

	return nil
}

func buildCsRouteV1CreateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	v, err := navigateValue(opts, []string{"destination"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["destination"] = v
	}

	return params, nil
}

func sendCsRouteV1CreateRequest(d *schema.ResourceData, params interface{},
	client *golangsdk.ServiceClient) (interface{}, error) {
	url, err := replaceVars(d, "reserved_cluster/{cluster_id}/peering/{peering_id}/route", nil)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	r := golangsdk.Result{}
	_, r.Err = client.Post(url, params, &r.Body, &golangsdk.RequestOpts{
		OkCodes: successHTTPCodes,
	})
	if r.Err != nil {
		return nil, fmtp.Errorf("Error running api(create), err=%s", r.Err)
	}
	return r.Body, nil
}

func fetchCsRouteV1ByList(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	link, err := replaceVars(d, "reserved_cluster/{cluster_id}/peering/{peering_id}/route", nil)
	if err != nil {
		return nil, err
	}
	link = client.ServiceURL(link)

	return findCsRouteV1ByList(client, link, d.Id())
}

func findCsRouteV1ByList(client *golangsdk.ServiceClient, link, resourceID string) (interface{}, error) {
	r, err := sendCsRouteV1ListRequest(client, link)
	if err != nil {
		return nil, err
	}
	for _, item := range r.([]interface{}) {
		val, ok := item.(map[string]interface{})["id"]
		if ok && resourceID == convertToStr(val) {
			return item, nil
		}
	}

	return nil, fmtp.Errorf("Error finding the resource by list api")
}

func sendCsRouteV1ListRequest(client *golangsdk.ServiceClient, url string) (interface{}, error) {
	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json"}})
	if r.Err != nil {
		return nil, fmtp.Errorf("Error running api(list) for resource(CsRouteV1), err=%s", r.Err)
	}

	v, err := navigateValue(r.Body, []string{"routes"}, nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func fillCsRouteV1ListRespBody(body interface{}) interface{} {
	result := make(map[string]interface{})
	val, ok := body.(map[string]interface{})
	if !ok {
		val = make(map[string]interface{})
	}

	if v, ok := val["destination"]; ok {
		result["destination"] = v
	} else {
		result["destination"] = nil
	}

	return result
}

func flattenCsRouteV1Options(response map[string]interface{}) (map[string]interface{}, error) {
	opts := make(map[string]interface{})

	v, err := navigateValue(response, []string{"list", "destination"}, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Route:destination, err: %s", err)
	}
	opts["destination"] = v

	return opts, nil
}

func setCsRouteV1States(d *schema.ResourceData, opts map[string]interface{}) error {
	for k, v := range opts {
		//lintignore:R001
		if err := d.Set(k, v); err != nil {
			return fmtp.Errorf("Error setting CsRouteV1:%s, err: %s", k, err)
		}
	}
	return nil
}
