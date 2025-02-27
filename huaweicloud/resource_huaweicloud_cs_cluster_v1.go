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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceCsClusterV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCsClusterV1Create,
		Read:   resourceCsClusterV1Read,
		Update: resourceCsClusterV1Update,
		Delete: resourceCsClusterV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"max_spu_num": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"subnet_cidr": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"subnet_gateway": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"vpc_cidr": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"manager_node_spu_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"used_spu_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceCsClusterV1UserInputParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"terraform_resource_data": d,
		"description":             d.Get("description"),
		"max_spu_num":             d.Get("max_spu_num"),
		"name":                    d.Get("name"),
		"subnet_cidr":             d.Get("subnet_cidr"),
		"subnet_gateway":          d.Get("subnet_gateway"),
		"vpc_cidr":                d.Get("vpc_cidr"),
	}
}

func resourceCsClusterV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CloudStreamV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	opts := resourceCsClusterV1UserInputParams(d)

	params, err := buildCsClusterV1CreateParameters(opts, nil)
	if err != nil {
		return fmtp.Errorf("Error building the request body of api(create), err=%s", err)
	}
	r, err := sendCsClusterV1CreateRequest(d, params, client)
	if err != nil {
		return fmtp.Errorf("Error creating CsClusterV1, err=%s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)

	obj, err := asyncWaitCsClusterV1Create(d, config, r, client, timeout)
	if err != nil {
		return err
	}
	id, err := navigateValue(obj, []string{"payload", "cluster_id"}, nil)
	if err != nil {
		return fmtp.Errorf("Error constructing id, err=%s", err)
	}
	d.SetId(convertToStr(id))

	return resourceCsClusterV1Read(d, meta)
}

func resourceCsClusterV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CloudStreamV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	res := make(map[string]interface{})

	v, err := sendCsClusterV1ReadRequest(d, client)
	if err != nil {
		return err
	}
	res["read"] = fillCsClusterV1ReadRespBody(v)

	states, err := flattenCsClusterV1Options(res)
	if err != nil {
		return err
	}

	return setCsClusterV1States(d, states)
}

func resourceCsClusterV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	client, err := config.CloudStreamV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	opts := resourceCsClusterV1UserInputParams(d)

	params, err := buildCsClusterV1UpdateParameters(opts, nil)
	if err != nil {
		return fmtp.Errorf("Error building the request body of api(update), err=%s", err)
	}
	if e, _ := isEmptyValue(reflect.ValueOf(params)); !e {
		_, err = sendCsClusterV1UpdateRequest(d, params, client)
		if err != nil {
			return fmtp.Errorf("Error updating (CsClusterV1: %v), err=%s", d.Id(), err)
		}
	}

	return resourceCsClusterV1Read(d, meta)
}

func resourceCsClusterV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CloudStreamV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	url, err := replaceVars(d, "reserved_cluster/{id}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	logp.Printf("[DEBUG] Deleting Cluster %q", d.Id())
	r := golangsdk.Result{}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:      successHTTPCodes,
		JSONBody:     nil,
		JSONResponse: &r.Body,
		MoreHeaders:  map[string]string{"Content-Type": "application/json"},
	})
	if r.Err != nil {
		return fmtp.Errorf("Error deleting Cluster %q, err=%s", d.Id(), r.Err)
	}

	_, err = asyncWaitCsClusterV1Delete(d, config, r.Body, client, d.Timeout(schema.TimeoutDelete))
	return err
}

func buildCsClusterV1CreateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	v, err := navigateValue(opts, []string{"description"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["desc"] = v
	}

	v, err = navigateValue(opts, []string{"max_spu_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["maximum_spu_quota"] = v
	}

	v, err = navigateValue(opts, []string{"name"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["name"] = v
	}

	v, err = navigateValue(opts, []string{"subnet_cidr"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["subnet_cidr"] = v
	}

	v, err = navigateValue(opts, []string{"subnet_gateway"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["subnet_gateway"] = v
	}

	v, err = navigateValue(opts, []string{"vpc_cidr"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["vpc_cidr"] = v
	}

	return params, nil
}

func sendCsClusterV1CreateRequest(d *schema.ResourceData, params interface{},
	client *golangsdk.ServiceClient) (interface{}, error) {
	url := client.ServiceURL("reserved_cluster")

	r := golangsdk.Result{}
	_, r.Err = client.Post(url, params, &r.Body, &golangsdk.RequestOpts{
		OkCodes: successHTTPCodes,
	})
	if r.Err != nil {
		return nil, fmtp.Errorf("Error running api(create), err=%s", r.Err)
	}
	return r.Body, nil
}

func asyncWaitCsClusterV1Create(d *schema.ResourceData, config *config.Config, result interface{},
	client *golangsdk.ServiceClient, timeout time.Duration) (interface{}, error) {

	data := make(map[string]interface{})
	pathParameters := map[string][]string{
		"cluster_id": []string{"payload", "cluster_id"},
	}
	for key, path := range pathParameters {
		value, err := navigateValue(result, path, nil)
		if err != nil {
			return nil, fmtp.Errorf("Error retrieving async operation path parameter, err=%s", err)
		}
		data[key] = value
	}

	url, err := replaceVars(d, "reserved_cluster/{cluster_id}", data)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	return waitToFinish(
		[]string{"3", "5", "8"},
		[]string{"1", "4"},
		timeout, 1*time.Second,
		func() (interface{}, string, error) {
			r := golangsdk.Result{}
			_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
				MoreHeaders: map[string]string{"Content-Type": "application/json"}})
			if r.Err != nil {
				return nil, "", nil
			}

			status, err := navigateValue(r.Body, []string{"payload", "status_code"}, nil)
			if err != nil {
				return nil, "", nil
			}
			return r.Body, convertToStr(status), nil
		},
	)
}

func buildCsClusterV1UpdateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	v, err := navigateValue(opts, []string{"description"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["desc"] = v
	}

	v, err = navigateValue(opts, []string{"max_spu_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["maximum_spu_quota"] = v
	}

	v, err = navigateValue(opts, []string{"name"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["name"] = v
	}

	return params, nil
}

func sendCsClusterV1UpdateRequest(d *schema.ResourceData, params interface{},
	client *golangsdk.ServiceClient) (interface{}, error) {
	url, err := replaceVars(d, "reserved_cluster/{id}", nil)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	r := golangsdk.Result{}
	_, r.Err = client.Patch(url, params, &r.Body, &golangsdk.RequestOpts{
		OkCodes: successHTTPCodes,
	})
	if r.Err != nil {
		return nil, fmtp.Errorf("Error running api(update), err=%s", r.Err)
	}
	return r.Body, nil
}

func asyncWaitCsClusterV1Delete(d *schema.ResourceData, config *config.Config, result interface{},
	client *golangsdk.ServiceClient, timeout time.Duration) (interface{}, error) {

	url, err := replaceVars(d, "reserved_cluster/{id}", nil)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	return waitToFinish(
		[]string{"Done"}, []string{"Pending"}, timeout, 1*time.Second,
		func() (interface{}, string, error) {
			r := golangsdk.Result{}
			_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
				OkCodes:     []int{200, 400},
				MoreHeaders: map[string]string{"Content-Type": "application/json"}})
			if r.Err != nil {
				return nil, "", nil
			}

			if checkCsClusterV1DeleteFinished(r.Body) {
				return r.Body, "Done", nil
			}
			return r.Body, "Pending", nil
		},
	)
}

func sendCsClusterV1ReadRequest(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	url, err := replaceVars(d, "reserved_cluster/{id}", nil)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json"}})
	if r.Err != nil {
		return nil, fmtp.Errorf("Error running api(read) for resource(CsClusterV1), err=%s", r.Err)
	}

	v, err := navigateValue(r.Body, []string{"payload"}, nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func fillCsClusterV1ReadRespBody(body interface{}) interface{} {
	result := make(map[string]interface{})
	val, ok := body.(map[string]interface{})
	if !ok {
		val = make(map[string]interface{})
	}

	if v, ok := val["created_at"]; ok {
		result["created_at"] = v
	} else {
		result["created_at"] = nil
	}

	if v, ok := val["desc"]; ok {
		result["desc"] = v
	} else {
		result["desc"] = nil
	}

	if v, ok := val["manager_node_spu"]; ok {
		result["manager_node_spu"] = v
	} else {
		result["manager_node_spu"] = nil
	}

	if v, ok := val["maximum_spu_quota"]; ok {
		result["maximum_spu_quota"] = v
	} else {
		result["maximum_spu_quota"] = nil
	}

	if v, ok := val["name"]; ok {
		result["name"] = v
	} else {
		result["name"] = nil
	}

	if v, ok := val["spu_used"]; ok {
		result["spu_used"] = v
	} else {
		result["spu_used"] = nil
	}

	if v, ok := val["subnet_cidr"]; ok {
		result["subnet_cidr"] = v
	} else {
		result["subnet_cidr"] = nil
	}

	if v, ok := val["subnet_gateway"]; ok {
		result["subnet_gateway"] = v
	} else {
		result["subnet_gateway"] = nil
	}

	if v, ok := val["vpc_cidr"]; ok {
		result["vpc_cidr"] = v
	} else {
		result["vpc_cidr"] = nil
	}

	return result
}

func flattenCsClusterV1Options(response map[string]interface{}) (map[string]interface{}, error) {
	opts := make(map[string]interface{})

	v, err := flattenCsClusterV1CreatedAT(response, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Cluster:created_at, err: %s", err)
	}
	opts["created_at"] = v

	v, err = navigateValue(response, []string{"read", "desc"}, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Cluster:description, err: %s", err)
	}
	opts["description"] = v

	v, err = flattenCsClusterV1ManagerNodeSpuNum(response, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Cluster:manager_node_spu_num, err: %s", err)
	}
	opts["manager_node_spu_num"] = v

	v, err = flattenCsClusterV1MaxSpuNum(response, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Cluster:max_spu_num, err: %s", err)
	}
	opts["max_spu_num"] = v

	v, err = navigateValue(response, []string{"read", "name"}, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Cluster:name, err: %s", err)
	}
	opts["name"] = v

	v, err = navigateValue(response, []string{"read", "subnet_cidr"}, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Cluster:subnet_cidr, err: %s", err)
	}
	opts["subnet_cidr"] = v

	v, err = navigateValue(response, []string{"read", "subnet_gateway"}, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Cluster:subnet_gateway, err: %s", err)
	}
	opts["subnet_gateway"] = v

	v, err = flattenCsClusterV1UsedSpuNum(response, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Cluster:used_spu_num, err: %s", err)
	}
	opts["used_spu_num"] = v

	v, err = navigateValue(response, []string{"read", "vpc_cidr"}, nil)
	if err != nil {
		return nil, fmtp.Errorf("Error flattening Cluster:vpc_cidr, err: %s", err)
	}
	opts["vpc_cidr"] = v

	return opts, nil
}

func flattenCsClusterV1CreatedAT(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"read", "created_at"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	v1 := int64(v.(float64) / 1000)
	return convertSeconds2Str(v1), nil
}

func flattenCsClusterV1ManagerNodeSpuNum(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"read", "manager_node_spu"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	return convertToInt(v)
}

func flattenCsClusterV1MaxSpuNum(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"read", "maximum_spu_quota"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	return convertToInt(v)
}

func flattenCsClusterV1UsedSpuNum(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"read", "spu_used"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	return convertToInt(v)
}

func setCsClusterV1States(d *schema.ResourceData, opts map[string]interface{}) error {
	for k, v := range opts {
		//lintignore:R001
		if err := d.Set(k, v); err != nil {
			return fmtp.Errorf("Error setting CsClusterV1:%s, err: %s", k, err)
		}
	}
	return nil
}
