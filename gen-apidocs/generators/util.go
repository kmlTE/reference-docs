/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generators

import (
	"fmt"
	"github.com/kmlTE/reference-docs/gen-apidocs/generators/api"
	"strings"
	"log"
	"gopkg.in/yaml.v2"
)

func PrintInfo(config *api.Config) {
	definitions := config.Definitions

	hasOrphaned := false
	for name, d := range definitions.All {
		if !d.FoundInField && !d.FoundInOperation {
			if !strings.Contains(name, "meta.v1.APIVersions") && !strings.Contains(name, "meta.v1.Patch") {
				hasOrphaned = true
			}
		}
	}
	if hasOrphaned {
		fmt.Printf("----------------------------------\n")
		fmt.Printf("Orphaned Definitions:\n")
		for name, d := range definitions.All {
			if !d.FoundInField && !d.FoundInOperation {
				if !strings.Contains(name, "meta.v1.APIVersions") && !strings.Contains(name, "meta.v1.Patch") {
					fmt.Printf("[%s]\n", name)
				}
			}
		}
		if !*api.AllowErrors {
			panic("Orphaned definitions found.")
		}
	}

	missingFromToc := false
	for _, d := range definitions.All {
		if !d.InToc && len(d.OperationCategories) > 0 && !d.IsOldVersion && !d.IsInlined {
			missingFromToc = true
		}
	}

	if missingFromToc {
		fmt.Printf("----------------------------------\n")
		fmt.Printf("Definitions with Operations Missing from Toc (Excluding old version):\n")
		for name, d := range definitions.All {
			if !d.InToc && len(d.OperationCategories) > 0 && !d.IsOldVersion && !d.IsInlined {
				fmt.Printf("[%s]\n", name)
				for _, oc := range d.OperationCategories {
					for _, o := range oc.Operations {
						fmt.Printf("\t [%s]\n", o.ID)
					}
				}
			}
		}
	}

	//fmt.Printf("Old definitions:\n")
	//for name, d := range definitions.All {
	//	if !d.InToc && len(d.OperationCategories) > 0 && d.IsOldVersion && !d.IsInlined {
	//		fmt.Printf("[%s]\n", name)
	//		for _, oc := range d.OperationCategories {
	//			for _, o := range oc.Operations {
	//				fmt.Printf("\t [%s]\n", o.ID)
	//			}
	//		}
	//	}
	//}
}


func PrintToscaInfo(config *api.Config) {
	definitions := config.Definitions

	fmt.Printf("-+-+-+-+-+-+-+-+-+-+-+-+-+-+-\n")

	for _, kind := range config.IncludedObjects {
		PrintDefitionInfo(definitions.ByKind[kind][0])
	}
	
	fmt.Printf("-+-+-+-+-+-+-+-+-+-+-+-+-+-+-\n")
}

func PrintDefitionInfo(def *api.Definition) {
	fmt.Printf("===\n")
	fmt.Printf("[%s]\n", def.Name)
	PrintFieldsInfo(def.Fields, 0)
	fmt.Printf("RequiredFields: %+v\n", def.RequiredFields)
	fmt.Printf("===\n")
}

func PrintFieldsInfo(fields api.Fields, indent int) {
	tab := strings.Repeat("\t", indent)
	for _, field := range fields {
		fmt.Printf("%s---\n", tab)
		fmt.Printf("%sName: %s\n", tab, field.Name)
		fmt.Printf("%sType: %s\n", tab, field.Type)
		fmt.Printf("%sComplex: %t\n", tab, field.HasComplexType())
		fmt.Printf("%sRequired: %t\n", tab, field.Required)
		fmt.Printf("%s---\n", tab)
		if field.HasComplexType() {
			fmt.Printf("%sRetrieve complex type for %s:\n", tab, field.Type)
			if len(field.Definition.Fields) == 0 {
				fmt.Printf("%sComplex type %s has no fields. Data type derived from %s:\n", tab, field.Type, field.Definition.Type)
			} else {
				PrintFieldsInfo(field.Definition.Fields, indent + 1)
			}
		}
	}
}

func DumpToscaYAML(tosca *ToscaTypes) {

	fmt.Printf("-+-+-+-+-+-+-+-+-+-+-+-+-+-+-\n")

	t, err := yaml.Marshal(&tosca)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    fmt.Printf("--- tosca dump:\n%s\n\n", string(t))

    fmt.Printf("-+-+-+-+-+-+-+-+-+-+-+-+-+-+-\n")
}