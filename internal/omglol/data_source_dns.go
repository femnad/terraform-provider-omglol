package omglol

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const baseUrl = "https://api.omg.lol"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsDataSource{}
)

func NewDNSDataSource() datasource.DataSource {
	return dnsDataSource{}
}

type dnsDataSourceModel struct {
}

type dnsModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Name types.String `tfsdk:"name"`
	Data types.String `tfsdk:"data"`
	TTL  types.Int64  `tfsdk:"ttl"`
}

type dnsDataSource struct {
	client *Client
}

func (d dnsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected data source configure type",
			fmt.Sprintf("Expected Client, got: %T", req.ProviderData),
		)
	}

	d.client = client
}

func (d dnsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns"
}

func (d dnsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"dns": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{Computed: true},
					"data": schema.StringAttribute{Computed: true},
					"ttl":  schema.StringAttribute{Computed: true},
				},
			},
		},
	},
	}
}

func (d dnsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsModel
	records, err := d.client.getRecords()
	if err != nil {
		resp.Diagnostics.AddError("Unable to read DNS records", err.Error())
		return
	}

	for _, record := range records {
		state = dnsModel{
			ID:   types.Int64Value(int64(record.ID)),
			Type: types.StringValue(record.Type),
			Name: types.StringValue(record.Name),
			Data: types.StringValue(record.Data),
			TTL:  types.Int64Value(int64(record.TTL)),
		}
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
