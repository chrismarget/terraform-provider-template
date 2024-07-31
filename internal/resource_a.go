package internal

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ json.Marshaler   = (*modelResourceA)(nil)
	_ json.Unmarshaler = (*modelResourceA)(nil)
)

type modelResourceA struct {
	Id         types.String `tfsdk:"id"`
	StringAttr types.String `tfsdk:"string_attr"`
}

func (m modelResourceA) MarshalJSON() ([]byte, error) {
	var raw struct {
		StringAttr string `json:"string_attr"`
	}

	raw.StringAttr = m.StringAttr.ValueString()

	return json.Marshal(raw)
}

func (m *modelResourceA) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		StringAttr string `json:"string_attr"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	m.StringAttr = types.StringValue(raw.StringAttr)

	return nil
}

type resourceA struct{}

func (o *resourceA) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_a"
}

func (o *resourceA) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"string_attr": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (o *resourceA) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan modelResourceA
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	f, err := os.CreateTemp(os.TempDir(), "terraform-provider-template-*.json")
	if err != nil {
		resp.Diagnostics.AddError("failed creating output file", err.Error())
		return
	}
	defer f.Close()

	plan.Id = types.StringValue(f.Name())

	err = json.NewEncoder(f).Encode(plan)
	if err != nil {
		resp.Diagnostics.AddError("failed writing file contents", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (o *resourceA) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state modelResourceA
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	f, err := os.Open(state.Id.ValueString())
	if err != nil {
		if os.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("failed to open file for reading", err.Error())
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&state)
	if err != nil {
		resp.Diagnostics.AddError("failed reading file", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (o *resourceA) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan modelResourceA
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	f, err := os.OpenFile(plan.Id.ValueString(), os.O_WRONLY, 0o600)
	if err != nil {
		resp.Diagnostics.AddError("failed to open file for writing", err.Error())
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(plan)
	if err != nil {
		resp.Diagnostics.AddError("failed writing file contents", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (o *resourceA) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state modelResourceA
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := os.Remove(state.Id.ValueString())
	if err != nil {
		if os.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("failed removing file", err.Error())
	}
}
