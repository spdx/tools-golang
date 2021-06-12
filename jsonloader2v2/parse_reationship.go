package jsonloader2v2

import (
	"fmt"
	"reflect"

	"github.com/spdx/tools-golang/spdx"
)

func (spec JSONSpdxDocument) parseJsonRelationships2_2(key string, value interface{}, doc *spdxDocument2_2) error {

	//FIXME : NOASSERTION and NONE in relationship B value not compatible
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		relationships := reflect.ValueOf(value)
		for i := 0; i < relationships.Len(); i++ {
			relationship := relationships.Index(i).Interface().(map[string]interface{})
			rel := spdx.Relationship2_2{}
			aid, err := extractDocElementID(relationship["spdxElementId"].(string))
			if err != nil {
				return fmt.Errorf("%s", err)
			}
			rel.RefA = aid

			bid, err := extractDocElementID(relationship["relatedSpdxElement"].(string))
			if err != nil {
				return fmt.Errorf("%s", err)
			}
			rel.RefB = bid

			if relationship["relationshipType"] == nil {
				return fmt.Errorf("%s , %d", "RelationshipType propty missing in relationship number", i)
			}
			rel.Relationship = relationship["relationshipType"].(string)

			if relationship["comment"] != nil {
				rel.RelationshipComment = relationship["comment"].(string)
			}

			doc.Relationships = append(doc.Relationships, &rel)
		}

	}
	return nil
}
