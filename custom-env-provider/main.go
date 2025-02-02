package main

import (
	"context"
	"fmt"
	"os"

	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	provider := NewProvider()
	openfeature.SetProvider(provider)

	client := openfeature.NewClient("my-app")

	os.Setenv("FT_MY_FEATURE", "true")
	os.Setenv("FT_COUNT", "42")
	os.Setenv("FT_NAME", "test")

	ctx := context.Background()

	evalCtx := openfeature.NewEvaluationContext(
		"user-123",
		map[string]interface{}{},
	)

	boolValue, _ := client.BooleanValue(ctx, "my_feature", false, evalCtx)
	intValue, _ := client.IntValue(ctx, "count", 0, evalCtx)
	stringValue, _ := client.StringValue(ctx, "name", "", evalCtx)

	fmt.Println(boolValue, intValue, stringValue)

	evalCtx = openfeature.NewEvaluationContext(
		"user-123", // targetingKey
		map[string]interface{}{ // attributes
			"my_feature": true,
			"count":      1000,
		},
	)

	evaluatedBoolValue, _ := client.BooleanValue(ctx, "my_feature", false, evalCtx)
	evaluatedIntValue, _ := client.IntValue(ctx, "count", 0, evalCtx)
	evaluatedStringValue, _ := client.StringValue(ctx, "name", "", evalCtx)

	fmt.Println(evaluatedBoolValue, evaluatedIntValue, evaluatedStringValue)
}
