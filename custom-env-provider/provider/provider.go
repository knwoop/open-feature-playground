package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/open-feature/go-sdk/openfeature"
)

const (
	DefaultPrefix = "FT_"
	ReasonEnv     = "env"
	ReasonCtx     = "evaluation_context"
)

type SimpleEnvProvider struct {
	prefix string
}

func NewSimpleEnvProvider(opts ...ProviderOption) *SimpleEnvProvider {
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

	if ctxVal, ok := p.getFromContext(flagKey, evalCtx); ok {
		if boolVal, ok := ctxVal.(bool); ok {
			return openfeature.BoolResolutionDetail{
				Value: boolVal,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			}
		}

		// If value exists but type is wrong
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewTypeMismatchResolutionError(fmt.Sprintf("context value for %s is not a boolean", flagKey)),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

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
	if ctxVal, ok := p.getFromContext(flagKey, evalCtx); ok {
		if strVal, ok := ctxVal.(string); ok {
			return openfeature.StringResolutionDetail{
				Value: strVal,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			}
		}
		// If value exists but type is wrong
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewTypeMismatchResolutionError(fmt.Sprintf("context value for %s is not a string", flagKey)),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

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
	if ctxVal, ok := p.getFromContext(flagKey, evalCtx); ok {
		switch v := ctxVal.(type) {
		case int:
			return openfeature.IntResolutionDetail{
				Value: int64(v),
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			}
		case int64:
			return openfeature.IntResolutionDetail{
				Value: v,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			}
		case float64:
			return openfeature.IntResolutionDetail{
				Value: int64(v),
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			}
		default:
			return openfeature.IntResolutionDetail{
				Value: defaultValue,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewTypeMismatchResolutionError(fmt.Sprintf("context value for %s is not a number", flagKey)),
					Reason:          openfeature.ErrorReason,
				},
			}
		}
	}

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
	if ctxVal, ok := p.getFromContext(flagKey, evalCtx); ok {
		switch v := ctxVal.(type) {
		case float64:
			return openfeature.FloatResolutionDetail{
				Value: v,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			}
		case int:
			return openfeature.FloatResolutionDetail{
				Value: float64(v),
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			}
		case int64:
			return openfeature.FloatResolutionDetail{
				Value: float64(v),
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: ReasonCtx,
				},
			}
		default:
			return openfeature.FloatResolutionDetail{
				Value: defaultValue,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewTypeMismatchResolutionError(fmt.Sprintf("context value for %s is not a number", flagKey)),
					Reason:          openfeature.ErrorReason,
				},
			}
		}
	}

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

func (p *SimpleEnvProvider) getFromContext(flagKey string, evalCtx openfeature.FlattenedContext) (interface{}, bool) {
	if evalCtx == nil {
		return nil, false
	}
	// First try exact flag key
	if val, ok := evalCtx[flagKey]; ok {
		return val, true
	}

	// Then try with prefix
	prefixedKey := p.prefix + strings.ToUpper(flagKey)
	if val, ok := evalCtx[prefixedKey]; ok {
		return val, true
	}

	return nil, false
}
