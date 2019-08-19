package rdf2v1

import (
	"strconv"

	"github.com/spdx/tools-golang/v0/spdx"
)

func TransferDocument(spdxdoc *Document, sp *Snippet) *spdx.Document2_1 {

	stdDoc := spdx.Document2_1{

		CreationInfo:  transferCreationInfo(spdxdoc),
		Packages:      transferPackages(spdxdoc, sp),
		OtherLicenses: transferOtherLicenses(spdxdoc),
		Relationships: transferRelationships(spdxdoc),
		Annotations:   transferAnnotation(spdxdoc),
		Reviews:       transferReview(spdxdoc),
	}
	return &stdDoc
}

func transferCreationInfo(spdxdoc *Document) *spdx.CreationInfo2_1 {

	var listExtDocRef []string
	listExtDocRef = append(listExtDocRef, spdxdoc.ExternalDocumentRef.ExternalDocumentId.Val)
	listExtDocRef = append(listExtDocRef, spdxdoc.ExternalDocumentRef.SPDXDocument.Val)
	listExtDocRef = append(listExtDocRef, ExtractChecksumAlgo(spdxdoc.ExternalDocumentRef.Checksum.Algorithm.Val))
	listExtDocRef = append(listExtDocRef, spdxdoc.ExternalDocumentRef.Checksum.ChecksumValue.Val)

	stdCi := spdx.CreationInfo2_1{

		SPDXVersion:                spdxdoc.SPDXVersion.Val,
		DataLicense:                spdxdoc.License.LicenseId.Val,
		SPDXIdentifier:             spdxdoc.SPDXID.Val,
		DocumentName:               spdxdoc.DocumentName.Val,
		DocumentNamespace:          spdxdoc.DocumentNamespace.Val,
		ExternalDocumentReferences: listExtDocRef,
		LicenseListVersion:         spdxdoc.CreationInfo.LicenseListVersion.Val,
		CreatorPersons:             ExtractCreator(spdxdoc.CreationInfo, "Person"),
		CreatorOrganizations:       ExtractCreator(spdxdoc.CreationInfo, "Organization"),
		CreatorTools:               ExtractCreator(spdxdoc.CreationInfo, "Tool"),
		Created:                    spdxdoc.CreationInfo.Create.Val,
		CreatorComment:             spdxdoc.CreationInfo.Comment.Val,
		DocumentComment:            spdxdoc.DocumentComment.Val,
	}
	return &stdCi
}

func transferPackages(spdxdoc *Document, sp *Snippet) []*spdx.Package2_1 {
	var arrPkg []*spdx.Package2_1
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.Package != nil {
				for _, b := range a.Package {

					if b != nil {
						stdPkg := spdx.Package2_1{
							IsUnpackaged:          b.PackageName.Val == "",
							PackageName:           b.PackageName.Val,
							PackageSPDXIdentifier: b.PackageSPDXIdentifier.Val,
							PackageVersion:        b.PackageVersionInfo.Val,
							PackageFileName:       b.PackageFileName.Val,

							PackageSupplierPerson:       ExtractValueType(b.PackageSupplier.Val, "Person"),
							PackageSupplierOrganization: ExtractValueType(b.PackageSupplier.Val, "Organization"),
							PackageSupplierNOASSERTION:  b.PackageSupplier.Val == "NOASSERTION",

							PackageOriginatorPerson:       ExtractValueType(b.PackageOriginator.Val, "Person"),
							PackageOriginatorOrganization: ExtractValueType(b.PackageOriginator.Val, "Organization"),
							PackageOriginatorNOASSERTION:  b.PackageOriginator.Val == "NOASSERTION",

							PackageDownloadLocation:             b.PackageDownloadLocation.Val,
							FilesAnalyzed:                       !(b.PackageName.Val == ""),
							IsFilesAnalyzedTagPresent:           b.PackageName.Val == "",
							PackageVerificationCode:             b.PackageVerificationCode.PackageVerificationCode.Val,
							PackageVerificationCodeExcludedFile: b.PackageVerificationCode.PackageVerificationCodeExcludedFile.Val,

							PackageChecksumSHA1:   AlgoValue(b.PackageChecksum, "SHA1"),
							PackageChecksumSHA256: AlgoValue(b.PackageChecksum, "SHA256"),
							PackageChecksumMD5:    AlgoValue(b.PackageChecksum, "MD5"),

							PackageHomePage:   b.PackageHomepage.Val,
							PackageSourceInfo: b.PackageSourceInfo.Val,

							PackageLicenseConcluded:     PackageLicenseConcluded(b),
							PackageLicenseInfoFromFiles: ValueList(b.PackageLicenseInfoFromFiles),
							PackageLicenseDeclared:      b.PackageLicenseDeclared.Val,
							PackageLicenseComments:      b.PackageLicenseComments.Val,

							PackageCopyrightText:      b.PackageCopyrightText.Val,
							PackageSummary:            b.PackageSummary.Val,
							PackageDescription:        b.PackageDescription.Val,
							PackageComment:            b.PackageComment.Val,
							PackageExternalReferences: transferPkgExternalRef(b),
							Files:                     transferFilefromPackages(b, sp),
						}

						pointer := &stdPkg
						arrPkg = append(arrPkg, pointer)
					}
				}
			}
		}
	}
	return arrPkg
}

func transferOtherLicenses(spdxdoc *Document) []*spdx.OtherLicense2_1 {
	var arrOl []*spdx.OtherLicense2_1
	for _, a := range spdxdoc.ExtractedLicensingInfo {
		if a != nil {
			stdOl := spdx.OtherLicense2_1{
				LicenseIdentifier:      a.LicenseIdentifier.Val,
				ExtractedText:          a.ExtractedText.Val,
				LicenseName:            a.LicenseName.Val,
				LicenseCrossReferences: ValueList(a.LicenseSeeAlso),
				LicenseComment:         a.LicenseComment.Val,
			}
			pointer := &stdOl
			arrOl = append(arrOl, pointer)
		}
	}
	return arrOl
}

func transferRelationships(spdxdoc *Document) []*spdx.Relationship2_1 {
	var arrRel []*spdx.Relationship2_1
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			stdRel := spdx.Relationship2_1{
				RefA:                spdxdoc.SPDXID.Val,
				RefB:                RelationshipRef(a),
				Relationship:        ExtractRelType(a.RelationshipType.Val),
				RelationshipComment: a.RelationshipComment.Val,
			}
			pointer := &stdRel
			arrRel = append(arrRel, pointer)
		}
	}

	return arrRel
}

func transferAnnotation(spdxdoc *Document) []*spdx.Annotation2_1 {
	var arrAnn []*spdx.Annotation2_1
	// Annotations from Document
	for _, an := range spdxdoc.Annotation {
		stdAnn := spdx.Annotation2_1{
			Annotator:                ExtractKeyValue(an.Annotator.Val, "subvalue"),
			AnnotatorType:            ExtractKeyValue(an.Annotator.Val, "subkey"),
			AnnotationType:           an.AnnotationType.Val,
			AnnotationDate:           an.AnnotationDate.Val,
			AnnotationComment:        an.AnnotationComment.Val,
			AnnotationSPDXIdentifier: spdxdoc.SPDXID.Val,
		}
		pointer := &stdAnn
		arrAnn = append(arrAnn, pointer)
	}
	// Annotations from Packages
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.Package != nil {
				for _, b := range a.Package {
					if b != nil {
						for _, an := range b.Annotation {

							stdAnn := spdx.Annotation2_1{
								Annotator:                ExtractKeyValue(an.Annotator.Val, "subvalue"),
								AnnotatorType:            ExtractKeyValue(an.Annotator.Val, "subkey"),
								AnnotationType:           an.AnnotationType.Val,
								AnnotationDate:           an.AnnotationDate.Val,
								AnnotationComment:        an.AnnotationComment.Val,
								AnnotationSPDXIdentifier: b.PackageSPDXIdentifier.Val,
							}
							pointer := &stdAnn
							arrAnn = append(arrAnn, pointer)
						}
					}
				}
			}
		}

	}
	// Annotations from Files in Packages
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.Package != nil {
				for _, b := range a.Package {
					if b != nil {
						for _, c := range b.File {
							if c != nil {
								for _, an := range b.Annotation {

									stdAnn := spdx.Annotation2_1{
										Annotator:                ExtractKeyValue(an.Annotator.Val, "subvalue"),
										AnnotatorType:            ExtractKeyValue(an.Annotator.Val, "subkey"),
										AnnotationType:           an.AnnotationType.Val,
										AnnotationDate:           an.AnnotationDate.Val,
										AnnotationComment:        an.AnnotationComment.Val,
										AnnotationSPDXIdentifier: b.PackageSPDXIdentifier.Val,
									}
									pointer := &stdAnn
									arrAnn = append(arrAnn, pointer)
								}
							}
						}
					}
				}
			}
		}

	}
	// Annotations from Files in Relationship
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.File != nil {
				for _, b := range a.File {
					if b != nil {
						for _, an := range b.Annotation {

							stdAnn := spdx.Annotation2_1{
								Annotator:                ExtractKeyValue(an.Annotator.Val, "subvalue"),
								AnnotatorType:            ExtractKeyValue(an.Annotator.Val, "subkey"),
								AnnotationType:           an.AnnotationType.Val,
								AnnotationDate:           an.AnnotationDate.Val,
								AnnotationComment:        an.AnnotationComment.Val,
								AnnotationSPDXIdentifier: b.FileSPDXIdentifier.Val,
							}
							pointer := &stdAnn
							arrAnn = append(arrAnn, pointer)
						}
					}
				}
			}
		}

	}
	return arrAnn
}

func transferReview(spdxdoc *Document) []*spdx.Review2_1 {
	var arrRev []*spdx.Review2_1
	for _, a := range spdxdoc.Review {
		if a != nil {
			stdRev := spdx.Review2_1{
				Reviewer:      a.Reviewer.Val,
				ReviewerType:  ExtractKey(a.Reviewer.Val),
				ReviewDate:    a.ReviewDate.Val,
				ReviewComment: a.ReviewComment.Val,
			}
			pointer := &stdRev
			arrRev = append(arrRev, pointer)
		}
	}

	return arrRev
}

func transferFilefromRel(spdxdoc *Document, sp *Snippet) []*spdx.File2_1 {
	var arrFile []*spdx.File2_1
	var dependencyList []ValueStr
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.File != nil {
				for _, b := range a.File {
					if b != nil {
						for _, c := range b.FileDependency {
							if c != nil {
								dependencyList = append(dependencyList, c.FileName)
							}
						}
						stdFile := spdx.File2_1{

							FileName:           b.FileName.Val,
							FileSPDXIdentifier: b.FileSPDXIdentifier.Val,
							FileType:           ValueList(b.FileType),
							FileChecksumSHA1:   AlgoValue(b.FileChecksum, "SHA1"),
							FileChecksumSHA256: AlgoValue(b.FileChecksum, "SHA256"),
							FileChecksumMD5:    AlgoValue(b.FileChecksum, "MD5"),
							LicenseConcluded:   FileLicenseConcluded(b),
							LicenseInfoInFile:  ValueList(b.LicenseInfoInFile),
							LicenseComments:    b.FileLicenseComments.Val,
							FileCopyrightText:  b.FileCopyrightText.Val,
							ArtifactOfProjects: transferArtifactOfProject(b),
							FileComment:        b.FileComment.Val,
							FileNotice:         b.FileNoticeText.Val,
							FileContributor:    ValueList(b.FileContributor),
							FileDependencies:   ValueList(dependencyList),
							Snippets:           transferSnippets(sp),
						}
						pointer := &stdFile
						arrFile = append(arrFile, pointer)
					}
				}
			}
		}
	}
	return arrFile
}

func transferFilefromPackages(pkg *Package, sp *Snippet) []*spdx.File2_1 {
	var arrFile []*spdx.File2_1
	var dependencyList []ValueStr
	for _, b := range pkg.File {
		if b != nil {
			for _, c := range b.FileDependency {
				if c != nil {
					dependencyList = append(dependencyList, c.FileName)
				}
			}
			stdFile := spdx.File2_1{

				FileName:           b.FileName.Val,
				FileSPDXIdentifier: b.FileSPDXIdentifier.Val,
				FileType:           ValueList(b.FileType),
				FileChecksumSHA1:   AlgoValue(b.FileChecksum, "SHA1"),
				FileChecksumSHA256: AlgoValue(b.FileChecksum, "SHA256"),
				FileChecksumMD5:    AlgoValue(b.FileChecksum, "MD5"),
				LicenseConcluded:   FileLicenseConcluded(b),
				LicenseInfoInFile:  ValueList(b.LicenseInfoInFile),
				LicenseComments:    b.FileLicenseComments.Val,
				FileCopyrightText:  b.FileCopyrightText.Val,
				ArtifactOfProjects: transferArtifactOfProject(b),
				FileComment:        b.FileComment.Val,
				FileNotice:         b.FileNoticeText.Val,
				FileContributor:    ValueList(b.FileContributor),
				FileDependencies:   ValueList(dependencyList),
				Snippets:           transferSnippets(sp),
			}
			pointer := &stdFile
			arrFile = append(arrFile, pointer)
		}
	}

	return arrFile
}

func transferFilefromSnippets(sp *Snippet) *spdx.File2_1 {
	var dependencyList []ValueStr

	if sp.SnippetFromFile != nil {
		dependencyList = append(dependencyList, sp.SnippetFromFile.FileName)
	}

	stdFile := spdx.File2_1{

		FileName:           sp.SnippetFromFile.FileName.Val,
		FileSPDXIdentifier: sp.SnippetFromFile.FileSPDXIdentifier.Val,
		FileType:           ValueList(sp.SnippetFromFile.FileType),
		FileChecksumSHA1:   AlgoValue(sp.SnippetFromFile.FileChecksum, "SHA1"),
		FileChecksumSHA256: AlgoValue(sp.SnippetFromFile.FileChecksum, "SHA256"),
		FileChecksumMD5:    AlgoValue(sp.SnippetFromFile.FileChecksum, "MD5"),
		LicenseConcluded:   FileLicenseConcluded(sp.SnippetFromFile),
		LicenseInfoInFile:  ValueList(sp.SnippetFromFile.LicenseInfoInFile),
		LicenseComments:    sp.SnippetFromFile.FileLicenseComments.Val,
		FileCopyrightText:  sp.SnippetFromFile.FileCopyrightText.Val,
		ArtifactOfProjects: transferArtifactOfProject(sp.SnippetFromFile),
		FileComment:        sp.SnippetFromFile.FileComment.Val,
		FileNotice:         sp.SnippetFromFile.FileNoticeText.Val,
		FileContributor:    ValueList(sp.SnippetFromFile.FileContributor),
		FileDependencies:   ValueList(dependencyList),
		Snippets:           transferSnippets(sp),
	}

	return &stdFile
}

func transferArtifactOfProject(file *File) []*spdx.ArtifactOfProject2_1 {
	var arrAop []*spdx.ArtifactOfProject2_1
	for _, c := range file.Project {
		stdAop := spdx.ArtifactOfProject2_1{
			Name:     c.Name.Val,
			HomePage: c.HomePage.Val,
			URI:      c.URI.Val,
		}

		pointer := &stdAop
		arrAop = append(arrAop, pointer)
	}

	return arrAop
}

func transferSnippets(sp *Snippet) []*spdx.Snippet2_1 {
	var arrSn []*spdx.Snippet2_1
	if sp != nil {
		stdSn := spdx.Snippet2_1{
			SnippetSPDXIdentifier:         sp.SnippetSPDXIdentifier.Val,
			SnippetFromFileSPDXIdentifier: sp.SnippetFromFile.FileSPDXIdentifier.Val,
			SnippetName:                   sp.SnippetName.Val,
			SnippetLicenseComments:        sp.SnippetLicenseComments.Val,
			SnippetCopyrightText:          sp.SnippetCopyrightText.Val,
			SnippetLicenseConcluded:       sp.SnippetLicenseConcluded.Val,
			SnippetComment:                sp.SnippetComment.Val,
			LicenseInfoInSnippet:          ValueList(sp.LicenseInfoInSnippet),
		}
		pointer := &stdSn
		arrSn = append(arrSn, pointer)
		return arrSn
	}
	return nil
}

func transferPkgExternalRef(pkg *Package) []*spdx.PackageExternalReference2_1 {
	var arrPer []*spdx.PackageExternalReference2_1
	for _, a := range pkg.PackageExternalRef {
		if a != nil {

			stdPer := spdx.PackageExternalReference2_1{
				Category:           a.ReferenceCategory.Val,
				Locator:            a.ReferenceLocator.Val,
				RefType:            a.ReferenceType.ReferenceType.Val,
				ExternalRefComment: a.ReferenceComment.Val,
			}
			pointer := &stdPer
			arrPer = append(arrPer, pointer)
		}
	}

	return arrPer
}

func FileLicenseConcluded(file *File) string {
	var lc string
	if file.DisjunctiveLicenseSet != nil {
		var lcl []string
		lcl = ValueList(file.DisjunctiveLicenseSet.Member)
		for _, i := range lcl {
			if lc != "" {
				lc = lc + "or"
			}
			lc = lc + i
		}
	}
	if file.ConjunctiveLicenseSet != nil {
		lc = file.ConjunctiveLicenseSet.License.LicenseId.Val
	}
	if file.ExtractedLicensingInfo != nil {
		lc = file.ExtractedLicensingInfo.LicenseIdentifier.Val
	}
	if file.SnippetLicense != nil {
		lc = file.SnippetLicense.LicenseId.Val
	}
	return lc
}

func PackageLicenseConcluded(pkg *Package) string {
	var lc string
	if pkg.DisjunctiveLicenseSet != nil {
		var lcl []string
		lcl = ValueList(pkg.DisjunctiveLicenseSet.Member)
		for _, i := range lcl {
			if lc != "" {
				lc = lc + "or"
			}
			lc = lc + i
		}
		return lc
	}
	if pkg.ConjunctiveLicenseSet != nil {
		lc = pkg.ConjunctiveLicenseSet.License.LicenseId.Val
		return lc

	}

	if pkg.PackageLicense != nil {
		lc = pkg.PackageLicense.LicenseId.Val
		return lc

	}
	return lc
}

func InsertSupplier(a *spdx.Package2_1) string {
	if a.PackageSupplierPerson != "" {
		return ("Person: " + a.PackageSupplierPerson)
	}
	if a.PackageSupplierOrganization != "" {
		return ("Organization: " + a.PackageSupplierPerson)
	}
	return ""
}

func InsertOriginator(a *spdx.Package2_1) string {
	if a.PackageOriginatorPerson != "" {
		return ("Person: " + a.PackageOriginatorPerson)
	}
	if a.PackageOriginatorOrganization != "" {
		return ("Organization: " + a.PackageOriginatorPerson)
	}
	return ""
}

func RelationshipRef(rel *Relationship) string {
	var ref string

	for _, a := range rel.Package {
		if a != nil {
			ref = a.PackageSPDXIdentifier.Val
		}
	}
	for _, a := range rel.File {
		if a != nil {
			ref = a.FileSPDXIdentifier.Val
		}
	}
	if rel.SpdxElement != nil {
		ref = rel.SpdxElement.SpdxElement.Val
	}
	if rel.RelatedSpdxElement.Val != "" {
		ref = rel.RelatedSpdxElement.Val
	}

	return ref
}

func PkgChecksumAlgo(pkg2_1 *spdx.Package2_1) string {
	if pkg2_1.PackageChecksumSHA1 != "" {
		return "SHA1"
	}
	if pkg2_1.PackageChecksumSHA256 != "" {
		return "SHA256"
	}
	if pkg2_1.PackageChecksumMD5 != "" {
		return "MD5"
	}
	return ""
}

func PkgChecksumValue(pkg2_1 *spdx.Package2_1) string {
	if pkg2_1.PackageChecksumSHA1 != "" {
		return pkg2_1.PackageChecksumSHA1
	}
	if pkg2_1.PackageChecksumSHA256 != "" {
		return pkg2_1.PackageChecksumSHA256
	}
	if pkg2_1.PackageChecksumMD5 != "" {
		return pkg2_1.PackageChecksumMD5
	}
	return ""
}

func FileChecksumAlgo(File2_1 *spdx.File2_1) string {
	if File2_1.FileChecksumSHA1 != "" {
		return "SHA1"
	}
	if File2_1.FileChecksumSHA256 != "" {
		return "SHA256"
	}
	if File2_1.FileChecksumMD5 != "" {
		return "MD5"
	}
	return ""
}

func FileChecksumValue(File2_1 *spdx.File2_1) string {
	if File2_1.FileChecksumSHA1 != "" {
		return File2_1.FileChecksumSHA1
	}
	if File2_1.FileChecksumSHA256 != "" {
		return File2_1.FileChecksumSHA256
	}
	if File2_1.FileChecksumMD5 != "" {
		return File2_1.FileChecksumMD5
	}
	return ""
}

func convertPointertoInt(str string) int {
	value, _ := strconv.Atoi(str)
	return value
}
