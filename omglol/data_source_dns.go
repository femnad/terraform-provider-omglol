package omglol

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const baseUrl = "https://api.omg.lol"

func dataSourceDns() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDnsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceDnsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	a := m.(auth)
	records, err := getRecords(a)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	for _, record := range records {
		if record.Name == name {
			recordId := strconv.Itoa(record.Id)

			d.SetId(recordId)

			setErr := setResourceData(d, record)
			if setErr != nil {
				return diag.FromErr(setErr)
			}

			return diags
		}
	}

	return diag.Errorf("Unable to find record with name %s", name)
}

func setResourceData(d *schema.ResourceData, record dnsRecord) error {
	mapping := map[string]any{
		"data": &record.Data,
		"ttl":  &record.TTL,
		"type": &record.Type,
	}

	for key, attr := range mapping {
		err := d.Set(key, attr)
		if err != nil {
			return err
		}
	}

	return nil
}
