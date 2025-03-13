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
	emptyResponse      = tfModelResp(string1, nil)
	responseWithValues = tfModelResp(string1, &TFModelRootComputedOptional{
		Computed: types.StringValue("computed"),
		Optional: types.StringValue("optional"),
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

func tfModelReq(id types.String, rootComputedOptional *TFModelRootComputedOptional) *TFModel {
	return tfModel(id, rootComputedOptional, true)
}

func tfModelResp(id types.String, rootComputedOptional *TFModelRootComputedOptional) *TFModel {
	return tfModel(id, rootComputedOptional, false)
}

func tfModel(id types.String, rootComputedOptional *TFModelRootComputedOptional, useUnknownForNull bool) *TFModel {
	computedOptional := asObjectValue(ctx, rootComputedOptional, ModelRootComputedOptionalObjectType.AttrTypes)
	if useUnknownForNull && computedOptional.IsNull() {
		computedOptional = types.ObjectUnknown(ModelRootComputedOptionalObjectType.AttrTypes)
	}
	return &TFModel{
		Id:                   id,
		RootComputedOptional: computedOptional,
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
			CreateObject: tfModelReq(types.StringUnknown(), nil),
		}),
	}))
}
func TestAccNestedComputedOptionalNonEmptyPlanWhenResponseIsSet(t *testing.T) {
	resource.Test(t, testCase(resource.TestStep{
		PreConfig: preConfig(func() {
			APIBehavior.CreateResponse = responseWithValues
			APIBehavior.ReadResponse = responseWithValues
		}),
		Config: exampleEmpty,
	}))
}
