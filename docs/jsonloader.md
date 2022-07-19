SPDX-License-Identifier: CC-BY-4.0

## Usage

A json.Unmarshal function on the v2_2.Document struct is defined so that when the JSON is unmarshalled, the function is called and the JSON can be processed in a custom way. Then a new map[string]interface{} is defined which temporarily holds the unmarshalled JSON. The map is then parsed into the v2_2.Document using functions defined for each different section.

JSON => map[string]interface{} => v2_2.Document

## Some Key Points 

- The packages have a property "hasFiles" defined in the schema which is an array of the SPDX Identifiers of the files of that package. The parser first parses all the files into the UnpackagedFiles map of the document and then when it parses the Packages, it removes the respective files from the UnpackagedFiles map and places them inside the Files map of the corresponding package.

- The snippets have a property "snippetFromFile" which has the SPDX identifier of the file to which the snippet is related. Thus the snippets require the files to be parsed before them. Then the snippets are parsed one by one and inserted into the respective files using this property.

## Assumptions

The json file loader in `package jsonloader` makes the following assumptions:

### Order of appearance of the properties 
* The parser does not make any assumptions based on the order in which the properties appear.

### Annotations
* The JSON SPDX schema does not define the SPDX Identifier property for the annotation object. The parser assumes the SPDX Identifier of the parent property of the currently-being-parsed annotation array to be the SPDX Identifer for all the annotation objects of that array.
