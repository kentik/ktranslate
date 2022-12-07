package enrich

import (
	"fmt"
	"regexp"

	"go.starlark.net/starlark"
)

// catch(f) evaluates f() and returns its evaluation error message
// if it failed or None if it succeeded.
func catch(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var fn starlark.Callable
	if err := starlark.UnpackArgs("catch", args, kwargs, "fn", &fn); err != nil {
		return nil, err
	}
	if _, err := starlark.Call(thread, fn, nil, nil); err != nil {
		return starlark.String(err.Error()), nil
	}
	return starlark.None, nil
}

// Keep a cache of seen regexps seen to speed things up. Note -- this isn't guarded with a lock, should it be?
var reCache = map[string]*regexp.Regexp{}

// findAllSubmatch(re, target) will compile re as a regexp and then run findAllSubmatch(target)
func findAllSubmatch(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var re starlark.String
	var target starlark.String
	if err := starlark.UnpackArgs("findAllSubmatch", args, kwargs, "re", &re, "target", &target); err != nil {
		return nil, err
	}

	// Use cache, assume there's lots of hits and only a few res so don't need a lock here to populate map.
	reg := re.GoString()
	r := reCache[reg]
	if r == nil {
		rn, err := regexp.Compile(reg)
		if err != nil {
			return starlark.None, err
		}
		reCache[reg] = rn
		r = rn
	}

	outputs := []starlark.Value{}
	for _, match := range r.FindAllSubmatch([]byte(target.GoString()), -1) {
		outputs = append(outputs, starlark.NewList([]starlark.Value{
			starlark.String(match[0]), starlark.String(match[1]),
		}))
	}

	return starlark.NewList(outputs), nil
}

type builtinMethod func(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

func builtinAttr(recv starlark.Value, name string, method builtinMethod) (starlark.Value, error) {
	if method == nil {
		return starlark.None, fmt.Errorf("no such method '%s'", name)
	}

	// Allocate a closure over 'method'.
	impl := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		return method(b, args, kwargs)
	}
	return starlark.NewBuiltin(name, impl).BindReceiver(recv), nil
}
