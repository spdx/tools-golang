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

// ToJSON first processes the document by:
//   - setting each Element's CreationInfo property to the SpdxDocument's CreationInfo if nil
//   - collecting all element references to the top-level Elements slice
//
// ... and after this initial processing, outputs the document as compact JSON LD,
// including accounting for empty IDs by outputting blank node spdxId values
func (d *Document) ToJSON(writer io.Writer) error {
	// all Elements need to have creationInfo set...
	d.setCreationInfo(d.SpdxDocument.CreationInfo, &d.SpdxDocument)

	// The Elements list should not be serialized - the graph of the SpdxDocument includes all other properties, such as RootElements
	elements := d.Elements
	defer func() { d.Elements = elements }()
	d.Elements = nil

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

			var allElements []AnyElement
			for _, o := range graph {
				// collect all graph elements except SpdxDocument itself
				if el, ok := o.(AnyElement); ok && el != doc {
					allElements = append(allElements, el)
				}
			}
			d.Elements = allElements

			return nil
		}
	}
	return fmt.Errorf("no SPDX document found")
}

// all elements in the Document should be available in the SpdxDocument.Elements proeprty --
// on JSON LD deserialization, move all elements from @graph to the Elements property
func collectAllElements(d *SpdxDocument) map[reflect.Value]AnyElement {
	all := map[reflect.Value]AnyElement{}
	all[reflect.ValueOf(d)] = d
	_ = ld.VisitObjectGraph(d, func(path []any, value reflect.Value) error {
		if value.Kind() == reflect.Pointer {
			if _, ok := all[value]; !ok {
				if e, ok := value.Interface().(AnyElement); ok {
					all[value] = e
				}
			}
		}
		return nil
	})
	return all
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
