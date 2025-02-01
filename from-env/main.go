package main

import (
	"context"
	"fmt"

	fromEnv "github.com/open-feature/go-sdk-contrib/providers/from-env/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	// register the provider against the go-sdk
	openfeature.SetProvider(&fromEnv.FromEnvProvider{})
	// create a client from via the go-sdk
	client := openfeature.NewClient("am-i-yellow-client")

	// we are now able to evaluate our stored flags
	resB, err := client.BooleanValueDetails(
		context.Background(),
		"AM_I_YELLOW",
		false,
		openfeature.NewEvaluationContext(
			"",
			map[string]interface{}{
				"color": "yellow",
			},
		),
	)
	fmt.Println(resB, err)

	resB, err = client.BooleanValueDetails(
		context.Background(),
		"AM_I_YELLOW",
		false,
		openfeature.NewEvaluationContext(
			"user",
			map[string]interface{}{
				"color": "yellow",
			},
		),
	)
	fmt.Println(resB, err)

	resS, err := client.StringValueDetails(
		context.Background(),
		"AM_I_YELLOW",
		"i am a default value",
		openfeature.NewEvaluationContext(
			"",
			map[string]interface{}{
				"color": "not yellow",
			},
		),
	)
	fmt.Println(resS, err)
}
