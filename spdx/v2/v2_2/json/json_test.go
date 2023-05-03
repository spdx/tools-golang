// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package json

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"

	"github.com/spdx/tools-golang/json"
	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
	"github.com/spdx/tools-golang/spdx/v2/v2_2/example"
)

func TestLoad(t *testing.T) {
	want := example.Copy()

	want.Relationships = append(want.Relationships, []*spdx.Relationship{
		{
			RefA:         common.DocElementID{ElementRefID: "DOCUMENT"},
			RefB:         common.DocElementID{ElementRefID: "File"},
			Relationship: "DESCRIBES",
		},
		{
			RefA:         common.DocElementID{ElementRefID: "DOCUMENT"},
			RefB:         common.DocElementID{ElementRefID: "Package"},
			Relationship: "DESCRIBES",
		},
		{
			RefA:         common.DocElementID{ElementRefID: "Package"},
			RefB:         common.DocElementID{ElementRefID: "CommonsLangSrc"},
			Relationship: "CONTAINS",
		},
		{
			RefA:         common.DocElementID{ElementRefID: "Package"},
			RefB:         common.DocElementID{ElementRefID: "JenaLib"},
			Relationship: "CONTAINS",
		},
		{
			RefA:         common.DocElementID{ElementRefID: "Package"},
			RefB:         common.DocElementID{ElementRefID: "DoapSource"},
			Relationship: "CONTAINS",
		},
	}...)

	file, err := os.Open("../../../../examples/sample-docs/json/SPDXJSONExample-v2.2.spdx.json")
	if err != nil {
		panic(fmt.Errorf("error opening File: %s", err))
	}

	var got spdx.Document
	err = json.ReadInto(file, &got)
	if err != nil {
		t.Errorf("json.parser.Load() error = %v", err)
		return
	}

	if !cmp.Equal(want, got, cmpopts.IgnoreUnexported(spdx.Package{})) {
		t.Errorf("got incorrect struct after parsing YAML example: %s", cmp.Diff(want, got, cmpopts.IgnoreUnexported(spdx.Package{})))
		return
	}
}

func Test_Write(t *testing.T) {
	want := example.Copy()

	// we always output FilesAnalyzed, even though we handle reading files where it is omitted
	for _, p := range want.Packages {
		p.IsFilesAnalyzedTagPresent = true
	}

	w := &bytes.Buffer{}

	if err := json.Write(&want, w); err != nil {
		t.Errorf("Write() error = %v", err.Error())
		return
	}

	// we should be able to parse what the writer wrote, and it should be identical to the original struct we wrote
	var got spdx.Document
	err := json.ReadInto(bytes.NewReader(w.Bytes()), &got)
	if err != nil {
		t.Errorf("failed to parse written document: %v", err.Error())
		return
	}

	if !cmp.Equal(want, got, cmpopts.IgnoreUnexported(spdx.Package{})) {
		t.Errorf("got incorrect struct after writing and re-parsing JSON example: %s", cmp.Diff(want, got, cmpopts.IgnoreUnexported(spdx.Package{})))
		return
	}
}

func Test_ShorthandFields(t *testing.T) {
	contents := `{
		"spdxVersion": "SPDX-2.2",
		"dataLicense": "CC0-1.0",
		"SPDXID": "SPDXRef-DOCUMENT",
		"name": "SPDX-Tools-v2.0",
		"documentDescribes": [
			"SPDXRef-Container"
		],
		"packages": [
			{
				"name": "Container",
				"SPDXID": "SPDXRef-Container"
			},
			{
				"name": "Package-1",
				"SPDXID": "SPDXRef-Package-1",
				"versionInfo": "1.1.1",
				"hasFiles": [
					"SPDXRef-File-1",
					"SPDXRef-File-2"
				]
			},
			{
				"name": "Package-2",
				"SPDXID": "SPDXRef-Package-2",
				"versionInfo": "2.2.2"
			}
		],
		"files": [
			{
				"fileName": "./f1",
				"SPDXID": "SPDXRef-File-1"
			},
			{
				"fileName": "./f2",
				"SPDXID": "SPDXRef-File-2"
			}
		]
	}`

	doc := spdx.Document{}
	err := json.ReadInto(strings.NewReader(contents), &doc)

	require.NoError(t, err)

	id := func(s string) common.DocElementID {
		return common.DocElementID{
			ElementRefID: common.ElementID(s),
		}
	}

	require.Equal(t, spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: "DOCUMENT",
		DocumentName:   "SPDX-Tools-v2.0",
		Packages: []*spdx.Package{
			{
				PackageName:           "Container",
				PackageSPDXIdentifier: "Container",
				FilesAnalyzed:         true,
			},
			{
				PackageName:           "Package-1",
				PackageSPDXIdentifier: "Package-1",
				PackageVersion:        "1.1.1",
				FilesAnalyzed:         true,
			},
			{
				PackageName:           "Package-2",
				PackageSPDXIdentifier: "Package-2",
				PackageVersion:        "2.2.2",
				FilesAnalyzed:         true,
			},
		},
		Files: []*spdx.File{
			{
				FileName:           "./f1",
				FileSPDXIdentifier: "File-1",
			},
			{
				FileName:           "./f2",
				FileSPDXIdentifier: "File-2",
			},
		},
		Relationships: []*spdx.Relationship{
			{
				RefA:         id("DOCUMENT"),
				RefB:         id("Container"),
				Relationship: common.TypeRelationshipDescribe,
			},
			{
				RefA:         id("Package-1"),
				RefB:         id("File-1"),
				Relationship: common.TypeRelationshipContains,
			},
			{
				RefA:         id("Package-1"),
				RefB:         id("File-2"),
				Relationship: common.TypeRelationshipContains,
			},
		},
	}, doc)
}

func Test_JsonEnums(t *testing.T) {
	contents := `{
        "spdxVersion": "SPDX-2.2",
        "dataLicense": "CC0-1.0",
        "SPDXID": "SPDXRef-DOCUMENT",
        "name": "SPDX-Tools-v2.0",
        "documentDescribes": [
            "SPDXRef-Container"
        ],
        "packages": [
            {
                "name": "Container",
                "SPDXID": "SPDXRef-Container"
            },
            {
                "name": "Package-1",
                "SPDXID": "SPDXRef-Package-1",
                "versionInfo": "1.1.1",
                "externalRefs": [{
                    "referenceCategory": "PACKAGE_MANAGER",
                    "referenceLocator": "pkg:somepkg/ns/name1",
                    "referenceType": "purl"
                }]
            },
            {
                "name": "Package-2",
                "SPDXID": "SPDXRef-Package-2",
                "versionInfo": "2.2.2",
                "externalRefs": [{
                    "referenceCategory": "PACKAGE-MANAGER",
                    "referenceLocator": "pkg:somepkg/ns/name2",
                    "referenceType": "purl"
                }]
            },
			{
                "name": "Package-3",
                "SPDXID": "SPDXRef-Package-3",
                "versionInfo": "3.3.3",
                "externalRefs": [{
                    "referenceCategory": "PERSISTENT_ID",
                    "referenceLocator": "gitoid:blob:sha1:261eeb9e9f8b2b4b0d119366dda99c6fd7d35c64",
                    "referenceType": "gitoid"
                }]
            },
			{
                "name": "Package-4",
                "SPDXID": "SPDXRef-Package-4",
                "versionInfo": "4.4.4",
                "externalRefs": [{
                    "referenceCategory": "PERSISTENT-ID",
                    "referenceLocator": "gitoid:blob:sha1:261eeb9e9f8b2b4b0d119366dda99c6fd7d35c64",
                    "referenceType": "gitoid"
                }]
            }
        ]
    }`

	doc := spdx.Document{}
	err := json.ReadInto(strings.NewReader(contents), &doc)

	require.NoError(t, err)

	id := func(s string) common.DocElementID {
		return common.DocElementID{
			ElementRefID: common.ElementID(s),
		}
	}

	require.Equal(t, spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: "DOCUMENT",
		DocumentName:   "SPDX-Tools-v2.0",
		Packages: []*spdx.Package{
			{
				PackageName:           "Container",
				PackageSPDXIdentifier: "Container",
				FilesAnalyzed:         true,
			},
			{
				PackageName:           "Package-1",
				PackageSPDXIdentifier: "Package-1",
				PackageVersion:        "1.1.1",
				PackageExternalReferences: []*spdx.PackageExternalReference{
					{
						Category: common.CategoryPackageManager,
						RefType:  common.TypePackageManagerPURL,
						Locator:  "pkg:somepkg/ns/name1",
					},
				},
				FilesAnalyzed: true,
			},
			{
				PackageName:           "Package-2",
				PackageSPDXIdentifier: "Package-2",
				PackageVersion:        "2.2.2",
				PackageExternalReferences: []*spdx.PackageExternalReference{
					{
						Category: common.CategoryPackageManager,
						RefType:  common.TypePackageManagerPURL,
						Locator:  "pkg:somepkg/ns/name2",
					},
				},
				FilesAnalyzed: true,
			},
			{
				PackageName:           "Package-3",
				PackageSPDXIdentifier: "Package-3",
				PackageVersion:        "3.3.3",
				PackageExternalReferences: []*spdx.PackageExternalReference{
					{
						Category: common.CategoryPersistentId,
						RefType:  common.TypePersistentIdGitoid,
						Locator:  "gitoid:blob:sha1:261eeb9e9f8b2b4b0d119366dda99c6fd7d35c64",
					},
				},
				FilesAnalyzed: true,
			},
			{
				PackageName:           "Package-4",
				PackageSPDXIdentifier: "Package-4",
				PackageVersion:        "4.4.4",
				PackageExternalReferences: []*spdx.PackageExternalReference{
					{
						Category: common.CategoryPersistentId,
						RefType:  common.TypePersistentIdGitoid,
						Locator:  "gitoid:blob:sha1:261eeb9e9f8b2b4b0d119366dda99c6fd7d35c64",
					},
				},
				FilesAnalyzed: true,
			},
		},
		Relationships: []*spdx.Relationship{
			{
				RefA:         id("DOCUMENT"),
				RefB:         id("Container"),
				Relationship: common.TypeRelationshipDescribe,
			},
		},
	}, doc)
}
