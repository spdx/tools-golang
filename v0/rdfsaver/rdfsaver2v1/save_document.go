// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"errors"

	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) Document(doc *rdf2v1.Document) (docId goraptor.Term, err error) {

	if doc == nil {
		return nil, errors.New("Nil document.")
	}

	docId = rdf2v1.Blank("doc")

	if err = f.setNodeType(docId, rdf2v1.TypeDocument); err != nil {
		return
	}

	if err = f.addLiteral(docId, "specVersion", doc.SPDXVersion.Val); err != nil {
		return
	}

	if doc.DataLicense.Val != "" {
		if err = f.addTerm(docId, "dataLicense", rdf2v1.Uri(rdf2v1.LicenseUri+doc.DataLicense.Val)); err != nil {
			return
		}
	}
	if doc.DocumentName.Val != "" {
		if err = f.addTerm(docId, "name", rdf2v1.Uri(rdf2v1.LicenseUri+doc.DataLicense.Val)); err != nil {
			return
		}
	}
	if err = f.addLiteral(docId, "rdfs:comment", doc.DocumentComment.Val); err != nil {
		return
	}

	if id, err := f.CreationInfo(doc.CreationInfo); err == nil {
		if err = f.addTerm(docId, "creationInfo", id); err != nil {
			return docId, err
		}
	} else {
		return docId, err
	}

	if id, err := f.License(doc.License); err == nil {
		if err = f.addTerm(docId, "dataLicense", id); err != nil {
			return docId, err
		}
	} else {
		return docId, err
	}

	if id, err := f.ExternalDocumentRef(doc.ExternalDocumentRef); err == nil {
		if err = f.addTerm(docId, "externalDocumentRef", id); err != nil {
			return docId, err
		}
	} else {
		return docId, err
	}
	if err = f.Relationships(docId, "relationship", doc.Relationship); err != nil {
		return
	}

	if err = f.Reviews(docId, "reviewed", doc.Review); err != nil {
		return
	}
	if err = f.Annotations(docId, "annotation", doc.Annotation); err != nil {
		return
	}
	if err = f.ExtractedLicInfos(docId, "hasExtractedLicensingInfo", doc.ExtractedLicensingInfo); err != nil {
		return
	}
	return docId, nil
}
func (f *Formatter) ExternalDocumentRef(edr *rdf2v1.ExternalDocumentRef) (id goraptor.Term, err error) {
	id = f.NodeId("edr")

	if edr != nil {

		if err = f.setNodeType(id, rdf2v1.TypeExternalDocumentRef); err != nil {
			return
		}

		err = f.addPairs(id,
			Pair{"externalDocumentId", edr.ExternalDocumentId.Val},
			Pair{"spdxDocument", edr.SPDXDocument.Val},
		)

		if err != nil {
			return
		}

		if edr.Checksum != nil {
			cksumId, err := f.Checksum(edr.Checksum)
			if err != nil {
				return id, err
			}
			if err = f.addTerm(id, "checksum", cksumId); err != nil {
				return id, err
			}
		}
	}
	return id, nil
}
