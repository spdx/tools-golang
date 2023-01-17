// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"reflect"
	"sort"
	"testing"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
)

func Test_rdfParser2_3_getAnyLicenseFromNode(t *testing.T) {
	// since this function is a mux, we just have to make sure that with each
	// type of input, it is able to redirect the request to an appropriate
	// license getter.

	// TestCase 1: input node is just a node string without any associated
	//			   triple (either a NONE|NOASSERTION) because for other case,
	//			   the license should've been associated with other triples
	parser, _ := parserFromBodyContent(``)
	inputNode := &gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       NS_SPDX + "NONE",
	}
	lic, err := parser.getAnyLicenseFromNode(inputNode)
	if err != nil {
		t.Errorf("error parsing a valid license input: %v", err)
	}
	// checking if the return type is a SpecialLicense
	switch lic.(type) {
	case SpecialLicense:
	default:
		t.Errorf("expected license to be of type SpecialLicense, found %v", reflect.TypeOf(lic))
	}

	// TestCase 2: DisjunctiveLicenseSet:
	parser, _ = parserFromBodyContent(`
		<spdx:DisjunctiveLicenseSet>
			<spdx:member rdf:resource="http://spdx.org/licenses/Nokia"/>
			<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
		</spdx:DisjunctiveLicenseSet>
	`)
	inputNode = parser.gordfParserObj.Triples[0].Subject
	lic, err = parser.getAnyLicenseFromNode(inputNode)
	if err != nil {
		t.Errorf("error parsing a valid license input: %v", err)
	}
	// checking if the return type is a DisjunctiveLicenseSet
	switch lic.(type) {
	case DisjunctiveLicenseSet:
	default:
		t.Errorf("expected license to be of type DisjunctiveLicenseSet, found %v", reflect.TypeOf(lic))
	}

	// TestCase 3: ConjunctiveLicenseSet:
	parser, _ = parserFromBodyContent(`
		<spdx:ConjunctiveLicenseSet>
			<spdx:member rdf:resource="http://spdx.org/licenses/Nokia"/>
			<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
		</spdx:ConjunctiveLicenseSet>
	`)
	inputNode = parser.gordfParserObj.Triples[0].Subject
	lic, err = parser.getAnyLicenseFromNode(inputNode)
	if err != nil {
		t.Errorf("error parsing a valid license input: %v", err)
	}
	// checking if the return type is a ConjunctiveLicenseSet
	switch lic.(type) {
	case ConjunctiveLicenseSet:
	default:
		t.Errorf("expected license to be of type ConjunctiveLicenseSet, found %v", reflect.TypeOf(lic))
	}

	// TestCase 4: ExtractedLicensingInfo
	parser, _ = parserFromBodyContent(`
		<spdx:ExtractedLicensingInfo rdf:about="http://spdx.dev/spdx.rdf#LicenseRef-Freeware">
			<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
			<spdx:name>freeware</spdx:name>
			<spdx:extractedText><![CDATA[...]]></spdx:extractedText>
	  	</spdx:ExtractedLicensingInfo>
	`)
	inputNode = parser.gordfParserObj.Triples[0].Subject
	lic, err = parser.getAnyLicenseFromNode(inputNode)
	if err != nil {
		t.Errorf("error parsing a valid license input: %v", err)
	}
	// checking if the return type is a ExtractedLicensingInfo
	switch lic.(type) {
	case ExtractedLicensingInfo:
	default:
		t.Errorf("expected license to be of type ExtractedLicensingInfo, found %v", reflect.TypeOf(lic))
	}

	// TestCase 4: ExtractedLicensingInfo
	parser, _ = parserFromBodyContent(`
		<spdx:ExtractedLicensingInfo rdf:about="http://spdx.dev/spdx.rdf#LicenseRef-Freeware">
			<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
			<spdx:name>freeware</spdx:name>
			<spdx:extractedText><![CDATA[...]]></spdx:extractedText>
	  	</spdx:ExtractedLicensingInfo>
	`)
	inputNode = parser.gordfParserObj.Triples[0].Subject
	lic, err = parser.getAnyLicenseFromNode(inputNode)
	if err != nil {
		t.Errorf("error parsing a valid license input: %v", err)
	}
	// checking if the return type is a ExtractedLicensingInfo
	switch lic.(type) {
	case ExtractedLicensingInfo:
	default:
		t.Errorf("expected license to be of type ExtractedLicensingInfo, found %v", reflect.TypeOf(lic))
	}

	// TestCase 5: License
	parser, _ = parserFromBodyContent(`
		<spdx:License rdf:about="http://spdx.org/licenses/Apache-2.0">
			<spdx:standardLicenseTemplate>&lt;&gt; Apache License Version 2.0, January 2004 http://www.apache.org/licenses/&lt;&gt;&lt;&gt; TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION&lt;&gt; &lt;&gt; Definitions. "License" shall mean the terms and conditions for use, reproduction, and distribution as defined by Sections 1 through 9 of this document. "Licensor" shall mean the copyright owner or entity authorized by the copyright owner that is granting the License. "Legal Entity" shall mean the union of the acting entity and all other entities that control, are controlled by, or are under common control with that entity. For the purposes of this definition, "control" means (i) the power, direct or indirect, to cause the direction or management of such entity, whether by contract or otherwise, or (ii) ownership of fifty percent (50%) or more of the outstanding shares, or (iii) beneficial ownership of such entity. "You" (or "Your") shall mean an individual or Legal Entity exercising permissions granted by this License. "Source" form shall mean the preferred form for making modifications, including but not limited to software source code, documentation source, and configuration files. "Object" form shall mean any form resulting from mechanical transformation or translation of a Source form, including but not limited to compiled object code, generated documentation, and conversions to other media types. "Work" shall mean the work of authorship, whether in Source or Object form, made available under the License, as indicated by a copyright notice that is included in or attached to the work (an example is provided in the Appendix below). "Derivative Works" shall mean any work, whether in Source or Object form, that is based on (or derived from) the Work and for which the editorial revisions, annotations, elaborations, or other modifications represent, as a whole, an original work of authorship. For the purposes of this License, Derivative Works shall not include works that remain separable from, or merely link (or bind by name) to the interfaces of, the Work and Derivative Works thereof. "Contribution" shall mean any work of authorship, including the original version of the Work and any modifications or additions to that Work or Derivative Works thereof, that is intentionally submitted to Licensor for inclusion in the Work by the copyright owner or by an individual or Legal Entity authorized to submit on behalf of the copyright owner. For the purposes of this definition, "submitted" means any form of electronic, verbal, or written communication sent to the Licensor or its representatives, including but not limited to communication on electronic mailing lists, source code control systems, and issue tracking systems that are managed by, or on behalf of, the Licensor for the purpose of discussing and improving the Work, but excluding communication that is conspicuously marked or otherwise designated in writing by the copyright owner as "Not a Contribution." "Contributor" shall mean Licensor and any individual or Legal Entity on behalf of whom a Contribution has been received by Licensor and subsequently incorporated within the Work. &lt;&gt; Grant of Copyright License. Subject to the terms and conditions of this License, each Contributor hereby grants to You a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare Derivative Works of, publicly display, publicly perform, sublicense, and distribute the Work and such Derivative Works in Source or Object form. &lt;&gt; Grant of Patent License. Subject to the terms and conditions of this License, each Contributor hereby grants to You a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable (except as stated in this section) patent license to make, have made, use, offer to sell, sell, import, and otherwise transfer the Work, where such license applies only to those patent claims licensable by such Contributor that are necessarily infringed by their Contribution(s) alone or by combination of their Contribution(s) with the Work to which such Contribution(s) was submitted. If You institute patent litigation against any entity (including a cross-claim or counterclaim in a lawsuit) alleging that the Work or a Contribution incorporated within the Work constitutes direct or contributory patent infringement, then any patent licenses granted to You under this License for that Work shall terminate as of the date such litigation is filed. &lt;&gt; Redistribution. You may reproduce and distribute copies of the Work or Derivative Works thereof in any medium, with or without modifications, and in Source or Object form, provided that You meet the following conditions: &lt;&gt; You must give any other recipients of the Work or Derivative Works a copy of this License; and &lt;&gt; You must cause any modified files to carry prominent notices stating that You changed the files; and &lt;&gt; You must retain, in the Source form of any Derivative Works that You distribute, all copyright, patent, trademark, and attribution notices from the Source form of the Work, excluding those notices that do not pertain to any part of the Derivative Works; and &lt;&gt; If the Work includes a "NOTICE" text file as part of its distribution, then any Derivative Works that You distribute must include a readable copy of the attribution notices contained within such NOTICE file, excluding those notices that do not pertain to any part of the Derivative Works, in at least one of the following places: within a NOTICE text file distributed as part of the Derivative Works; within the Source form or documentation, if provided along with the Derivative Works; or, within a display generated by the Derivative Works, if and wherever such third-party notices normally appear. The contents of the NOTICE file are for informational purposes only and do not modify the License. You may add Your own attribution notices within Derivative Works that You distribute, alongside or as an addendum to the NOTICE text from the Work, provided that such additional attribution notices cannot be construed as modifying the License. You may add Your own copyright statement to Your modifications and may provide additional or different license terms and conditions for use, reproduction, or distribution of Your modifications, or for any such Derivative Works as a whole, provided Your use, reproduction, and distribution of the Work otherwise complies with the conditions stated in this License. &lt;&gt; Submission of Contributions. Unless You explicitly state otherwise, any Contribution intentionally submitted for inclusion in the Work by You to the Licensor shall be under the terms and conditions of this License, without any additional terms or conditions. Notwithstanding the above, nothing herein shall supersede or modify the terms of any separate license agreement you may have executed with Licensor regarding such Contributions. &lt;&gt; Trademarks. This License does not grant permission to use the trade names, trademarks, service marks, or product names of the Licensor, except as required for reasonable and customary use in describing the origin of the Work and reproducing the content of the NOTICE file. &lt;&gt; Disclaimer of Warranty. Unless required by applicable law or agreed to in writing, Licensor provides the Work (and each Contributor provides its Contributions) on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied, including, without limitation, any warranties or conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE. You are solely responsible for determining the appropriateness of using or redistributing the Work and assume any risks associated with Your exercise of permissions under this License. &lt;&gt; Limitation of Liability. In no event and under no legal theory, whether in tort (including negligence), contract, or otherwise, unless required by applicable law (such as deliberate and grossly negligent acts) or agreed to in writing, shall any Contributor be liable to You for damages, including any direct, indirect, special, incidental, or consequential damages of any character arising as a result of this License or out of the use or inability to use the Work (including but not limited to damages for loss of goodwill, work stoppage, computer failure or malfunction, or any and all other commercial damages or losses), even if such Contributor has been advised of the possibility of such damages. &lt;&gt; Accepting Warranty or Additional Liability. While redistributing the Work or Derivative Works thereof, You may choose to offer, and charge a fee for, acceptance of support, warranty, indemnity, or other liability obligations and/or rights consistent with this License. However, in accepting such obligations, You may act only on Your own behalf and on Your sole responsibility, not on behalf of any other Contributor, and only if You agree to indemnify, defend, and hold each Contributor harmless for any liability incurred by, or claims asserted against, such Contributor by reason of your accepting any such warranty or additional liability.&lt;&gt; END OF TERMS AND CONDITIONS APPENDIX: How to apply the Apache License to your work. To apply the Apache License to your work, attach the following boilerplate notice, with the fields enclosed by brackets "[]" replaced with your own identifying information. (Don't include the brackets!) The text should be enclosed in the appropriate comment syntax for the file format. We also recommend that a file or class name and description of purpose be included on the same "printed page" as the copyright notice for easier identification within third-party archives. Copyright &lt;&gt; Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0 Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.&lt;&gt;</spdx:standardLicenseTemplate>
			<rdfs:seeAlso>http://www.apache.org/licenses/LICENSE-2.0</rdfs:seeAlso>
			<spdx:name>Apache License 2.0</spdx:name>
			<spdx:licenseId>Apache-2.0</spdx:licenseId>
			<spdx:isOsiApproved>true</spdx:isOsiApproved>
			<rdfs:seeAlso>http://www.opensource.org/licenses/Apache-2.0</rdfs:seeAlso>
			<spdx:licenseText>...</spdx:licenseText>
			<spdx:standardLicenseHeader>...</spdx:standardLicenseHeader>
	  </spdx:License>
	`)
	inputNode = parser.gordfParserObj.Triples[0].Subject
	lic, err = parser.getAnyLicenseFromNode(inputNode)
	if err != nil {
		t.Errorf("error parsing a valid license input: %v", err)
	}
	// checking if the return type is a License
	switch lic.(type) {
	case License:
	default:
		t.Errorf("expected license to be of type License, found %v", reflect.TypeOf(lic))
	}

	// TestCase 5: WithExceptionOperator
	parser, _ = parserFromBodyContent(`
		<spdx:WithExceptionOperator>
			<spdx:licenseException>
				<spdx:LicenseException rdf:nodeID="A1">
					<spdx:example></spdx:example>
					<spdx:licenseExceptionId>Libtool-exception</spdx:licenseExceptionId>
					<rdfs:comment></rdfs:comment>
				</spdx:LicenseException>
			</spdx:licenseException>
			<spdx:member rdf:resource="http://spdx.org/licenses/GPL-2.0-or-later"/>
		</spdx:WithExceptionOperator>
	`)
	inputNode = parser.gordfParserObj.Triples[0].Subject
	lic, err = parser.getAnyLicenseFromNode(inputNode)
	if err != nil {
		t.Errorf("error parsing a valid license input: %v", err)
	}
	// checking if the return type is a WithExceptionOperator
	switch lic.(type) {
	case WithExceptionOperator:
	default:
		t.Errorf("expected license to be of type WithExceptionOperator, found %v", reflect.TypeOf(lic))
	}

	// TestCase 6: OrLaterOperator
	parser, _ = parserFromBodyContent(`
		<spdx:OrLaterOperator>
			<spdx:member>
				<spdx:SimpleLicensingInfo>
					<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
					<spdx:name>freeware</spdx:name>
				</spdx:SimpleLicensingInfo>
			</spdx:member>
		</spdx:OrLaterOperator>
	`)
	inputNode = parser.gordfParserObj.Triples[0].Subject
	lic, err = parser.getAnyLicenseFromNode(inputNode)
	if err != nil {
		t.Errorf("error parsing a valid license input: %v", err)
	}
	// checking if the return type is a OrLaterOperator
	switch lic.(type) {
	case OrLaterOperator:
	default:
		t.Errorf("expected license to be of type OrLaterOperator, found %v", reflect.TypeOf(lic))
	}

	// TestCase 7: checking if an unknown license raises an error.
	parser, _ = parserFromBodyContent(`
		<spdx:UnknownLicense>
			<spdx:unknownTag />
		</spdx:UnknownLicense>
	`)
	node := parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getAnyLicenseFromNode(node)
	t.Log(err)
	if err == nil {
		t.Errorf("should've raised an error for invalid input")
	}

	// TestCase 8: cyclic dependent license must raise an error.
	parser, _ = parserFromBodyContent(`
		<spdx:ConjunctiveLicenseSet rdf:about="#SPDXRef-RecursiveLicense">
			<spdx:member rdf:resource="http://spdx.org/licenses/GPL-2.0-or-later"/>
			<spdx:member>
				<spdx:ConjunctiveLicenseSet rdf:about="#SPDXRef-RecursiveLicense">
					<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
					<spdx:member rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-RecursiveLicense"/>
				</spdx:ConjunctiveLicenseSet>
			</spdx:member>
		</spdx:ConjunctiveLicenseSet>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getAnyLicenseFromNode(node)
	if err == nil {
		t.Errorf("expected an error due to cyclic dependent license. found %v", err)
	}
}

func Test_rdfParser2_3_getConjunctiveLicenseSetFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var err error
	var licenseNode *gordfParser.Node
	var license ConjunctiveLicenseSet

	// TestCase 1: invalid license member
	parser, _ = parserFromBodyContent(`
		<spdx:ConjunctiveLicenseSet>
			<spdx:member rdf:resource="http://spdx.org/licenses/Unknown"/>
			<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
		</spdx:ConjunctiveLicenseSet>
	`)
	licenseNode = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getConjunctiveLicenseSetFromNode(licenseNode)
	if err == nil {
		t.Errorf("expected an error saying invalid license member, found <nil>")
	}

	// TestCase 2: invalid predicate in the licenseSet.
	parser, _ = parserFromBodyContent(`
		<spdx:ConjunctiveLicenseSet>
			<spdx:member rdf:resource="http://spdx.org/licenses/CC0-1.0"/>
			<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
			<spdx:unknownTag />
		</spdx:ConjunctiveLicenseSet>
	`)
	licenseNode = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getConjunctiveLicenseSetFromNode(licenseNode)
	if err == nil {
		t.Errorf("expected an error saying invalid predicate found")
	}

	// TestCase 3: valid example.
	parser, _ = parserFromBodyContent(`
		<spdx:ConjunctiveLicenseSet>
			<spdx:member rdf:resource="http://spdx.org/licenses/Nokia"/>
			<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
		</spdx:ConjunctiveLicenseSet>
	`)
	licenseNode = parser.gordfParserObj.Triples[0].Subject
	license, err = parser.getConjunctiveLicenseSetFromNode(licenseNode)
	if err != nil {
		t.Errorf("unexpected error parsing licenseSet: %v", err)
	}
	nMembers := len(license.members)
	if nMembers != 2 {
		t.Errorf("expected licenseSet to have 2 members, found %d", nMembers)
	}
	licenseMembers := mapLicensesToStrings(license.members)
	expectedLicenseMembers := []string{"LGPL-2.0", "Nokia"}
	sort.Strings(licenseMembers)
	if !reflect.DeepEqual(licenseMembers, expectedLicenseMembers) {
		t.Errorf("expected %v, found %v", expectedLicenseMembers, licenseMembers)
	}
}

func Test_rdfParser2_3_getDisjunctiveLicenseSetFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var err error
	var licenseNode *gordfParser.Node
	var license DisjunctiveLicenseSet

	// TestCase 1: invalid license member
	parser, _ = parserFromBodyContent(`
		<spdx:DisjunctiveLicenseSet>
			<spdx:member rdf:resource="http://spdx.org/licenses/Unknown"/>
			<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
		</spdx:DisjunctiveLicenseSet>
	`)
	licenseNode = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getDisjunctiveLicenseSetFromNode(licenseNode)
	if err == nil {
		t.Errorf("expected an error saying invalid license member, found <nil>")
	}

	// TestCase 2: invalid predicate in the licenseSet.
	parser, _ = parserFromBodyContent(`
		<spdx:DisjunctiveLicenseSet>
			<spdx:member rdf:resource="http://spdx.org/licenses/Unknown"/>
			<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
			<spdx:unknownTag />
		</spdx:DisjunctiveLicenseSet>
	`)
	licenseNode = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getDisjunctiveLicenseSetFromNode(licenseNode)
	if err == nil {
		t.Errorf("expected an error saying invalid predicate found")
	}

	// TestCase 3: valid example.
	parser, _ = parserFromBodyContent(`
		<spdx:DisjunctiveLicenseSet>
			<spdx:member rdf:resource="http://spdx.org/licenses/Nokia"/>
			<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
		</spdx:DisjunctiveLicenseSet>
	`)
	licenseNode = parser.gordfParserObj.Triples[0].Subject
	license, err = parser.getDisjunctiveLicenseSetFromNode(licenseNode)
	if err != nil {
		t.Errorf("unexpected error parsing licenseSet: %v", err)
	}
	nMembers := len(license.members)
	if nMembers != 2 {
		t.Errorf("expected licenseSet to have 2 members, found %d", nMembers)
	}
	licenseMembers := mapLicensesToStrings(license.members)
	expectedLicenseMembers := []string{"LGPL-2.0", "Nokia"}
	sort.Strings(licenseMembers)
	if !reflect.DeepEqual(licenseMembers, expectedLicenseMembers) {
		t.Errorf("expected %v, found %v", expectedLicenseMembers, licenseMembers)
	}
}

func Test_rdfParser2_3_getLicenseExceptionFromNode(t *testing.T) {
	var licenseException LicenseException
	var err error
	var node *gordfParser.Node
	var parser *rdfParser2_3

	// TestCase 1: invalid value for rdf:seeAlso
	parser, _ = parserFromBodyContent(`
		<spdx:LicenseException>
			<rdfs:seeAlso>see-also</rdfs:seeAlso>
			<spdx:example></spdx:example>
			<spdx:licenseExceptionId>Libtool-exception</spdx:licenseExceptionId>
			<rdfs:comment></rdfs:comment>
		</spdx:LicenseException>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getLicenseExceptionFromNode(node)
	if err == nil {
		t.Errorf("should've raised an error due to invalid uri for rdfs:seeAlso")
	}

	// TestCase 2: invalid predicate for licenseException
	// TestCase 1: invalid value for rdf:seeAlso
	parser, _ = parserFromBodyContent(`
		<spdx:LicenseException>
			<spdx:example></spdx:example>
			<spdx:licenseExceptionId>Libtool-exception</spdx:licenseExceptionId>
			<rdfs:unknown></rdfs:unknown>
		</spdx:LicenseException>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getLicenseExceptionFromNode(node)
	if err == nil {
		t.Errorf("should've raised an error due to invalid predicate")
	}

	// TestCase 3: everything valid
	// TestCase 1: invalid value for rdf:seeAlso
	parser, _ = parserFromBodyContent(`
		<spdx:LicenseException>
			<rdfs:seeAlso rdf:resource="http://www.opensource.org/licenses/GPL-3.0"/>
			<spdx:example>no example</spdx:example>
			<spdx:licenseExceptionId>Libtool-exception</spdx:licenseExceptionId>
			<rdfs:comment>no comments</rdfs:comment>
			<spdx:licenseExceptionText>text</spdx:licenseExceptionText>
			<spdx:name>name</spdx:name>
		</spdx:LicenseException>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	licenseException, err = parser.getLicenseExceptionFromNode(node)
	if err != nil {
		t.Fatalf("unexpected error while parsing a valid licenseException")
	}
	expectedCrossReference := "http://www.opensource.org/licenses/GPL-3.0"
	if licenseException.seeAlso != expectedCrossReference {
		t.Errorf("expected: %s, found: %s", expectedCrossReference, licenseException.seeAlso)
	}
	expectedExample := "no example"
	if licenseException.example != expectedExample {
		t.Errorf("expected: %s, got: %s", expectedExample, licenseException.example)
	}
	if licenseException.licenseExceptionId != "Libtool-exception" {
		t.Errorf("expected: %s, got: %s", "Libtool-exception", licenseException.licenseExceptionId)
	}
	if licenseException.comment != "no comments" {
		t.Errorf("expected: %s, got: %s", "no comments", licenseException.comment)
	}
	if licenseException.licenseExceptionText != "text" {
		t.Errorf("expected: '%s', got: '%s'", "text", licenseException.licenseExceptionText)
	}
	if licenseException.name != "name" {
		t.Errorf("expected: '%s', got: '%s'", "name", licenseException.name)
	}
}

func Test_rdfParser2_3_getLicenseFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var license License
	var err error

	// TestCase 1: isOsiApproved is not a valid boolean
	parser, _ = parserFromBodyContent(`
		<spdx:License>
			<spdx:isOsiApproved>no</spdx:isOsiApproved>
		</spdx:License>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	license, err = parser.getLicenseFromNode(node)
	if err == nil {
		t.Errorf("expected function to raise an error stating isOsiApproved should be a valid boolean type")
	}

	// TestCase 2: rdf:seeAlso not a valid uri must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:License>
			<rdfs:seeAlso>uri</rdfs:seeAlso>
		</spdx:License>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	license, err = parser.getLicenseFromNode(node)
	if err == nil {
		t.Errorf("expected function to raise an error stating invalid uri for rdfs:seeAlso")
	}

	// TestCase 3: isDeprecatedLicenseId is not a valid boolean
	parser, _ = parserFromBodyContent(`
		<spdx:License>
			<spdx:isDeprecatedLicenseId>yes</spdx:isDeprecatedLicenseId>
		</spdx:License>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	license, err = parser.getLicenseFromNode(node)
	if err == nil {
		t.Errorf("expected function to raise an error stating isDeprecatedLicenseId should be a valid boolean type")
	}

	// TestCase 4: isFsfLibre is not a valid boolean
	parser, _ = parserFromBodyContent(`
		<spdx:License>
			<spdx:isFsfLibre>no</spdx:isFsfLibre>
		</spdx:License>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	license, err = parser.getLicenseFromNode(node)
	if err == nil {
		t.Errorf("expected function to raise an error stating isFsfLibre should be a valid boolean type")
	}

	// TestCase 5: invalid triple for License:
	parser, _ = parserFromBodyContent(`
		<spdx:License>
			<spdx:unknown />
		</spdx:License>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	license, err = parser.getLicenseFromNode(node)
	if err == nil {
		t.Errorf("invalid predicate should've raised an error, got <nil>")
	}

	// TestCase 5: everything valid:
	parser, _ = parserFromBodyContent(`
		<spdx:License rdf:about="http://spdx.org/licenses/GPL-3.0-or-later">
			<rdfs:seeAlso>http://www.opensource.org/licenses/GPL-3.0</rdfs:seeAlso>
			<spdx:isOsiApproved>true</spdx:isOsiApproved>
			<spdx:licenseText>GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007</spdx:licenseText>
			<spdx:name>GNU General Public License v3.0 or later</spdx:name>
			<spdx:standardLicenseHeaderTemplate>...</spdx:standardLicenseHeaderTemplate>
			<spdx:licenseId>GPL-3.0-or-later</spdx:licenseId>
			<rdfs:comment>This license was released: 29 June 2007</rdfs:comment>
			<spdx:isFsfLibre>true</spdx:isFsfLibre>
			<spdx:standardLicenseHeader>...</spdx:standardLicenseHeader>
			<spdx:standardLicenseTemplate>....</spdx:standardLicenseTemplate>
		</spdx:License>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	license, err = parser.getLicenseFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid input: %v", err)
	}
	expectedSeeAlso := "http://www.opensource.org/licenses/GPL-3.0"
	if len(license.seeAlso) != 1 {
		t.Fatalf("expected seeAlso to have 1 element, got %d", len(license.seeAlso))
	}
	if license.seeAlso[len(license.seeAlso)-1] != expectedSeeAlso {
		t.Errorf("expected %s, got %s", expectedSeeAlso, license.seeAlso)
	}
	if license.isOsiApproved != true {
		t.Errorf("expected %t, got %t", true, license.isOsiApproved)
	}
	expectedLicenseText := "GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007"
	if license.licenseText != expectedLicenseText {
		t.Errorf("expected %s, got %s", expectedSeeAlso, license.licenseText)
	}
	expectedName := "GNU General Public License v3.0 or later"
	if license.name != expectedName {
		t.Errorf("expected %s, got %s", expectedName, license.name)
	}
	expectedstdLicHeader := "..."
	if license.standardLicenseHeader != expectedstdLicHeader {
		t.Errorf("expected %s, got %s", expectedstdLicHeader, license.standardLicenseHeader)
	}
	expectedLicenseId := "GPL-3.0-or-later"
	if expectedLicenseId != license.licenseID {
		t.Errorf("expected %s, got %s", expectedLicenseId, license.licenseID)
	}
	expectedLicenseComment := "This license was released: 29 June 2007"
	if expectedLicenseComment != license.comment {
		t.Errorf("expected %s, got %s", expectedLicenseComment, license.comment)
	}
	expectedstdLicTemplate := "..."
	if license.standardLicenseHeader != expectedstdLicTemplate {
		t.Errorf("expected %s, got %s", expectedstdLicTemplate, license.standardLicenseTemplate)
	}
	expectedstdLicHeaderTemplate := "..."
	if license.standardLicenseHeaderTemplate != expectedstdLicHeaderTemplate {
		t.Errorf("expected %s, got %s", expectedstdLicHeaderTemplate, license.standardLicenseHeaderTemplate)
	}
	if license.isFsfLibre != true {
		t.Errorf("expected %t, got %t", true, license.isFsfLibre)
	}
}

func Test_rdfParser2_3_getOrLaterOperatorFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var err error

	// TestCase 1: more than one member in the OrLaterOperator tag must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:OrLaterOperator>
			<spdx:member>
				<spdx:SimpleLicensingInfo>
					<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
					<spdx:name>freeware</spdx:name>
				</spdx:SimpleLicensingInfo>
			</spdx:member>
			<spdx:member rdf:resource="http://spdx.org/licenses/Nokia"/>
		</spdx:OrLaterOperator>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getOrLaterOperatorFromNode(node)
	if err == nil {
		t.Error("expected an error due to more than one members, got <nil>")
	}

	// TestCase 2: Invalid predicate must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:OrLaterOperator>
			<spdx:members>
				<spdx:SimpleLicensingInfo>
					<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
					<spdx:name>freeware</spdx:name>
				</spdx:SimpleLicensingInfo>
			</spdx:members>
		</spdx:OrLaterOperator>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getOrLaterOperatorFromNode(node)
	if err == nil {
		t.Error("expected an error due to invalid predicate, got <nil>")
	}

	// TestCase 5: invalid member
	parser, _ = parserFromBodyContent(`
		<spdx:OrLaterOperator>
			<spdx:member>
				<spdx:SimpleLicensingInfo>
					<spdx:invalidTag />
					<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
					<spdx:name>freeware</spdx:name>
				</spdx:SimpleLicensingInfo>
			</spdx:member>
		</spdx:OrLaterOperator>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getOrLaterOperatorFromNode(node)
	if err == nil {
		t.Errorf("expected an error parsing invalid license member, got %v", err)
	}

	// TestCase 4: valid input
	parser, _ = parserFromBodyContent(`
		<spdx:OrLaterOperator>
			<spdx:member>
				<spdx:SimpleLicensingInfo>
					<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
					<spdx:name>freeware</spdx:name>
				</spdx:SimpleLicensingInfo>
			</spdx:member>
		</spdx:OrLaterOperator>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getOrLaterOperatorFromNode(node)
	if err != nil {
		t.Errorf("unexpected error parsing a valid input: %v", err)
	}
}

func Test_rdfParser2_3_getSimpleLicensingInfoFromNode(t *testing.T) {
	// nothing to test. The just provides an interface to call function that
	// uses triples to render a SimpleLicensingInfo.
	parser, _ := parserFromBodyContent(`
		<spdx:SimpleLicensingInfo>
			<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
			<spdx:name>freeware</spdx:name>
		</spdx:SimpleLicensingInfo>
	`)
	node := parser.gordfParserObj.Triples[0].Subject
	_, err := parser.getSimpleLicensingInfoFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid input: %v", err)
	}
}

func Test_rdfParser2_3_getSimpleLicensingInfoFromTriples(t *testing.T) {
	var parser *rdfParser2_3
	var err error
	var license SimpleLicensingInfo

	// TestCase 1: invalid rdf:seeAlso attribute
	parser, _ = parserFromBodyContent(`
		<spdx:SimpleLicensingInfo>
			<rdfs:seeAlso>an invalid uri</rdfs:seeAlso>
		</spdx:SimpleLicensingInfo>
    `)
	_, err = parser.getSimpleLicensingInfoFromTriples(parser.gordfParserObj.Triples)
	if err == nil {
		t.Error("expected an error reporting invalid uri for rdf:seeAlso, got <nil>")
	}

	// TestCase 2: invalid predicate must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:SimpleLicensingInfo>
			<rdfs:invalidPredicate />
		</spdx:SimpleLicensingInfo>
    `)
	_, err = parser.getSimpleLicensingInfoFromTriples(parser.gordfParserObj.Triples)
	if err == nil {
		t.Error("expected an error reporting invalid predicate, got <nil>")
	}

	// TestCase 3: valid example
	parser, _ = parserFromBodyContent(`
		<spdx:SimpleLicensingInfo>
			<rdfs:comment>comment</rdfs:comment>
			<spdx:licenseId>lid</spdx:licenseId>
			<spdx:name>name</spdx:name>
			<rdfs:seeAlso>https://opensource.org/licenses/MPL-1.0</rdfs:seeAlso>
			<spdx:example>example</spdx:example>
		</spdx:SimpleLicensingInfo>
    `)
	license, err = parser.getSimpleLicensingInfoFromTriples(parser.gordfParserObj.Triples)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedComment := "comment"
	expectedLicenseId := "lid"
	expectedName := "name"
	expectedSeeAlso := "https://opensource.org/licenses/MPL-1.0"
	expectedExample := "example"
	if expectedComment != license.comment {
		t.Errorf("expected %v, got %v", expectedComment, license.comment)
	}
	if expectedLicenseId != license.licenseID {
		t.Errorf("expected %v, got %v", expectedLicenseId, license.licenseID)
	}
	if expectedName != license.name {
		t.Errorf("expected %v, got %v", expectedName, license.name)
	}
	if len(license.seeAlso) != 1 {
		t.Fatalf("expected seeAlso to have 1 element, found %d", len(license.seeAlso))
	}
	if license.seeAlso[0] != expectedSeeAlso {
		t.Errorf("expected %v, got %v", expectedSeeAlso, license.seeAlso[0])
	}
	if license.example != expectedExample {
		t.Errorf("expected %v, got %v", expectedExample, license.example)
	}
}

func Test_rdfParser2_3_getSpecialLicenseFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var license SpecialLicense

	// TestCase 1: NONE
	parser, _ = parserFromBodyContent(``)
	node = &gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       NS_SPDX + "NONE",
	}
	license, err := parser.getSpecialLicenseFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid node: %v", err)
	}
	if license.value != "NONE" {
		t.Errorf("expected %s, got %s", "NONE", license.value)
	}

	// TestCase 2: NOASSERTION
	parser, _ = parserFromBodyContent(``)
	node = &gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       NS_SPDX + "NOASSERTION",
	}
	license, err = parser.getSpecialLicenseFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid node: %v", err)
	}
	if license.value != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", license.value)
	}

	// TestCase 4: undefined standard license
	parser, _ = parserFromBodyContent(``)
	node = &gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       "https://opensource.org/licenses/unknown",
	}
	_, err = parser.getSpecialLicenseFromNode(node)
	if err == nil {
		t.Errorf("expected an error saying invalid license")
	}

	// TestCase 4: valid standard license
	parser, _ = parserFromBodyContent(``)
	node = &gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       "https://opensource.org/licenses/MPL-1.0",
	}
	license, err = parser.getSpecialLicenseFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid node: %v", err)
	}
	if license.value != "MPL-1.0" {
		t.Errorf("expected %s, got %s", "MPL-1.0", license.value)
	}
}

func Test_rdfParser2_3_getWithExceptionOperatorFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var err error

	// TestCase 1: more than one member in the OrLaterOperator tag must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:WithExceptionOperator>
			<spdx:member>
				<spdx:SimpleLicensingInfo>
					<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
					<spdx:name>freeware</spdx:name>
				</spdx:SimpleLicensingInfo>
			</spdx:member>
			<spdx:member rdf:resource="http://spdx.org/licenses/Nokia"/>
		</spdx:WithExceptionOperator>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getWithExceptionOperatorFromNode(node)
	if err == nil {
		t.Error("expected an error due to more than one members, got <nil>")
	}

	// TestCase 2: Invalid predicate must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:WithExceptionOperator>
			<spdx:members>
				<spdx:SimpleLicensingInfo>
					<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
					<spdx:name>freeware</spdx:name>
				</spdx:SimpleLicensingInfo>
			</spdx:members>
		</spdx:WithExceptionOperator>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getWithExceptionOperatorFromNode(node)
	if err == nil {
		t.Error("expected an error due to invalid predicate, got <nil>")
	}

	// TestCase 3: Invalid member
	parser, _ = parserFromBodyContent(`
		<spdx:WithExceptionOperator>
			<spdx:member>
				<spdx:License rdf:about="http://spdx.org/licenses/GPL-2.0-or-later">
					<spdx:unknownTag />
				</spdx:License>
			</spdx:member>
		</spdx:WithExceptionOperator>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getWithExceptionOperatorFromNode(node)
	if err == nil {
		t.Error("expected an error due to error parsing a member, got <nil>")
	}

	// TestCase 4: Invalid licenseException
	parser, _ = parserFromBodyContent(`
		<spdx:WithExceptionOperator>
			<spdx:member>
				<spdx:License rdf:about="http://spdx.org/licenses/GPL-2.0-or-later"/>
			</spdx:member>
			<spdx:licenseException>
				<spdx:LicenseException>
					<spdx:invalidTag />
					<spdx:example>example</spdx:example>
					<spdx:licenseExceptionId>Libtool-exception</spdx:licenseExceptionId>
					<rdfs:comment>comment</rdfs:comment>
				</spdx:LicenseException>
			</spdx:licenseException>
		</spdx:WithExceptionOperator>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getWithExceptionOperatorFromNode(node)
	if err == nil {
		t.Error("expected an error due to invalid licenseException, got <nil>")
	}

	// TestCase 5: valid input
	parser, _ = parserFromBodyContent(`
		<spdx:WithExceptionOperator>
			<spdx:member>
				<spdx:License rdf:about="http://spdx.org/licenses/GPL-2.0-or-later"/>
			</spdx:member>
			<spdx:licenseException>
				<spdx:LicenseException>
					<spdx:example>example</spdx:example>
					<spdx:licenseExceptionId>Libtool-exception</spdx:licenseExceptionId>
					<rdfs:comment>comment</rdfs:comment>
				</spdx:LicenseException>
			</spdx:licenseException>
		</spdx:WithExceptionOperator>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getWithExceptionOperatorFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid input: %v", err)
	}
}
