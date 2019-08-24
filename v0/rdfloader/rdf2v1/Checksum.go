// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"fmt"
	"strings"

	"github.com/deltamobile/goraptor"
)

type Checksum struct {
	Algorithm     ValueStr
	ChecksumValue ValueStr
}

func (p *Parser) requestChecksum(node goraptor.Term) (*Checksum, error) {
	obj, err := p.requestElementType(node, TypeChecksum)
	if err != nil {
		return nil, err
	}
	return obj.(*Checksum), err
}

func (p *Parser) MapChecksum(cksum *Checksum) *builder {
	builder := &builder{t: TypeChecksum, ptr: cksum}
	key := false
	builder.updaters = map[string]updater{
		"algorithm": func(obj goraptor.Term) error {
			if key {
				return fmt.Errorf("Algorithm set already.")
			}
			algostr := termStr(obj)
			cksum.Algorithm.Val = ExtractChecksumAlgo(algostr)
			key = true
			return nil
		},
		"checksumValue": update(&cksum.ChecksumValue),
	}
	return builder
}

func ExtractChecksumAlgo(str string) string {
	str = strings.Replace(str, "http://spdx.org/rdf/terms#checksumAlgorithm_", "", 1)
	str = strings.ToUpper(str)
	return str
}

func InsertChecksumAlgo(str string) string {
	str = strings.ToLower(str)
	str = "http://spdx.org/rdf/terms#checksumAlgorithm_" + str
	return str
}

// Takes in the checksum, compares it's algo with a string, if matches returns the Value
func AlgoValue(cksum *Checksum, t string) string {
	algo := ExtractChecksumAlgo(cksum.Algorithm.Val)
	if strings.Contains(algo, t) {
		return cksum.ChecksumValue.Val
	}
	return ""
}
