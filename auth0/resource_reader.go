package auth0

import "github.com/hashicorp/terraform/helper/schema"

func readStringFromResource(d *schema.ResourceData, key string) string {
	if attr, ok := d.GetOk(key); ok {
		return attr.(string)
	}
	return ""
}

func readBoolFromResource(d *schema.ResourceData, key string) bool {
	if attr, ok := d.GetOk(key); ok {
		return attr.(bool)
	}
	return false
}

func readMapFromResource(d *schema.ResourceData, key string) map[string]interface{} {

	if attr, ok := d.GetOk(key); ok {
		result := attr.(map[string]interface{})
		return result
	}

	return nil
}
