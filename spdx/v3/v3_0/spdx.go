package v3_0

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kzantow/go-ld"
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
	LDContext ld.Context
}

func LDContext() ld.Context {
	return context()
}

func NewDocument(conformance ProfileIdentifierType, name string, createdBy AnyAgent, createdUsing AnyTool) *Document {
	ci := &CreationInfo{
		SpecVersion:  "3.0.1", // TODO is there a way to ascertain this version from generated code programmatically?
		Created:      time.Now(),
		CreatedBy:    AgentList{createdBy},
		CreatedUsing: ToolList{createdUsing},
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
		LDContext: context(),
	}
}

func (d *Document) Validate(setCreationInfo bool) error {
	if setCreationInfo {
		// all Elements need to have creationInfo set...
		d.setCreationInfo(d.SpdxDocument.CreationInfo, d.SpdxDocument)
	}
	return ld.ValidateGraph(d.SpdxDocument)
}

func (d *Document) Append(e ...AnyElement) {
	d.SpdxDocument.RootElements = append(d.SpdxDocument.RootElements, e...)
}

// ToJSON first processes the document by:
//   - setting each Element's CreationInfo property to the SpdxDocument's CreationInfo if nil
//   - collecting all element references to the top-level Elements slice
//
// ... and after this initial processing, outputs the document as compact JSON LD,
// including accounting for empty IDs by outputting blank node spdxId values
func (d *Document) ToJSON(writer io.Writer) error {
	if d.SpdxDocument == nil {
		return fmt.Errorf("no document object created")
	}

	// all Elements need to have creationInfo set...
	d.setCreationInfo(d.SpdxDocument.CreationInfo, d.SpdxDocument)

	// ensure the Elements
	d.ensureAllDocumentElements()

	return d.LDContext.ToJSON(writer, d.SpdxDocument)
}

func (d *Document) setCreationInfo(creationInfo AnyCreationInfo, doc *SpdxDocument) {
	creationInfoInterfaceType := reflect.TypeOf((*AnyCreationInfo)(nil)).Elem()
	ci := reflect.ValueOf(creationInfo)
	_ = ld.VisitObjectGraph(doc, func(path []any, value reflect.Value) error {
		t := value.Type()
		if t == creationInfoInterfaceType && value.IsNil() {
			value.Set(ci)
		}
		return nil
	})
}

func (d *Document) FromJSON(reader io.Reader) error {
	graph, err := d.LDContext.FromJSON(reader)
	if err != nil {
		return err
	}
	for _, e := range graph {
		if doc, ok := e.(*SpdxDocument); ok {
			d.SpdxDocument = doc
			return nil
		}
	}
	return fmt.Errorf("no SPDX document found")
}

func (d *Document) ensureAllDocumentElements() {
	all := map[reflect.Value]struct{}{}
	for _, e := range d.Elements {
		v := reflect.ValueOf(e)
		if v.Kind() != reflect.Pointer {
			panic("non-pointer type in elements: %v" + spew.Sdump(v))
		}
		all[v] = struct{}{}
	}
	all[reflect.ValueOf(d.SpdxDocument)] = struct{}{}
	_ = ld.VisitObjectGraph(d.SpdxDocument, func(path []any, value reflect.Value) error {
		if value.Kind() == reflect.Pointer {
			if _, ok := all[value]; ok {
				return nil
			}
			if e, ok := value.Interface().(AnyElement); ok {
				all[value] = struct{}{}
				d.Elements = append(d.Elements, e)
			}
		}
		return nil
	})
}
