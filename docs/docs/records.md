---
id: records
title: Records
---

![Screenshot of the record viewer](/img/records.png)

A record is roughly analogous to a row in a relational database. Logically, they're a tuple of named fields.

In JSON, they're an Object with string keys and non-Array or Object values. 

In memory, they're stored as structs manipulated using reflection. On disk, they're Gob encoded structs.

Unlike a relational database, Aft treats relationships between records specially; they're not stored on the record itself, but in a separate "links" index. This is, however, opaque to the end-user.
