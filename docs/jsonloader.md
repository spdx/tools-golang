SPDX-License-Identifier: CC-BY-4.0

## Working

A UnmarshallJSON function on the spdx.Document2_2 struct is defined so that when the JSON is unmarshalled in it the function is called and we can implement the process in a custom way . Then  a new map[string]interface{} is deifined which temporarily holds the unmarshalled json . The map is then parsed into the spdx.Document2_2 using functions defined for it’s different sections .

JSON  →  map[string]interface{}  → spdx.Document2_2

## Some Key Points 

- The packages have a property "hasFiles" defined in the schema which is an array of the SPDX Identifiers of the files of that pacakge . The parses first parses all the files into the Unpackaged files map of the document and then when it parses the packages , it removes the respective files from the unpackaged files map and places it inside the files map of that package .

- The snippets have a property "snippetFromFile" which has the SPDX identiifer of the file to which the snippet is related . Thus the snippets require the files to be parsed before them . Then the snippets are parsed one by one and inserted into the respective files using this property .


The json file loader in `package jsonloader` makes the following assumptions:


### Order of appearance of the properties 
* The parser does not make any pre-assumptions based on the order in which the properties appear . 


### Annotations
* The json spdx schema does not define the SPDX Identifier property for the annotation object . The parser assumes the spdx Identifier of the parent property of the currently being parsed annotation array to be the SPDX Identifer for all the annotation objects of that array.