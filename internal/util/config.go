package util

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type LookupValid interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string |
		~bool
}

func LookupEnv[T LookupValid](env string, defaultValue T) T {
	value, ok := os.LookupEnv(env)
	if !ok {
		return defaultValue
	}

	var ret any
	switch t := any(defaultValue).(type) {
	case int, int8, int16, int32, int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Failed to lookup env %s as %v: %v", env, t, err))
		}
		ret = v
	case uint, uint8, uint16, uint32, uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Failed to lookup env %s as %v: %v", env, t, err))
		}
		ret = v
	case float32, float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(fmt.Sprintf("Failed to lookup env %s as %v: %v", env, t, err))
		}
		ret = v
	case string:
		ret = value
	case time.Duration:
		v, err := time.ParseDuration(value)
		if err != nil {
			panic(fmt.Sprintf("Failed to lookup env %s as %v: %v", env, t, err))
		}
		ret = v
	case bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			panic(fmt.Sprintf("Failed to lookup env %s as %v: %v", env, t, err))
		}
		ret = v
	}
	return ret.(T)
}
