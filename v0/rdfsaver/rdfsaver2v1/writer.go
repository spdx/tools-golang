// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"os"
	"strconv"

	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func Write(output *os.File, doc *rdf2v1.Document, sn *rdf2v1.Snippet) error {
	f := NewFormatter(output, "rdfxml-abbrev")
	_, snippeterr := f.Snippet(sn)
	_ = snippeterr
	_, docerr := f.Document(doc)
	if docerr != nil {
		return nil
	}
	f.Close()
	return nil
}

// Formatter struct to write the output
type Formatter struct {
	serializer *goraptor.Serializer
	nodeIds    map[string]int
	fileIds    map[string]goraptor.Term
}

type Pair struct {
	key, val string
}

// NewFormatter initialses a new Formatter Interface
func NewFormatter(output *os.File, format string) *Formatter {

	// a new goraptor.NewSerializer
	s := goraptor.NewSerializer(format)

	s.StartStream(output, rdf2v1.BaseUri)

	// goraptor.NamespaceHandler:
	// handler function to be called when the parser encounters a namespace.
	s.SetNamespace("rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns#")
	s.SetNamespace("spdx", "http://spdx.org/rdf/terms#")
	s.SetNamespace("rdfs", "http://www.w3.org/2000/01/rdf-schema#")
	s.SetNamespace("doap:", "http://usefulinc.com/ns/doap#")
	s.SetNamespace("j.0:", "http://www.w3.org/2009/pointers#")

	return &Formatter{
		serializer: s,
		nodeIds:    make(map[string]int),
		fileIds:    make(map[string]goraptor.Term),
	}
}

// NodeId method to set new node ID for a particular prefix
func (f *Formatter) NodeId(prefix string) *goraptor.Blank {

	f.nodeIds[prefix]++

	id := goraptor.Blank(prefix + strconv.Itoa(f.nodeIds[prefix]))
	return &id
}

// Sets node Type
func (f *Formatter) setNodeType(node, t goraptor.Term) error {
	return f.add(node, rdf2v1.Prefix("ns:Type"), t)
}

// Add 'keys' to 'values' for subject 'to'
func (f *Formatter) add(to, key, value goraptor.Term) error {
	return f.serializer.Add(&goraptor.Statement{
		Subject:   to,
		Predicate: key,
		Object:    value,
	})
}

func (f *Formatter) addTerm(to goraptor.Term, key string, value goraptor.Term) error {
	return f.add(to, rdf2v1.Prefix(key), value)
}

func (f *Formatter) addLiteral(to goraptor.Term, key, value string) error {
	if value == "" {
		return nil
	}
	return f.add(to, rdf2v1.Prefix(key), &goraptor.Literal{Value: value})
}

func (f *Formatter) addPairs(to goraptor.Term, Pairs ...Pair) error {
	for _, p := range Pairs {
		if err := f.addLiteral(to, p.key, p.val); err != nil {
			return err
		}
	}
	return nil
}

// Close to free the serializer
func (f *Formatter) Close() {
	f.serializer.EndStream()
	f.serializer.Free()
}
