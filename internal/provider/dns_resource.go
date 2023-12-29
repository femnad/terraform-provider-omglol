package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/femnad/terraform-provider-omglol/internal/omglol"
)

const (
	defaultTtl = 3600
)

var (
	_ resource.Resource              = &dnsResource{}
	_ resource.Resource              = &dnsResource{}
	_ resource.ResourceWithConfigure = &dnsResource{}
)

type dnsRecordModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Name types.String `tfsdk:"name"`
	Data types.String `tfsdk:"data"`
	TTL  types.Int64  `tfsdk:"ttl"`
}

func NewDNSResource() resource.Resource {
	return &dnsResource{}
}

type dnsResource struct {
	client *omglol.Client
}

func (d dnsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*omglol.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Expected *omglol.Client, got: %T", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d dnsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns"
}

func (d dnsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{Computed: true},
			"type": schema.StringAttribute{Computed: true},
			"data": schema.StringAttribute{Computed: true},
			"ttl":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d dnsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dnsRecordModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	planRecord := omglol.DNSRecord{
		Type: plan.Type.ValueString(),
		Name: plan.Name.ValueString(),
		Data: plan.Data.ValueString(),
		TTL:  int(plan.TTL.ValueInt64()),
	}

	record, err := d.client.CreateRecord(planRecord)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating record",
			fmt.Sprintf("Could not create DNS record: %s", err),
		)
	}

	plan.ID = types.Int64Value(int64(record.ID))
	plan.Type = types.StringValue(record.Type)
	plan.Name = types.StringValue(record.Name)
	plan.Data = types.StringValue(record.Data)
	plan.TTL = types.Int64Value(int64(record.TTL))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d dnsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dnsRecordModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	record, err := d.client.GetRecord(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading record",
			fmt.Sprintf("Could not read DNS record with ID %s: %v", state.ID, err),
		)
	}

	state.ID = types.Int64Value(int64(record.ID))
	state.Type = types.StringValue(record.Type)
	state.Name = types.StringValue(record.Name)
	state.Data = types.StringValue(record.Data)
	state.TTL = types.Int64Value(int64(record.TTL))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d dnsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (d dnsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
