package db

func (op FindOneOperation) Apply(db DB) interface{} {
	val := findOne(db, op.ModelName, op.UniqueQuery)
	return val
}
