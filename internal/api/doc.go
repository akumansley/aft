/*
Package api implements a JSON-based API for data in aft. The implementation
is split into three sub-packages: handlers, parsers, and operations.

Operations

The API supports the following operations

  * findOne
  * findMany
  * create
  * update
  * updateMany
  * upsert
  * delete
  * deleteMany
  * count

And several of those support nested operations as well. Update, for example,
supports connect, disconnect, set, create, update or delete.

Many operations also return a result set. For example, create returns the record
created, and findOne returns the record that was found. The API supports a
notion of inclusion and selection on a given result set.

For interface types, the API allows for differentiated selection and inclusion
clauses via the `case` statement.

Records in Aft always have two system attributes: `id` and `type`. `id` is a
UUID that (uniquely) identifies the object, and `type` is a string that
identifies the concrete model by name. These fields cannot be set by clients,
and will error if they are included in operations that modify data.


*/
package api
