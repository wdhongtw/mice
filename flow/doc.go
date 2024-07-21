/*
Package flow is a helper library around "iter.Seq" types.

The library is intend to provide the lacking wheels from standard and/or
"x/exp/xiter" library.

For example, "Empty" and "Pack" is provided to build a sequence of
zero and one item, "Any" and "All" boolean short-circuit is also provided.
But "Map", "Filter" and "Reduce" is not provided since that is planned to be
in "x/exp/xiter". Also transformation from/to slice/map is in the standard
library "slices" and "maps".

All function Xxx comes with a Xxx2 version to address the usage between
"iter.Seq" and "iter.Seq2", if reasonable.

Function with immediate transformation, e.g. key, is not provided, since that
users can already achieve with another "Map" operation.

Wish someday we can use tuple as primitive generic type, so we don't have
to write these Xxx2 stuffs.
*/
package flow
