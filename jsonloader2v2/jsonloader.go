// Package jsonloader is used to load and parse SPDX JSON documents
// into tools-golang data structures.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package jsonloader2v2

import (
	"encoding/json"
	"fmt"
)

//TODO : return spdx.Document2_2
func Load2_2(content []byte) (*spdxDocument2_2, error) {
	// check whetehr the Json is valid or not
	if !json.Valid(content) {
		return nil, fmt.Errorf("%s", "Invalid JSON Specification")
	}
	result := spdxDocument2_2{}
	// unmarshall the json into the result struct
	err := json.Unmarshal(content, &result)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return &result, nil
}

func (doc *spdxDocument2_2) UnmarshalJSON(data []byte) error {
	var specs JSONSpdxDocument
	//unmarshall the json into the intermediate stricture map[string]interface{}
	err := json.Unmarshal(data, &specs)
	if err != nil {
		return err
	}
	// parse the data from the intermediate structure to the spdx.Document2_2{}
	err = specs.newDocument(doc)
	if err != nil {
		return err
	}
	return nil
}

func (spec JSONSpdxDocument) newDocument(doc *spdxDocument2_2) error {
	// raneg through all the keys in the map and send them to appropriate arsing functions
	for key, val := range spec {
		switch key {
		case "dataLicense", "spdxVersion", "SPDXID", "documentNamespace", "name", "comment", "creationInfo", "externalDocumentRefs":
			err := spec.parseJsonCreationInfo2_2(key, val, doc)
			if err != nil {
				return err
			}
		// case "annotations":
		// 	err := spec.parseJsonAnnotations2_2(key, val, doc)
		// 	if err != nil {
		// 		return err
		// 	}
		case "relationships":
			err := spec.parseJsonRelationships2_2(key, val, doc)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unrecognized key here %v", key)

		}

	}
	return nil
}
