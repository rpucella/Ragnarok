package primitives

import (
	"rpucella.net/ragnarok/internal/util"
	"rpucella.net/ragnarok/internal/value"
	"strings"
)

var PrimStringAppend = mkPrimitive(
	"string-append", 0, -1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		v := ""
		for _, arg := range args {
			if err := checkArgType(name, arg, IsString); err != nil {
				return nil, err
			}
			v += arg.GetString()
		}
		return value.NewString(v), nil
	})

var PrimStringLength = mkPrimitive("string-length", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		if err := checkArgType(name, args[0], IsString); err != nil {
			return nil, err
		}
		return value.NewInteger(len(args[0].GetString())), nil
	})

var PrimStringLower = mkPrimitive("string-lower", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		if err := checkArgType(name, args[0], IsString); err != nil {
			return nil, err
		}
		return value.NewString(strings.ToLower(args[0].GetString())), nil
	})

var PrimStringUpper = mkPrimitive("string-upper", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		if err := checkArgType(name, args[0], IsString); err != nil {
			return nil, err
		}
		return value.NewString(strings.ToUpper(args[0].GetString())), nil
	})

var PrimStringSubstring = mkPrimitive("string-substring", 1, 3,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		if err := checkArgType(name, args[0], IsString); err != nil {
			return nil, err
		}
		start := 0
		end := len(args[0].GetString())
		if len(args) > 2 {
			if err := checkArgType(name, args[2], isInt); err != nil {
				return nil, err
			}
			end = util.MinInt(args[2].GetInt(), end)
		}
		if len(args) > 1 {
			if err := checkArgType(name, args[1], isInt); err != nil {
				return nil, err
			}
			start = util.MaxInt(args[1].GetInt(), start)
		}
		// or perhaps raise an exception
		if end < start {
			return value.NewString(""), nil
		}
		return value.NewString(args[0].GetString()[start:end]), nil
	})

