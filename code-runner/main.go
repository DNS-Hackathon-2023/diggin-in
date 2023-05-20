package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	slJSON "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"gopkg.in/yaml.v3"
	"io"
	"os/exec"
	"time"
)

// This is hack and very inefficient
func anyToStarlark(
	thread *starlark.Thread,
	v interface{},
) (starlark.Value, error) {
	// Convert to JSON
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	data := string(b)

	// Decode using starlark internal json decoder.
	decode := slJSON.Module.Members["decode"].(*starlark.Builtin)
	res, err := decode.CallInternal(thread, starlark.Tuple{starlark.String(data)}, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func execDigVerbose(thread *starlark.Thread, domain string) (starlark.Value, error) {
	dig := exec.Command("dig", domain)
	jc := exec.Command("jc", "dig")

	rx, tx := io.Pipe()
	dig.Stdout = tx
	jc.Stdin = rx

	var buf bytes.Buffer
	jc.Stdout = &buf

	dig.Start()
	jc.Start()
	dig.Wait()
	tx.Close()
	jc.Wait()

	// Decode using starlark internal json decoder.
	decode := slJSON.Module.Members["decode"].(*starlark.Builtin)
	res, err := decode.CallInternal(thread, starlark.Tuple{starlark.String(buf.String())}, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func execDig(domain string) (interface{}, error) {
	dig := exec.Command("dig", "+yaml", domain)
	res, err := dig.Output()
	if err != nil {
		return nil, err
	}
	// Decode yaml
	var data interface{}
	err = yaml.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func apiDig(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	// Parse arguments
	var (
		domain string
	)
	if err := starlark.UnpackArgs("dig", args, kwargs, "domain", &domain); err != nil {
		return nil, err
	}

	// Run dig
	/*
		result, err := execDig(domain)
		if err != nil {
			return nil, err
		}

		res, err := anyToStarlark(thread, result)
		if err != nil {
			return nil, err
		}
	*/
	res, err := execDigVerbose(thread, domain)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// State: Set thread local value
func apiStateSet(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	// Parse arguments
	var (
		key   string
		value starlark.Value
	)
	if err := starlark.UnpackArgs("set", args, kwargs, "key", &key, "value", &value); err != nil {
		return nil, err
	}

	// Set thread local value
	thread.SetLocal(key, value)

	return starlark.None, nil
}

// State: Get thread local value
func apiStateGet(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	// Parse arguments
	var (
		key          string
		defaultValue starlark.Value = starlark.None
	)
	if err := starlark.UnpackArgs("get", args, kwargs, "key", &key, "default?", &defaultValue); err != nil {
		return nil, err
	}

	// Get thread local value
	value := thread.Local(key)
	if value == nil {
		return defaultValue, nil
	}

	return value.(starlark.Value), nil
}

func main() {
	// Execute Starlark program in a file.
	thread := &starlark.Thread{Name: "main"}

	// Globals
	var modMeasure = &starlarkstruct.Module{
		Name: "measure",
		Members: starlark.StringDict{
			"dig": starlark.NewBuiltin("dig", apiDig),
		},
	}

	var modState = &starlarkstruct.Module{
		Name: "state",
		Members: starlark.StringDict{
			"get": starlark.NewBuiltin("get", apiStateGet),
			"set": starlark.NewBuiltin("set", apiStateSet),
		},
	}

	env := starlark.StringDict{
		"measure": modMeasure,
		"state":   modState,
	}

	globals, err := starlark.ExecFile(thread, "hej.star", nil, env)
	if err != nil {
		panic(err)
	}

	// Retrieve a module global.

	// Entry point
	main := globals["loop"]
	for {
		fmt.Println("executing loop")
		_, err := starlark.Call(thread, main, nil, nil)
		if err != nil {
			panic(err)
		}

		time.Sleep(10 * time.Second)
	}
}
