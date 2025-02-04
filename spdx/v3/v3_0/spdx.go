package v3_0

import (
	"reflect"
	"time"
)

/*
SPDX 3 models and serialization code is generated from some different prototype golang support for shacl2code, in: https://github.com/kzantow-anchore/shacl2code

To regenerate, use something like this command:
.venv/bin/python -m shacl2code generate -i https://spdx.org/rdf/3.0.1/spdx-model.ttl -i https://spdx.org/rdf/3.0.1/spdx-json-serialize-annotations.ttl -x https://spdx.org/rdf/3.0.1/spdx-context.jsonld golang
--package v3_0 --license MIT
--output $HOME/projects/tools-golang/spdx/v3/v3_0/model.go
--remap-props element=elements,extension=extensions,externalIdentifier=externalIdentifiers,externalRef=externalRefs,rootElement=rootElements
*/

type Document struct {
	*SpdxDocument
}

func NewDocument(conformance profileIdentifierType, name string, createdBy AnyAgent, createdUsing AnyTool) *Document {
	ci := &CreationInfo{
		Created:       time.Now(),
		CreatedBys:    AgentList{createdBy},
		CreatedUsings: ToolList{createdUsing},
	}
	return &Document{
		SpdxDocument: &SpdxDocument{
			ElementCollection: ElementCollection{
				Element: Element{
					Name:         name,
					CreationInfo: ci,
				},
				ProfileConformances: []profileIdentifierType{conformance},
			},
		},
		//LDContext: LDContext(),
	}
}

func (d *Document) Append(e ...AnyElement) {
	d.SpdxDocument.RootElements = append(d.SpdxDocument.RootElements, e...)
	d.SpdxDocument.Elements = append(d.SpdxDocument.Elements, e...)
}

func (d *Document) setCreationInfo(creationInfo AnyCreationInfo, doc *SpdxDocument) {
	iCreationInfoType := reflect.TypeOf((*AnyCreationInfo)(nil)).Elem()
	ci := reflect.ValueOf(creationInfo)
	_ = visitObjectGraph(map[reflect.Value]struct{}{}, reflect.ValueOf(doc), func(v reflect.Value) error {
		if v.IsZero() {
			return nil
		}
		t := v.Type()
		if t.Kind() == reflect.Interface && v.IsNil() && t.Implements(iCreationInfoType) {
			v.Set(ci)
		}
		return nil
	})
}

func (d *Document) ensureSpdxIDs(doc *SpdxDocument, idGen idGenerator) {
	iElementType := reflect.TypeOf((*AnyElement)(nil)).Elem()
	_ = visitObjectGraph(map[reflect.Value]struct{}{}, reflect.ValueOf(doc), func(v reflect.Value) error {
		if v.Type().Implements(iElementType) {
			el, ok := v.Interface().(AnyElement)
			if ok {
				e := el.asElement()
				if e != nil && e.ID == "" {
					e.ID = idGen(el)
				}
			}
		}
		return nil
	})
}

type idGenerator func(e any) string

func baseType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func visitObjectGraph(visited map[reflect.Value]struct{}, v reflect.Value, visitor func(reflect.Value) error) error {
	if _, ok := visited[v]; ok {
		return nil
	}
	visited[v] = struct{}{}
	if !v.IsValid() {
		return nil
	}
	err := visitor(v)
	if err != nil {
		return err
	}
	switch v.Kind() {
	case reflect.Interface:
		if !v.IsNil() {
			return visitObjectGraph(visited, v.Elem(), visitor)
		}
	case reflect.Pointer:
		if v.IsNil() {
			return nil
		}
		return visitObjectGraph(visited, v.Elem(), visitor)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			err = visitObjectGraph(visited, v.Field(i), visitor)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			err = visitObjectGraph(visited, v.Index(i), visitor)
			if err != nil {
				return err
			}
		}
	default:
	}
	return nil
}
