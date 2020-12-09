---
title: Datastore
order: 7
---

Aft stores [records](#records) in the Aft Datastore. 

The datastore is split into two main components:

1. An in-memory "hold"
2. A durable log

Hold
---

The Aft hold is an in-memory immutable patricia trie with string keys and a few different index patterns: by ID, by [interface](#interfaces), and by relationships.

Log
---

The Aft log stores a binary encoding of every write made to the database, and provides durability. It is currently `fsync`'d after every transaction. Upon startup, the log is replayed into the Hold to restore the previous state of the datastore.

