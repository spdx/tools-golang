// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"errors"

	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) Snippet(snip *rdf2v1.Snippet) (snipId goraptor.Term, err error) {

	if snip == nil {
		return nil, errors.New("Nil Snippet.")
	}

	snipId = rdf2v1.Blank("snip")

	if err = f.setNodeType(snipId, rdf2v1.TypeSnippet); err != nil {
		return
	}
	if err = f.addLiteral(snipId, "name", snip.SnippetName.Val); err != nil {
		return
	}
	if err = f.addLiteral(snipId, "copyrightText", snip.SnippetCopyrightText.Val); err != nil {
		return
	}
	if err = f.addLiteral(snipId, "licenseComments", snip.SnippetLicenseComments.Val); err != nil {
		return
	}
	if err = f.addLiteral(snipId, "rdfs:comment", snip.SnippetComment.Val); err != nil {
		return
	}

	if snip.SnippetLicenseConcluded.Val != "" {
		if err = f.addTerm(snipId, "licenseConcluded", rdf2v1.Prefix(snip.SnippetLicenseConcluded.Val)); err != nil {
			return
		}
	}
	for _, li := range snip.LicenseInfoInSnippet {
		if err = f.addLiteral(snipId, "licenseInfoInSnippet", li.Val); err != nil {
			return
		}
	}
	if err = f.SnippetStartEndPointers(snipId, "range", snip.SnippetStartEndPointer); err != nil {

		return

	}
	if snip.SnippetFromFile != nil {
		sfId, err := f.File(snip.SnippetFromFile)
		if err != nil {
			return snipId, err
		}
		if err = f.addTerm(snipId, "snippetFromFile", sfId); err != nil {
			return snipId, err
		}
	}
	return snipId, nil
}

func (f *Formatter) SnippetStartEndPointer(se *rdf2v1.SnippetStartEndPointer) (id goraptor.Term, err error) {
	id = f.NodeId("ssep")

	if err = f.setNodeType(id, rdf2v1.TypeSnippetStartEndPointer); err != nil {
		return
	}

	if err = f.ByteOffsetPointers(id, "j.0:endPointer", se.ByteOffsetPointer); err != nil {
		if err = f.LineCharPointers(id, "j.0:endPointer", se.LineCharPointer); err != nil {
			return
		}
	}
	if err = f.LineCharPointers(id, "j.0:startPointer", se.LineCharPointer); err != nil {
		if err = f.ByteOffsetPointers(id, "j.0:startPointer", se.ByteOffsetPointer); err != nil {
			return
		}
	}

	return id, nil
}

func (f *Formatter) LineCharPointer(lcp *rdf2v1.LineCharPointer) (id goraptor.Term, err error) {
	id = f.NodeId("lc")

	if err = f.setNodeType(id, rdf2v1.TypeLineCharPointer); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"j.0:reference", lcp.Reference.Val},
		Pair{"j.0:lineNumber", lcp.LineNumber.Val},
	)

	if err != nil {
		return
	}

	return id, nil
}
func (f *Formatter) ByteOffsetPointer(bop *rdf2v1.ByteOffsetPointer) (id goraptor.Term, err error) {
	id = f.NodeId("bo")

	if err = f.setNodeType(id, rdf2v1.TypeByteOffsetPointer); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"j.0:reference", bop.Reference.Val},
		Pair{"j.0:offset", bop.Offset.Val},
	)

	if err != nil {
		return
	}

	return id, nil
}

func (f *Formatter) SnippetStartEndPointers(parent goraptor.Term, element string, ses []*rdf2v1.SnippetStartEndPointer) error {

	if len(ses) == 0 {
		return nil
	}

	for _, se := range ses {
		if se != nil {
		}
		sepId, err := f.SnippetStartEndPointer(se)

		if err != nil {
			return err
		}
		if sepId == nil {
			continue
		}
		if err = f.addTerm(parent, element, sepId); err != nil {
			return err
		}
	}

	return nil
}

func (f *Formatter) ByteOffsetPointers(parent goraptor.Term, element string, bos []*rdf2v1.ByteOffsetPointer) error {

	if len(bos) == 0 {
		return nil
	}

	for _, bo := range bos {
		if bo != nil {
			bopId, err := f.ByteOffsetPointer(bo)
			if err != nil {
				return err
			}
			if bopId == nil {
				continue
			}
			if err = f.addTerm(parent, element, bopId); err != nil {
				return err
			}
		}
	}
	return nil
}
func (f *Formatter) LineCharPointers(parent goraptor.Term, element string, lcs []*rdf2v1.LineCharPointer) error {

	if len(lcs) == 0 {
		return nil
	}

	for _, lc := range lcs {
		if lc != nil {
			lcId, err := f.LineCharPointer(lc)
			if err != nil {
				return err
			}
			if lcId == nil {
				continue
			}
			if err = f.addTerm(parent, element, lcId); err != nil {
				return err
			}
		}
	}
	return nil
}
