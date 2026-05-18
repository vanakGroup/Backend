package logging_help

import "github.com/Vanak/Backend/internal/domain/logging/logtypes"

func MapToZapParams(extra map[logtypes.ExtraKey]interface{}) []interface{} {
	params := make([]interface{}, 0)

	for k, v := range extra {
		params = append(params, string(k))
		params = append(params, v)
	}

	return params
}

func MapToZeroParams(extra map[logtypes.ExtraKey]interface{}) map[string]interface{} {
	params := make(map[string]interface{})

	for k, v := range extra {
		params[string(k)] = v
	}

	return params
}
