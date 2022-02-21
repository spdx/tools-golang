SPDX-License-Identifier: CC-BY-4.0

## Working

The spdx document is converted to map[string]interface{} and then the entire map is converted to json using a single json Marshall function call . The saver uses a tempoarary storage to store all the files (Paackaged and Unpackaged) together in a single data structure in order to comply with the json schema defined by spdx .

spdx.Document2_2  →  map[string]interface{}  → JSON

## Some Key Points

- The packages have a property "hasFiles" defined in the schema which is an array of the SPDX Identifiers of the files of that pacakge . The saver iterates through the files of a package and inserted all the SPDX Identifiers of the files in the "hasFiles" array . In addition it adds the file to a temporary storage map to store all the files of the entire document at a single place .

- The files require the packages to be saved before them in order to ensure that the packaged files are added to the temporary storage before the files are saved .

- The snippets are saved after the files and a property "snippetFromFile" identifies the file of the snippets.

The json file loader in `package jsonsaver` makes the following assumptions:


### Order of appearance of the properties
* The saver does not make any pre-assumptions based on the order in which the properties are saved . 


### Annotations
* The json spdx schema does not define the SPDX Identifier property for the annotation object . The saver inserts the annotation inside the element who spdx identifier mathches the annotation SPDX identifier .

### Indentation
* The jsonsaver uses the marshall indent function with "" as he prefix and "\t" as the indent character  , passed as funtion parameters .