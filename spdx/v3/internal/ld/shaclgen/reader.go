package shaclgen

import (
	"bytes"
	"slices"
	"strconv"
	"strings"

	"github.com/deiu/rdf2go"
)

func ParseSHACLFromUrl(url string) (map[string]*Class, []*Individual) {
	ttl := fetch(url)
	return ParseSHACL(url, ttl)
}

func ParseSHACL(url string, ttl []byte) (map[string]*Class, []*Individual) {
	g := rdf2go.NewGraph(url)
	must(g.Parse(bytes.NewReader(ttl), "text/turtle"))

	iriToClass := map[string]*Class{}
	individuals := []*Individual{}

	var allUsedProps []*rdf2go.Triple
	used := func(triples ...*rdf2go.Triple) []*rdf2go.Triple {
		allUsedProps = append(allUsedProps, triples...)
		return triples
	}

	classes := used(g.All(nil, rdfType, owlClass)...)
	for _, class := range sorted(classes, bySubject) {
		log("Class", class)

		out := &Class{
			IRI:     cleanIRI(class.Subject.String()),
			Comment: getComment(used, g, class.Subject),
			Kind:    getNodeKind(used, g, class.Subject),
		}
		if _, ok := iriToClass[out.IRI]; ok {
			panic("duplicate type definition: " + out.IRI)
		}
		iriToClass[out.IRI] = out

		superclass := oneOptional(used, g.All(class.Subject, rdfSubclassOf, nil))
		if superclass != nil {
			out.ParentIRI = cleanIRI(superclass.Object.String())
			log("  extends: ", out.ParentIRI)
		}

		properties := used(g.All(class.Subject, shaclProperty, nil)...)

		for _, property := range sorted(properties, byObject) {
			path := oneRequired(used, g.All(property.Object, shaclPath, nil))
			used(path)
			log("  property:", class.Subject.String(), "path", path.Object.String())

			if path.Object.Equal(rdfType) {
				log("marking class %v as abstract due to path requirement:", out.GoName, path)
				out.Abstract = true
				continue
			}

			nodeKind := oneOptional(used, g.All(property.Object, shaclNodeKind, nil))
			if nodeKind == nil {
				continue // only include properties that have nodeKind, there are some properties specified for the path we need to exclude
			}

			prop := &Property{
				IRI:     cleanIRI(path.Object.String()),
				Comment: getComment(used, g, path.Object),
			}

			// get the data type
			typeIRI := oneOptional(used, g.All(property.Object, shaclClass, nil))
			if typeIRI == nil {
				typeIRI = oneOptional(used, g.All(property.Object, shaclDatatype, nil))
			}
			if typeIRI == nil {
				typeIRI = oneOptional(used, g.All(path.Object, rdfRange, nil))
			}
			if typeIRI != nil {
				prop.TypeIRI = cleanIRI(typeIRI.Object.String())
			} else {
				if oneOptional(used, g.All(property.Object, shaclNot, nil)) != nil {
					// SPDX TTL has some path references for types disallowed don't have a typeIRI,
					// they just include some constraints on the same path since they are not supposed to be instantiated directly
					continue
				}
				log("WARNING: No type IRI for: ", property.Object)
				panic("no typeIRI found for: " + property.Object.String())
			}

			out.Properties = append(out.Properties, prop)

			minCount := oneOptional(used, g.All(property.Object, shaclMinCount, nil))
			used(minCount)
			if minCount != nil {
				prop.MinCount = parseIntegerValue(minCount)
			}

			maxCount := oneOptional(used, g.All(property.Object, shaclMaxCount, nil))
			if maxCount != nil {
				prop.MaxCount = parseIntegerValue(maxCount)
			} else {
				prop.MaxCount = -1 // how is * represented?
			}

			allowedValues := oneOptional(used, g.All(property.Object, shaclIn, nil))
			var usedPropertyNodes []*rdf2go.Triple
			for allowedValues != nil {
				usedPropertyNodes = append(usedPropertyNodes, allowedValues)
				validation := oneOptional(used, g.All(allowedValues.Object, rdfFirst, nil))
				if validation != nil {
					used(validation)
					log("    validation:", nodeDisplay(validation))
					prop.Validations = append(prop.Validations, AllowedIRIValidation(validation.Object.String()))
				}
				allowedValues = oneOptional(used, g.All(allowedValues.Object, rdfRest, nil))
			}

			pattern := oneOptional(used, g.All(property.Object, shaclPattern, nil))
			if pattern != nil {
				prop.Validations = append(prop.Validations, MatchPatternValidation(cleanText(pattern.Object.String())))
			}

			//allProps := g.All(property.Object, nil, nil)
			//for _, p := range sorted(allProps, byObject) {
			//	if slices.Contains(append(usedPropertyNodes, path, property), p) {
			//		continue
			//	}
			//	log("  ... extra prop prop:", p)
			//}
		}

		//allprops := g.All(class.Subject, nil, nil)
		//for _, p := range sorted(allprops, byObject) {
		//	if slices.Contains(allUsedProps, p) {
		//		continue
		//	}
		//	log("  ... unused class prop:", p)
		//}

	}

	// read all the named individuals
	namedIndividuals := used(g.All(nil, rdfType, owlNamedIndividual)...)
	for _, namedIndividual := range sorted(namedIndividuals, bySubject) {
		entries := used(g.All(namedIndividual.Subject, rdfType, nil)...)
		for _, entry := range entries {
			if entry == namedIndividual {
				continue
			}
			typeIRI := cleanIRI(entry.Object.String())
			if c := iriToClass[typeIRI]; c != nil {
				label := ""
				if l := oneOptional(used, g.All(namedIndividual.Subject, rdfLabel, nil)); l != nil {
					label = cleanText(l.Object.String())
				}
				ni := Individual{
					IRI:     cleanIRI(namedIndividual.Subject.String()),
					Label:   label,
					Comment: getComment(used, g, namedIndividual.Subject),
					TypeIRI: typeIRI,
				}
				individuals = append(individuals, &ni)
			}
		}
	}

	log("------------------- UNUSED PROPS -------------------")
	for p := range g.IterTriples() {
		if slices.Contains(allUsedProps, p) {
			continue
		}
		log("  ... unused triple:", p)
	}

	return iriToClass, individuals
}

func getNodeKind(used usedFunc, g *rdf2go.Graph, subject rdf2go.Term) string {
	propNodeKind := oneOptional(used, g.All(subject, shaclNodeKind, nil))
	if propNodeKind != nil {
		return cleanIRI(propNodeKind.Object.RawValue())
	}
	return ""
}

func getComment(used usedFunc, g *rdf2go.Graph, subject rdf2go.Term) string {
	allComments := g.All(subject, rdfComment, nil)
	var comment *rdf2go.Triple
	for _, c := range allComments {
		value := c.Object.String()
		comment = c
		if strings.HasSuffix(value, "@en") {
			// use English comment
			break
		}
	}
	if comment != nil {
		used(comment)
		value := comment.Object.String()
		parts := strings.Split(value, "@")
		value = strings.Join(parts[:len(parts)-1], "@")
		return strings.TrimSpace(strings.Trim(value, "\""))
	}
	return ""
}

func parseIntegerValue(count *rdf2go.Triple) int {
	if count == nil {
		return 0
	}
	val := count.Object.String()
	val = strings.Split(val, "^")[0]
	val = strings.Trim(val, "\"")
	return get(strconv.Atoi(val))
}

var (
	rdfType            = rdf2go.NewResource("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	rdfFirst           = rdf2go.NewResource("http://www.w3.org/1999/02/22-rdf-syntax-ns#first")
	rdfRest            = rdf2go.NewResource("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest")
	rdfSubclassOf      = rdf2go.NewResource("http://www.w3.org/2000/01/rdf-schema#subClassOf")
	rdfComment         = rdf2go.NewResource("http://www.w3.org/2000/01/rdf-schema#comment")
	rdfLabel           = rdf2go.NewResource("http://www.w3.org/2000/01/rdf-schema#label")
	rdfRange           = rdf2go.NewResource("http://www.w3.org/2000/01/rdf-schema#range")
	owlClass           = rdf2go.NewResource("http://www.w3.org/2002/07/owl#Class")
	owlNamedIndividual = rdf2go.NewResource("http://www.w3.org/2002/07/owl#NamedIndividual")
	//owlObjectProperty   = rdf2go.NewResource("http://www.w3.org/2002/07/owl#ObjectProperty")
	//owlDatatypeProperty = rdf2go.NewResource("http://www.w3.org/2002/07/owl#DatatypeProperty")
	//shaclNodeShape    = rdf2go.NewResource("http://www.w3.org/ns/shacl#NodeShape")
	//shaclIRI            = rdf2go.NewResource("http://www.w3.org/ns/shacl#IRI")
	//shaclBlankNodeOrIRI = rdf2go.NewResource("http://www.w3.org/ns/shacl#BlankNodeOrIRI")
	shaclNodeKind = rdf2go.NewResource("http://www.w3.org/ns/shacl#nodeKind")
	shaclProperty = rdf2go.NewResource("http://www.w3.org/ns/shacl#property")
	shaclClass    = rdf2go.NewResource("http://www.w3.org/ns/shacl#class")
	shaclPath     = rdf2go.NewResource("http://www.w3.org/ns/shacl#path")
	shaclDatatype = rdf2go.NewResource("http://www.w3.org/ns/shacl#datatype")
	shaclIn       = rdf2go.NewResource("http://www.w3.org/ns/shacl#in")
	shaclNot      = rdf2go.NewResource("http://www.w3.org/ns/shacl#not")
	shaclMinCount = rdf2go.NewResource("http://www.w3.org/ns/shacl#minCount")
	shaclMaxCount = rdf2go.NewResource("http://www.w3.org/ns/shacl#maxCount")
	shaclPattern  = rdf2go.NewResource("http://www.w3.org/ns/shacl#pattern")
)
