package v3_0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/spdx/tools-golang/spdx/v3/internal"
	"github.com/spdx/tools-golang/spdx/v3/internal/ld"
)

/*
SPDX 3 models and serialization code is generated from internal/generate/main.go
To regenerate all models, run: make generate
*/

const Version = "3.0.1" // TODO is there a way to ascertain this version from generated code programmatically?

type Document struct {
	SpdxDocument
	LDContext ld.Context
}

func (d *Document) UnmarshalJSON(data []byte) error {
	if d.LDContext == nil {
		d.LDContext = context()
	}
	err := d.FromJSON(bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}

func (d *Document) MarshalJSON() ([]byte, error) {
	if d.LDContext == nil {
		d.LDContext = context()
	}
	buf := bytes.Buffer{}
	err := d.Write(&buf)
	return buf.Bytes(), err
}

func (d *Document) Write(w io.Writer) error {
	return d.ToJSON(w)
}

func NewDocument(conformance ProfileIdentifierType, documentName string, createdBy AnyAgent, createdUsing AnyTool) *Document {
	if createdBy == nil {
		createdBy = &SoftwareAgent{
			Comment: "Created with github.com/spdx/tools-golang",
			Name:    "tools-golang",
		}
	}
	ci := &CreationInfo{
		SpecVersion:  Version,
		Created:      time.Now(),
		CreatedBy:    notNil(AgentList{createdBy}),
		CreatedUsing: notNil(ToolList{createdUsing}),
	}
	id := ""
	name := documentName
	if internal.IsURI(name) {
		id = name
		name = ""
	}
	return &Document{
		SpdxDocument: SpdxDocument{
			ID:                  id,
			Name:                name,
			CreationInfo:        ci,
			ProfileConformances: conformanceFrom(conformance),
		},
		LDContext: context(),
	}
}

func conformanceFrom(conformance ProfileIdentifierType) []ProfileIdentifierType {
	out := []ProfileIdentifierType{ProfileIdentifierType_Core}
	switch conformance {
	case ProfileIdentifierType_Core:
	case ProfileIdentifierType_Software:
		out = append(out, conformance)
	case ProfileIdentifierType_Ai:
		out = append(out, ProfileIdentifierType_Software, conformance)
	case ProfileIdentifierType_Dataset:
		out = append(out, ProfileIdentifierType_Software, ProfileIdentifierType_Ai, conformance)
	}
	return out
}

func (d *Document) Validate(setCreationInfo bool) error {
	if setCreationInfo {
		// all Elements need to have creationInfo set...
		d.setCreationInfo(d.SpdxDocument.CreationInfo, &d.SpdxDocument)
	}
	return ld.ValidateGraph(d.SpdxDocument)
}

//func (d *Document) Append(e ...AnyElement) {
//	d.SpdxDocument.Elements = append(d.SpdxDocument.Elements, e...)
//	d.SpdxDocument.RootElements = append(d.SpdxDocument.RootElements, e...)
//}

// ToJSON first processes the document by:
//   - setting each Element's CreationInfo property to the SpdxDocument's CreationInfo if nil
//   - collecting all element references to the top-level Elements slice
//
// ... and after this initial processing, outputs the document as compact JSON LD,
// including accounting for empty IDs by outputting blank node spdxId values
func (d *Document) ToJSON(writer io.Writer) error {
	// all Elements need to have creationInfo set...
	d.setCreationInfo(d.SpdxDocument.CreationInfo, &d.SpdxDocument)

	// ensure the Elements are in the root Element list
	d.ensureAllDocumentElements()

	if d.LDContext == nil {
		d.LDContext = context()
	}

	// our default behavior is to ensure a URI for the document prefix, defaulting to a sub-URI of the document ID
	documentPrefix := d.ID
	if !internal.IsURI(documentPrefix) {
		name := d.ID
		if name == "" {
			name = d.Name
		}
		documentPrefix = internal.NewDocumentID(name)
		if d.ID == "" {
			d.ID = documentPrefix
		}
	}

	namespaceMap := map[string]string{}
	defaultPrefix := internal.DefaultSpdxNamespace
	for _, mapEntry := range d.NamespaceMaps {
		namespaceMap[string(mapEntry.GetNamespace())] = mapEntry.GetPrefix()
		if strings.HasPrefix(string(mapEntry.GetNamespace()), documentPrefix) {
			// if the user has provided an ID and included a reference in the namespace map, we will just use the namespace prefix
			documentPrefix = mapEntry.GetPrefix()
			break
		}
		// if we need to use the default prefix, avoid clashes with existing namespace map entries
		for strings.HasPrefix(mapEntry.GetPrefix(), defaultPrefix) {
			defaultPrefix += "2"
		}
	}

	// If the prefix is still a URI, we add a namespace map for the default prefix to refer to this document
	if internal.IsURI(documentPrefix) {
		// at this point we have a URI, need to ensure a separator character so expansion makes sense
		if !strings.HasSuffix(documentPrefix, "/") && !strings.HasSuffix(documentPrefix, internal.DefaultSpdxNamespaceSeparator) {
			documentPrefix += internal.DefaultSpdxNamespaceSeparator
		}
		ns := documentPrefix
		d.NamespaceMaps = append(d.NamespaceMaps, &NamespaceMap{
			Prefix:    defaultPrefix,
			Namespace: URI(ns),
		})
		namespaceMap[ns] = defaultPrefix
		documentPrefix = defaultPrefix // each element will have a unique URI based on the spdx document namespace
	}

	return internal.ToJSON("https://spdx.org/rdf/3.0.1/spdx-context.jsonld", d.LDContext, &d.SpdxDocument, internal.PrefixedIdGenerator(documentPrefix, namespaceMap), writer)
}

func (d *Document) setCreationInfo(creationInfo AnyCreationInfo, doc *SpdxDocument) {
	if creationInfo == nil {
		return
	}
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
	if d.LDContext == nil {
		d.LDContext = context()
	}
	graph, err := d.LDContext.FromJSON(reader)
	if err != nil {
		return err
	}
	for _, e := range graph {
		if doc, ok := e.(*SpdxDocument); ok {
			d.SpdxDocument = *doc
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
			panic(fmt.Sprintf("non-pointer type in elements: %#v", v))
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

var _ interface {
	json.Marshaler
	json.Unmarshaler
} = (*Document)(nil)

func notNil[T any, ListType ~[]T](values ListType) ListType {
	var out ListType
	for _, v := range values {
		if isNil(v) {
			continue
		}
		out = append(out, v)
	}
	return out
}
