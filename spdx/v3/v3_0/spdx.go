package v3_0

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"time"
)

/*
SPDX 3 models and serialization code is generated from prototype golang support for shacl2code,
https://github.com/kzantow-anchore/shacl2code/tree/golang-bindings (with contributions from Nisha and Keith)

To regenerate, use something like this command:
.venv/bin/python -m shacl2code generate -i https://spdx.org/rdf/3.0.0/spdx-model.ttl -i https://spdx.org/rdf/3.0.0/spdx-json-serialize-annotations.ttl -x https://spdx.org/rdf/3.0.0/spdx-context.jsonld golang --package v3_0 --license MIT --output $HOME/projects/tools-golang/spdx/v3/v3_0/model.go --remap-props element=elements,extension=extensions,externalIdentifier=externalIdentifiers,externalRef=externalRefs,rootElement=rootElements
*/

type Document struct {
	creationInfo *CreationInfo
	document     *SpdxDocument
	graph        []any
	ldc          ldContext
}

func NewDocument(creator IAgent) *Document {
	ci := &CreationInfo{
		Created: time.Now().Format(time.RFC3339),
		CreatedBy: []IAgent{
			creator,
		},
	}
	creator.SetCreationInfo(ci)
	return &Document{
		creationInfo: ci,
		document: &SpdxDocument{
			CreationInfo: ci,
		},
		graph: []any{ci, creator},
		ldc:   ldGlobal,
	}
}

func (d *Document) CreationInfo() ICreationInfo {
	return d.creationInfo
}

func (d *Document) AddElement(e ...IElement) {
	d.document.RootElements = append(d.document.RootElements, e...)
	d.document.Elements = append(d.document.Elements, e...)
}

func (d *Document) GetElements() []IElement {
	return d.document.RootElements
}

func (d *Document) Document() *SpdxDocument {
	return d.document
}

func (d *Document) Packages() []IPackage {
	return get[IPackage](d)
}

func (d *Document) Relationships() Relationships {
	return Relationships{get[IRelationship](d)}
}

func (d *Document) Files() []IFile {
	return get[IFile](d)
}

func (d *Document) ToJSON(writer io.Writer) error {
	if d.document == nil {
		return fmt.Errorf("no document object created")
	}
	// all IElement need to have creationInfo set...
	if d.creationInfo != nil {
		d.setCreationInfo(d.creationInfo, d.document)
	}
	// all IElement need to have spdxID...
	if makeIdGenerator != nil {
		idGen := makeIdGenerator(d.document)
		if d.document.GetSpdxId() == "" {
			d.document.SetSpdxId(idGen(d.document))
		}
		d.ensureSpdxIDs(d.document, idGen)
	}

	maps, err := d.ldc.toMaps(d.document)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(maps)
}

func (d *Document) setCreationInfo(creationInfo ICreationInfo, doc ISpdxDocument) {
	iCreationInfoType := reflect.TypeOf((*ICreationInfo)(nil)).Elem()
	ci := reflect.ValueOf(creationInfo)
	_ = visitObjectGraph(map[reflect.Value]struct{}{}, reflect.ValueOf(doc), func(v reflect.Value) error {
		t := v.Type()
		if t.Kind() == reflect.Interface && v.IsNil() && t.Implements(iCreationInfoType) {
			v.Set(ci)
		}
		return nil
	})
}

func (d *Document) ensureSpdxIDs(doc ISpdxDocument, idGen idGenerator) {
	iElementType := reflect.TypeOf((*IElement)(nil)).Elem()
	_ = visitObjectGraph(map[reflect.Value]struct{}{}, reflect.ValueOf(doc), func(v reflect.Value) error {
		if v.Type().Implements(iElementType) {
			el, ok := v.Interface().(IElement)
			if ok && el.GetSpdxId() == "" {
				el.SetSpdxId(idGen(el))
			}
		}
		return nil
	})
}

type idGenerator func(e IElement) string

var makeIdGenerator = func(doc ISpdxDocument) idGenerator {
	nextID := map[reflect.Type]uint{}
	return func(e IElement) string {
		if _, ok := e.(ISpdxDocument); ok {
			return fmt.Sprintf("%v", rand.Uint64())
		}
		t := baseType(reflect.TypeOf(e))
		// should these be blank nodes?
		id := nextID[t] + 1
		nextID[t] = id
		return fmt.Sprintf("_:%v-%v-%v", doc.GetSpdxId(), t.Name(), id)
	}
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

func (d *Document) FromJSON(reader io.Reader) error {
	graph, err := d.ldc.FromJSON(reader)
	if err != nil {
		return err
	}
	d.graph = append(d.graph, graph)
	for _, e := range graph {
		if doc, ok := e.(*SpdxDocument); ok {
			d.document = doc
			return nil
		}
	}
	return fmt.Errorf("no document found")
}

func get[T any](ctx *Document) []T {
	var out []T
	for _, i := range ctx.graph {
		if i, ok := i.(T); ok {
			out = append(out, i)
		}
	}
	return out
}
