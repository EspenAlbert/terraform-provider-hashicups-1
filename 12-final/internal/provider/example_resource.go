// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ExampleResource{}
var _ resource.ResourceWithImportState = &ExampleResource{}

func NewExampleResource() resource.Resource {
	return &ExampleResource{}
}

// ExampleResource defines the resource implementation.
type ExampleResource struct {
	client *hashicups.Client
}

// TFModel describes the resource data model.
type TFModel struct {
	RootComputedOptional types.Object `tfsdk:"root_computed_optional"`
	Id                   types.String `tfsdk:"id"`
}

type TFModelRootComputedOptional struct {
	Computed types.String `tfsdk:"computed"`
	Optional types.String `tfsdk:"optional"`
}

var ModelRootComputedOptionalObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"computed": types.StringType,
	"optional": types.StringType,
}}

type APIBehaviorStruct struct {
	CurrentOperation string
	CreateObject     *TFModel
	CreateResponse   *TFModel
	UpdateObject     *TFModel
	UpdateResponse   *TFModel
	ReadObjects      []*TFModel
	ReadResponse     *TFModel
}

func (a *APIBehaviorStruct) AddReadObject(data TFModel) {
	a.ReadObjects = append(a.ReadObjects, &data)
}

func (a *APIBehaviorStruct) Reset() {
	a.CreateObject = nil
	a.CreateResponse = nil
	a.UpdateObject = nil
	a.UpdateResponse = nil
	a.ReadObjects = nil
	a.ReadResponse = nil
}

var APIBehavior = APIBehaviorStruct{}

func (r *ExampleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (r *ExampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"root_computed_optional": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"computed": schema.StringAttribute{
						Computed: true,
					},
					"optional": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ExampleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hashicups.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ExampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TFModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	APIBehavior.CreateObject = &data
	if APIBehavior.CreateResponse == nil {
		resp.Diagnostics.AddError("APIBehavior.CreateResponse is nil", "")
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &APIBehavior.CreateResponse)...)
}

func (r *ExampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TFModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	APIBehavior.AddReadObject(data)
	if resp.Diagnostics.HasError() {
		return
	}
	if APIBehavior.ReadResponse == nil {
		resp.Diagnostics.AddError("APIBehavior.ReadResponse is nil", "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, APIBehavior.ReadResponse)...)
}

func (r *ExampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	APIBehavior.UpdateObject = &data
	if APIBehavior.UpdateResponse == nil {
		resp.Diagnostics.AddError("APIBehavior.UpdateResponse is nil", "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, APIBehavior.UpdateResponse)...)
}

func (r *ExampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ExampleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
