package primitives

import (
	"fmt"
	"io/ioutil"
	"rpucella.net/ragnarok/internal/ragnarok"
	"rpucella.net/ragnarok/internal/value"
	"sort"
	"time"
)

var PrimQuit = mkPrimitive(
	"quit", 0, 0,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		context, ok := ctxt.(*ragnarok.Context)
		if !ok {
			return nil, fmt.Errorf("Problem understanding context")
		}
		context.Bail()
		return value.NewNil(), nil
	})

var PrimEnv = mkPrimitive(
	"env", 0, 0,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		context, ok := ctxt.(*ragnarok.Context)
		if !ok {
			return nil, fmt.Errorf("Problem understanding context")
		}
		bindings := context.CurrentEnv.Bindings()
		maxWidth := 0
		keys := make([]string, len(bindings))
		i := 0
		for name := range bindings {
			if len(name) > maxWidth {
				maxWidth = len(name)
			}
			keys[i] = name
			i += 1
		}
		sort.Strings(keys)
		for _, name := range keys {
			context.Report(fmt.Sprintf("%*s %s", -maxWidth-2, name, bindings[name].Str()))
		}
		return value.NewNil(), nil
	})

var PrimGo = mkPrimitive(
	"go", 0, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		context, ok := ctxt.(*ragnarok.Context)
		if !ok {
			return nil, fmt.Errorf("Problem understanding context")
		}
		if len(args) == 0 {
			context.NextCurrentModule = context.HomeModule
			return value.NewNil(), nil
		}
		if err := checkArgType(name, args[0], IsSymbol); err != nil {
			return nil, err
		}
		context.NextCurrentModule = args[0].GetString()
		return value.NewNil(), nil
	})

var PrimModules = mkPrimitive(
	"modules", 0, 0,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		context, ok := ctxt.(*ragnarok.Context)
		if !ok {
			return nil, fmt.Errorf("Problem understanding context")
		}
		var result value.Value = value.NewEmpty()
		for m := range context.Ecosystem.Modules() {
			result = value.NewCons(value.NewSymbol(m), result)
		}
		return result, nil
	})

var PrimHelp = mkPrimitive(
	"help", 0, 0,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		context, ok := ctxt.(*ragnarok.Context)
		if !ok {
			return nil, fmt.Errorf("Problem understanding context")
		}
		context.Report("Some help about the system")
		context.Report("")
		context.Report("      (quit)   bail out")
		context.Report("   (modules)   see available modules")
		context.Report("  (go 'buff)   navigate to a particular buffer")
		context.Report("")
		return value.NewNil(), nil
	})

var PrimPrint = mkPrimitive(
	"print", 0, -1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		for _, arg := range args {
			fmt.Print(arg.Display(), " ")
		}
		fmt.Println()
		return value.NewNil(), nil
	})

var PrimLoad = mkPrimitive(
	"load", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		context, ok := ctxt.(*ragnarok.Context)
		if !ok {
			return nil, fmt.Errorf("Problem understanding context")
		}
		if err := checkArgType(name, args[0], IsString); err != nil {
			return nil, err
		}
		filename := args[0].GetString()
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		str := string(b)
		if err := context.ReadAll(str, context); err != nil {
			return nil, err
		}
		return value.NewNil(), nil
	})

var PrimTimedApply = mkPrimitive(
	"timed-apply", 2, 2,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		context, ok := ctxt.(*ragnarok.Context)
		if !ok {
			return nil, fmt.Errorf("Problem understanding context")
		}
		timeTrack := func(start time.Time) {
			elapsed := time.Since(start)
			context.Report(fmt.Sprintf("Time: %s", elapsed))
		}
		defer timeTrack(time.Now())
		if err := checkArgType(name, args[0], IsFunction); err != nil {
			return nil, err
		}
		if err := checkArgType(name, args[1], isList); err != nil {
			return nil, err
		}
		arguments := make([]value.Value, listLength(args[1]))
		current := args[1]
		for i := range arguments {
			arguments[i] = current.GetHead()
			current = current.GetTail()
		}
		if !value.IsEmpty(current) {
			return nil, fmt.Errorf("%s: malformed list", name)
		}
		return args[0].Apply(arguments, ctxt)
	})
