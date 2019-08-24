// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"strconv"

	"github.com/spdx/tools-golang/v0/spdx"
)

func CollectDocument(doc2v1 *spdx.Document2_1) *Document {

	stdDoc := Document{
		SPDXVersion:            Str(doc2v1.CreationInfo.SPDXVersion),
		DataLicense:            Str(doc2v1.CreationInfo.DataLicense),
		Review:                 collectReview(doc2v1),
		DocumentName:           Str(doc2v1.CreationInfo.DocumentName),
		DocumentNamespace:      Str(doc2v1.CreationInfo.DocumentNamespace),
		SPDXID:                 Str(doc2v1.CreationInfo.SPDXIdentifier),
		DocumentComment:        Str(doc2v1.CreationInfo.DocumentComment),
		ExtractedLicensingInfo: collectExtractedLicInfo(doc2v1),
		Relationship:           collectRelationships(doc2v1),
		CreationInfo:           collectCreationInfo(doc2v1),
		Annotation:             collectDocAnnotation(doc2v1),
		ExternalDocumentRef:    collectExternalDocumentRef(doc2v1),
		License:                collectLicense(doc2v1),
	}
	return &stdDoc
}

func collectCreationInfo(doc2v1 *spdx.Document2_1) *CreationInfo {

	stdCi := CreationInfo{

		SPDXIdentifier:     Str(doc2v1.CreationInfo.SPDXIdentifier),
		LicenseListVersion: Str(doc2v1.CreationInfo.LicenseListVersion),
		Creator:            InsertCreator(doc2v1.CreationInfo),
		Create:             Str(doc2v1.CreationInfo.Created),
		Comment:            Str(doc2v1.CreationInfo.CreatorComment),
	}
	return &stdCi
}

func collectExternalDocumentRef(doc2v1 *spdx.Document2_1) *ExternalDocumentRef {
	
  if doc2v1.CreationInfo.ExternalDocumentReferences != nil {
		stdEdr := ExternalDocumentRef{

			ExternalDocumentId: Str(doc2v1.CreationInfo.ExternalDocumentReferences[0]),
			SPDXDocument:       Str(doc2v1.CreationInfo.ExternalDocumentReferences[1]),
			Checksum:           collectDocChecksum(doc2v1),
		}
		return &stdEdr
	}
	return nil
}

func collectLicense(doc2v1 *spdx.Document2_1) *License {
	stdEdr := License{
		LicenseId: Str(LicenseUri + doc2v1.CreationInfo.DataLicense),
	}
	return &stdEdr
}

func collectVerificationCode(pkg2_1 *spdx.Package2_1) *PackageVerificationCode {

	stdVc := PackageVerificationCode{

		PackageVerificationCode:             Str(pkg2_1.PackageVerificationCode),
		PackageVerificationCodeExcludedFile: Str(pkg2_1.PackageVerificationCodeExcludedFile),
	}
	return &stdVc
}

func collectPackageChecksum(pkg2_1 *spdx.Package2_1) *Checksum {

	stdPc := Checksum{

		Algorithm:     Str(PkgChecksumAlgo(pkg2_1)),
		ChecksumValue: Str(PkgChecksumValue(pkg2_1)),
	}
	return &stdPc
}

func collectFileChecksum(File2_1 *spdx.File2_1) *Checksum {

	stdFc := Checksum{

		Algorithm:     Str(FileChecksumAlgo(File2_1)),
		ChecksumValue: Str(FileChecksumValue(File2_1)),
	}
	return &stdFc
}

func collectDocChecksum(doc2v1 *spdx.Document2_1) *Checksum {

	stdFc := Checksum{

		Algorithm:     Str(doc2v1.CreationInfo.ExternalDocumentReferences[2]),
		ChecksumValue: Str(doc2v1.CreationInfo.ExternalDocumentReferences[3]),
	}
	return &stdFc
}

func collectDocAnnotation(doc2v1 *spdx.Document2_1) []*Annotation {
	var arrAnn []*Annotation
	for _, an := range doc2v1.Annotations {
		if an.AnnotationSPDXIdentifier == doc2v1.CreationInfo.SPDXIdentifier {
			stdAnn := Annotation{
				Annotator:                Str(an.AnnotatorType + an.Annotator),
				AnnotationType:           Str(an.AnnotationType),
				AnnotationDate:           Str(an.AnnotationDate),
				AnnotationComment:        Str(an.AnnotationComment),
				AnnotationSPDXIdentifier: Str(an.AnnotationSPDXIdentifier),
			}
			pointer := &stdAnn
			arrAnn = append(arrAnn, pointer)
		}
	}
	return arrAnn
}

func collectFileAnnotation(doc2v1 *spdx.Document2_1, f *spdx.File2_1) []*Annotation {
	var arrAnn []*Annotation
	for _, pkg := range doc2v1.Packages {
		for _, file := range pkg.Files {
			if file != nil {
				for _, an := range doc2v1.Annotations {
					if an.AnnotationSPDXIdentifier == f.FileSPDXIdentifier {
						stdAnn := Annotation{
							Annotator:                Str(an.AnnotatorType + an.Annotator),
							AnnotationType:           Str(an.AnnotationType),
							AnnotationDate:           Str(an.AnnotationDate),
							AnnotationComment:        Str(an.AnnotationComment),
							AnnotationSPDXIdentifier: Str(an.AnnotationSPDXIdentifier),
						}
						pointer := &stdAnn
						arrAnn = append(arrAnn, pointer)
					}
				}
			}
		}
	}
	return arrAnn
}

func collectPackageAnnotation(doc2v1 *spdx.Document2_1) []*Annotation {
	var arrAnn []*Annotation
	for _, pkg := range doc2v1.Packages {
		for _, an := range doc2v1.Annotations {
			if an.AnnotationSPDXIdentifier == pkg.PackageSPDXIdentifier {
				stdAnn := Annotation{
					Annotator:                Str(an.AnnotatorType + an.Annotator),
					AnnotationType:           Str(an.AnnotationType),
					AnnotationDate:           Str(an.AnnotationDate),
					AnnotationComment:        Str(an.AnnotationComment),
					AnnotationSPDXIdentifier: Str(an.AnnotationSPDXIdentifier),
				}
				pointer := &stdAnn
				arrAnn = append(arrAnn, pointer)
			}
		}
	}
	return arrAnn
}

func collectReview(doc2v1 *spdx.Document2_1) []*Review {
	var arrRev []*Review
	for _, a := range doc2v1.Reviews {
		if a != nil {
			stdRev := Review{
				Reviewer:      Str(a.Reviewer),
				ReviewDate:    Str(a.ReviewDate),
				ReviewComment: Str(a.ReviewComment),
			}
			pointer := &stdRev
			arrRev = append(arrRev, pointer)
		}
	}

	return arrRev
}

func collectRelationships(doc2v1 *spdx.Document2_1) []*Relationship {
	var arrRel []*Relationship
	for _, a := range doc2v1.Relationships {

		if a != nil {
			stdRel := Relationship{
				RelationshipType:    Str(a.Relationship),
				RelationshipComment: Str(a.RelationshipComment),
				Package:             collectPackagesfromRelationship(a, doc2v1),
			}
			pointer := &stdRel
			arrRel = append(arrRel, pointer)
		}

	}

	return arrRel
}

func collectFilesfromPackages(pkg *spdx.Package2_1, doc2v1 *spdx.Document2_1) []*File {
	var arrFile []*File
	if pkg != nil {
		for _, b := range pkg.Files {
			if b != nil {
				stdFile := File{
					FileName:            Str(b.FileName),
					FileSPDXIdentifier:  Str(b.FileSPDXIdentifier),
					FileType:            ValueStrList(b.FileType),
					FileChecksum:        collectFileChecksum(b),
					LicenseInfoInFile:   ValueStrList(b.LicenseInfoInFile),
					FileLicenseComments: Str(b.LicenseComments),
					FileCopyrightText:   Str(b.FileCopyrightText),
					Project:             collectArtifactOfProject(b),
					FileComment:         Str(b.FileComment),
					FileNoticeText:      Str(b.FileNotice),
					FileContributor:     ValueStrList(b.FileContributor),
					Annotation:          collectFileAnnotation(doc2v1, b),
				}
				pointer := &stdFile
				arrFile = append(arrFile, pointer)
			}
		}
	}
	return arrFile
}

func collectPackagesfromRelationship(rel *spdx.Relationship2_1, doc2v1 *spdx.Document2_1) []*Package {
	var arrPkg []*Package
	for _, a := range doc2v1.Packages {
		if a != nil {
			if a.PackageSPDXIdentifier == rel.RefB {
				stdPkg := Package{
					PackageName:                 Str(a.PackageName),
					PackageVersionInfo:          Str(a.PackageVersion),
					PackageFileName:             Str(a.PackageFileName),
					PackageSPDXIdentifier:       Str(a.PackageSPDXIdentifier),
					PackageDownloadLocation:     Str(a.PackageDownloadLocation),
					PackageVerificationCode:     collectVerificationCode(a),
					PackageComment:              Str(a.PackageComment),
					PackageChecksum:             collectPackageChecksum(a),
					PackageLicenseComments:      Str(a.PackageLicenseComments),
					PackageLicenseInfoFromFiles: ValueStrList(a.PackageLicenseInfoFromFiles),
					PackageLicenseDeclared:      Str(a.PackageLicenseDeclared),
					PackageCopyrightText:        Str(a.PackageCopyrightText),
					PackageHomepage:             Str(a.PackageHomePage),
					PackageSupplier:             Str(InsertSupplier(a)),
					PackageExternalRef:          collectPkgExternalRef(a),
					PackageOriginator:           Str(InsertOriginator(a)),
					PackageSourceInfo:           Str(a.PackageSummary),
					FilesAnalyzed:               Str(strconv.FormatBool(a.FilesAnalyzed)),
					PackageSummary:              Str(a.PackageSummary),
					PackageDescription:          Str(a.PackageDescription),
					Annotation:                  collectPackageAnnotation(doc2v1),
					File:                        collectFilesfromPackages(a, doc2v1),
				}
				pointer := &stdPkg
				arrPkg = append(arrPkg, pointer)
			}
		}
	}
	return arrPkg
}

func collectExtractedLicInfo(doc2v1 *spdx.Document2_1) []*ExtractedLicensingInfo {
	var arrEl []*ExtractedLicensingInfo
	for _, a := range doc2v1.OtherLicenses {
		if a != nil {
			stdEl := ExtractedLicensingInfo{
				LicenseIdentifier: Str(a.LicenseIdentifier),
				LicenseName:       Str(a.LicenseName),
				ExtractedText:     Str(a.ExtractedText),
				LicenseComment:    Str(a.LicenseComment),
				LicenseSeeAlso:    ValueStrList(a.LicenseCrossReferences),
			}
			pointer := &stdEl
			arrEl = append(arrEl, pointer)
		}
	}
	return arrEl
}

func collectArtifactOfProject(a *spdx.File2_1) []*Project {
	var arrp []*Project

	for _, c := range a.ArtifactOfProjects {
		stdp := Project{
			Name:     Str(c.Name),
			HomePage: Str(c.HomePage),
		}

		pointer := &stdp
		arrp = append(arrp, pointer)
	}

	return arrp
}

func CollectSnippets(doc2v1 *spdx.Document2_1) *Snippet {
	for _, pkg := range doc2v1.Packages {
		if pkg != nil {
			for _, file := range pkg.Files {
				if file != nil {
					for _, sp := range file.Snippets {
						if sp != nil {
							stdSn := Snippet{
								SnippetName:             Str(sp.SnippetName),
								SnippetSPDXIdentifier:   Str(sp.SnippetSPDXIdentifier),
								SnippetLicenseComments:  Str(sp.SnippetLicenseComments),
								SnippetCopyrightText:    Str(sp.SnippetCopyrightText),
								SnippetLicenseConcluded: Str(sp.SnippetLicenseConcluded),
								SnippetComment:          Str(sp.SnippetComment),
								LicenseInfoInSnippet:    ValueStrList(sp.LicenseInfoInSnippet),
							}

							return &stdSn
						}
					}
				}
			}
		}
	}
	return nil
}

func collectPkgExternalRef(pkg *spdx.Package2_1) []*ExternalRef {
	var arrPer []*ExternalRef
	for _, a := range pkg.PackageExternalReferences {
		if a != nil {
			stdEl := ExternalRef{
				ReferenceLocator:  Str(a.Locator),
				ReferenceType:     collectReferenceType(a),
				ReferenceCategory: Str(a.Category),
				ReferenceComment:  Str(a.ExternalRefComment),
			}
			pointer := &stdEl
			arrPer = append(arrPer, pointer)
		}
	}
	return arrPer
}

func collectReferenceType(pkger *spdx.PackageExternalReference2_1) *ReferenceType {

	stdRt := ReferenceType{
		ReferenceType: Str(pkger.RefType),
	}
	return &stdRt
}
