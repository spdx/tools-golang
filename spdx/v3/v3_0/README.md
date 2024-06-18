# SPDX 3 Model Patterns

This is a non-exhaustive list of some patterns we are able to use in Golang to work with an SPDX 3 data model.

I've included some variations I think we should consider. An explanation is as follows, along with some pros/cons:

- [embedded first](embedded_first) - lean into modifying go structs directly with exported fields, using struct embedding.
This requires using interfaces, but the interfaces are simple single-functions to get an editable embedded struct.
  - pro: interfaces are simple
  - pro: modifying structs is simple and idiomatic
  - pro: minimal surface area: very few additional functions necessary
  - con: construction is confusing -- which nesting has field X?
  - con: naming is hard, need `Package` and `IPackage`


- [interfaces only](interfaces_only) - only export interfaces and use interfaces to interact with structs.
NOTE: this also includes a variant that setters return `error`, which could be applied to other options.
This option only exports the interface types and simple no-arg constructor functions.
  - pro: idiomatic getter naming (without `Get`)
  - pro: single type exported, data structs are not
  - pro: minimal surface area: no extra helper functions/structs that we may decide to change later
  - con: absolutely no JSON/etc. output
  - con: getter/setter still not especially idiomatic
  - con: duplicated code because no embedding
  - con: construction is tedious and could be error-prone if referencing wrong variable name


- [interfaces only with constructors](interfaces_only/constructors)
  - (generally same pros/cons as above), but:
  - pro: much simpler construction pattern for users
  - pro: single validation spot during construction (with returning `error` pattern)
  - con: more structs and functions exported which cannot be changed without breaking backwards compatibility


- [inverted](inverted) - since we know the entire data model and it does not need to be extended,
  it is possible to implement an uber-struct without the need to use interfaces at all.
I thought I'd mention this for some semblance of completeness, but this is probably a very bad idea. 
  - pro: single set of types exported
  - pro: fairly easy to interact with existing documents
  - con: easy to construct incorrectly
  - con: potential unneecssary memory usage


- [embedded with interfaces and constructors](embedded_with_interfaces_and_constructors) - takes the approach of both 
embedding and interfaces, along with constructor functions. This makes creation reasonably simple and working with
existing documents reasonably simple.
  - pro: probably the most idiomatic Go of these options
  - pro: simple construction
  - pro: reasonably simple to work with existing documents
  - pro: eliminates some code duplication seen with `interfaces_only`
  - con: getter/setter still not especially idiomatic
  - con: larger surface area: at least 4 exported parts for each type
  - con: absolutely no JSON/etc. output
