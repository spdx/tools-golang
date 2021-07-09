// Package saver2v2 contains functions to render and write a json
// formatted version of an in-memory SPDX document and its sections
// (version 2.2).
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package saver2v2

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestRenderDocument2_2(t *testing.T) {

	test1 := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:       "SPDX-2.2",
			DataLicense:       "CC0-1.0",
			SPDXIdentifier:    spdx.ElementID("DOCUMENT"),
			DocumentName:      "tools-golang-0.0.1.abcdef",
			DocumentNamespace: "https://github.com/spdx/spdx-docs/tools-golang/tools-golang-0.0.1.abcdef.whatever",
			CreatorPersons: []string{
				"John Doe ()",
			},
			CreatorOrganizations: []string{"ExampleCodeInspect"},
			CreatorTools:         []string{"LicenseFind-1.0"},
			DocumentComment:      "This document was created using SPDX 2.0 using licenses from the web site.",
			LicenseListVersion:   "3.8",
			CreatorComment:       "This package has been shipped in source and binary form.\nThe binaries were created with gcc 4.5.1 and expect to link to\ncompatible system run time libraries.",
			Created:              "2018-10-10T06:20:00Z",
			ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{
				"spdx-tool-1.2": {
					DocumentRefID: "spdx-tool-1.2",
					URI:           "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
					Alg:           "SHA1",
					Checksum:      "d6a770ba38583ed4bb4525bd96e50461655d2759",
				},
			},
		},
		OtherLicenses: []*spdx.OtherLicense2_2{
			{
				ExtractedText:     "\"THE BEER-WARE LICENSE\" (Revision 42):\nphk@FreeBSD.ORG wrote this file. As long as you retain this notice you\ncan do whatever you want with this stuff. If we meet some day, and you think this stuff is worth it, you can buy me a beer in return Poul-Henning Kamp  </\nLicenseName: Beer-Ware License (Version 42)\nLicenseCrossReference:  http://people.freebsd.org/~phk/\nLicenseComment: \nThe beerware license has a couple of other standard variants.",
				LicenseIdentifier: "LicenseRef-Beerware-4.2",
			},
			{
				ExtractedText:     "/*\n * (c) Copyright 2009 University of Bristol\n * All rights reserved.\n *\n * Redistribution and use in source and binary forms, with or without\n * modification, are permitted provided that the following conditions\n * are met:\n * 1. Redistributions of source code must retain the above copyright\n *    notice, this list of conditions and the following disclaimer.\n * 2. Redistributions in binary form must reproduce the above copyright\n *    notice, this list of conditions and the following disclaimer in the\n *    documentation and/or other materials provided with the distribution.\n * 3. The name of the author may not be used to endorse or promote products\n *    derived from this software without specific prior written permission.\n *\n * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR\n * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\n * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.\n * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,\n * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT\n * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\n * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\n * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF\n * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n*/",
				LicenseIdentifier: "LicenseRef-4",
			},
			{
				LicenseComment:         "This is tye CyperNeko License",
				ExtractedText:          "The CyberNeko Software License, Version 1.0\n\n \n(C) Copyright 2002-2005, Andy Clark.  All rights reserved.\n \nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions\nare met:\n\n1. Redistributions of source code must retain the above copyright\n   notice, this list of conditions and the following disclaimer. \n\n2. Redistributions in binary form must reproduce the above copyright\n   notice, this list of conditions and the following disclaimer in\n   the documentation and/or other materials provided with the\n   distribution.\n\n3. The end-user documentation included with the redistribution,\n   if any, must include the following acknowledgment:  \n     \"This product includes software developed by Andy Clark.\"\n   Alternately, this acknowledgment may appear in the software itself,\n   if and wherever such third-party acknowledgments normally appear.\n\n4. The names \"CyberNeko\" and \"NekoHTML\" must not be used to endorse\n   or promote products derived from this software without prior \n   written permission. For written permission, please contact \n   andyc@cyberneko.net.\n\n5. Products derived from this software may not be called \"CyberNeko\",\n   nor may \"CyberNeko\" appear in their name, without prior written\n   permission of the author.\n\nTHIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESSED OR IMPLIED\nWARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\nOF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\nDISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR OTHER CONTRIBUTORS\nBE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, \nOR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT \nOF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR \nBUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, \nWHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE \nOR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, \nEVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
				LicenseIdentifier:      "LicenseRef-3",
				LicenseName:            "CyberNeko License",
				LicenseCrossReferences: []string{"http://people.apache.org/~andyc/neko/LICENSE", "http://justasample.url.com"},
			},
			{
				ExtractedText:     "This package includes the GRDDL parser developed by Hewlett Packard under the following license:\n� Copyright 2007 Hewlett-Packard Development Company, LP\n\nRedistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met: \n\nRedistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer. \nRedistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution. \nThe name of the author may not be used to endorse or promote products derived from this software without specific prior written permission. \nTHIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
				LicenseIdentifier: "LicenseRef-2",
			},
			{
				ExtractedText:     "/*\n * (c) Copyright 2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009 Hewlett-Packard Development Company, LP\n * All rights reserved.\n *\n * Redistribution and use in source and binary forms, with or without\n * modification, are permitted provided that the following conditions\n * are met:\n * 1. Redistributions of source code must retain the above copyright\n *    notice, this list of conditions and the following disclaimer.\n * 2. Redistributions in binary form must reproduce the above copyright\n *    notice, this list of conditions and the following disclaimer in the\n *    documentation and/or other materials provided with the distribution.\n * 3. The name of the author may not be used to endorse or promote products\n *    derived from this software without specific prior written permission.\n *\n * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR\n * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\n * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.\n * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,\n * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT\n * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\n * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\n * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF\n * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n*/",
				LicenseIdentifier: "LicenseRef-1",
			},
		},
		Annotations: []*spdx.Annotation2_2{
			{
				AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "File"},
				AnnotationDate:           "2011-01-29T18:30:22Z",
				AnnotationType:           "OTHER",
				AnnotatorType:            "Person",
				Annotator:                "File Commenter",
				AnnotationComment:        "File level annotation",
			},
			{
				AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "Package"},
				AnnotationDate:           "2011-01-29T18:30:22Z",
				AnnotationType:           "OTHER",
				AnnotatorType:            "Person",
				Annotator:                "Package Commenter",
				AnnotationComment:        "Package level annotation",
			},
			{
				AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
				AnnotationDate:           "2010-02-10T00:00:00Z",
				AnnotationType:           "REVIEW",
				AnnotatorType:            "Person",
				Annotator:                "Joe Reviewer",
				AnnotationComment:        "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses",
			},
			{
				AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
				AnnotationDate:           "2011-03-13T00:00:00Z",
				AnnotationType:           "REVIEW",
				AnnotatorType:            "Person",
				Annotator:                "Suzanne Reviewer",
				AnnotationComment:        "Another example reviewer.",
			},
			{
				AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
				AnnotationDate:           "2010-01-29T18:30:22Z",
				AnnotationType:           "OTHER",
				AnnotatorType:            "Person",
				Annotator:                "Jane Doe ()",
				AnnotationComment:        "Document level annotation",
			},
		},
	}
	var b []byte

	type args struct {
		doc *spdx.Document2_2
		buf *bytes.Buffer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success test ",
			args: args{
				doc: test1,
				buf: bytes.NewBuffer(b),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RenderDocument2_2(tt.args.doc, tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("RenderDocument2_2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
