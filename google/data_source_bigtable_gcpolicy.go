package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleBigtableGCPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleBigtableGCPolicyRead,

		Schema: map[string]*schema.Schema{
			"max_num_versions": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"max_age", "union", "intersection"},
			},
			"max_age": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"max_num_versions", "union", "intersection"},
			},
			"union": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:      true,
				ConflictsWith: []string{"max_num_versions", "max_age", "intersection"},
			},
			"intersection": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:      true,
				ConflictsWith: []string{"max_num_versions", "max_age", "union"},
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleBigtableGCPolicyRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("bigtable-gc-policy")

	if maxNumVersions, ok := d.GetOk("max_num_versions"); ok {
		d.Set("json", fmt.Sprintf("{\"maxNumVersions\":%d}", maxNumVersions))
		return nil
	}

	if maxAge, ok := d.GetOk("max_age"); ok {
		d.Set("json", fmt.Sprintf("{\"maxAge\":%s}", maxAge))
		return nil
	}

	policies := convertStringSet(d.Get("intersection").(*schema.Set))
	if len(policies) > 0 {
		json := "{\"intersection\":["
		for i, policy := range policies {
			json = json + policy
			if i < (len(policies) - 1) {
				json = json + ","
			}
		}
		d.Set("json", json+"]}")
		return nil
	}

	policies = convertStringSet(d.Get("union").(*schema.Set))
	if len(policies) > 0 {
		json := "{\"union\":["
		for i, policy := range policies {
			json = json + policy
			if i < (len(policies) - 1) {
				json = json + ","
			}
		}
		d.Set("json", json+"]}")
		return nil
	}

	return nil
}
