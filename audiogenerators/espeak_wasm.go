package audiogenerators

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/sillsdev/go-aeneas/datatypes"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// WASM... oh how I would love to like you. Unfortunately, in order to
// run properly, you need WASI, and there are 3 main WASI implementations in
// Go, wazero, wasmer-go, and gowasm. However, wasmer-go and gowasm use CGo,
// which could be problematic on Windows, and wazero doesn't actually have a
// full emscripten WASI implementation, meaning that one has to be brought to
// the table. The approach I took was to try to understand which of the WASI
// functions were being called, but it was getting into compiler intrinsics
// (e.g., __cxz_allocate_exception)

// Once WASM is fully used to load espeak-ng.wasm, a new fun task comes about
// in the form of reimplementing the BenLubar espeak library
// (https://github.com/BenLubar/espeak) with references to functions loaded
// from WASM instead of being linked at compile time.

// The code below uses wazero. This may or may not be the right choice, but it
// is an example implementation using wazero

// One problem with wazero: it uses a different definition of a function than
// emscripten compiled espeak with, so on top of providing custom emscripten
// functions, the only working functions provided through WASI might have to
// be reimplemented. Cannot tell if this is a fault of emscripten and the
// resulting WASM file, or if it's an error on the part of WASI implementation
// in wazero

//go:embed espeak-ng.wasm
var espeakWasm []byte

type EspeakWasmGenerator struct {
	ctx     context.Context
	runtime wazero.Runtime
	config  wazero.ModuleConfig
	module  api.Module
}

func (gen EspeakWasmGenerator) GenerateAudioFile(parameters *datatypes.Parameters, phrase string, outputPath string) error {
	log.Fatal("Not implemented!")
	return nil
}

func (gen EspeakWasmGenerator) Close() {
	gen.runtime.Close(gen.ctx)
	gen.module.Close(gen.ctx)
}

type EspeakWasmGeneratorFactory struct {
}

func (fac EspeakWasmGeneratorFactory) GetName() string {
	return "espeak-wasm"
}

func (fac EspeakWasmGeneratorFactory) GetAudioGenerator() (*datatypes.AudioGenerator, error) {
	ctx := context.Background()

	r := wazero.NewRuntime(ctx)

	config := wazero.
		NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithFS(os.DirFS("/"))

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// how to load custom functions in the env module
	// modBuilder := r.NewHostModuleBuilder("env")
	// _, err := modBuilder.NewFunctionBuilder().WithFunc(wasiexit).Export("exit").Instantiate(ctx)

	finalMod, err := r.InstantiateWithConfig(ctx, espeakWasm, config)
	if err != nil {
		fmt.Println("error: could not instantiate module ", err)
		return nil, err
	}

	// to get a function:
	// f := finalMod.ExportedFunction("espeak_ListVoices")

	var gen datatypes.AudioGenerator = EspeakWasmGenerator{ctx, r, config, finalMod}
	return &gen, nil
}
