package jsonloader2v2

import (
	"fmt"
	"reflect"

	"github.com/spdx/tools-golang/spdx"
)

func (spec JSONSpdxDocument) parseJsonSnippets2_2(key string, value interface{}, doc *spdxDocument2_2) error {

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		snippets := reflect.ValueOf(value)
		for i := 0; i < snippets.Len(); i++ {
			snippetmap := snippets.Index(i).Interface().(map[string]interface{})
			// create a new package
			snippet := &spdx.Snippet2_2{}
			//extract the SPDXID of the package
			eID, err := extractElementID(snippetmap["SPDXID"].(string))
			if err != nil {
				return fmt.Errorf("%s", err)
			}
			snippet.SnippetSPDXIdentifier = eID
			//range over all other properties now
			for k, v := range snippetmap {
				switch k {
				case "SPDXID", "snippetFromFile":
					//redundant case
				case "name":
					snippet.SnippetName = v.(string)
				case "copyrightText":
					snippet.SnippetCopyrightText = v.(string)
				case "licenseComments":
					snippet.SnippetLicenseComments = v.(string)
				case "licenseConcluded":
					snippet.SnippetLicenseConcluded = v.(string)
				case "licenseInfoInSnippets":
					if reflect.TypeOf(v).Kind() == reflect.Slice {
						info := reflect.ValueOf(v)
						for i := 0; i < info.Len(); i++ {
							snippet.LicenseInfoInSnippet = append(snippet.LicenseInfoInSnippet, info.Index(i).Interface().(string))
						}
					}
				case "attributionTexts":
					if reflect.TypeOf(v).Kind() == reflect.Slice {
						info := reflect.ValueOf(v)
						for i := 0; i < info.Len(); i++ {
							snippet.SnippetAttributionTexts = append(snippet.SnippetAttributionTexts, info.Index(i).Interface().(string))
						}
					}
				case "comment":
					snippet.SnippetComment = v.(string)
				case "ranges":
					//TODO: optimise this logic
					if reflect.TypeOf(v).Kind() == reflect.Slice {
						info := reflect.ValueOf(v)
						lineRanges := info.Index(0).Interface().(map[string]interface{})
						lineRangeStart := lineRanges["endPointer"].(map[string]interface{})
						lineRangeEnd := lineRanges["startPointer"].(map[string]interface{})
						snippet.SnippetLineRangeStart = int(lineRangeStart["lineNumber"].(float64))
						snippet.SnippetLineRangeEnd = int(lineRangeEnd["lineNumber"].(float64))

						byteRanges := info.Index(1).Interface().(map[string]interface{})
						byteRangeStart := byteRanges["endPointer"].(map[string]interface{})
						byteRangeEnd := byteRanges["startPointer"].(map[string]interface{})
						snippet.SnippetLineRangeStart = int(byteRangeStart["offset"].(float64))
						snippet.SnippetLineRangeEnd = int(byteRangeEnd["offset"].(float64))
					}
				default:
					return fmt.Errorf("received unknown tag %v in files section", k)
				}
			}
			fileID, err2 := extractElementID(snippetmap["snippetFromFile"].(string))
			if err2 != nil {
				return fmt.Errorf("%s", err2)
			}
			doc.UnpackagedFiles[fileID].Snippets[eID] = snippet
		}

	}
	return nil
}
