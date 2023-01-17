// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"reflect"
	"testing"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
)

func Test_rdfParser2_3_getExtractedLicensingInfoFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var err error
	var node *gordfParser.Node

	// TestCase 1: invalid predicate must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:ExtractedLicensingInfo rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#LicenseRef-Freeware">
			<spdx:licenseID>LicenseRef-Freeware</spdx:licenseID>
			<spdx:name>freeware</spdx:name>
			<spdx:extractedText><![CDATA[Software classified as freeware is licensed at no cost and is either fully functional for an unlimited time; or has only basic functions enabled with a fully functional version available commercially or as shareware.[8] In contrast to free software, the author usually restricts one or more rights of the user, including the rights to use, copy, distribute, modify and make derivative works of the software or extract the source code.[1][2][9][10] The software license may impose various additional restrictions on the type of use, e.g. only for personal use, private use, individual use, non-profit use, non-commercial use, academic use, educational use, use in charity or humanitarian organizations, non-military use, use by public authorities or various other combinations of these type of restrictions.[11] For instance, the license may be "free for private, non-commercial use". The software license may also impose various other restrictions, such as restricted use over a network, restricted use on a server, restricted use in a combination with some types of other software or with some hardware devices, prohibited distribution over the Internet other than linking to author's website, restricted distribution without author's consent, restricted number of copies, etc.]]></spdx:extractedText>
		</spdx:ExtractedLicensingInfo>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getExtractedLicensingInfoFromNode(node)
	if err == nil {
		t.Errorf("expected an error saying invalid predicate, got <nil>")
	}

	// TestCase 2: valid input
	parser, _ = parserFromBodyContent(`
		<spdx:ExtractedLicensingInfo rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#LicenseRef-Freeware">
			<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
			<spdx:name>freeware</spdx:name>
			<spdx:extractedText><![CDATA[Software classified as freeware is licensed at no cost and is either fully functional for an unlimited time; or has only basic functions enabled with a fully functional version available commercially or as shareware.[8] In contrast to free software, the author usually restricts one or more rights of the user, including the rights to use, copy, distribute, modify and make derivative works of the software or extract the source code.[1][2][9][10] The software license may impose various additional restrictions on the type of use, e.g. only for personal use, private use, individual use, non-profit use, non-commercial use, academic use, educational use, use in charity or humanitarian organizations, non-military use, use by public authorities or various other combinations of these type of restrictions.[11] For instance, the license may be "free for private, non-commercial use". The software license may also impose various other restrictions, such as restricted use over a network, restricted use on a server, restricted use in a combination with some types of other software or with some hardware devices, prohibited distribution over the Internet other than linking to author's website, restricted distribution without author's consent, restricted number of copies, etc.]]></spdx:extractedText>
		</spdx:ExtractedLicensingInfo>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getExtractedLicensingInfoFromNode(node)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_rdfParser2_3_extractedLicenseToOtherLicense(t *testing.T) {
	// nothing to test for this function.
	parser, _ := parserFromBodyContent(`
		<spdx:ExtractedLicensingInfo rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#LicenseRef-Freeware">
			<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
			<spdx:name>freeware</spdx:name>
			<spdx:extractedText><![CDATA[Software classified as freeware is licensed at no cost and is either fully functional for an unlimited time; or has only basic functions enabled with a fully functional version available commercially or as shareware.[8] In contrast to free software, the author usually restricts one or more rights of the user, including the rights to use, copy, distribute, modify and make derivative works of the software or extract the source code.[1][2][9][10] The software license may impose various additional restrictions on the type of use, e.g. only for personal use, private use, individual use, non-profit use, non-commercial use, academic use, educational use, use in charity or humanitarian organizations, non-military use, use by public authorities or various other combinations of these type of restrictions.[11] For instance, the license may be "free for private, non-commercial use". The software license may also impose various other restrictions, such as restricted use over a network, restricted use on a server, restricted use in a combination with some types of other software or with some hardware devices, prohibited distribution over the Internet other than linking to author's website, restricted distribution without author's consent, restricted number of copies, etc.]]></spdx:extractedText>
		</spdx:ExtractedLicensingInfo>
	`)
	node := parser.gordfParserObj.Triples[0].Subject
	extLicense, _ := parser.getExtractedLicensingInfoFromNode(node)
	othLic := parser.extractedLicenseToOtherLicense(extLicense)

	if othLic.LicenseIdentifier != extLicense.licenseID {
		t.Errorf("expected %v, got %v", othLic.LicenseIdentifier, extLicense.licenseID)
	}
	if othLic.ExtractedText != extLicense.extractedText {
		t.Errorf("expected %v, got %v", othLic.ExtractedText, extLicense.extractedText)
	}
	if othLic.LicenseComment != extLicense.comment {
		t.Errorf("expected %v, got %v", othLic.LicenseComment, extLicense.comment)
	}
	if !reflect.DeepEqual(othLic.LicenseCrossReferences, extLicense.seeAlso) {
		t.Errorf("expected %v, got %v", othLic.LicenseCrossReferences, extLicense.seeAlso)
	}
	if othLic.LicenseName != extLicense.name {
		t.Errorf("expected %v, got %v", othLic.LicenseName, extLicense.name)
	}
}
