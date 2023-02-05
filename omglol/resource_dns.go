package omglol

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultTtl = 3600
)

func resourceDns() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDnsCreate,
		ReadContext:   resourceDnsRead,
		UpdateContext: resourceDnsUpdate,
		DeleteContext: resourceDnsDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Default:  defaultTtl,
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDnsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	a := m.(auth)
	record := recordFromData(d)

	record, err := createRecord(a, record)
	if err != nil {
		return diag.FromErr(err)
	}
	id := strconv.Itoa(record.Id)
	d.SetId(id)

	return diag.Diagnostics{}
}

func resourceDnsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	a := m.(auth)
	recordId := d.Id()

	record, err := getRecord(a, recordId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = setResourceData(d, record)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func recordFromData(d *schema.ResourceData) dnsRecordIdInt {
	name := d.Get("name").(string)
	rType := d.Get("type").(string)
	data := d.Get("data").(string)
	ttl := d.Get("ttl").(int)

	record := dnsRecordIdInt{
		Name: name,
		Type: rType,
		Data: data,
		TTL:  ttl,
	}

	return record
}

func resourceDnsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceDnsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	a := m.(auth)
	id := d.Id()

	err := deleteRecord(a, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
