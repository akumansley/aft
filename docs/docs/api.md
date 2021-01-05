---
id: api
title: API
---


Aft's API supports reading and mutating data in the datastore using a rich JSON API.

Every call uses an HTTP POST with a JSON body (and 'application/json' Content-Type).

The routes are exposed in the following URL format:

```
https://$BASE_URL/api/$INTERFACE.$OPERATION
```

So for example, to perform a `findMany` on the `model` model, one would send an HTTP POST to

```
https://$BASE_URL/api/model.findMany
```


![Screenshot of making an API call](/img/api.png)


## findOne/findMany

`findOne` and `findMany` both return records from the Aft datastore. Both return an object with a single key, "data".


`findOne` returns as "data" a single record as a JSON object, whereas `findMany` returns an array.

Both operations accept a JSON object as a parameter describing the query. At the top level, the request must have a key, "where", which must be an object.

The where object may contain keys referencing the fields (i.e. attributes or relationships) of the model being queried. The values of those keys represent equality filters. For example, to query for a user named Andrew, one might send:

```
{
	"where": {
		"name": "Andrew"
	}
}
```

Where clauses may be applied to relationships as well, in which case Aft accepts a nested "where" object. For example, if user object has a related "profile" object with a "bio" field, one might query:

```
{
	"where": {
		"profile": {
			"bio": "This is my cool bio!"
		}
	}
}
```

For relationships with Mutli cardinality, an aggregate must be specified. For example, if a user has a related collection of posts:

```
{
	"where": {
		"posts": {
			"some": {"text": "some cool post"}
		}
	}
}
```

The supported aggregations are "some", "every" and "none".

### Includes

A typical client UI requires data not just from a single record, but also a graph of related records.

In order to support efficient and transactional loading of these record graphs, Aft supports "includes".

Includes are peers to the "where" object at the top level of a `findMany` or `findOne` request.


```
{
	"where": {
		"name": "Andrew"
	},
	"include": {
		"posts": true
	}
}
```

Includes may be similarly nested.

```
{
	"where": {
		"name": "Andrew"
	},
	"include": {
		"posts": {
			"include": {
				"comments": true
			}
		}
	}
}
```

### Selection

For efficiency, clients may request only a subset of fields on a record by using a "select" object. Selects may also specify records being included, though only one or the other can exist on the same level of nesting.


```
{
	"where": {
		"name": "Andrew"
	},
	"select": {
		"id": true,
		"name": true,
		"posts": true
	}
}
```

## count

Count accepts the same arguments as a `findMany` call, but returns an object with a single key, "count", that has as its value a number representing the number of records that would've been returned from the equivalent findMany.


## create, update, updateMany, delete, upsert

The mutation methods are similar to the read methods.

Create accepts a single key "data". All records in Aft have an implcit id field which is always a server-generated UUID4. Thus it is omitted from the create call.

Nested related objects may be mutated—created, updated, connected, disconnecting, set—in the same call.

So to create a user:

```
{
	"data": {
		"name": "Andrew",
		"posts": {
			"create": [{
				"text": "a new post!"
				}],
		}
	}
}
```

Updates take two arguments: a "where" object to locate the record to be updated, and "data" describing the mutation.

Updates accept the same set of nested mutations.


```
{
	"where": {
		"name": "Andrew",
	}
	"data": {
		"age": 32,
		"posts": {
			"create": [{
				"text": "a new post!"
				}],
		}
	}
}
```

Deletes just accept a where.


```
{
	"where": {
		"name": "Andrew",
	}
}
```

And upserts accept a where object and two possibilities. A create, taken in the event that no record maches the where object and an update in the opposite case.


```
{
	"where": {
		"name": "Andrew",
	},
	"create": {
		"name": "Andrew",
		"age": 32,
	},
	"update": {
		"age": 32,
	}
}
```

