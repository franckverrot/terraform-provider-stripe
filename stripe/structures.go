package stripe

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}

func expandMetadata(d *schema.ResourceData) map[string]string {
	old, new := d.GetChange("metadata")

	// Set the old values to empty string so that they can be removed
	expanded := expandStringMap(old.(map[string]interface{}))
	for key := range expanded {
		expanded[key] = ""
	}

	// Add entries for the new/updated fields defined in Terraform
	for key, value := range expandStringMap(new.(map[string]interface{})) {
		expanded[key] = value
	}

	return expanded
}

func expandStringList(d *schema.ResourceData, key string) []*string {
	elements := d.Get(key).([]interface{})

	if _, ok := d.GetOk(key); ok {
		expanded := make([]*string, len(elements))

		for i, element := range elements {
			tmp := element.(string)
			expanded[i] = &tmp
		}

		return expanded
	}

	return nil
}

func getMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
