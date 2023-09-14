// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	gogen "github.com/pulumi/pulumi/pkg/v3/codegen/go"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"

	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

const Tool = "pulumi-gen-pdec"

type Language string

const (
	Go     Language = "go"
	Schema Language = "schema"
)

func main() {
	printUsage := func() {
		fmt.Printf("Usage: %s <language> <out-dir> [schema-file] [version]\n", os.Args[0])
	}

	args := os.Args[1:]
	if len(args) < 2 {
		printUsage()
		os.Exit(1)
	}

	language, outdir := Language(args[0]), args[1]

	var schemaFile string
	var version string
	if language != Schema {
		if len(args) < 4 {
			printUsage()
			os.Exit(1)
		}
		schemaFile, version = args[2], args[3]
	}

	switch language {
	case Go:
		genGo(readSchema(schemaFile, version), outdir)
	case Schema:
		pkgSpec := generateSchema()
		mustWritePulumiSchema(pkgSpec, outdir)
	default:
		panic(fmt.Sprintf("Unrecognized language %q", language))
	}
}
func rawMessage(v interface{}) schema.RawMessage {
	bytes, err := json.Marshal(v)
	contract.Assert(err == nil)
	return bytes
}

func generateSchema() schema.PackageSpec {
	return schema.PackageSpec{
		Name:        "pdec",
		Description: "",
		License:     "Apache-2.0",
		Keywords:    []string{},
		Language: map[string]schema.RawMessage{
			"go": rawMessage(map[string]interface{}{
				"generateResourceContainerTypes": true,
				"importBasePath":                 "github.com/corymhall/pulumi-provider-pde/sdk/go/pdec",
				"liftSingleValueMethodReturns":   true,
			}),
		},
		Functions: map[string]schema.FunctionSpec{
			"pdec:local:Profile/addLines": {
				Inputs: &schema.ObjectTypeSpec{
					Properties: map[string]schema.PropertySpec{
						"__self__": {
							TypeSpec: schema.TypeSpec{Ref: "#/resources/pdec:local:Profile"},
						},
						"lines": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
						},
					},
					Required: []string{"__self__", "lines"},
				},
				Outputs: &schema.ObjectTypeSpec{
					Properties: map[string]schema.PropertySpec{},
				},
			},
			"pdec:local:Profile/getFileName": {
				Inputs: &schema.ObjectTypeSpec{
					Properties: map[string]schema.PropertySpec{
						"__self__": {
							TypeSpec: schema.TypeSpec{Ref: "#/resources/pdec:local:Profile"},
						},
					},
					Required: []string{"__self__"},
				},
				Outputs: &schema.ObjectTypeSpec{
					Properties: map[string]schema.PropertySpec{
						"result": {
							TypeSpec: schema.TypeSpec{Type: "string"},
						},
					},
					Required: []string{"result"},
				},
			},
		},
		Resources: map[string]schema.ResourceSpec{
			"pdec:local:Profile": {
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Description: "",
					Properties: map[string]schema.PropertySpec{
						"fileName": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
					},
					Required: []string{"fileName"},
				},
				InputProperties: map[string]schema.PropertySpec{
					"fileName": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
				},
				RequiredInputs: []string{"fileName"},
				Methods: map[string]string{
					"getFileName": "pdec:local:Profile/getFileName",
					"addLines":    "pdec:local:Profile/addLines",
				},
				IsComponent: true,
			},
		},
	}
}

func readSchema(schemaPath string, version string) *schema.Package {
	// Read in, decode, and import the schema.
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		panic(err)
	}

	var pkgSpec schema.PackageSpec
	if err = json.Unmarshal(schemaBytes, &pkgSpec); err != nil {
		panic(err)
	}
	pkgSpec.Version = version

	pkg, err := schema.ImportSpec(pkgSpec, nil)
	if err != nil {
		panic(err)
	}
	return pkg
}

func genGo(pkg *schema.Package, outdir string) {
	files, err := gogen.GeneratePackage(Tool, pkg)
	if err != nil {
		panic(err)
	}
	mustWriteFiles(outdir, files)
}
func mustWriteFiles(rootDir string, files map[string][]byte) {
	for filename, contents := range files {
		mustWriteFile(rootDir, filename, contents)
	}
}

func mustWriteFile(rootDir, filename string, contents []byte) {
	outPath := filepath.Join(rootDir, filename)

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		panic(err)
	}
	err := ioutil.WriteFile(outPath, contents, 0600)
	if err != nil {
		panic(err)
	}
}

func mustWritePulumiSchema(pkgSpec schema.PackageSpec, outdir string) {
	schemaJSON, err := json.MarshalIndent(pkgSpec, "", "    ")
	if err != nil {
		panic(errors.Wrap(err, "marshaling Pulumi schema"))
	}
	mustWriteFile(outdir, "schema.json", schemaJSON)
}
