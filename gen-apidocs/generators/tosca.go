/*
Copyright 2020 SODALITE EU Project.

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
	"github.com/kmlTE/reference-docs/gen-apidocs/generators/api"
	"strings"
)

func BuildToscaTypesFromDefinitions(config *api.Config, tosca *ToscaTypes) {
	definitions := config.Definitions
	for _, kind := range config.IncludedObjects {
		def := definitions.ByKind[kind][0]
		AddDefinitionToDataTypes(def, tosca)
		AddDefinitionToNodeTypes(def, tosca)
		PopulateToscaTypesFromComplexFields(def.Fields, tosca)
	}
}

func AddDefinitionToDataTypes(def *api.Definition, tosca *ToscaTypes) {
	if def.IsWrapper() { // for wrapper add only to tosca data types
		tosca.DataTypes[GetDataTypeName(def.Name)] = DataType{
			DerivedFrom: GetToscaTypeFromSpec(def.Type),
			Description: GetDescription(def.RawDescription),
		}
	} else {
		tosca.DataTypes[GetDataTypeName(def.Name)] = DataType{
			DerivedFrom: DataTypeBase,
			Description: GetDescription(def.RawDescription),
			Properties: GetDataTypeProperties(def.Fields),
		}
	}
}

func AddDefinitionToNodeTypes(def *api.Definition, tosca *ToscaTypes) {
	if !def.IsWrapper() { // do not include wrappers
		tosca.NodeTypes[GetNodeTypeName(def.Name)] = NodeType{
			DerivedFrom: NodeTypeBase,
			Description: GetDescription(def.RawDescription),
			Properties: GetNodeTypeProperties(GetDataTypeName(def.Name)),
		}
	}
}

func AddDefinitionToToscaTypes(def *api.Definition, tosca *ToscaTypes) {
	if def.IsWrapper() { // for wrapper add only to tosca data types
		tosca.DataTypes[GetDataTypeName(def.Name)] = DataType{
			DerivedFrom: GetToscaTypeFromSpec(def.Type),
			Description: GetDescription(def.RawDescription),
		}
	} else {
		tosca.DataTypes[GetDataTypeName(def.Name)] = DataType{
			DerivedFrom: DataTypeBase,
			Description: GetDescription(def.RawDescription),
			Properties: GetDataTypeProperties(def.Fields),
		}
		tosca.NodeTypes[GetNodeTypeName(def.Name)] = NodeType{
			DerivedFrom: NodeTypeBase,
			Description: GetDescription(def.RawDescription),
			Properties: GetNodeTypeProperties(GetDataTypeName(def.Name)),
		}
	}
}

func PopulateToscaTypesFromComplexFields(fields api.Fields, tosca *ToscaTypes) {
	for _, field := range fields {
		if !field.HasComplexType() || TypeExistsInTosca(GetDataTypeName(field.Name), tosca) {
			continue
		}
		AddDefinitionToDataTypes(field.Definition, tosca)
		PopulateToscaTypesFromComplexFields(field.Definition.Fields, tosca)
	}
}

func TypeExistsInTosca(dt_name string, tosca *ToscaTypes) bool {
	_, ok := tosca.DataTypes[dt_name]
	return ok
}

func GetDataTypeProperties(fields api.Fields) map[string]PropertyDefinition {
	properties := map[string]PropertyDefinition{}
	for _, field := range fields {
		var field_type string
		var entry_schema EntrySchemaDefinition

		if field.HasComplexType() {
			base_type := GetBaseType(field.Type)
			if field.Definition.IsWrapper() {
				field_type = GetDataTypeName(base_type)
			} else {
				field_type = GetEntrySchema(field.Type)
				entry_schema = EntrySchemaDefinition{
					Type: GetDataTypeName(base_type),
				}
			}
		} else {
			field_type = GetToscaTypeFromSpec(field.Type)
			if IsArray(field.Type) {
				field_type = ToscaArray
				entry_schema = EntrySchemaDefinition{
					Type: GetToscaTypeFromSpec(GetBaseType(field.Type)),
				}
			}
		}

		properties[field.Name] = PropertyDefinition{
			Type: field_type,
			Description: GetDescription(field.Description),
			Required: field.Required,
			EntrySchema: entry_schema,
		}
	}
	return properties
}

func GetNodeTypeProperties(dt_name string) map[string]PropertyDefinition {
	return map[string]PropertyDefinition{
		"definition": PropertyDefinition{
			Type: ToscaMap,
			Description: "Full definition can be found in " + dt_name,
			Required: true,
			EntrySchema: EntrySchemaDefinition{
				Type: dt_name,
			},
		},
	}
}

func IsArray(t string) bool {
	return strings.Contains(t, SpecArraySeparator)
}

func GetBaseType(t string) string {
	return strings.Split(t, SpecArraySeparator)[0]
}

func GetEntrySchema(t string) string {
	if IsArray(t) {
		return ToscaArray
	}
	return ToscaMap
}

func GetDescription(d string) string {
	n := 20
	if len(d) > n {
		return d[:n]
	}
	return d
}


/* TYPES */

// spec data type to tosca data type
const ToscaArray = "list"
const ToscaMap = "map"
const SpecArray = "array"
const SpecArraySeparator = " " + SpecArray
const SpecMap = "object"

func GetToscaTypeFromSpec(spec_type string) string {
	switch spec_type {
    case SpecArray:
        return ToscaArray
    case SpecMap:
        return ToscaMap
    default:
        return spec_type
    }
}

// data types
const DataTypeBase = "sodalite.datatypes.Kubernetes.Kind"

func GetDataTypeName(n string) string {
	return string(DataTypeBase + "." + n)
}

type DataType struct {
	DerivedFrom  string                             `yaml:"derived_from,omitempty"`
	Description  string                             `yaml:"description,omitempty"`
	Properties   map[string]PropertyDefinition      `yaml:"properties,omitempty"`
}

// node types
const NodeTypeBase = "sodalite.nodes.Kubernetes.Kind"

func GetNodeTypeName(n string) string {
	return string(NodeTypeBase + "." + n)
}

type NodeType struct {
	DerivedFrom  string                             `yaml:"derived_from,omitempty"`
	Description  string                             `yaml:"description,omitempty"`
	Properties   map[string]PropertyDefinition      `yaml:"properties,omitempty"`
	// Requirements []map[string]RequirementDefinition `yaml:"requirements,omitempty"`
	// Capabilities map[string]CapabilityDefinition    `yaml:"capabilities,omitempty"`
	// Attributes   map[string]AttributeDefinition     `yaml:"attributes,omitempty"`
	// Interfaces   map[string]InterfaceDefinition     `yaml:"interfaces,omitempty"`
}

type EntrySchemaDefinition struct {
	Type string `yaml:"type"`
}

type PropertyDefinition struct {
	Type        string             		`yaml:"type"`
	Description string             		`yaml:"description,omitempty"`
	Required    bool               		`yaml:"required"`
	Default     string             		`yaml:"default,omitempty"`
	EntrySchema EntrySchemaDefinition	`yaml:"entry_schema,omitempty"`
}

type ToscaTypes struct {
	Version   	string 				`yaml:"tosca_definitions_version,omitempty"`
	DataTypes 	map[string]DataType `yaml:"data_types,omitempty"`
	NodeTypes 	map[string]NodeType `yaml:"node_types,omitempty"` 
}