// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import "github.com/spdx/gordf/rdfloader/parser"

var (
	// NAMESPACES
	NS_SPDX = "http://spdx.org/rdf/terms#"
	NS_RDFS = "http://www.w3.org/2000/01/rdf-schema#"
	NS_RDF  = parser.RDFNS
	NS_PTR  = "http://www.w3.org/2009/pointers#"
	NS_DOAP = "http://usefulinc.com/ns/doap#"

	// SPDX properties
	SPDX_SPEC_VERSION                            = NS_SPDX + "specVersion"
	SPDX_DATA_LICENSE                            = NS_SPDX + "dataLicense"
	SPDX_NAME                                    = NS_SPDX + "name"
	SPDX_EXTERNAL_DOCUMENT_REF                   = NS_SPDX + "externalDocumentRef"
	SPDX_LICENSE_LIST_VERSION                    = NS_SPDX + "licenseListVersion"
	SPDX_CREATOR                                 = NS_SPDX + "creator"
	SPDX_CREATED                                 = NS_SPDX + "created"
	SPDX_REVIEWED                                = NS_SPDX + "reviewed"
	SPDX_DESCRIBES_PACKAGE                       = NS_SPDX + "describesPackage"
	SPDX_HAS_EXTRACTED_LICENSING_INFO            = NS_SPDX + "hasExtractedLicensingInfo"
	SPDX_RELATIONSHIP                            = NS_SPDX + "relationship"
	SPDX_ANNOTATION                              = NS_SPDX + "annotation"
	SPDX_COMMENT                                 = NS_SPDX + "comment"
	SPDX_CREATION_INFO                           = NS_SPDX + "creationInfo"
	SPDX_CHECKSUM_ALGORITHM_SHA1                 = NS_SPDX + "checksumAlgorithm_sha1"
	SPDX_CHECKSUM_ALGORITHM_SHA256               = NS_SPDX + "checksumAlgorithm_sha256"
	SPDX_CHECKSUM_ALGORITHM_MD5                  = NS_SPDX + "checksumAlgorithm_md5"
	SPDX_EXTERNAL_DOCUMENT_ID                    = NS_SPDX + "externalDocumentId"
	SPDX_SPDX_DOCUMENT                           = NS_SPDX + "spdxDocument"
	SPDX_SPDX_DOCUMENT_CAPITALIZED               = NS_SPDX + "SpdxDocument"
	SPDX_CHECKSUM                                = NS_SPDX + "checksum"
	SPDX_CHECKSUM_CAPITALIZED                    = NS_SPDX + "Checksum"
	SPDX_ANNOTATION_TYPE                         = NS_SPDX + "annotationType"
	SPDX_ANNOTATION_TYPE_OTHER                   = NS_SPDX + "annotationType_other"
	SPDX_ANNOTATION_TYPE_REVIEW                  = NS_SPDX + "annotationType_review"
	SPDX_LICENSE_INFO_IN_FILE                    = NS_SPDX + "licenseInfoInFile"
	SPDX_LICENSE_CONCLUDED                       = NS_SPDX + "licenseConcluded"
	SPDX_LICENSE_COMMENTS                        = NS_SPDX + "licenseComments"
	SPDX_COPYRIGHT_TEXT                          = NS_SPDX + "copyrightText"
	SPDX_ARTIFACT_OF                             = NS_SPDX + "artifactOf"
	SPDX_NOTICE_TEXT                             = NS_SPDX + "noticeText"
	SPDX_FILE_CONTRIBUTOR                        = NS_SPDX + "fileContributor"
	SPDX_FILE_DEPENDENCY                         = NS_SPDX + "fileDependency"
	SPDX_FILE_TYPE                               = NS_SPDX + "fileType"
	SPDX_FILE_NAME                               = NS_SPDX + "fileName"
	SPDX_EXTRACTED_TEXT                          = NS_SPDX + "extractedText"
	SPDX_LICENSE_ID                              = NS_SPDX + "licenseId"
	SPDX_FILE                                    = NS_SPDX + "File"
	SPDX_PACKAGE                                 = NS_SPDX + "Package"
	SPDX_SPDX_ELEMENT                            = NS_SPDX + "SpdxElement"
	SPDX_VERSION_INFO                            = NS_SPDX + "versionInfo"
	SPDX_PACKAGE_FILE_NAME                       = NS_SPDX + "packageFileName"
	SPDX_SUPPLIER                                = NS_SPDX + "supplier"
	SPDX_ORIGINATOR                              = NS_SPDX + "originator"
	SPDX_DOWNLOAD_LOCATION                       = NS_SPDX + "downloadLocation"
	SPDX_FILES_ANALYZED                          = NS_SPDX + "filesAnalyzed"
	SPDX_PACKAGE_VERIFICATION_CODE               = NS_SPDX + "packageVerificationCode"
	SPDX_SOURCE_INFO                             = NS_SPDX + "sourceInfo"
	SPDX_LICENSE_INFO_FROM_FILES                 = NS_SPDX + "licenseInfoFromFiles"
	SPDX_LICENSE_DECLARED                        = NS_SPDX + "licenseDeclared"
	SPDX_SUMMARY                                 = NS_SPDX + "summary"
	SPDX_DESCRIPTION                             = NS_SPDX + "description"
	SPDX_EXTERNAL_REF                            = NS_SPDX + "externalRef"
	SPDX_HAS_FILE                                = NS_SPDX + "hasFile"
	SPDX_PRIMARY_PACKAGE_PURPOSE                 = NS_SPDX + "primaryPackagePurpose"
	SPDX_RELEASE_DATE                            = NS_SPDX + "releaseDate"
	SPDX_BUILT_DATE                              = NS_SPDX + "builtDate"
	SPDX_VALID_UNTIL_DATE                        = NS_SPDX + "validUntilDate"
	SPDX_ATTRIBUTION_TEXT                        = NS_SPDX + "attributionText"
	SPDX_PACKAGE_VERIFICATION_CODE_VALUE         = NS_SPDX + "packageVerificationCodeValue"
	SPDX_PACKAGE_VERIFICATION_CODE_EXCLUDED_FILE = NS_SPDX + "packageVerificationCodeExcludedFile"
	SPDX_RELATED_SPDX_ELEMENT                    = NS_SPDX + "relatedSpdxElement"
	SPDX_RELATIONSHIP_TYPE                       = NS_SPDX + "relationshipType"
	SPDX_SNIPPET_FROM_FILE                       = NS_SPDX + "snippetFromFile"
	SPDX_LICENSE_INFO_IN_SNIPPET                 = NS_SPDX + "licenseInfoInSnippet"
	SPDX_RANGE                                   = NS_SPDX + "range"
	SPDX_REVIEWER                                = NS_SPDX + "reviewer"
	SPDX_REVIEW_DATE                             = NS_SPDX + "reviewDate"
	SPDX_SNIPPET                                 = NS_SPDX + "Snippet"
	SPDX_ALGORITHM                               = NS_SPDX + "algorithm"
	SPDX_CHECKSUM_VALUE                          = NS_SPDX + "checksumValue"
	SPDX_REFERENCE_CATEGORY                      = NS_SPDX + "referenceCategory"
	SPDX_REFERENCE_CATEGORY_PACKAGE_MANAGER      = NS_SPDX + "referenceCategory_packageManager"
	SPDX_REFERENCE_CATEGORY_SECURITY             = NS_SPDX + "referenceCategory_security"
	SPDX_REFERENCE_CATEGORY_OTHER                = NS_SPDX + "referenceCategory_other"

	SPDX_REFERENCE_TYPE                   = NS_SPDX + "referenceType"
	SPDX_REFERENCE_LOCATOR                = NS_SPDX + "referenceLocator"
	SPDX_ANNOTATION_DATE                  = NS_SPDX + "annotationDate"
	SPDX_ANNOTATOR                        = NS_SPDX + "annotator"
	SPDX_MEMBER                           = NS_SPDX + "member"
	SPDX_DISJUNCTIVE_LICENSE_SET          = NS_SPDX + "DisjunctiveLicenseSet"
	SPDX_CONJUNCTIVE_LICENSE_SET          = NS_SPDX + "ConjunctiveLicenseSet"
	SPDX_EXTRACTED_LICENSING_INFO         = NS_SPDX + "ExtractedLicensingInfo"
	SPDX_SIMPLE_LICENSING_INFO            = NS_SPDX + "SimpleLicensingInfo"
	SPDX_NONE_CAPS                        = NS_SPDX + "NONE"
	SPDX_NOASSERTION_CAPS                 = NS_SPDX + "NOASSERTION"
	SPDX_NONE_SMALL                       = NS_SPDX + "none"
	SPDX_NOASSERTION_SMALL                = NS_SPDX + "noassertion"
	SPDX_LICENSE                          = NS_SPDX + "License"
	SPDX_LISTED_LICENSE                   = NS_SPDX + "ListedLicense"
	SPDX_EXAMPLE                          = NS_SPDX + "example"
	SPDX_IS_OSI_APPROVED                  = NS_SPDX + "isOsiApproved"
	SPDX_STANDARD_LICENSE_TEMPLATE        = NS_SPDX + "standardLicenseTemplate"
	SPDX_IS_DEPRECATED_LICENSE_ID         = NS_SPDX + "isDeprecatedLicenseId"
	SPDX_IS_FSF_LIBRE                     = NS_SPDX + "isFsfLibre"
	SPDX_LICENSE_TEXT                     = NS_SPDX + "licenseText"
	SPDX_STANDARD_LICENSE_HEADER          = NS_SPDX + "standardLicenseHeader"
	SPDX_LICENSE_EXCEPTION_ID             = NS_SPDX + "licenseExceptionId"
	SPDX_LICENSE_EXCEPTION_TEXT           = NS_SPDX + "licenseExceptionText"
	SPDX_LICENSE_EXCEPTION                = NS_SPDX + "licenseException"
	SPDX_WITH_EXCEPTION_OPERATOR          = NS_SPDX + "WithExceptionOperator"
	SPDX_OR_LATER_OPERATOR                = NS_SPDX + "OrLaterOperator"
	SPDX_STANDARD_LICENSE_HEADER_TEMPLATE = NS_SPDX + "standardLicenseHeaderTemplate"

	// RDFS properties
	RDFS_COMMENT  = NS_RDFS + "comment"
	RDFS_SEE_ALSO = NS_RDFS + "seeAlso"

	// RDF properties
	RDF_TYPE = NS_RDF + "type"

	// DOAP properties
	DOAP_HOMEPAGE = NS_DOAP + "homepage"
	DOAP_NAME     = NS_DOAP + "name"

	// PTR properties
	PTR_START_END_POINTER   = NS_PTR + "StartEndPointer"
	PTR_START_POINTER       = NS_PTR + "startPointer"
	PTR_BYTE_OFFSET_POINTER = NS_PTR + "ByteOffsetPointer"
	PTR_LINE_CHAR_POINTER   = NS_PTR + "LineCharPointer"
	PTR_REFERENCE           = NS_PTR + "reference"
	PTR_OFFSET              = NS_PTR + "offset"
	PTR_LINE_NUMBER         = NS_PTR + "lineNumber"
	PTR_END_POINTER         = NS_PTR + "endPointer"

	// prefixes
	PREFIX_RELATIONSHIP_TYPE = "relationshipType_"
)

func AllRelationshipTypes() []string {
	return []string{
		"amendment", "ancestorOf", "buildDependencyOf", "buildToolOf",
		"containedBy", "contains", "copyOf", "dataFile", "dataFileOf",
		"dependencyManifestOf", "dependencyOf", "dependsOn", "descendantOf",
		"describedBy", "describes", "devDependencyOf", "devToolOf",
		"distributionArtifact", "documentation", "dynamicLink", "exampleOf",
		"expandedFromArchive", "fileAdded", "fileDeleted", "fileModified",
		"generatedFrom", "generates", "hasPrerequisite", "metafileOf",
		"optionalComponentOf", "optionalDependencyOf", "other", "packageOf",
		"patchApplied", "patchFor", "prerequisiteFor", "providedDependencyOf",
		"runtimeDependencyOf", "staticLink", "testDependencyOf", "testOf",
		"testToolOf", "testcaseOf", "variantOf",
	}
}

func AllStandardLicenseIDS() []string {
	return []string{
		"0BSD", "389-exception", "AAL", "Abstyles", "Adobe-2006", "Adobe-Glyph",
		"ADSL", "AFL-1.1", "AFL-1.2", "AFL-2.0", "AFL-2.1", "AFL-3.0", "Afmparse",
		"AGPL-1.0-only", "AGPL-1.0-or-later", "AGPL-1.0", "AGPL-3.0-only",
		"AGPL-3.0-or-later", "AGPL-3.0", "Aladdin", "AMDPLPA", "AML", "AMPAS",
		"ANTLR-PD", "Apache-1.0", "Apache-1.1", "Apache-2.0", "APAFML", "APL-1.0",
		"APSL-1.0", "APSL-1.1", "APSL-1.2", "APSL-2.0", "Artistic-1.0-cl8",
		"Artistic-1.0-Perl", "Artistic-1.0", "Artistic-2.0", "",
		"Autoconf-exception-2.0", "Autoconf-exception-3.0", "Bahyph", "Barr",
		"Beerware", "Bison-exception-2.2", "BitTorrent-1.0", "BitTorrent-1.1",
		"blessing", "BlueOak-1.0.0", "Bootloader-exception", "Borceux", "BSD-1-Clause",
		"BSD-2-Clause-FreeBSD", "BSD-2-Clause-NetBSD", "BSD-2-Clause-Patent",
		"BSD-2-Clause-Views", "BSD-2-Clause", "BSD-3-Clause-Attribution",
		"BSD-3-Clause-Clear", "BSD-3-Clause-LBNL",
		"BSD-3-Clause-No-Nuclear-License-2014", "BSD-3-Clause-No-Nuclear-License",
		"BSD-3-Clause-No-Nuclear-Warranty", "BSD-3-Clause-Open-MPI", "BSD-3-Clause",
		"BSD-4-Clause-UC", "BSD-4-Clause", "BSD-Protection", "BSD-Source-Code",
		"BSL-1.0", "bzip2-1.0.5", "bzip2-1.0.6", "CAL-1.0-Combined-Work-Exception",
		"CAL-1.0", "Caldera", "CATOSL-1.1", "CC-BY-1.0", "CC-BY-2.0", "CC-BY-2.5",
		"CC-BY-3.0-AT", "CC-BY-3.0", "CC-BY-4.0", "CC-BY-NC-1.0", "CC-BY-NC-2.0",
		"CC-BY-NC-2.5", "CC-BY-NC-3.0", "CC-BY-NC-4.0", "CC-BY-NC-ND-1.0",
		"CC-BY-NC-ND-2.0", "CC-BY-NC-ND-2.5", "CC-BY-NC-ND-3.0-IGO", "CC-BY-NC-ND-3.0",
		"CC-BY-NC-ND-4.0", "CC-BY-NC-SA-1.0", "CC-BY-NC-SA-2.0", "CC-BY-NC-SA-2.5",
		"CC-BY-NC-SA-3.0", "CC-BY-NC-SA-4.0", "CC-BY-ND-1.0", "CC-BY-ND-2.0",
		"CC-BY-ND-2.5", "CC-BY-ND-3.0", "CC-BY-ND-4.0", "CC-BY-SA-1.0", "CC-BY-SA-2.0",
		"CC-BY-SA-2.5", "CC-BY-SA-3.0-AT", "CC-BY-SA-3.0", "CC-BY-SA-4.0", "CC-PDDC",
		"CC0-1.0", "CDDL-1.0", "CDDL-1.1", "CDLA-Permissive-1.0", "CDLA-Sharing-1.0",
		"CECILL-1.0", "CECILL-1.1", "CECILL-2.0", "CECILL-2.1", "CECILL-B", "CECILL-C",
		"CERN-OHL-1.1", "CERN-OHL-1.2", "CERN-OHL-P-2.0", "CERN-OHL-S-2.0",
		"CERN-OHL-W-2.0", "ClArtistic", "Classpath-exception-2.0",
		"CLISP-exception-2.0", "CNRI-Jython", "CNRI-Python-GPL-Compatible",
		"CNRI-Python", "Condor-1.1", "copyleft-next-0.3.0", "copyleft-next-0.3.1",
		"CPAL-1.0", "CPL-1.0", "CPOL-1.02", "Crossword", "CrystalStacker",
		"CUA-OPL-1.0", "Cube", "curl", "D-FSL-1.0", "diffmark",
		"DigiRule-FOSS-exception", "DOC", "Dotseqn", "DSDP", "dvipdfm", "ECL-1.0",
		"ECL-2.0", "eCos-2.0", "eCos-exception-2.0", "EFL-1.0", "EFL-2.0", "eGenix",
		"Entessa", "EPICS", "EPL-1.0", "EPL-2.0", "ErlPL-1.1", "etalab-2.0",
		"EUDatagrid", "EUPL-1.0", "EUPL-1.1", "EUPL-1.2", "Eurosym", "Fair",
		"Fawkes-Runtime-exception", "FLTK-exception", "Font-exception-2.0",
		"Frameworx-1.0", "FreeImage", "freertos-exception-2.0", "FSFAP", "FSFUL",
		"FSFULLR", "FTL", "GCC-exception-2.0", "GCC-exception-3.1",
		"GFDL-1.1-invariants-only", "GFDL-1.1-invariants-or-later",
		"GFDL-1.1-no-invariants-only", "GFDL-1.1-no-invariants-or-later",
		"GFDL-1.1-only", "GFDL-1.1-or-later", "GFDL-1.1", "GFDL-1.2-invariants-only",
		"GFDL-1.2-invariants-or-later", "GFDL-1.2-no-invariants-only",
		"GFDL-1.2-no-invariants-or-later", "GFDL-1.2-only", "GFDL-1.2-or-later",
		"GFDL-1.2", "GFDL-1.3-invariants-only", "GFDL-1.3-invariants-or-later",
		"GFDL-1.3-no-invariants-only", "GFDL-1.3-no-invariants-or-later",
		"GFDL-1.3-only", "GFDL-1.3-or-later", "GFDL-1.3", "Giftware", "GL2PS", "Glide",
		"Glulxe", "GLWTPL", "gnu-javamail-exception", "gnuplot", "GPL-1.0+",
		"GPL-1.0-only", "GPL-1.0-or-later", "GPL-1.0", "GPL-2.0+", "GPL-2.0-only",
		"GPL-2.0-or-later", "GPL-2.0-with-autoconf-exception",
		"GPL-2.0-with-bison-exception", "GPL-2.0-with-classpath-exception",
		"GPL-2.0-with-font-exception", "GPL-2.0-with-GCC-exception", "GPL-2.0",
		"GPL-3.0+", "GPL-3.0-linking-exception", "GPL-3.0-linking-source-exception",
		"GPL-3.0-only", "GPL-3.0-or-later", "GPL-3.0-with-autoconf-exception",
		"GPL-3.0-with-GCC-exception", "GPL-3.0", "GPL-CC-1.0", "gSOAP-1.3b",
		"HaskellReport", "Hippocratic-2.1", "HPND-sell-variant", "HPND",
		"i2p-gpl-java-exception", "IBM-pibs", "ICU", "IJG", "ImageMagick", "iMatix",
		"Imlib2", "Info-ZIP", "Intel-ACPI", "Intel", "Interbase-1.0", "IPA", "IPL-1.0",
		"ISC", "JasPer-2.0", "JPNIC", "JSON", "LAL-1.2", "LAL-1.3", "Latex2e",
		"Leptonica", "LGPL-2.0+", "LGPL-2.0-only", "LGPL-2.0-or-later", "LGPL-2.0",
		"LGPL-2.1+", "LGPL-2.1-only", "LGPL-2.1-or-later", "LGPL-2.1", "LGPL-3.0+",
		"LGPL-3.0-linking-exception", "LGPL-3.0-only", "LGPL-3.0-or-later", "LGPL-3.0",
		"LGPLLR", "libpng-2.0", "Libpng", "libselinux-1.0", "libtiff",
		"Libtool-exception", "licenses", "LiLiQ-P-1.1", "LiLiQ-R-1.1",
		"LiLiQ-Rplus-1.1", "Linux-OpenIB", "Linux-syscall-note", "LLVM-exception",
		"LPL-1.0", "LPL-1.02", "LPPL-1.0", "LPPL-1.1", "LPPL-1.2", "LPPL-1.3a",
		"LPPL-1.3c", "LZMA-exception", "MakeIndex", "mif-exception", "MirOS", "MIT-0",
		"MIT-advertising", "MIT-CMU", "MIT-enna", "MIT-feh", "MIT", "MITNFA",
		"Motosoto", "mpich2", "MPL-1.0", "MPL-1.1", "MPL-2.0-no-copyleft-exception",
		"MPL-2.0", "MS-PL", "MS-RL", "MTLL", "MulanPSL-1.0", "MulanPSL-2.0", "Multics",
		"Mup", "NASA-1.3", "Naumen", "NBPL-1.0", "NCGL-UK-2.0", "NCSA", "Net-SNMP",
		"NetCDF", "Newsletr", "NGPL", "NIST-PD-fallback", "NIST-PD", "NLOD-1.0",
		"NLPL", "Nokia-Qt-exception-1.1", "Nokia", "NOSL", "Noweb", "NPL-1.0",
		"NPL-1.1", "NPOSL-3.0", "NRL", "NTP-0", "NTP", "Nunit", "O-UDA-1.0",
		"OCaml-LGPL-linking-exception", "OCCT-exception-1.0", "OCCT-PL", "OCLC-2.0",
		"ODbL-1.0", "ODC-By-1.0", "OFL-1.0-no-RFN", "OFL-1.0-RFN", "OFL-1.0",
		"OFL-1.1-no-RFN", "OFL-1.1-RFN", "OFL-1.1", "OGC-1.0", "OGL-Canada-2.0",
		"OGL-UK-1.0", "OGL-UK-2.0", "OGL-UK-3.0", "OGTSL", "OLDAP-1.1", "OLDAP-1.2",
		"OLDAP-1.3", "OLDAP-1.4", "OLDAP-2.0.1", "OLDAP-2.0", "OLDAP-2.1",
		"OLDAP-2.2.1", "OLDAP-2.2.2", "OLDAP-2.2", "OLDAP-2.3", "OLDAP-2.4",
		"OLDAP-2.5", "OLDAP-2.6", "OLDAP-2.7", "OLDAP-2.8", "OML", "",
		"OpenJDK-assembly-exception-1.0", "OpenSSL", "openvpn-openssl-exception",
		"OPL-1.0", "OSET-PL-2.1", "OSL-1.0", "OSL-1.1", "OSL-2.0", "OSL-2.1",
		"OSL-3.0", "Parity-6.0.0", "Parity-7.0.0", "PDDL-1.0", "PHP-3.0", "PHP-3.01",
		"Plexus", "PolyForm-Noncommercial-1.0.0", "PolyForm-Small-Business-1.0.0",
		"PostgreSQL", "PS-or-PDF-font-exception-20170817", "PSF-2.0", "psfrag",
		"psutils", "Python-2.0", "Qhull", "QPL-1.0", "Qt-GPL-exception-1.0",
		"Qt-LGPL-exception-1.1", "Qwt-exception-1.0", "Rdisc", "RHeCos-1.1", "RPL-1.1",
		"RPL-1.5", "RPSL-1.0", "RSA-MD", "RSCPL", "Ruby", "SAX-PD", "Saxpath", "SCEA",
		"Sendmail-8.23", "Sendmail", "SGI-B-1.0", "SGI-B-1.1", "SGI-B-2.0", "SHL-0.5",
		"SHL-0.51", "SHL-2.0", "SHL-2.1", "SimPL-2.0", "SISSL-1.2", "SISSL",
		"Sleepycat", "SMLNJ", "SMPPL", "SNIA", "Spencer-86", "Spencer-94",
		"Spencer-99", "SPL-1.0", "SSH-OpenSSH", "SSH-short", "SSPL-1.0",
		"StandardML-NJ", "SugarCRM-1.1.3", "Swift-exception", "SWL", "TAPR-OHL-1.0",
		"TCL", "TCP-wrappers", "TMate", "TORQUE-1.1", "TOSL", "TU-Berlin-1.0",
		"TU-Berlin-2.0", "u-boot-exception-2.0", "UCL-1.0", "Unicode-DFS-2015",
		"Unicode-DFS-2016", "Unicode-TOU", "Universal-FOSS-exception-1.0", "Unlicense",
		"UPL-1.0", "Vim", "VOSTROM", "VSL-1.0", "W3C-19980720", "W3C-20150513", "W3C",
		"Watcom-1.0", "Wsuipa", "WTFPL", "WxWindows-exception-3.1", "wxWindows", "X11",
		"Xerox", "XFree86-1.1", "xinetd", "Xnet", "xpp", "XSkat", "YPL-1.0", "YPL-1.1",
		"Zed", "Zend-2.0", "Zimbra-1.3", "Zimbra-1.4", "zlib-acknowledgement", "Zlib",
		"ZPL-1.1", "ZPL-2.0", "ZPL-2.1",
	}
}
