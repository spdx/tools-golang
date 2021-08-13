// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestJSONSpdxDocument_parseJsonOtherLicenses2_2(t *testing.T) {

	jsonData, err := ioutil.ReadFile("jsonfiles/otherlicensestest.json") // b has type []byte
	if err != nil {
		log.Fatal(err)
	}

	otherLicensestest1 := []*spdx.OtherLicense2_2{
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
	}

	otherLicensestest2 := []byte(`{
		"hasExtractedLicensingInfos":[{
			"extractedText" : "\"THE BEER-WARE LICENSE\" (Revision 42):\nphk@FreeBSD.ORG wrote this file. As long as you retain this notice you\ncan do whatever you want with this stuff. If we meet some day, and you think this stuff is worth it, you can buy me a beer in return Poul-Henning Kamp  </\nLicenseName: Beer-Ware License (Version 42)\nLicenseCrossReference:  http://people.freebsd.org/~phk/\nLicenseComment: \nThe beerware license has a couple of other standard variants.",
			"licensd" : "LicenseRef-Beerware-4.2"
		  }]
		}
  	`)

	var specs JSONSpdxDocument
	json.Unmarshal(jsonData, &specs)

	var specs2 JSONSpdxDocument
	json.Unmarshal(otherLicensestest2, &specs2)

	type args struct {
		key   string
		value interface{}
		doc   *spdxDocument2_2
	}
	tests := []struct {
		name    string
		spec    JSONSpdxDocument
		args    args
		wantErr bool
		want    []*spdx.OtherLicense2_2
	}{
		// TODO: Add test cases.
		{
			name: "successTest",
			spec: specs,
			args: args{
				key:   "hasExtractedLicensingInfos",
				value: specs["hasExtractedLicensingInfos"],
				doc:   &spdxDocument2_2{},
			},
			want:    otherLicensestest1,
			wantErr: false,
		},
		{
			name: "failure test --> unknown tag",
			spec: specs2,
			args: args{
				key:   "hasExtractedLicensingInfos",
				value: specs2["hasExtractedLicensingInfos"],
				doc:   &spdxDocument2_2{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.spec.parseJsonOtherLicenses2_2(tt.args.key, tt.args.value, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("JSONSpdxDocument.parseJsonOtherLicenses2_2() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				for i := 0; i < len(tt.want); i++ {
					if !reflect.DeepEqual(tt.args.doc.OtherLicenses[i], tt.want[i]) {
						t.Errorf("Load2_2() = %v, want %v", tt.args.doc.OtherLicenses[i], tt.want[i])
					}
				}
			}
		})
	}
}
