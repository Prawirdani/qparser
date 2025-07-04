package qparser

import (
	"reflect"
	"strconv"
)

func setStr(value string) reflect.Value {
	return reflect.ValueOf(value)
}

func setBool(value string) reflect.Value {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return invalidReflect
	}
	return reflect.ValueOf(v)
}

func setInt(bitSize int) Setter {
	return func(value string) reflect.Value {
		v, err := strconv.ParseInt(value, 10, bitSize)
		if err != nil {
			return invalidReflect
		}
		return reflect.ValueOf(v)
	}
}

func setUint(bitSize int) Setter {
	return func(value string) reflect.Value {
		v, err := strconv.ParseUint(value, 10, bitSize)
		if err != nil {
			return invalidReflect
		}
		return reflect.ValueOf(v)
	}
}

func setFloat(bitSize int) Setter {
	return func(value string) reflect.Value {
		v, err := strconv.ParseFloat(value, bitSize)
		if err != nil {
			return invalidReflect
		}
		return reflect.ValueOf(v)
	}
}
