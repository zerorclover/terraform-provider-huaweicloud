package utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
)

// UpdateResourceTags is a helper to update the tags for a resource.
// It expects the tags field to be named "tags"
func UpdateResourceTags(conn *golangsdk.ServiceClient, d *schema.ResourceData, resourceType, id string) error {
	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		// remove old tags
		if len(oMap) > 0 {
			taglist := ExpandResourceTags(oMap)
			err := tags.Delete(conn, resourceType, id, taglist).ExtractErr()
			if err != nil {
				return err
			}
		}

		// set new tags
		if len(nMap) > 0 {
			taglist := ExpandResourceTags(nMap)
			err := tags.Create(conn, resourceType, id, taglist).ExtractErr()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// TagsToMap returns the list of tags into a map.
func TagsToMap(tags []tags.ResourceTag) map[string]string {
	result := make(map[string]string)
	for _, val := range tags {
		result[val.Key] = val.Value
	}

	return result
}

// ExpandResourceTags returns the tags for the given map of data.
func ExpandResourceTags(tagmap map[string]interface{}) []tags.ResourceTag {
	var taglist []tags.ResourceTag

	for k, v := range tagmap {
		tag := tags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	return taglist
}

// GetDNSZoneTagType returns resource tag type of DNS zone by zoneType
func GetDNSZoneTagType(zoneType string) (string, error) {
	if zoneType == "public" {
		return "DNS-public_zone", nil
	} else if zoneType == "private" {
		return "DNS-private_zone", nil
	}
	return "", fmt.Errorf("invalid zone type: %s", zoneType)
}

// GetDNSRecordSetTagType returns resource tag type of DNS record set by zoneType
func GetDNSRecordSetTagType(zoneType string) (string, error) {
	if zoneType == "public" {
		return "DNS-public_recordset", nil
	} else if zoneType == "private" {
		return "DNS-private_recordset", nil
	}
	return "", fmt.Errorf("invalid zone type: %s", zoneType)
}
