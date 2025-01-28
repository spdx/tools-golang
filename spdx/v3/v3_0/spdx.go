package v3_0

import (
	"fmt"
	"math/rand"
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
	graph []any
	//ldc          ldContext
}

func NewDocument(conformance ProfileIdentifierType, name string, createdBy AnyAgent, createdUsing AnyTool) *Document {
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
				ProfileConformances: []ProfileIdentifierType{conformance},
			},
		},
		graph: []any{ci},
		//ldc:   ldGlobal,
	}
}

func (d *Document) Append(e ...AnyElement) {
	d.SpdxDocument.RootElements = append(d.SpdxDocument.RootElements, e...)
	d.SpdxDocument.Elements = append(d.SpdxDocument.Elements, e...)
}

//func (d *Document) ToJSON(writer io.Writer) error {
//	if d.SpdxDocument == nil {
//		return fmt.Errorf("no document object created")
//	}
//	// all IElement need to have creationInfo set...
//	d.setCreationInfo(d.SpdxDocument.CreationInfo, d.SpdxDocument)
//
//	// all IElement need to have spdxID...
//	if makeIdGenerator != nil {
//		idGen := makeIdGenerator(d.SpdxDocument)
//		if d.SpdxDocument.ID == "" {
//			d.SpdxDocument.ID = idGen(d.SpdxDocument)
//		}
//		d.ensureSpdxIDs(d.SpdxDocument, idGen)
//	}
//
//	maps, err := d.ldc.toMaps(d.SpdxDocument)
//	if err != nil {
//		return err
//	}
//	enc := json.NewEncoder(writer)
//	enc.SetEscapeHTML(false)
//	enc.SetIndent("", "  ")
//	return enc.Encode(maps)
//}

func (d *Document) setCreationInfo(creationInfo AnyCreationInfo, doc *SpdxDocument) {
	iCreationInfoType := reflect.TypeOf((*AnyCreationInfo)(nil)).Elem()
	ci := reflect.ValueOf(creationInfo)
	_ = visitObjectGraph(map[reflect.Value]struct{}{}, reflect.ValueOf(doc), func(v reflect.Value) error {
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

var makeIdGenerator = func(doc *SpdxDocument) idGenerator {
	nextID := map[reflect.Type]uint{}
	return func(e any) string {
		if _, ok := e.(*SpdxDocument); ok {
			return fmt.Sprintf("%v", rand.Uint64())
		}
		t := baseType(reflect.TypeOf(e))
		// should these be blank nodes?
		id := nextID[t] + 1
		nextID[t] = id
		return fmt.Sprintf("_:%v-%v-%v", doc.ID, t.Name(), id)
	}
}

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

//func (d *Document) FromJSON(reader io.Reader) error {
//	graph, err := d.ldc.FromJSON(reader)
//	if err != nil {
//		return err
//	}
//	d.graph = append(d.graph, graph)
//	for _, e := range graph {
//		if doc, ok := e.(*SpdxDocument); ok {
//			d.SpdxDocument = doc
//			return nil
//		}
//	}
//	return fmt.Errorf("no document found")
//}
