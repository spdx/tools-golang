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
	obj, err := p.requestElementType(node, typeChecksum)
	if err != nil {
		return nil, err
	}
	return obj.(*Checksum), err
}

func (p *Parser) MapChecksum(cksum *Checksum) *builder {
	builder := &builder{t: typeChecksum, ptr: cksum}
	key := false
	builder.updaters = map[string]updater{
		"algorithm": func(obj goraptor.Term) error {
			if key {
				return fmt.Errorf("Algorithm defined already.")
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

// Takes in the checksum, compares it's algo with a string, if matches returns the algo
func AlgoIdentifier(cksum *Checksum, t string) string {
	algo := ExtractChecksumAlgo(cksum.Algorithm.Val)
	if strings.Contains(algo, t) {
		return t
	}
	return ""
}
