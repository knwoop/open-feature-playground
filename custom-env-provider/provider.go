package main

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/open-feature/go-sdk/openfeature"
)

const (
	DefaultPrefix = "FT_"
	ReasonEnv     = "environment_variable"
)

type SimpleEnvProvider struct {
	prefix string
}

func NewProvider(opts ...ProviderOption) *SimpleEnvProvider {
	p := &SimpleEnvProvider{
		prefix: DefaultPrefix,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

type ProviderOption func(*SimpleEnvProvider)

func WithPrefix(prefix string) ProviderOption {
	return func(p *SimpleEnvProvider) {
		p.prefix = prefix
	}
}

func (p *SimpleEnvProvider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{
		Name: "simple-env-flag-evaluator",
	}
}

func (p *SimpleEnvProvider) Hooks() []openfeature.Hook {
	return []openfeature.Hook{}
}

func (p *SimpleEnvProvider) BooleanEvaluation(ctx context.Context, flagKey string, defaultValue bool, evalCtx openfeature.FlattenedContext) openfeature.BoolResolutionDetail {
	val := os.Getenv(p.prefix + strings.ToUpper(flagKey))
	if val == "" {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason: openfeature.DefaultReason,
			},
		}
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewParseErrorResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	return openfeature.BoolResolutionDetail{
		Value: boolVal,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: ReasonEnv,
		},
	}
}

func (p *SimpleEnvProvider) StringEvaluation(ctx context.Context, flagKey string, defaultValue string, evalCtx openfeature.FlattenedContext) openfeature.StringResolutionDetail {
	val := os.Getenv(p.prefix + strings.ToUpper(flagKey))
	if val == "" {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason: openfeature.DefaultReason,
			},
		}
	}

	return openfeature.StringResolutionDetail{
		Value: val,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: ReasonEnv,
		},
	}
}

func (p *SimpleEnvProvider) IntEvaluation(ctx context.Context, flagKey string, defaultValue int64, evalCtx openfeature.FlattenedContext) openfeature.IntResolutionDetail {
	val := os.Getenv(p.prefix + strings.ToUpper(flagKey))
	if val == "" {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason: openfeature.DefaultReason,
			},
		}
	}

	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewParseErrorResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	return openfeature.IntResolutionDetail{
		Value: intVal,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: ReasonEnv,
		},
	}
}

func (p *SimpleEnvProvider) FloatEvaluation(ctx context.Context, flagKey string, defaultValue float64, evalCtx openfeature.FlattenedContext) openfeature.FloatResolutionDetail {
	val := os.Getenv(p.prefix + strings.ToUpper(flagKey))
	if val == "" {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason: openfeature.DefaultReason,
			},
		}
	}

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewParseErrorResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	return openfeature.FloatResolutionDetail{
		Value: floatVal,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: ReasonEnv,
		},
	}
}

func (p *SimpleEnvProvider) ObjectEvaluation(ctx context.Context, flagKey string, defaultValue interface{}, evalCtx openfeature.FlattenedContext) openfeature.InterfaceResolutionDetail {
	// not support for object type
	return openfeature.InterfaceResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewGeneralResolutionError("Object type is not supported in simple env provider"),
			Reason:          openfeature.ErrorReason,
		},
	}
}
