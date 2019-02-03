package stripe

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	log.Printf("expandedmetadata: %#v\n", result)
	return result
}

func expandMetadata(d *schema.ResourceData) map[string]string {
	return expandStringMap(d.Get("metadata").(map[string]interface{}))
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
