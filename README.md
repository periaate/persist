# PartDB

PartDB is a database built around a hash map implementation made with Go generics, allowing for type safety, native performance, but couples it with persistence.
PartDB is being specifically developed to tackle partite datasets, primarily bipartite datasets. This means that the focus is on memory contiguity and memory latency as the most common operations are set membership operations.

PartDB is closer to a data structure persistence system than a typical database, as the focus here isn't on abstraction, but on control.
