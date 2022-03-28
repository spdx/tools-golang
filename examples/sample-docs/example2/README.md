# Example 2

## Description

```
content
├── build
│   └── hello
└── src
    ├── Makefile
    └── hello.c
```

The content is identical to [example1](../example1): one [C source file](content/src/hello.c) with a simple "hello world" program, compiled into a [single binary](content/build/hello) with no dependencies via a [Makefile](content/src/Makefile).

However, where example1 had a single SPDX document containing both source and binary, example2 instead has separate SPDX documents for [source](spdx/example2-src.spdx) and [binary](spdx/example2-bin.spdx).

This describes a scenario where binary and source are distributed separately, but where we want to be able to reflect the relationships between binary and source packages.

## Comments

Substantively, this is the same software as in [example1](../example).
However, here we are representing the sources and binaries as two separate Packages, on the assumption that we're distributing them separately.
Because of this, the source Package and binary Package are described in two separate SPDX documents.

Note that these do not _have_ to be in separate documents.
A single SPDX document can contain multiple Packages.
However, because the assumption in this scenario is that the binaries and the sources are distributed separately, it makes sense here that separate SPDX documents could accompany the binaries and the sources.

Relationships across separate documents are handled via `DocumentRef-` tags, defined via external document references in the Document Creation Info section.
Note that these external document references and relationships cannot be circular: one document can refer to the other, but (to my knowledge) they cannot refer circularly to each other.
To reference another document in an ExternalDocumentRef definition, you need to specify its hash, so it isn't possible for two documents to refer to one another; each would need to modify its own contents based on the other's hash value.

In the [SPDX document for the binary](spdx/example2-bin.spdx), note how the Relationships at the end of the document include `DocumentRef-hello-src:` as a prefix.
This uses the `DocumentRef-` defined in the `ExternalDocumentRef` tag at the top of the document.
This is the mechanism used to refer to SPDX identifiers for elements defined in other SPDX documents.
