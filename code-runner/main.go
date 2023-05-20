package main

import (
	"fmt"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func apiDig(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	fmt.Println("running dig:", args, kwargs)

	// Parse arguments
	var (
		domain string
	)
	if err := starlark.UnpackArgs("dig", args, kwargs, "domain", &domain); err != nil {
		return nil, err
	}

	// Run dig

	// Return result

	return starlark.String("result..."), nil
}

func main() {
	// Execute Starlark program in a file.
	thread := &starlark.Thread{Name: "main"}

	// Globals
	env := starlark.StringDict{
		"measure": starlarkstruct.FromStringDict(
			starlark.String("measure"),
			starlark.StringDict{
				"dig": starlark.NewBuiltin("dig", apiDig),
			}),
	}

	globals, err := starlark.ExecFile(thread, "hej.star", nil, env)
	if err != nil {
		panic(err)
	}

	// Retrieve a module global.

	// Entry point
	main := globals["loop"]
	v, err := starlark.Call(thread, main, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("v", v)
}
