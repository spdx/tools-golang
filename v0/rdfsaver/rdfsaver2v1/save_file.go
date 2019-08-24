// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) File(file *rdf2v1.File) (id goraptor.Term, err error) {
	id, ok := f.fileIds[file.FileName.Val]
	if ok {
		return
	}

	id = f.NodeId("file")
	f.fileIds[file.FileName.Val] = id

	if err = f.setNodeType(id, rdf2v1.TypeFile); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"fileName", file.FileName.Val},
		Pair{"licenseComments", file.FileLicenseComments.Val},
		Pair{"copyrightText", file.FileCopyrightText.Val},
		Pair{"rdfs:comment", file.FileComment.Val},
		Pair{"noticeText", file.FileNoticeText.Val},
	)

	if err != nil {
		return
	}
	if file.FileChecksum != nil {
		cksumId, err := f.Checksum(file.FileChecksum)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "checksum", cksumId); err != nil {
			return id, err
		}
	}

	if file.ExtractedLicensingInfo != nil {
		exlicId, err := f.ExtractedLicInfo(file.ExtractedLicensingInfo)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "hasExtractedLicensingInfo", exlicId); err != nil {
			return id, err
		}
	}
	if file.DisjunctiveLicenseSet != nil {
		dlsId, err := f.DisjunctiveLicenseSet(file.DisjunctiveLicenseSet)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "member", dlsId); err != nil {
			return id, err
		}
	}
	if file.ConjunctiveLicenseSet != nil {
		clsId, err := f.ConjunctiveLicenseSet(file.ConjunctiveLicenseSet)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "member", clsId); err != nil {
			return id, err
		}
	}

	for _, fc := range file.FileContributor {
		if err = f.addLiteral(id, "fileContributor", fc.Val); err != nil {
			return
		}
	}
	for _, lif := range file.LicenseInfoInFile {
		if err = f.addTerm(id, "licenseInfoInFile", rdf2v1.Prefix(lif.Val)); err != nil {
			return
		}
	}
	for _, ft := range file.FileType {
		if err = f.addTerm(id, "fileType", rdf2v1.Prefix(ft.Val)); err != nil {
			return
		}
	}

	if file.SnippetLicense != nil {
		filelicId, err := f.License(file.SnippetLicense)
		if err != nil {
			filelicId, err = f.DisjunctiveLicenseSet(file.DisjunctiveLicenseSet)
			if err != nil {
				filelicId, err = f.ExtractedLicInfo(file.ExtractedLicensingInfo)
				if err != nil {
					return id, err
				}
			}
		}
		if err = f.addTerm(id, "licenseConcluded", filelicId); err != nil {
			return id, err

		}
	}

	for _, dep := range file.FileDependency {
		if file.FileDependency != nil {
			fdId, err := f.File(file.FileDependency)
			if err != nil {
				return id, err
			}
			if err = f.addTerm(id, "fileDependency", fdId); err != nil {
				return id, err
			}
		}
	}

  if file.FileRelationship != nil {
		frId, err := f.Relationship(file.FileRelationship)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "relationship", frId); err != nil {
			return id, err
		}
	}

	if err = f.Annotations(id, "annotation", file.Annotation); err != nil {
		return
	}
	if err = f.Projects(id, "artifactOf", file.Project); err != nil {
		return
	}
	return id, err

}

func (f *Formatter) Files(parent goraptor.Term, element string, files []*rdf2v1.File) error {
	if len(files) == 0 {
		return nil
	}
	for _, file := range files {
		fId, err := f.File(file)
		if err != nil {
			return err
		}
		if fId == nil {
			continue
		}
		if err = f.addTerm(parent, element, fId); err != nil {
			return err
		}
	}
	return nil
}
func (f *Formatter) Project(pro *rdf2v1.Project) (id goraptor.Term, err error) {
	id = f.NodeId("pro")

	if err = f.setNodeType(id, rdf2v1.TypeProject); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"homepage", pro.HomePage.Val},
		Pair{"name", pro.Name.Val},
	)

	return id, err
}
func (f *Formatter) Projects(parent goraptor.Term, element string, pros []*rdf2v1.Project) error {

	if len(pros) == 0 {
		return nil
	}

	for _, pro := range pros {
		proId, err := f.Project(pro)
		if err != nil {
			return err
		}
		if proId == nil {
			continue
		}
		if err = f.addTerm(parent, element, proId); err != nil {
			return err
		}
	}
	return nil
}
