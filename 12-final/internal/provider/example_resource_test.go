package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	resourceType = "hashicups_example"
	resourceName = "test"
	resourceID   = "hashicups_example.test"
)

const exampleEmpty = `resource "hashicups_example" "test" {}`

var (
	ctx                = context.Background()
	string1            = types.StringValue("1")
	emptyResponse      = tfModelResp(tfModelDef{id: string1})
	responseWithValues = tfModelResp(tfModelDef{
		id: string1,
		rootComputedOptional: &TFModelRootComputedOptional{
			Computed: types.StringValue("computed"),
			Optional: types.StringValue("optional"),
		},
	})
	responseWithDefault = tfModelResp(tfModelDef{
		id: string1,
		rootComputedDefault: &TFModelRootComputedDefault{
			Computed: types.StringValue("computed"),
			Default:  types.StringValue(DefaultValue),
		},
	})
	responseWithNonDefault = tfModelResp(tfModelDef{
		id: string1,
		rootComputedDefault: &TFModelRootComputedDefault{
			Computed: types.StringValue("computed"),
			Default:  types.StringValue("non-default"),
		},
	})
)

func assertGlobalState(t *testing.T, expectedState APIBehaviorStruct) func(_ *terraform.State) error {
	return func(s *terraform.State) error {
		actual := APIBehavior
		if expectedState.CreateObject != nil {
			assert.Equal(t, expectedState.CreateObject, actual.CreateObject)
		}
		return nil
	}
}
func asObjectValue[T any](ctx context.Context, t T, attrs map[string]attr.Type) types.Object {
	objType, diagsLocal := types.ObjectValueFrom(ctx, attrs, t)
	if diagsLocal.HasError() {
		panic("failed to convert object to model")
	}
	return objType
}

type tfModelDef struct {
	id                   types.String
	rootComputedOptional *TFModelRootComputedOptional
	rootComputedDefault  *TFModelRootComputedDefault
}

func tfModelReq(model tfModelDef) *TFModel {
	return tfModel(model, true)
}

func tfModelResp(model tfModelDef) *TFModel {
	return tfModel(model, false)
}

func tfModel(model tfModelDef, useUnknownForNull bool) *TFModel {
	computedOptional := asObjectValue(ctx, model.rootComputedOptional, ModelRootComputedOptionalObjectType.AttrTypes)
	if useUnknownForNull && computedOptional.IsNull() {
		computedOptional = types.ObjectUnknown(ModelRootComputedOptionalObjectType.AttrTypes)
	}
	computedDefault := asObjectValue(ctx, model.rootComputedDefault, ModelRootComputedDefaultObjectType.AttrTypes)
	if useUnknownForNull && computedDefault.IsNull() {
		computedDefault = types.ObjectUnknown(ModelRootComputedDefaultObjectType.AttrTypes)
	}
	return &TFModel{
		Id:                   model.id,
		RootComputedOptional: computedOptional,
		RootComputedDefault:  computedDefault,
	}
}

func preConfig(call func()) func() {
	return func() {
		APIBehavior.Reset()
		call()
	}
}

func testCase(steps ...resource.TestStep) resource.TestCase {
	return resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	}
}

func TestAccNestedComputedOptionalOK(t *testing.T) {
	resource.Test(t, testCase(resource.TestStep{
		PreConfig: preConfig(func() {
			APIBehavior.CreateResponse = emptyResponse
			APIBehavior.ReadResponse = emptyResponse
		}),
		Config: exampleEmpty,
		Check: assertGlobalState(t, APIBehaviorStruct{
			CreateObject: tfModelReq(tfModelDef{id: types.StringUnknown()}),
		}),
	}))
}
func TestAccErrorNestedComputedOptionalNonEmptyPlanWhenResponseIsSet(t *testing.T) {
	resource.Test(t, testCase(resource.TestStep{
		PreConfig: preConfig(func() {
			APIBehavior.CreateResponse = responseWithValues
			APIBehavior.ReadResponse = responseWithValues
		}),
		Config: exampleEmpty,
	}))
}

func TestAccNestedComputedDefaultAPIReturnDefaultOK(t *testing.T) {
	resource.Test(t, testCase(resource.TestStep{
		PreConfig: preConfig(func() {
			APIBehavior.CreateResponse = responseWithDefault
			APIBehavior.ReadResponse = responseWithDefault
		}),
		Config: exampleEmpty,
		Check: assertGlobalState(t, APIBehaviorStruct{
			CreateObject: tfModelReq(tfModelDef{id: types.StringUnknown()}),
		}),
	}))
}

func TestAccErrorNestedComputedDefaultNonEmptyPlanWhenResponseIsNonDefault(t *testing.T) {
	resource.Test(t, testCase(resource.TestStep{
		PreConfig: preConfig(func() {
			APIBehavior.CreateResponse = responseWithNonDefault
			APIBehavior.ReadResponse = responseWithNonDefault
		}),
		Config: exampleEmpty,
	}))
}
