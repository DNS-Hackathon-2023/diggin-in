package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	sldataframe "github.com/qri-io/starlib/dataframe"
	slbase64 "github.com/qri-io/starlib/encoding/base64"
	slcsv "github.com/qri-io/starlib/encoding/csv"
	slyaml "github.com/qri-io/starlib/encoding/yaml"
	slhash "github.com/qri-io/starlib/hash"
	slre "github.com/qri-io/starlib/re"

	slJSON "go.starlark.net/lib/json"
	slmath "go.starlark.net/lib/math"
	sltime "go.starlark.net/lib/time"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/exec"
	"time"
)

// Cli Flags
var (
	scriptFile string
	waitSec    int
	outFile    string
	probeId    string
)

func init() {
	// Get current hostname
	hostname, _ := os.Hostname()

	flag.StringVar(&scriptFile, "script", "", "Script file to execute")
	flag.IntVar(&waitSec, "wait", 10, "Wait time between executions")
	flag.StringVar(&outFile, "out", "", "Output file, gets truncated on start")
	flag.StringVar(&probeId, "probeid", hostname, "Probe ID")
}

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

type Result struct {
	Timestamp time.Time       `json:"timestamp"`
	ProbeID   string          `json:"probe_id"`
	Result    json.RawMessage `json:"result"`
	Tag       string          `json:"tag"`
}

// Collect: write data to output file
func apiCollect(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	// Get file from thread
	file := thread.Local("outputFile").(*os.File)

	encodeArgs := starlark.Tuple{args}
	tag := ""
	if args.Len() == 1 {
		encodeArgs = starlark.Tuple{args.Index(0)}
	} else if args.Len() == 2 {
		encodeArgs = starlark.Tuple{args.Index(1)}
		tag = args.Index(0).(starlark.String).GoString()
	}

	// Encode args as json
	encode := slJSON.Module.Members["encode"].(*starlark.Builtin)
	encodeRes, err := encode.CallInternal(thread, encodeArgs, nil)
	if err != nil {
		return nil, err
	}

	data := encodeRes.(starlark.String).GoString()

	// Encode payload
	result := Result{
		ProbeID:   probeId,
		Timestamp: time.Now().UTC(),
		Result:    []byte(data),
		Tag:       tag,
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	if _, err := file.Write(encoded); err != nil {
		return nil, err
	}
	file.WriteString("\n")

	return starlark.None, nil
}

func execDigVerbose(thread *starlark.Thread, args string) (starlark.Value, error) {
	argv := strings.Split(args, " ")

	dig := exec.Command("dig", argv...)
	jc := exec.Command("jc", "--dig")

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

	state := thread.Local("state")
	if state == nil {
		state = make(map[string]interface{})
	}
	state.(map[string]interface{})[key] = value

	// Set thread local value
	thread.SetLocal("state", state)

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
	state := thread.Local("state")
	value := state.(map[string]interface{})[key]
	if value == nil {
		return defaultValue, nil
	}

	return value.(starlark.Value), nil
}

func loadModule(module string, loader func() (starlark.StringDict, error)) starlark.Value {
	dict, err := loader()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return dict[module]
}

func libraryModule(module string) starlark.Value {
	switch module {
	case "base64":
		return loadModule(module, slbase64.LoadModule)
	case "csv":
		return loadModule(module, slcsv.LoadModule)
	case "dataframe":
		return sldataframe.Module
	case "hash":
		return loadModule(module, slhash.LoadModule)
	case "math":
		return slmath.Module
	case "measure":
		return &starlarkstruct.Module{
			Name: module,
			Members: starlark.StringDict{
				"dig": starlark.NewBuiltin("dig", apiDig),
			},
		}
	case "re":
		return loadModule(module, slre.LoadModule)
	case "state":
		return &starlarkstruct.Module{
			Name: module,
			Members: starlark.StringDict{
				"get": starlark.NewBuiltin("get", apiStateGet),
				"set": starlark.NewBuiltin("set", apiStateSet),
			},
		}
	case "time":
		return sltime.Module
	case "yaml":
		return loadModule(module, slyaml.LoadModule)
	}
	log.Fatalf("invalid module %q", module)
	return nil
}

func main() {
	flag.Parse()
	if scriptFile == "" {
		fmt.Println("No script file specified")
		flag.Usage()
		return
	}

	if outFile == "" {
		fmt.Println("Output filename missing")
		flag.Usage()
		return
	}
	file, err := os.Create(outFile)
	if err != nil {
		fmt.Println("Failed to open output file:", err)
		return
	}
	defer file.Close()

	// Execute Starlark program in a file.
	thread := &starlark.Thread{Name: "main"}
	thread.SetLocal("outputFile", file)
	thread.SetLocal("state", map[string]interface{}{})

	// Globals
	env := starlark.StringDict{
		"collect": starlark.NewBuiltin("collect", apiCollect),
	}
	for _, mod := range []string{
		"re",
		"base64",
		"csv",
		"dataframe",
		"hash",
		"math",
		"measure",
		"re",
		"state",
		"time",
		"yaml",
	} {
		env[mod] = libraryModule(mod)
	}

	if scriptFile == "" {
		fmt.Println("No script file specified")
		flag.Usage()
		return
	}

	globals, err := starlark.ExecFile(thread, scriptFile, nil, env)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Retrieve a module global.

	// Entry point
	main := globals["loop"]
	for {
		fmt.Println("executing loop")
		_, err := starlark.Call(thread, main, nil, nil)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		time.Sleep(time.Duration(waitSec) * time.Second)
	}
}
