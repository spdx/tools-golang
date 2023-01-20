// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package convert

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_1"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

func Test_ConvertSPDXDocuments(t *testing.T) {
	tests := []struct {
		name     string
		source   interface{}
		expected interface{}
	}{
		{
			name: "basic v2_2 to v2_3",
			source: v2_2.Document{
				SPDXVersion: v2_2.Version,
				Packages: []*v2_2.Package{
					{
						PackageName: "Pkg 1",
						Files: []*v2_2.File{
							{
								FileName: "File 1",
							},
							{
								FileName: "File 2",
							},
						},
						PackageVerificationCode: common.PackageVerificationCode{
							Value: "verification code value",
							ExcludedFiles: []string{
								"a",
								"b",
							},
						},
					},
				},
			},
			expected: spdx.Document{
				SPDXVersion: spdx.Version,
				Packages: []*spdx.Package{
					{
						PackageName: "Pkg 1",
						Files: []*spdx.File{
							{
								FileName: "File 1",
							},
							{
								FileName: "File 2",
							},
						},
						PackageVerificationCode: &common.PackageVerificationCode{
							Value: "verification code value",
							ExcludedFiles: []string{
								"a",
								"b",
							},
						},
					},
				},
			},
		},
		{
			name: "full 2.1 -> 2.3 document",
			source: v2_1.Document{
				SPDXVersion:       "SPDX-2.2",
				DataLicense:       "data license",
				SPDXIdentifier:    "spdx id",
				DocumentName:      "doc name",
				DocumentNamespace: "doc namespace",
				ExternalDocumentReferences: []v2_1.ExternalDocumentRef{
					{
						DocumentRefID: "doc ref id 1",
						URI:           "uri 1",
						Checksum: common.Checksum{
							Algorithm: "algo 1",
							Value:     "value 1",
						},
					},
					{
						DocumentRefID: "doc ref id 2",
						URI:           "uri 2",
						Checksum: common.Checksum{
							Algorithm: "algo 2",
							Value:     "value 2",
						},
					},
				},
				DocumentComment: "doc comment",
				CreationInfo: &v2_1.CreationInfo{
					LicenseListVersion: "license list version",
					Creators: []common.Creator{
						{
							Creator:     "creator 1",
							CreatorType: "type 1",
						},
						{
							Creator:     "creator 2",
							CreatorType: "type 2",
						},
					},
					Created:        "created date",
					CreatorComment: "creator comment",
				},
				Packages: []*v2_1.Package{
					{
						PackageName:               "package name 1",
						PackageSPDXIdentifier:     "id 1",
						PackageVersion:            "version 1",
						PackageFileName:           "file 1",
						PackageSupplier:           nil,
						PackageOriginator:         nil,
						PackageDownloadLocation:   "",
						FilesAnalyzed:             true,
						IsFilesAnalyzedTagPresent: true,
						PackageVerificationCode: common.PackageVerificationCode{
							Value:         "value 1",
							ExcludedFiles: []string{"a", "b"},
						},
						PackageChecksums: []common.Checksum{
							{
								Algorithm: "alg 1",
								Value:     "val 1",
							},
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
						},
						PackageHomePage:             "home page 1",
						PackageSourceInfo:           "source info 1",
						PackageLicenseConcluded:     "license concluded 1",
						PackageLicenseInfoFromFiles: []string{"a", "b"},
						PackageLicenseDeclared:      "license declared 1",
						PackageLicenseComments:      "license comments 1",
						PackageCopyrightText:        "copyright text 1",
						PackageSummary:              "summary 1",
						PackageDescription:          "description 1",
						PackageComment:              "comment 1",
						PackageExternalReferences: []*v2_1.PackageExternalReference{
							{
								Category:           "cat 1",
								RefType:            "type 1",
								Locator:            "locator 1",
								ExternalRefComment: "comment 1",
							},
							{
								Category:           "cat 2",
								RefType:            "type 2",
								Locator:            "locator 2",
								ExternalRefComment: "comment 2",
							},
						},
						Files: []*v2_1.File{
							{
								FileName:           "file 1",
								FileSPDXIdentifier: "id 1",
								FileTypes:          []string{"a", "b"},
								Checksums: []common.Checksum{
									{
										Algorithm: "alg 1",
										Value:     "val 1",
									},
									{
										Algorithm: "alg 2",
										Value:     "val 2",
									},
								},
								LicenseConcluded:   "license concluded 1",
								LicenseInfoInFiles: []string{"f1", "f2", "f3"},
								LicenseComments:    "comments 1",
								FileCopyrightText:  "copy text 1",
								ArtifactOfProjects: []*v2_1.ArtifactOfProject{
									{
										Name:     "name 1",
										HomePage: "home 1",
										URI:      "uri 1",
									},
									{
										Name:     "name 2",
										HomePage: "home 2",
										URI:      "uri 2",
									},
								},
								FileComment:      "comment 1",
								FileNotice:       "notice 1",
								FileContributors: []string{"c1", "c2"},
								FileDependencies: []string{"dep1", "dep2", "dep3"},
								Snippets: map[common.ElementID]*v2_1.Snippet{
									common.ElementID("e1"): {
										SnippetSPDXIdentifier:         "id1",
										SnippetFromFileSPDXIdentifier: "file1",
										Ranges: []common.SnippetRange{
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             1,
													LineNumber:         2,
													FileSPDXIdentifier: "f1",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             3,
													LineNumber:         4,
													FileSPDXIdentifier: "f2",
												},
											},
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             5,
													LineNumber:         6,
													FileSPDXIdentifier: "f3",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             7,
													LineNumber:         8,
													FileSPDXIdentifier: "f4",
												},
											},
										},
										SnippetLicenseConcluded: "license 1",
										LicenseInfoInSnippet:    []string{"a", "b"},
										SnippetLicenseComments:  "license comment 1",
										SnippetCopyrightText:    "copy 1",
										SnippetComment:          "comment 1",
										SnippetName:             "name 1",
									},
									common.ElementID("e2"): {
										SnippetSPDXIdentifier:         "id2",
										SnippetFromFileSPDXIdentifier: "file2",
										Ranges: []common.SnippetRange{
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             5,
													LineNumber:         6,
													FileSPDXIdentifier: "f3",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             7,
													LineNumber:         8,
													FileSPDXIdentifier: "f4",
												},
											},
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             9,
													LineNumber:         10,
													FileSPDXIdentifier: "f13",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             11,
													LineNumber:         12,
													FileSPDXIdentifier: "f14",
												},
											},
										},
										SnippetLicenseConcluded: "license 1",
										LicenseInfoInSnippet:    []string{"a", "b"},
										SnippetLicenseComments:  "license comment 1",
										SnippetCopyrightText:    "copy 1",
										SnippetComment:          "comment 1",
										SnippetName:             "name 1",
									},
								},
								Annotations: []v2_1.Annotation{
									{
										Annotator: common.Annotator{
											Annotator:     "ann 1",
											AnnotatorType: "typ 1",
										},
										AnnotationDate: "date 1",
										AnnotationType: "type 1",
										AnnotationSPDXIdentifier: common.DocElementID{
											DocumentRefID: "doc 1",
											ElementRefID:  "elem 1",
											SpecialID:     "spec 1",
										},
										AnnotationComment: "comment 1",
									},
									{
										Annotator: common.Annotator{
											Annotator:     "ann 2",
											AnnotatorType: "typ 2",
										},
										AnnotationDate: "date 2",
										AnnotationType: "type 2",
										AnnotationSPDXIdentifier: common.DocElementID{
											DocumentRefID: "doc 2",
											ElementRefID:  "elem 2",
											SpecialID:     "spec 2",
										},
										AnnotationComment: "comment 2",
									},
								},
							},
						},
						Annotations: []v2_1.Annotation{
							{
								Annotator: common.Annotator{
									Annotator:     "ann 1",
									AnnotatorType: "typ 1",
								},
								AnnotationDate: "date 1",
								AnnotationType: "type 1",
								AnnotationSPDXIdentifier: common.DocElementID{
									DocumentRefID: "doc 1",
									ElementRefID:  "elem 1",
									SpecialID:     "spec 1",
								},
								AnnotationComment: "comment 1",
							},
							{
								Annotator: common.Annotator{
									Annotator:     "ann 2",
									AnnotatorType: "typ 2",
								},
								AnnotationDate: "date 2",
								AnnotationType: "type 2",
								AnnotationSPDXIdentifier: common.DocElementID{
									DocumentRefID: "doc 2",
									ElementRefID:  "elem 2",
									SpecialID:     "spec 2",
								},
								AnnotationComment: "comment 2",
							},
						},
					},
				},
				Files: []*v2_1.File{
					{
						FileName:           "file 1",
						FileSPDXIdentifier: "id 1",
						FileTypes:          []string{"t1", "t2"},
						Checksums: []common.Checksum{
							{
								Algorithm: "alg 1",
								Value:     "val 1",
							},
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
						},
						LicenseConcluded:   "concluded 1",
						LicenseInfoInFiles: []string{"f1", "f2", "f3"},
						LicenseComments:    "comments 1",
						FileCopyrightText:  "copy 1",
						ArtifactOfProjects: []*v2_1.ArtifactOfProject{
							{
								Name:     "name 1",
								HomePage: "home 1",
								URI:      "uri 1",
							},
							{
								Name:     "name 2",
								HomePage: "home 2",
								URI:      "uri 2",
							},
						},
						FileComment:      "comment 1",
						FileNotice:       "notice 1",
						FileContributors: []string{"c1", "c2"},
						FileDependencies: []string{"d1", "d2", "d3", "d4"},
						Snippets:         nil, // already have snippets elsewhere
						Annotations:      nil, // already have annotations elsewhere
					},
					{
						FileName:           "file 2",
						FileSPDXIdentifier: "id 2",
						FileTypes:          []string{"t3", "t4"},
						Checksums: []common.Checksum{
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
							{
								Algorithm: "alg 3",
								Value:     "val 3",
							},
						},
						LicenseConcluded:   "concluded 2",
						LicenseInfoInFiles: []string{"f1", "f2", "f3"},
						LicenseComments:    "comments 2",
						FileCopyrightText:  "copy 2",
						ArtifactOfProjects: []*v2_1.ArtifactOfProject{
							{
								Name:     "name 2",
								HomePage: "home 2",
								URI:      "uri 2",
							},
							{
								Name:     "name 4",
								HomePage: "home 4",
								URI:      "uri 4",
							},
						},
						FileComment:      "comment 2",
						FileNotice:       "notice 2",
						FileContributors: []string{"c1", "c2"},
						FileDependencies: []string{"d1", "d2", "d3", "d4"},
						Snippets:         nil, // already have snippets elsewhere
						Annotations:      nil, // already have annotations elsewhere
					},
				},
				OtherLicenses: []*v2_1.OtherLicense{
					{
						LicenseIdentifier:      "id 1",
						ExtractedText:          "text 1",
						LicenseName:            "name 1",
						LicenseCrossReferences: []string{"x1", "x2", "x3"},
						LicenseComment:         "comment 1",
					},
					{
						LicenseIdentifier:      "id 2",
						ExtractedText:          "text 2",
						LicenseName:            "name 2",
						LicenseCrossReferences: []string{"x4", "x5", "x6"},
						LicenseComment:         "comment 2",
					},
				},
				Relationships: []*v2_1.Relationship{
					{
						RefA: common.DocElementID{
							DocumentRefID: "doc 1",
							ElementRefID:  "elem 1",
							SpecialID:     "spec 1",
						},
						RefB: common.DocElementID{
							DocumentRefID: "doc 2",
							ElementRefID:  "elem 2",
							SpecialID:     "spec 2",
						},
						Relationship:        "type 1",
						RelationshipComment: "comment 1",
					},
					{
						RefA: common.DocElementID{
							DocumentRefID: "doc 3",
							ElementRefID:  "elem 3",
							SpecialID:     "spec 3",
						},
						RefB: common.DocElementID{
							DocumentRefID: "doc 4",
							ElementRefID:  "elem 4",
							SpecialID:     "spec 4",
						},
						Relationship:        "type 2",
						RelationshipComment: "comment 2",
					},
				},
				Annotations: []*v2_1.Annotation{
					{
						Annotator: common.Annotator{
							Annotator:     "annotator 1",
							AnnotatorType: "annotator type 1",
						},
						AnnotationDate: "date 1",
						AnnotationType: "type 1",
						AnnotationSPDXIdentifier: common.DocElementID{
							DocumentRefID: "doc 1",
							ElementRefID:  "elem 1",
							SpecialID:     "spec 1",
						},
						AnnotationComment: "comment 1",
					},
					{
						Annotator: common.Annotator{
							Annotator:     "annotator 2",
							AnnotatorType: "annotator type 2",
						},
						AnnotationDate: "date 2",
						AnnotationType: "type 2",
						AnnotationSPDXIdentifier: common.DocElementID{
							DocumentRefID: "doc 2",
							ElementRefID:  "elem 2",
							SpecialID:     "spec 2",
						},
						AnnotationComment: "comment 2",
					},
				},
				Snippets: []v2_1.Snippet{
					{
						SnippetSPDXIdentifier:         "id1",
						SnippetFromFileSPDXIdentifier: "file1",
						Ranges: []common.SnippetRange{
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             1,
									LineNumber:         2,
									FileSPDXIdentifier: "f1",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             3,
									LineNumber:         4,
									FileSPDXIdentifier: "f2",
								},
							},
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             5,
									LineNumber:         6,
									FileSPDXIdentifier: "f3",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             7,
									LineNumber:         8,
									FileSPDXIdentifier: "f4",
								},
							},
						},
						SnippetLicenseConcluded: "license 1",
						LicenseInfoInSnippet:    []string{"a", "b"},
						SnippetLicenseComments:  "license comment 1",
						SnippetCopyrightText:    "copy 1",
						SnippetComment:          "comment 1",
						SnippetName:             "name 1",
					},
					{
						SnippetSPDXIdentifier:         "id2",
						SnippetFromFileSPDXIdentifier: "file2",
						Ranges: []common.SnippetRange{
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             5,
									LineNumber:         6,
									FileSPDXIdentifier: "f3",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             7,
									LineNumber:         8,
									FileSPDXIdentifier: "f4",
								},
							},
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             9,
									LineNumber:         10,
									FileSPDXIdentifier: "f13",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             11,
									LineNumber:         12,
									FileSPDXIdentifier: "f14",
								},
							},
						},
						SnippetLicenseConcluded: "license 1",
						LicenseInfoInSnippet:    []string{"a", "b"},
						SnippetLicenseComments:  "license comment 1",
						SnippetCopyrightText:    "copy 1",
						SnippetComment:          "comment 1",
						SnippetName:             "name 1",
					},
				},
				Reviews: []*v2_1.Review{
					{
						Reviewer:      "reviewer 1",
						ReviewerType:  "type 1",
						ReviewDate:    "date 1",
						ReviewComment: "comment 1",
					},
					{
						Reviewer:      "reviewer 2",
						ReviewerType:  "type 2",
						ReviewDate:    "date 2",
						ReviewComment: "comment 2",
					},
				},
			},
			expected: spdx.Document{
				SPDXVersion:       "SPDX-2.3", // ConvertFrom updates this value
				DataLicense:       "data license",
				SPDXIdentifier:    "spdx id",
				DocumentName:      "doc name",
				DocumentNamespace: "doc namespace",
				ExternalDocumentReferences: []spdx.ExternalDocumentRef{
					{
						DocumentRefID: "doc ref id 1",
						URI:           "uri 1",
						Checksum: common.Checksum{
							Algorithm: "algo 1",
							Value:     "value 1",
						},
					},
					{
						DocumentRefID: "doc ref id 2",
						URI:           "uri 2",
						Checksum: common.Checksum{
							Algorithm: "algo 2",
							Value:     "value 2",
						},
					},
				},
				DocumentComment: "doc comment",
				CreationInfo: &spdx.CreationInfo{
					LicenseListVersion: "license list version",
					Creators: []common.Creator{
						{
							Creator:     "creator 1",
							CreatorType: "type 1",
						},
						{
							Creator:     "creator 2",
							CreatorType: "type 2",
						},
					},
					Created:        "created date",
					CreatorComment: "creator comment",
				},
				Packages: []*spdx.Package{
					{
						IsUnpackaged:              true,
						PackageName:               "package name 1",
						PackageSPDXIdentifier:     "id 1",
						PackageVersion:            "version 1",
						PackageFileName:           "file 1",
						PackageSupplier:           nil,
						PackageOriginator:         nil,
						PackageDownloadLocation:   "",
						FilesAnalyzed:             true,
						IsFilesAnalyzedTagPresent: true,
						PackageVerificationCode: &common.PackageVerificationCode{
							Value:         "value 1",
							ExcludedFiles: []string{"a", "b"},
						},
						PackageChecksums: []common.Checksum{
							{
								Algorithm: "alg 1",
								Value:     "val 1",
							},
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
						},
						PackageHomePage:             "home page 1",
						PackageSourceInfo:           "source info 1",
						PackageLicenseConcluded:     "license concluded 1",
						PackageLicenseInfoFromFiles: []string{"a", "b"},
						PackageLicenseDeclared:      "license declared 1",
						PackageLicenseComments:      "license comments 1",
						PackageCopyrightText:        "copyright text 1",
						PackageSummary:              "summary 1",
						PackageDescription:          "description 1",
						PackageComment:              "comment 1",
						PackageExternalReferences: []*spdx.PackageExternalReference{
							{
								Category:           "cat 1",
								RefType:            "type 1",
								Locator:            "locator 1",
								ExternalRefComment: "comment 1",
							},
							{
								Category:           "cat 2",
								RefType:            "type 2",
								Locator:            "locator 2",
								ExternalRefComment: "comment 2",
							},
						},
						Files: []*spdx.File{
							{
								FileName:           "file 1",
								FileSPDXIdentifier: "id 1",
								FileTypes:          []string{"a", "b"},
								Checksums: []common.Checksum{
									{
										Algorithm: "alg 1",
										Value:     "val 1",
									},
									{
										Algorithm: "alg 2",
										Value:     "val 2",
									},
								},
								LicenseConcluded:   "license concluded 1",
								LicenseInfoInFiles: []string{"f1", "f2", "f3"},
								LicenseComments:    "comments 1",
								FileCopyrightText:  "copy text 1",
								ArtifactOfProjects: []*spdx.ArtifactOfProject{
									{
										Name:     "name 1",
										HomePage: "home 1",
										URI:      "uri 1",
									},
									{
										Name:     "name 2",
										HomePage: "home 2",
										URI:      "uri 2",
									},
								},
								FileComment:      "comment 1",
								FileNotice:       "notice 1",
								FileContributors: []string{"c1", "c2"},
								FileDependencies: []string{"dep1", "dep2", "dep3"},
								Snippets: map[common.ElementID]*spdx.Snippet{
									common.ElementID("e1"): {
										SnippetSPDXIdentifier:         "id1",
										SnippetFromFileSPDXIdentifier: "file1",
										Ranges: []common.SnippetRange{
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             1,
													LineNumber:         2,
													FileSPDXIdentifier: "f1",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             3,
													LineNumber:         4,
													FileSPDXIdentifier: "f2",
												},
											},
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             5,
													LineNumber:         6,
													FileSPDXIdentifier: "f3",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             7,
													LineNumber:         8,
													FileSPDXIdentifier: "f4",
												},
											},
										},
										SnippetLicenseConcluded: "license 1",
										LicenseInfoInSnippet:    []string{"a", "b"},
										SnippetLicenseComments:  "license comment 1",
										SnippetCopyrightText:    "copy 1",
										SnippetComment:          "comment 1",
										SnippetName:             "name 1",
									},
									common.ElementID("e2"): {
										SnippetSPDXIdentifier:         "id2",
										SnippetFromFileSPDXIdentifier: "file2",
										Ranges: []common.SnippetRange{
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             5,
													LineNumber:         6,
													FileSPDXIdentifier: "f3",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             7,
													LineNumber:         8,
													FileSPDXIdentifier: "f4",
												},
											},
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             9,
													LineNumber:         10,
													FileSPDXIdentifier: "f13",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             11,
													LineNumber:         12,
													FileSPDXIdentifier: "f14",
												},
											},
										},
										SnippetLicenseConcluded: "license 1",
										LicenseInfoInSnippet:    []string{"a", "b"},
										SnippetLicenseComments:  "license comment 1",
										SnippetCopyrightText:    "copy 1",
										SnippetComment:          "comment 1",
										SnippetName:             "name 1",
									},
								},
								Annotations: []spdx.Annotation{
									{
										Annotator: common.Annotator{
											Annotator:     "ann 1",
											AnnotatorType: "typ 1",
										},
										AnnotationDate: "date 1",
										AnnotationType: "type 1",
										AnnotationSPDXIdentifier: common.DocElementID{
											DocumentRefID: "doc 1",
											ElementRefID:  "elem 1",
											SpecialID:     "spec 1",
										},
										AnnotationComment: "comment 1",
									},
									{
										Annotator: common.Annotator{
											Annotator:     "ann 2",
											AnnotatorType: "typ 2",
										},
										AnnotationDate: "date 2",
										AnnotationType: "type 2",
										AnnotationSPDXIdentifier: common.DocElementID{
											DocumentRefID: "doc 2",
											ElementRefID:  "elem 2",
											SpecialID:     "spec 2",
										},
										AnnotationComment: "comment 2",
									},
								},
							},
						},
						Annotations: []spdx.Annotation{
							{
								Annotator: common.Annotator{
									Annotator:     "ann 1",
									AnnotatorType: "typ 1",
								},
								AnnotationDate: "date 1",
								AnnotationType: "type 1",
								AnnotationSPDXIdentifier: common.DocElementID{
									DocumentRefID: "doc 1",
									ElementRefID:  "elem 1",
									SpecialID:     "spec 1",
								},
								AnnotationComment: "comment 1",
							},
							{
								Annotator: common.Annotator{
									Annotator:     "ann 2",
									AnnotatorType: "typ 2",
								},
								AnnotationDate: "date 2",
								AnnotationType: "type 2",
								AnnotationSPDXIdentifier: common.DocElementID{
									DocumentRefID: "doc 2",
									ElementRefID:  "elem 2",
									SpecialID:     "spec 2",
								},
								AnnotationComment: "comment 2",
							},
						},
					},
				},
				Files: []*spdx.File{
					{
						FileName:           "file 1",
						FileSPDXIdentifier: "id 1",
						FileTypes:          []string{"t1", "t2"},
						Checksums: []common.Checksum{
							{
								Algorithm: "alg 1",
								Value:     "val 1",
							},
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
						},
						LicenseConcluded:   "concluded 1",
						LicenseInfoInFiles: []string{"f1", "f2", "f3"},
						LicenseComments:    "comments 1",
						FileCopyrightText:  "copy 1",
						ArtifactOfProjects: []*spdx.ArtifactOfProject{
							{
								Name:     "name 1",
								HomePage: "home 1",
								URI:      "uri 1",
							},
							{
								Name:     "name 2",
								HomePage: "home 2",
								URI:      "uri 2",
							},
						},
						FileComment:      "comment 1",
						FileNotice:       "notice 1",
						FileContributors: []string{"c1", "c2"},
						FileDependencies: []string{"d1", "d2", "d3", "d4"},
						Snippets:         nil, // already have snippets elsewhere
						Annotations:      nil, // already have annotations elsewhere
					},
					{
						FileName:           "file 2",
						FileSPDXIdentifier: "id 2",
						FileTypes:          []string{"t3", "t4"},
						Checksums: []common.Checksum{
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
							{
								Algorithm: "alg 3",
								Value:     "val 3",
							},
						},
						LicenseConcluded:   "concluded 2",
						LicenseInfoInFiles: []string{"f1", "f2", "f3"},
						LicenseComments:    "comments 2",
						FileCopyrightText:  "copy 2",
						ArtifactOfProjects: []*spdx.ArtifactOfProject{
							{
								Name:     "name 2",
								HomePage: "home 2",
								URI:      "uri 2",
							},
							{
								Name:     "name 4",
								HomePage: "home 4",
								URI:      "uri 4",
							},
						},
						FileComment:      "comment 2",
						FileNotice:       "notice 2",
						FileContributors: []string{"c1", "c2"},
						FileDependencies: []string{"d1", "d2", "d3", "d4"},
						Snippets:         nil, // already have snippets elsewhere
						Annotations:      nil, // already have annotations elsewhere
					},
				},
				OtherLicenses: []*spdx.OtherLicense{
					{
						LicenseIdentifier:      "id 1",
						ExtractedText:          "text 1",
						LicenseName:            "name 1",
						LicenseCrossReferences: []string{"x1", "x2", "x3"},
						LicenseComment:         "comment 1",
					},
					{
						LicenseIdentifier:      "id 2",
						ExtractedText:          "text 2",
						LicenseName:            "name 2",
						LicenseCrossReferences: []string{"x4", "x5", "x6"},
						LicenseComment:         "comment 2",
					},
				},
				Relationships: []*spdx.Relationship{
					{
						RefA: common.DocElementID{
							DocumentRefID: "doc 1",
							ElementRefID:  "elem 1",
							SpecialID:     "spec 1",
						},
						RefB: common.DocElementID{
							DocumentRefID: "doc 2",
							ElementRefID:  "elem 2",
							SpecialID:     "spec 2",
						},
						Relationship:        "type 1",
						RelationshipComment: "comment 1",
					},
					{
						RefA: common.DocElementID{
							DocumentRefID: "doc 3",
							ElementRefID:  "elem 3",
							SpecialID:     "spec 3",
						},
						RefB: common.DocElementID{
							DocumentRefID: "doc 4",
							ElementRefID:  "elem 4",
							SpecialID:     "spec 4",
						},
						Relationship:        "type 2",
						RelationshipComment: "comment 2",
					},
				},
				Annotations: []*spdx.Annotation{
					{
						Annotator: common.Annotator{
							Annotator:     "annotator 1",
							AnnotatorType: "annotator type 1",
						},
						AnnotationDate: "date 1",
						AnnotationType: "type 1",
						AnnotationSPDXIdentifier: common.DocElementID{
							DocumentRefID: "doc 1",
							ElementRefID:  "elem 1",
							SpecialID:     "spec 1",
						},
						AnnotationComment: "comment 1",
					},
					{
						Annotator: common.Annotator{
							Annotator:     "annotator 2",
							AnnotatorType: "annotator type 2",
						},
						AnnotationDate: "date 2",
						AnnotationType: "type 2",
						AnnotationSPDXIdentifier: common.DocElementID{
							DocumentRefID: "doc 2",
							ElementRefID:  "elem 2",
							SpecialID:     "spec 2",
						},
						AnnotationComment: "comment 2",
					},
				},
				Snippets: []spdx.Snippet{
					{
						SnippetSPDXIdentifier:         "id1",
						SnippetFromFileSPDXIdentifier: "file1",
						Ranges: []common.SnippetRange{
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             1,
									LineNumber:         2,
									FileSPDXIdentifier: "f1",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             3,
									LineNumber:         4,
									FileSPDXIdentifier: "f2",
								},
							},
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             5,
									LineNumber:         6,
									FileSPDXIdentifier: "f3",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             7,
									LineNumber:         8,
									FileSPDXIdentifier: "f4",
								},
							},
						},
						SnippetLicenseConcluded: "license 1",
						LicenseInfoInSnippet:    []string{"a", "b"},
						SnippetLicenseComments:  "license comment 1",
						SnippetCopyrightText:    "copy 1",
						SnippetComment:          "comment 1",
						SnippetName:             "name 1",
					},
					{
						SnippetSPDXIdentifier:         "id2",
						SnippetFromFileSPDXIdentifier: "file2",
						Ranges: []common.SnippetRange{
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             5,
									LineNumber:         6,
									FileSPDXIdentifier: "f3",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             7,
									LineNumber:         8,
									FileSPDXIdentifier: "f4",
								},
							},
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             9,
									LineNumber:         10,
									FileSPDXIdentifier: "f13",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             11,
									LineNumber:         12,
									FileSPDXIdentifier: "f14",
								},
							},
						},
						SnippetLicenseConcluded: "license 1",
						LicenseInfoInSnippet:    []string{"a", "b"},
						SnippetLicenseComments:  "license comment 1",
						SnippetCopyrightText:    "copy 1",
						SnippetComment:          "comment 1",
						SnippetName:             "name 1",
					},
				},
				Reviews: []*spdx.Review{
					{
						Reviewer:      "reviewer 1",
						ReviewerType:  "type 1",
						ReviewDate:    "date 1",
						ReviewComment: "comment 1",
					},
					{
						Reviewer:      "reviewer 2",
						ReviewerType:  "type 2",
						ReviewDate:    "date 2",
						ReviewComment: "comment 2",
					},
				},
			},
		},
		{
			name: "full 2.2 -> 2.3 document",
			source: v2_2.Document{
				SPDXVersion:       "SPDX-2.2",
				DataLicense:       "data license",
				SPDXIdentifier:    "spdx id",
				DocumentName:      "doc name",
				DocumentNamespace: "doc namespace",
				ExternalDocumentReferences: []v2_2.ExternalDocumentRef{
					{
						DocumentRefID: "doc ref id 1",
						URI:           "uri 1",
						Checksum: common.Checksum{
							Algorithm: "algo 1",
							Value:     "value 1",
						},
					},
					{
						DocumentRefID: "doc ref id 2",
						URI:           "uri 2",
						Checksum: common.Checksum{
							Algorithm: "algo 2",
							Value:     "value 2",
						},
					},
				},
				DocumentComment: "doc comment",
				CreationInfo: &v2_2.CreationInfo{
					LicenseListVersion: "license list version",
					Creators: []common.Creator{
						{
							Creator:     "creator 1",
							CreatorType: "type 1",
						},
						{
							Creator:     "creator 2",
							CreatorType: "type 2",
						},
					},
					Created:        "created date",
					CreatorComment: "creator comment",
				},
				Packages: []*v2_2.Package{
					{
						IsUnpackaged:              true,
						PackageName:               "package name 1",
						PackageSPDXIdentifier:     "id 1",
						PackageVersion:            "version 1",
						PackageFileName:           "file 1",
						PackageSupplier:           nil,
						PackageOriginator:         nil,
						PackageDownloadLocation:   "",
						FilesAnalyzed:             true,
						IsFilesAnalyzedTagPresent: true,
						PackageVerificationCode: common.PackageVerificationCode{
							Value:         "value 1",
							ExcludedFiles: []string{"a", "b"},
						},
						PackageChecksums: []common.Checksum{
							{
								Algorithm: "alg 1",
								Value:     "val 1",
							},
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
						},
						PackageHomePage:             "home page 1",
						PackageSourceInfo:           "source info 1",
						PackageLicenseConcluded:     "license concluded 1",
						PackageLicenseInfoFromFiles: []string{"a", "b"},
						PackageLicenseDeclared:      "license declared 1",
						PackageLicenseComments:      "license comments 1",
						PackageCopyrightText:        "copyright text 1",
						PackageSummary:              "summary 1",
						PackageDescription:          "description 1",
						PackageComment:              "comment 1",
						PackageExternalReferences: []*v2_2.PackageExternalReference{
							{
								Category:           "cat 1",
								RefType:            "type 1",
								Locator:            "locator 1",
								ExternalRefComment: "comment 1",
							},
							{
								Category:           "cat 2",
								RefType:            "type 2",
								Locator:            "locator 2",
								ExternalRefComment: "comment 2",
							},
						},
						PackageAttributionTexts: []string{"a", "b", "c"},
						Files: []*v2_2.File{
							{
								FileName:           "file 1",
								FileSPDXIdentifier: "id 1",
								FileTypes:          []string{"a", "b"},
								Checksums: []common.Checksum{
									{
										Algorithm: "alg 1",
										Value:     "val 1",
									},
									{
										Algorithm: "alg 2",
										Value:     "val 2",
									},
								},
								LicenseConcluded:   "license concluded 1",
								LicenseInfoInFiles: []string{"f1", "f2", "f3"},
								LicenseComments:    "comments 1",
								FileCopyrightText:  "copy text 1",
								ArtifactOfProjects: []*v2_2.ArtifactOfProject{
									{
										Name:     "name 1",
										HomePage: "home 1",
										URI:      "uri 1",
									},
									{
										Name:     "name 2",
										HomePage: "home 2",
										URI:      "uri 2",
									},
								},
								FileComment:          "comment 1",
								FileNotice:           "notice 1",
								FileContributors:     []string{"c1", "c2"},
								FileAttributionTexts: []string{"att1", "att2"},
								FileDependencies:     []string{"dep1", "dep2", "dep3"},
								Snippets: map[common.ElementID]*v2_2.Snippet{
									common.ElementID("e1"): {
										SnippetSPDXIdentifier:         "id1",
										SnippetFromFileSPDXIdentifier: "file1",
										Ranges: []common.SnippetRange{
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             1,
													LineNumber:         2,
													FileSPDXIdentifier: "f1",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             3,
													LineNumber:         4,
													FileSPDXIdentifier: "f2",
												},
											},
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             5,
													LineNumber:         6,
													FileSPDXIdentifier: "f3",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             7,
													LineNumber:         8,
													FileSPDXIdentifier: "f4",
												},
											},
										},
										SnippetLicenseConcluded: "license 1",
										LicenseInfoInSnippet:    []string{"a", "b"},
										SnippetLicenseComments:  "license comment 1",
										SnippetCopyrightText:    "copy 1",
										SnippetComment:          "comment 1",
										SnippetName:             "name 1",
										SnippetAttributionTexts: []string{"att1", "att2", "att3"},
									},
									common.ElementID("e2"): {
										SnippetSPDXIdentifier:         "id2",
										SnippetFromFileSPDXIdentifier: "file2",
										Ranges: []common.SnippetRange{
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             5,
													LineNumber:         6,
													FileSPDXIdentifier: "f3",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             7,
													LineNumber:         8,
													FileSPDXIdentifier: "f4",
												},
											},
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             9,
													LineNumber:         10,
													FileSPDXIdentifier: "f13",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             11,
													LineNumber:         12,
													FileSPDXIdentifier: "f14",
												},
											},
										},
										SnippetLicenseConcluded: "license 1",
										LicenseInfoInSnippet:    []string{"a", "b"},
										SnippetLicenseComments:  "license comment 1",
										SnippetCopyrightText:    "copy 1",
										SnippetComment:          "comment 1",
										SnippetName:             "name 1",
										SnippetAttributionTexts: []string{"att1", "att2", "att3"},
									},
								},
								Annotations: []v2_2.Annotation{
									{
										Annotator: common.Annotator{
											Annotator:     "ann 1",
											AnnotatorType: "typ 1",
										},
										AnnotationDate: "date 1",
										AnnotationType: "type 1",
										AnnotationSPDXIdentifier: common.DocElementID{
											DocumentRefID: "doc 1",
											ElementRefID:  "elem 1",
											SpecialID:     "spec 1",
										},
										AnnotationComment: "comment 1",
									},
									{
										Annotator: common.Annotator{
											Annotator:     "ann 2",
											AnnotatorType: "typ 2",
										},
										AnnotationDate: "date 2",
										AnnotationType: "type 2",
										AnnotationSPDXIdentifier: common.DocElementID{
											DocumentRefID: "doc 2",
											ElementRefID:  "elem 2",
											SpecialID:     "spec 2",
										},
										AnnotationComment: "comment 2",
									},
								},
							},
						},
						Annotations: []v2_2.Annotation{
							{
								Annotator: common.Annotator{
									Annotator:     "ann 1",
									AnnotatorType: "typ 1",
								},
								AnnotationDate: "date 1",
								AnnotationType: "type 1",
								AnnotationSPDXIdentifier: common.DocElementID{
									DocumentRefID: "doc 1",
									ElementRefID:  "elem 1",
									SpecialID:     "spec 1",
								},
								AnnotationComment: "comment 1",
							},
							{
								Annotator: common.Annotator{
									Annotator:     "ann 2",
									AnnotatorType: "typ 2",
								},
								AnnotationDate: "date 2",
								AnnotationType: "type 2",
								AnnotationSPDXIdentifier: common.DocElementID{
									DocumentRefID: "doc 2",
									ElementRefID:  "elem 2",
									SpecialID:     "spec 2",
								},
								AnnotationComment: "comment 2",
							},
						},
					},
				},
				Files: []*v2_2.File{
					{
						FileName:           "file 1",
						FileSPDXIdentifier: "id 1",
						FileTypes:          []string{"t1", "t2"},
						Checksums: []common.Checksum{
							{
								Algorithm: "alg 1",
								Value:     "val 1",
							},
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
						},
						LicenseConcluded:   "concluded 1",
						LicenseInfoInFiles: []string{"f1", "f2", "f3"},
						LicenseComments:    "comments 1",
						FileCopyrightText:  "copy 1",
						ArtifactOfProjects: []*v2_2.ArtifactOfProject{
							{
								Name:     "name 1",
								HomePage: "home 1",
								URI:      "uri 1",
							},
							{
								Name:     "name 2",
								HomePage: "home 2",
								URI:      "uri 2",
							},
						},
						FileComment:          "comment 1",
						FileNotice:           "notice 1",
						FileContributors:     []string{"c1", "c2"},
						FileAttributionTexts: []string{"att1", "att2"},
						FileDependencies:     []string{"d1", "d2", "d3", "d4"},
						Snippets:             nil, // already have snippets elsewhere
						Annotations:          nil, // already have annotations elsewhere
					},
					{
						FileName:           "file 2",
						FileSPDXIdentifier: "id 2",
						FileTypes:          []string{"t3", "t4"},
						Checksums: []common.Checksum{
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
							{
								Algorithm: "alg 3",
								Value:     "val 3",
							},
						},
						LicenseConcluded:   "concluded 2",
						LicenseInfoInFiles: []string{"f1", "f2", "f3"},
						LicenseComments:    "comments 2",
						FileCopyrightText:  "copy 2",
						ArtifactOfProjects: []*v2_2.ArtifactOfProject{
							{
								Name:     "name 2",
								HomePage: "home 2",
								URI:      "uri 2",
							},
							{
								Name:     "name 4",
								HomePage: "home 4",
								URI:      "uri 4",
							},
						},
						FileComment:          "comment 2",
						FileNotice:           "notice 2",
						FileContributors:     []string{"c1", "c2"},
						FileAttributionTexts: []string{"att1", "att2"},
						FileDependencies:     []string{"d1", "d2", "d3", "d4"},
						Snippets:             nil, // already have snippets elsewhere
						Annotations:          nil, // already have annotations elsewhere
					},
				},
				OtherLicenses: []*v2_2.OtherLicense{
					{
						LicenseIdentifier:      "id 1",
						ExtractedText:          "text 1",
						LicenseName:            "name 1",
						LicenseCrossReferences: []string{"x1", "x2", "x3"},
						LicenseComment:         "comment 1",
					},
					{
						LicenseIdentifier:      "id 2",
						ExtractedText:          "text 2",
						LicenseName:            "name 2",
						LicenseCrossReferences: []string{"x4", "x5", "x6"},
						LicenseComment:         "comment 2",
					},
				},
				Relationships: []*v2_2.Relationship{
					{
						RefA: common.DocElementID{
							DocumentRefID: "doc 1",
							ElementRefID:  "elem 1",
							SpecialID:     "spec 1",
						},
						RefB: common.DocElementID{
							DocumentRefID: "doc 2",
							ElementRefID:  "elem 2",
							SpecialID:     "spec 2",
						},
						Relationship:        "type 1",
						RelationshipComment: "comment 1",
					},
					{
						RefA: common.DocElementID{
							DocumentRefID: "doc 3",
							ElementRefID:  "elem 3",
							SpecialID:     "spec 3",
						},
						RefB: common.DocElementID{
							DocumentRefID: "doc 4",
							ElementRefID:  "elem 4",
							SpecialID:     "spec 4",
						},
						Relationship:        "type 2",
						RelationshipComment: "comment 2",
					},
				},
				Annotations: []*v2_2.Annotation{
					{
						Annotator: common.Annotator{
							Annotator:     "annotator 1",
							AnnotatorType: "annotator type 1",
						},
						AnnotationDate: "date 1",
						AnnotationType: "type 1",
						AnnotationSPDXIdentifier: common.DocElementID{
							DocumentRefID: "doc 1",
							ElementRefID:  "elem 1",
							SpecialID:     "spec 1",
						},
						AnnotationComment: "comment 1",
					},
					{
						Annotator: common.Annotator{
							Annotator:     "annotator 2",
							AnnotatorType: "annotator type 2",
						},
						AnnotationDate: "date 2",
						AnnotationType: "type 2",
						AnnotationSPDXIdentifier: common.DocElementID{
							DocumentRefID: "doc 2",
							ElementRefID:  "elem 2",
							SpecialID:     "spec 2",
						},
						AnnotationComment: "comment 2",
					},
				},
				Snippets: []v2_2.Snippet{
					{
						SnippetSPDXIdentifier:         "id1",
						SnippetFromFileSPDXIdentifier: "file1",
						Ranges: []common.SnippetRange{
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             1,
									LineNumber:         2,
									FileSPDXIdentifier: "f1",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             3,
									LineNumber:         4,
									FileSPDXIdentifier: "f2",
								},
							},
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             5,
									LineNumber:         6,
									FileSPDXIdentifier: "f3",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             7,
									LineNumber:         8,
									FileSPDXIdentifier: "f4",
								},
							},
						},
						SnippetLicenseConcluded: "license 1",
						LicenseInfoInSnippet:    []string{"a", "b"},
						SnippetLicenseComments:  "license comment 1",
						SnippetCopyrightText:    "copy 1",
						SnippetComment:          "comment 1",
						SnippetName:             "name 1",
						SnippetAttributionTexts: []string{"att1", "att2", "att3"},
					},
					{
						SnippetSPDXIdentifier:         "id2",
						SnippetFromFileSPDXIdentifier: "file2",
						Ranges: []common.SnippetRange{
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             5,
									LineNumber:         6,
									FileSPDXIdentifier: "f3",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             7,
									LineNumber:         8,
									FileSPDXIdentifier: "f4",
								},
							},
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             9,
									LineNumber:         10,
									FileSPDXIdentifier: "f13",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             11,
									LineNumber:         12,
									FileSPDXIdentifier: "f14",
								},
							},
						},
						SnippetLicenseConcluded: "license 1",
						LicenseInfoInSnippet:    []string{"a", "b"},
						SnippetLicenseComments:  "license comment 1",
						SnippetCopyrightText:    "copy 1",
						SnippetComment:          "comment 1",
						SnippetName:             "name 1",
						SnippetAttributionTexts: []string{"att1", "att2", "att3"},
					},
				},
				Reviews: []*v2_2.Review{
					{
						Reviewer:      "reviewer 1",
						ReviewerType:  "type 1",
						ReviewDate:    "date 1",
						ReviewComment: "comment 1",
					},
					{
						Reviewer:      "reviewer 2",
						ReviewerType:  "type 2",
						ReviewDate:    "date 2",
						ReviewComment: "comment 2",
					},
				},
			},
			expected: spdx.Document{
				SPDXVersion:       "SPDX-2.3", // ConvertFrom updates this value
				DataLicense:       "data license",
				SPDXIdentifier:    "spdx id",
				DocumentName:      "doc name",
				DocumentNamespace: "doc namespace",
				ExternalDocumentReferences: []spdx.ExternalDocumentRef{
					{
						DocumentRefID: "doc ref id 1",
						URI:           "uri 1",
						Checksum: common.Checksum{
							Algorithm: "algo 1",
							Value:     "value 1",
						},
					},
					{
						DocumentRefID: "doc ref id 2",
						URI:           "uri 2",
						Checksum: common.Checksum{
							Algorithm: "algo 2",
							Value:     "value 2",
						},
					},
				},
				DocumentComment: "doc comment",
				CreationInfo: &spdx.CreationInfo{
					LicenseListVersion: "license list version",
					Creators: []common.Creator{
						{
							Creator:     "creator 1",
							CreatorType: "type 1",
						},
						{
							Creator:     "creator 2",
							CreatorType: "type 2",
						},
					},
					Created:        "created date",
					CreatorComment: "creator comment",
				},
				Packages: []*spdx.Package{
					{
						IsUnpackaged:              true,
						PackageName:               "package name 1",
						PackageSPDXIdentifier:     "id 1",
						PackageVersion:            "version 1",
						PackageFileName:           "file 1",
						PackageSupplier:           nil,
						PackageOriginator:         nil,
						PackageDownloadLocation:   "",
						FilesAnalyzed:             true,
						IsFilesAnalyzedTagPresent: true,
						PackageVerificationCode: &common.PackageVerificationCode{
							Value:         "value 1",
							ExcludedFiles: []string{"a", "b"},
						},
						PackageChecksums: []common.Checksum{
							{
								Algorithm: "alg 1",
								Value:     "val 1",
							},
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
						},
						PackageHomePage:             "home page 1",
						PackageSourceInfo:           "source info 1",
						PackageLicenseConcluded:     "license concluded 1",
						PackageLicenseInfoFromFiles: []string{"a", "b"},
						PackageLicenseDeclared:      "license declared 1",
						PackageLicenseComments:      "license comments 1",
						PackageCopyrightText:        "copyright text 1",
						PackageSummary:              "summary 1",
						PackageDescription:          "description 1",
						PackageComment:              "comment 1",
						PackageExternalReferences: []*spdx.PackageExternalReference{
							{
								Category:           "cat 1",
								RefType:            "type 1",
								Locator:            "locator 1",
								ExternalRefComment: "comment 1",
							},
							{
								Category:           "cat 2",
								RefType:            "type 2",
								Locator:            "locator 2",
								ExternalRefComment: "comment 2",
							},
						},
						PackageAttributionTexts: []string{"a", "b", "c"},
						Files: []*spdx.File{
							{
								FileName:           "file 1",
								FileSPDXIdentifier: "id 1",
								FileTypes:          []string{"a", "b"},
								Checksums: []common.Checksum{
									{
										Algorithm: "alg 1",
										Value:     "val 1",
									},
									{
										Algorithm: "alg 2",
										Value:     "val 2",
									},
								},
								LicenseConcluded:   "license concluded 1",
								LicenseInfoInFiles: []string{"f1", "f2", "f3"},
								LicenseComments:    "comments 1",
								FileCopyrightText:  "copy text 1",
								ArtifactOfProjects: []*spdx.ArtifactOfProject{
									{
										Name:     "name 1",
										HomePage: "home 1",
										URI:      "uri 1",
									},
									{
										Name:     "name 2",
										HomePage: "home 2",
										URI:      "uri 2",
									},
								},
								FileComment:          "comment 1",
								FileNotice:           "notice 1",
								FileContributors:     []string{"c1", "c2"},
								FileAttributionTexts: []string{"att1", "att2"},
								FileDependencies:     []string{"dep1", "dep2", "dep3"},
								Snippets: map[common.ElementID]*spdx.Snippet{
									common.ElementID("e1"): {
										SnippetSPDXIdentifier:         "id1",
										SnippetFromFileSPDXIdentifier: "file1",
										Ranges: []common.SnippetRange{
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             1,
													LineNumber:         2,
													FileSPDXIdentifier: "f1",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             3,
													LineNumber:         4,
													FileSPDXIdentifier: "f2",
												},
											},
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             5,
													LineNumber:         6,
													FileSPDXIdentifier: "f3",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             7,
													LineNumber:         8,
													FileSPDXIdentifier: "f4",
												},
											},
										},
										SnippetLicenseConcluded: "license 1",
										LicenseInfoInSnippet:    []string{"a", "b"},
										SnippetLicenseComments:  "license comment 1",
										SnippetCopyrightText:    "copy 1",
										SnippetComment:          "comment 1",
										SnippetName:             "name 1",
										SnippetAttributionTexts: []string{"att1", "att2", "att3"},
									},
									common.ElementID("e2"): {
										SnippetSPDXIdentifier:         "id2",
										SnippetFromFileSPDXIdentifier: "file2",
										Ranges: []common.SnippetRange{
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             5,
													LineNumber:         6,
													FileSPDXIdentifier: "f3",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             7,
													LineNumber:         8,
													FileSPDXIdentifier: "f4",
												},
											},
											{
												StartPointer: common.SnippetRangePointer{
													Offset:             9,
													LineNumber:         10,
													FileSPDXIdentifier: "f13",
												},
												EndPointer: common.SnippetRangePointer{
													Offset:             11,
													LineNumber:         12,
													FileSPDXIdentifier: "f14",
												},
											},
										},
										SnippetLicenseConcluded: "license 1",
										LicenseInfoInSnippet:    []string{"a", "b"},
										SnippetLicenseComments:  "license comment 1",
										SnippetCopyrightText:    "copy 1",
										SnippetComment:          "comment 1",
										SnippetName:             "name 1",
										SnippetAttributionTexts: []string{"att1", "att2", "att3"},
									},
								},
								Annotations: []spdx.Annotation{
									{
										Annotator: common.Annotator{
											Annotator:     "ann 1",
											AnnotatorType: "typ 1",
										},
										AnnotationDate: "date 1",
										AnnotationType: "type 1",
										AnnotationSPDXIdentifier: common.DocElementID{
											DocumentRefID: "doc 1",
											ElementRefID:  "elem 1",
											SpecialID:     "spec 1",
										},
										AnnotationComment: "comment 1",
									},
									{
										Annotator: common.Annotator{
											Annotator:     "ann 2",
											AnnotatorType: "typ 2",
										},
										AnnotationDate: "date 2",
										AnnotationType: "type 2",
										AnnotationSPDXIdentifier: common.DocElementID{
											DocumentRefID: "doc 2",
											ElementRefID:  "elem 2",
											SpecialID:     "spec 2",
										},
										AnnotationComment: "comment 2",
									},
								},
							},
						},
						Annotations: []spdx.Annotation{
							{
								Annotator: common.Annotator{
									Annotator:     "ann 1",
									AnnotatorType: "typ 1",
								},
								AnnotationDate: "date 1",
								AnnotationType: "type 1",
								AnnotationSPDXIdentifier: common.DocElementID{
									DocumentRefID: "doc 1",
									ElementRefID:  "elem 1",
									SpecialID:     "spec 1",
								},
								AnnotationComment: "comment 1",
							},
							{
								Annotator: common.Annotator{
									Annotator:     "ann 2",
									AnnotatorType: "typ 2",
								},
								AnnotationDate: "date 2",
								AnnotationType: "type 2",
								AnnotationSPDXIdentifier: common.DocElementID{
									DocumentRefID: "doc 2",
									ElementRefID:  "elem 2",
									SpecialID:     "spec 2",
								},
								AnnotationComment: "comment 2",
							},
						},
					},
				},
				Files: []*spdx.File{
					{
						FileName:           "file 1",
						FileSPDXIdentifier: "id 1",
						FileTypes:          []string{"t1", "t2"},
						Checksums: []common.Checksum{
							{
								Algorithm: "alg 1",
								Value:     "val 1",
							},
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
						},
						LicenseConcluded:   "concluded 1",
						LicenseInfoInFiles: []string{"f1", "f2", "f3"},
						LicenseComments:    "comments 1",
						FileCopyrightText:  "copy 1",
						ArtifactOfProjects: []*spdx.ArtifactOfProject{
							{
								Name:     "name 1",
								HomePage: "home 1",
								URI:      "uri 1",
							},
							{
								Name:     "name 2",
								HomePage: "home 2",
								URI:      "uri 2",
							},
						},
						FileComment:          "comment 1",
						FileNotice:           "notice 1",
						FileContributors:     []string{"c1", "c2"},
						FileAttributionTexts: []string{"att1", "att2"},
						FileDependencies:     []string{"d1", "d2", "d3", "d4"},
						Snippets:             nil, // already have snippets elsewhere
						Annotations:          nil, // already have annotations elsewhere
					},
					{
						FileName:           "file 2",
						FileSPDXIdentifier: "id 2",
						FileTypes:          []string{"t3", "t4"},
						Checksums: []common.Checksum{
							{
								Algorithm: "alg 2",
								Value:     "val 2",
							},
							{
								Algorithm: "alg 3",
								Value:     "val 3",
							},
						},
						LicenseConcluded:   "concluded 2",
						LicenseInfoInFiles: []string{"f1", "f2", "f3"},
						LicenseComments:    "comments 2",
						FileCopyrightText:  "copy 2",
						ArtifactOfProjects: []*spdx.ArtifactOfProject{
							{
								Name:     "name 2",
								HomePage: "home 2",
								URI:      "uri 2",
							},
							{
								Name:     "name 4",
								HomePage: "home 4",
								URI:      "uri 4",
							},
						},
						FileComment:          "comment 2",
						FileNotice:           "notice 2",
						FileContributors:     []string{"c1", "c2"},
						FileAttributionTexts: []string{"att1", "att2"},
						FileDependencies:     []string{"d1", "d2", "d3", "d4"},
						Snippets:             nil, // already have snippets elsewhere
						Annotations:          nil, // already have annotations elsewhere
					},
				},
				OtherLicenses: []*spdx.OtherLicense{
					{
						LicenseIdentifier:      "id 1",
						ExtractedText:          "text 1",
						LicenseName:            "name 1",
						LicenseCrossReferences: []string{"x1", "x2", "x3"},
						LicenseComment:         "comment 1",
					},
					{
						LicenseIdentifier:      "id 2",
						ExtractedText:          "text 2",
						LicenseName:            "name 2",
						LicenseCrossReferences: []string{"x4", "x5", "x6"},
						LicenseComment:         "comment 2",
					},
				},
				Relationships: []*spdx.Relationship{
					{
						RefA: common.DocElementID{
							DocumentRefID: "doc 1",
							ElementRefID:  "elem 1",
							SpecialID:     "spec 1",
						},
						RefB: common.DocElementID{
							DocumentRefID: "doc 2",
							ElementRefID:  "elem 2",
							SpecialID:     "spec 2",
						},
						Relationship:        "type 1",
						RelationshipComment: "comment 1",
					},
					{
						RefA: common.DocElementID{
							DocumentRefID: "doc 3",
							ElementRefID:  "elem 3",
							SpecialID:     "spec 3",
						},
						RefB: common.DocElementID{
							DocumentRefID: "doc 4",
							ElementRefID:  "elem 4",
							SpecialID:     "spec 4",
						},
						Relationship:        "type 2",
						RelationshipComment: "comment 2",
					},
				},
				Annotations: []*spdx.Annotation{
					{
						Annotator: common.Annotator{
							Annotator:     "annotator 1",
							AnnotatorType: "annotator type 1",
						},
						AnnotationDate: "date 1",
						AnnotationType: "type 1",
						AnnotationSPDXIdentifier: common.DocElementID{
							DocumentRefID: "doc 1",
							ElementRefID:  "elem 1",
							SpecialID:     "spec 1",
						},
						AnnotationComment: "comment 1",
					},
					{
						Annotator: common.Annotator{
							Annotator:     "annotator 2",
							AnnotatorType: "annotator type 2",
						},
						AnnotationDate: "date 2",
						AnnotationType: "type 2",
						AnnotationSPDXIdentifier: common.DocElementID{
							DocumentRefID: "doc 2",
							ElementRefID:  "elem 2",
							SpecialID:     "spec 2",
						},
						AnnotationComment: "comment 2",
					},
				},
				Snippets: []spdx.Snippet{
					{
						SnippetSPDXIdentifier:         "id1",
						SnippetFromFileSPDXIdentifier: "file1",
						Ranges: []common.SnippetRange{
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             1,
									LineNumber:         2,
									FileSPDXIdentifier: "f1",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             3,
									LineNumber:         4,
									FileSPDXIdentifier: "f2",
								},
							},
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             5,
									LineNumber:         6,
									FileSPDXIdentifier: "f3",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             7,
									LineNumber:         8,
									FileSPDXIdentifier: "f4",
								},
							},
						},
						SnippetLicenseConcluded: "license 1",
						LicenseInfoInSnippet:    []string{"a", "b"},
						SnippetLicenseComments:  "license comment 1",
						SnippetCopyrightText:    "copy 1",
						SnippetComment:          "comment 1",
						SnippetName:             "name 1",
						SnippetAttributionTexts: []string{"att1", "att2", "att3"},
					},
					{
						SnippetSPDXIdentifier:         "id2",
						SnippetFromFileSPDXIdentifier: "file2",
						Ranges: []common.SnippetRange{
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             5,
									LineNumber:         6,
									FileSPDXIdentifier: "f3",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             7,
									LineNumber:         8,
									FileSPDXIdentifier: "f4",
								},
							},
							{
								StartPointer: common.SnippetRangePointer{
									Offset:             9,
									LineNumber:         10,
									FileSPDXIdentifier: "f13",
								},
								EndPointer: common.SnippetRangePointer{
									Offset:             11,
									LineNumber:         12,
									FileSPDXIdentifier: "f14",
								},
							},
						},
						SnippetLicenseConcluded: "license 1",
						LicenseInfoInSnippet:    []string{"a", "b"},
						SnippetLicenseComments:  "license comment 1",
						SnippetCopyrightText:    "copy 1",
						SnippetComment:          "comment 1",
						SnippetName:             "name 1",
						SnippetAttributionTexts: []string{"att1", "att2", "att3"},
					},
				},
				Reviews: []*spdx.Review{
					{
						Reviewer:      "reviewer 1",
						ReviewerType:  "type 1",
						ReviewDate:    "date 1",
						ReviewComment: "comment 1",
					},
					{
						Reviewer:      "reviewer 2",
						ReviewerType:  "type 2",
						ReviewDate:    "date 2",
						ReviewComment: "comment 2",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outType := reflect.TypeOf(test.expected)
			outInstance := reflect.New(outType).Interface()

			// convert the start document to the target document using the conversion chain
			err := Document(test.source, outInstance)
			if err != nil {
				t.Fatalf("error converting: %v", err)
			}
			outInstance = reflect.ValueOf(outInstance).Elem().Interface()

			// use JSONEq here because it is much easier to see differences
			require.JSONEq(t, toJSON(test.expected), toJSON(outInstance))
		})
	}
}

func toJSON(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
