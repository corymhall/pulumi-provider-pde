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

const Tool = "pulumi-gen-pde"

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
		Name:        "pde",
		Description: "",
		License:     "Apache-2.0",
		Keywords:    []string{},
		Language: map[string]schema.RawMessage{
			"go": rawMessage(map[string]interface{}{
				"generateResourceContainerTypes": true,
				"importBasePath":                 "github.com/corymhall/pulumi-provider-pde/sdk/go/pde",
				"liftSingleValueMethodReturns":   true,
			}),
		},
		// Functions: map[string]schema.FunctionSpec{
		// 	"pde:local:Profile/getFileName": {
		// 		Inputs: &schema.ObjectTypeSpec{
		// 			Properties: map[string]schema.PropertySpec{
		// 				"__self__": {
		// 					TypeSpec: schema.TypeSpec{Ref: "#/resources/pde:local:Profile"},
		// 				},
		// 			},
		// 			Required: []string{"__self__"},
		// 		},
		// 		Outputs: &schema.ObjectTypeSpec{
		// 			Properties: map[string]schema.PropertySpec{
		// 				"result": {
		// 					TypeSpec: schema.TypeSpec{Type: "string"},
		// 				},
		// 			},
		// 			Required: []string{"result"},
		// 		},
		// 	},
		// },
		Resources: map[string]schema.ResourceSpec{
			// "pde:local:Profile": {
			// 	ObjectTypeSpec: schema.ObjectTypeSpec{
			// 		Description: "",
			// 		Properties: map[string]schema.PropertySpec{
			// 			"fileName": {
			// 				TypeSpec:    schema.TypeSpec{Type: "string"},
			// 				Description: "",
			// 			},
			// 		},
			// 		Required: []string{"fileName"},
			// 	},
			// 	InputProperties: map[string]schema.PropertySpec{
			// 		"fileName": {
			// 			TypeSpec:    schema.TypeSpec{Type: "string"},
			// 			Description: "",
			// 		},
			// 	},
			// 	RequiredInputs: []string{"fileName"},
			// 	Methods: map[string]string{
			// 		"getFileName": "pde:local:Profile/getFileName",
			// 	},
			// 	IsComponent: true,
			// },
			"pde:installers:Shell": {
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Description: "",
					Properties: map[string]schema.PropertySpec{
						"interpreter": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"location": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"programName": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"binLocation": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"executable": {
							TypeSpec:    schema.TypeSpec{Type: "boolean"},
							Description: "",
						},
						"installCommands": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"uninstallCommands": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"updateCommands": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"versionCommand": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"downloadURL": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"environment": {
							TypeSpec: schema.TypeSpec{
								Type:                 "object",
								AdditionalProperties: &schema.TypeSpec{Type: "string"},
							},
						},
					},
					Required: []string{"installCommands", "programName", "downloadURL", "location"},
				},
				RequiredInputs: []string{"installCommands", "programName", "downloadURL"},
				InputProperties: map[string]schema.PropertySpec{
					"downloadURL": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"programName": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"binLocation": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"executable": {
						TypeSpec:    schema.TypeSpec{Type: "boolean"},
						Description: "",
					},
					"installCommands": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
					"uninstallCommands": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
					"updateCommands": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
				},
			},
			"pde:local:File": {
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Required:    []string{"path", "force", "content"},
					Description: "",
					Properties: map[string]schema.PropertySpec{
						"path": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"force": {
							TypeSpec:    schema.TypeSpec{Type: "boolean"},
							Description: "",
						},
						"content": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
					},
				},
				RequiredInputs: []string{"path", "content"},
				InputProperties: map[string]schema.PropertySpec{
					"path": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"force": {
						TypeSpec:    schema.TypeSpec{Type: "boolean"},
						Description: "",
					},
					"content": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
				},
			},
			"pde:local:Link": {
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Description: "",
					Properties: map[string]schema.PropertySpec{
						"source": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"target": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"overwrite": {
							TypeSpec:    schema.TypeSpec{Type: "boolean"},
							Description: "",
						},
						"retain": {
							TypeSpec:    schema.TypeSpec{Type: "boolean"},
							Description: "",
						},
						"recursive": {
							TypeSpec:    schema.TypeSpec{Type: "boolean"},
							Description: "",
						},
						"linked": {
							TypeSpec:    schema.TypeSpec{Type: "boolean"},
							Description: "",
						},
						"isDir": {
							TypeSpec:    schema.TypeSpec{Type: "boolean"},
							Description: "",
						},
						"targets": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
					},
					Required: []string{"source", "target", "linked", "isDir", "targets"},
				},
				InputProperties: map[string]schema.PropertySpec{
					"source": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"target": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"overwrite": {
						TypeSpec:    schema.TypeSpec{Type: "boolean"},
						Description: "",
					},
					"retain": {
						TypeSpec:    schema.TypeSpec{Type: "boolean"},
						Description: "",
					},
					"recursive": {
						TypeSpec:    schema.TypeSpec{Type: "boolean"},
						Description: "",
					},
				},
				RequiredInputs: []string{"source", "target"},
			},

			"pde:installers:GitHubRepo": {
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Description: "",
					Properties: map[string]schema.PropertySpec{
						"interpreter": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"locations": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"assetName": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"absFolderName": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"folderName": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"branch": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"repo": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"org": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"installCommands": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"uninstallCommands": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"updateCommands": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"version": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"environment": {
							TypeSpec: schema.TypeSpec{
								Type:                 "object",
								AdditionalProperties: &schema.TypeSpec{Type: "string"},
							},
						},
					},
					Required: []string{"org", "repo", "absFolderName"},
				},
				RequiredInputs: []string{"org", "repo"},
				InputProperties: map[string]schema.PropertySpec{
					"repo": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"org": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"branch": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"folderName": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"installCommands": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
					"uninstallCommands": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
					"updateCommands": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
				},
			},
			"pde:installers:GitHubRelease": {
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Description: "",
					Properties: map[string]schema.PropertySpec{
						"interpreter": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"locations": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"assetName": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"binFolder": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"repo": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"org": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"executable": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"installCommands": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"uninstallCommands": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"updateCommands": {
							TypeSpec: schema.TypeSpec{
								Type:  "array",
								Items: &schema.TypeSpec{Type: "string"},
							},
							Description: "",
						},
						"version": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"downloadURL": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "",
						},
						"environment": {
							TypeSpec: schema.TypeSpec{
								Type:                 "object",
								AdditionalProperties: &schema.TypeSpec{Type: "string"},
							},
						},
					},
					Required: []string{"org", "repo", "downloadURL"},
				},
				RequiredInputs: []string{"org", "repo"},
				InputProperties: map[string]schema.PropertySpec{
					"repo": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"org": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"assetName": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"binFolder": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"executable": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "",
					},
					"installCommands": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
					"uninstallCommands": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
					"updateCommands": {
						TypeSpec: schema.TypeSpec{
							Type:  "array",
							Items: &schema.TypeSpec{Type: "string"},
						},
						Description: "",
					},
				},
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
