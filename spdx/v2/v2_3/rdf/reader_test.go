package rdf

import (
	"io"
	"strings"
	"testing"
)

func Test_Read(t *testing.T) {
	var reader io.Reader
	var err error

	// TestCase 1: invalid rdf/xml must raise an error
	reader = strings.NewReader("")
	_, err = Read(reader)
	if err == nil {
		t.Errorf("expected an EOF error reading an empty file, got %v", err)
	}

	// TestCase 2: Valid rdf/xml but invalid spdx document must raise an error
	reader = strings.NewReader(`
		<rdf:RDF
			xmlns:spdx="http://spdx.org/rdf/terms#"
			xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
			xmlns="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#">
		</rdf:RDF>
	`)
	_, err = Read(reader)
	if err == nil {
		t.Errorf("expected an error due to no SpdxDocument Node in the document")
	}

	// TestCase 3: New SPDX package elements
	reader = strings.NewReader(`
		<rdf:RDF
			xmlns:spdx="http://spdx.org/rdf/terms#"
			xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
			xmlns="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#"
			xmlns:doap="http://usefulinc.com/ns/doap#"
			xmlns:j.0="http://www.w3.org/2009/pointers#"
			xmlns:rdfs="http://www.w3.org/2000/01/rdf-schema#">
			<spdx:SpdxDocument rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DOCUMENT">
				<spdx:specVersion>SPDX-2.0</spdx:specVersion>
				<spdx:relationship>
					<spdx:Relationship>
						<spdx:relationshipType rdf:resource="http://spdx.org/rdf/terms#relationshipType_describes"/>
						<spdx:relatedSpdxElement>
							<spdx:Package rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon">
								<spdx:name>Some-Package</spdx:name>
								<spdx:primaryPackagePurpose rdf:resource="packagePurpose_container" />
								<spdx:releaseDate>2021-10-15T02:38:00Z</spdx:releaseDate>
								<spdx:builtDate>2021-09-15T02:38:00Z</spdx:builtDate>
								<spdx:validUntilDate>2022-10-15T02:38:00Z</spdx:validUntilDate>
							</spdx:Package>
						</spdx:relatedSpdxElement>
					</spdx:Relationship>
				</spdx:relationship>
			</spdx:SpdxDocument>
		</rdf:RDF>
	`)

	doc, err := Read(reader)
	if err != nil {
		t.Errorf("expected valid SPDX document: %v", err)
	}

	if doc == nil {
		t.Fatalf("expected valid SPDX document but got nil")
	}

	if len(doc.Packages) == 0 {
		t.Errorf("expected packages but got none")
	}

	pkg := doc.Packages[0]
	if pkg.PackageName != "Some-Package" {
		t.Errorf("expected package nameof Some-Package but got: %s", pkg.PackageName)
	}
	if pkg.PrimaryPackagePurpose != "CONTAINER" {
		t.Errorf("expected package primary purpose of CONTAINER but got: %s", pkg.PrimaryPackagePurpose)
	}
	if pkg.ReleaseDate != "2021-10-15T02:38:00Z" {
		t.Errorf("expected release date of 2021-10-15T02:38:00Z but got: %s", pkg.ReleaseDate)
	}
	if pkg.BuiltDate != "2021-09-15T02:38:00Z" {
		t.Errorf("expected built date of 2021-09-15T02:38:00Z but got: %s", pkg.BuiltDate)
	}
	if pkg.ValidUntilDate != "2022-10-15T02:38:00Z" {
		t.Errorf("expected valid until date of 2022-10-15T02:38:00Z but got: %s", pkg.ValidUntilDate)
	}
}
