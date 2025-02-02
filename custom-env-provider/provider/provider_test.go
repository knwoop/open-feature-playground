package provider

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/open-feature/go-sdk/openfeature"
)

func TestNewSimpleEnvProvider(t *testing.T) {
	for name, test := range map[string]struct {
		prefix string
		want   string
	}{
		"default prefix": {
			prefix: "",
			want:   DefaultPrefix,
		},
		"custom prefix": {
			prefix: "TEST_",
			want:   "TEST_",
		},
	} {
		t.Run(name, func(t *testing.T) {
			provider := NewSimpleEnvProvider()
			if test.prefix != "" {
				provider = NewSimpleEnvProvider(WithPrefix(test.prefix))
			}

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(openfeature.ResolutionError{}),
			}
			if diff := cmp.Diff(test.want, provider.prefix, opts...); diff != "" {
				t.Errorf("prefix mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBooleanEvaluation(t *testing.T) {
	for name, test := range map[string]struct {
		envValue     string
		flagKey      string
		defaultValue bool
		evalCtx      openfeature.FlattenedContext
		want         openfeature.BoolResolutionDetail
	}{
		"environment value true": {
			envValue:     "true",
			flagKey:      "test_flag",
			defaultValue: false,
			want: openfeature.BoolResolutionDetail{
				Value: true,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonEnv,
				},
			},
		},
		"invalid environment value": {
			envValue:     "invalid",
			flagKey:      "test_flag",
			defaultValue: false,
			want: openfeature.BoolResolutionDetail{
				Value: false,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewParseErrorResolutionError("strconv.ParseBool: parsing \"invalid\": invalid syntax"),
					Reason:          openfeature.ErrorReason,
				},
			},
		},
		"context value": {
			flagKey:      "test_flag",
			defaultValue: false,
			evalCtx: openfeature.FlattenedContext{
				"test_flag": true,
			},
			want: openfeature.BoolResolutionDetail{
				Value: true,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			},
		},
		"invalid context value type": {
			flagKey:      "test_flag",
			defaultValue: false,
			evalCtx: openfeature.FlattenedContext{
				"test_flag": "not a bool",
			},
			want: openfeature.BoolResolutionDetail{
				Value: false,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewTypeMismatchResolutionError("context value for test_flag is not a boolean"),
					Reason:          openfeature.ErrorReason,
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			if test.envValue != "" {
				os.Setenv(DefaultPrefix+test.flagKey, test.envValue)
				defer os.Unsetenv(DefaultPrefix + test.flagKey)
			}

			provider := NewSimpleEnvProvider()
			result := provider.BooleanEvaluation(context.Background(), test.flagKey, test.defaultValue, test.evalCtx)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(openfeature.ResolutionError{}),
			}
			if diff := cmp.Diff(test.want, result, opts...); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStringEvaluation(t *testing.T) {
	for name, test := range map[string]struct {
		envValue     string
		flagKey      string
		defaultValue string
		evalCtx      openfeature.FlattenedContext
		want         openfeature.StringResolutionDetail
	}{
		"environment value": {
			envValue:     "test-value",
			flagKey:      "test_flag",
			defaultValue: "default",
			want: openfeature.StringResolutionDetail{
				Value: "test-value",
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonEnv,
				},
			},
		},
		"context value": {
			flagKey:      "test_flag",
			defaultValue: "default",
			evalCtx: openfeature.FlattenedContext{
				"test_flag": "context-value",
			},
			want: openfeature.StringResolutionDetail{
				Value: "context-value",
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			},
		},
		"invalid context value type": {
			flagKey:      "test_flag",
			defaultValue: "default",
			evalCtx: openfeature.FlattenedContext{
				"test_flag": 123,
			},
			want: openfeature.StringResolutionDetail{
				Value: "default",
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewTypeMismatchResolutionError("context value for test_flag is not a string"),
					Reason:          openfeature.ErrorReason,
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			if test.envValue != "" {
				os.Setenv(DefaultPrefix+test.flagKey, test.envValue)
				defer os.Unsetenv(DefaultPrefix + test.flagKey)
			}

			provider := NewSimpleEnvProvider()
			result := provider.StringEvaluation(context.Background(), test.flagKey, test.defaultValue, test.evalCtx)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(openfeature.ResolutionError{}),
			}
			if diff := cmp.Diff(test.want, result, opts...); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIntEvaluation(t *testing.T) {
	for name, test := range map[string]struct {
		envValue     string
		flagKey      string
		defaultValue int64
		evalCtx      openfeature.FlattenedContext
		want         openfeature.IntResolutionDetail
	}{
		"environment value": {
			envValue:     "123",
			flagKey:      "test_flag",
			defaultValue: 0,
			want: openfeature.IntResolutionDetail{
				Value: 123,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonEnv,
				},
			},
		},
		"context value int": {
			flagKey:      "test_flag",
			defaultValue: 0,
			evalCtx: openfeature.FlattenedContext{
				"test_flag": 123,
			},
			want: openfeature.IntResolutionDetail{
				Value: 123,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			},
		},
		"context value float64": {
			flagKey:      "test_flag",
			defaultValue: 0,
			evalCtx: openfeature.FlattenedContext{
				"test_flag": float64(123),
			},
			want: openfeature.IntResolutionDetail{
				Value: 123,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			},
		},
		"invalid environment value": {
			envValue:     "invalid",
			flagKey:      "test_flag",
			defaultValue: 0,
			want: openfeature.IntResolutionDetail{
				Value: 0,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewParseErrorResolutionError("strconv.ParseInt: parsing \"invalid\": invalid syntax"),
					Reason:          openfeature.ErrorReason,
				},
			},
		},
		"invalid context value type": {
			flagKey:      "test_flag",
			defaultValue: 0,
			evalCtx: openfeature.FlattenedContext{
				"test_flag": "not a number",
			},
			want: openfeature.IntResolutionDetail{
				Value: 0,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewTypeMismatchResolutionError("context value for test_flag is not a number"),
					Reason:          openfeature.ErrorReason,
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			if test.envValue != "" {
				os.Setenv(DefaultPrefix+test.flagKey, test.envValue)
				defer os.Unsetenv(DefaultPrefix + test.flagKey)
			}

			provider := NewSimpleEnvProvider()
			result := provider.IntEvaluation(context.Background(), test.flagKey, test.defaultValue, test.evalCtx)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(openfeature.ResolutionError{}),
			}
			if diff := cmp.Diff(test.want, result, opts...); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFloatEvaluation(t *testing.T) {
	for name, test := range map[string]struct {
		envValue     string
		flagKey      string
		defaultValue float64
		evalCtx      openfeature.FlattenedContext
		want         openfeature.FloatResolutionDetail
	}{
		"environment value": {
			envValue:     "123.45",
			flagKey:      "test_flag",
			defaultValue: 0,
			want: openfeature.FloatResolutionDetail{
				Value: 123.45,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonEnv,
				},
			},
		},
		"context value float64": {
			flagKey:      "test_flag",
			defaultValue: 0,
			evalCtx: openfeature.FlattenedContext{
				"test_flag": 123.45,
			},
			want: openfeature.FloatResolutionDetail{
				Value: 123.45,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			},
		},
		"context value int": {
			flagKey:      "test_flag",
			defaultValue: 0,
			evalCtx: openfeature.FlattenedContext{
				"test_flag": 123,
			},
			want: openfeature.FloatResolutionDetail{
				Value: 123,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			},
		},
		"invalid environment value": {
			envValue:     "invalid",
			flagKey:      "test_flag",
			defaultValue: 0,
			want: openfeature.FloatResolutionDetail{
				Value: 0,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewParseErrorResolutionError("strconv.ParseFloat: parsing \"invalid\": invalid syntax"),
					Reason:          openfeature.ErrorReason,
				},
			},
		},
		"invalid context value type": {
			flagKey:      "test_flag",
			defaultValue: 0,
			evalCtx: openfeature.FlattenedContext{
				"test_flag": "not a number",
			},
			want: openfeature.FloatResolutionDetail{
				Value: 0,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewTypeMismatchResolutionError("context value for test_flag is not a number"),
					Reason:          openfeature.ErrorReason,
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			if test.envValue != "" {
				os.Setenv(DefaultPrefix+test.flagKey, test.envValue)
				defer os.Unsetenv(DefaultPrefix + test.flagKey)
			}

			provider := NewSimpleEnvProvider()
			result := provider.FloatEvaluation(context.Background(), test.flagKey, test.defaultValue, test.evalCtx)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(openfeature.ResolutionError{}),
			}
			if diff := cmp.Diff(test.want, result, opts...); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestObjectEvaluation(t *testing.T) {
	for name, test := range map[string]struct {
		flagKey      string
		defaultValue interface{}
		evalCtx      openfeature.FlattenedContext
		want         openfeature.InterfaceResolutionDetail
	}{
		"object evaluation not supported": {
			flagKey:      "test_flag",
			defaultValue: map[string]interface{}{"key": "value"},
			want: openfeature.InterfaceResolutionDetail{
				Value: map[string]interface{}{"key": "value"},
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewGeneralResolutionError("Object type is not supported in simple env provider"),
					Reason:          openfeature.ErrorReason,
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			provider := NewSimpleEnvProvider()
			result := provider.ObjectEvaluation(context.Background(), test.flagKey, test.defaultValue, test.evalCtx)

			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(openfeature.ResolutionError{}),
			}
			if diff := cmp.Diff(test.want, result, opts...); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
