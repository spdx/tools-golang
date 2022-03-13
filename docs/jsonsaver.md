SPDX-License-Identifier: CC-BY-4.0

## Working

The SPDX document is converted to map[string]interface{} and then the entire map is converted to JSON using a single json.MarshalIndent function call. The saver uses temporary memory to store all the files (Packaged and Unpackaged) together in a single data structure in order to comply with the JSON schema defined by SPDX.

spdx.Document2_2 => map[string]interface{} => JSON

## Some Key Points

- The packages have a property "hasFiles" defined in the schema which is an array of the SPDX Identifiers of the files of that package. The saver iterates through the files of a package and inserts all the SPDX Identifiers of the files in the "hasFiles" array. In addition it adds each file to a temporary storage map to store all the files of the entire document at a single place.

- The files require the packages to be saved before them in order to ensure that the packaged files are added to the temporary storage before the files are saved.

- The snippets are saved after the files and a property "snippetFromFile" identifies the file containing each snippet.

## Assumptions

The json file loader in `package jsonsaver` makes the following assumptions:

### Order of appearance of the properties
* The saver does not make any pre-assumptions based on the order in which the properties are saved.

### Annotations
* The JSON SPDX schema does not define the SPDX Identifier property for the annotation object. The saver inserts the annotation inside the element whose SPDX Identifier matches the annotation's SPDX Identifier.

### Indentation
* The jsonsaver uses the MarshalIndent function with "" as the prefix and "\t" as the indent character, passed as function parameters.
