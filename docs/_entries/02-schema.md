---
title: Schema
order: 2
---

When you're developing an app with Aft you start by defining a schema.

![Screenshot of the schema page](/aft/img/schema.png)

Models
---

Models are roughly analogous to a table in a relational database. A model has a set of attributes, some relationships to other models, and it may implement some interfaces.

Attributes
---

If models are roughly tables, the attributes are roughly columns. An attribute has a name and a datatype. 

Datatypes
---

### Core Datatypes

Aft supports the typical core datatypes like string, int, or bytes. 

### Custom Datatypes

It's possible to create datatypes that additionally restrict a core datatype:

* Email address - stored as a string
* Phone number - stored as a string

A custom datatype has a "validator"—a function that parses incoming values from JSON and either accepts them, translated to a storage format, or rejects them as invalid.


Relationships
---

Relationships in Aft go from a source interface to a target interface. They can have a cardinality of either single or multiple.

You may also define "reverse" relationships to enable backreferences from target to source.

Interfaces
---

An interface acts as a union of one or more models. Interfaces can be queried using the columns on the interface, or can be queried using "case" statements to filter on properties unique to one model that implements the interface.

Models need to explicitly implement interfaces.

Functions
---

In addition to storing data, Aft is able to store functions. Currently functions can either be "native" Go functions, or scripts written in [Starlark](https://github.com/google/starlark-go).