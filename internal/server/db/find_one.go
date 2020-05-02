package db

func (op FindOneOperation) Apply(db DB) (st interface{}, err error) {
	return db.FindOne(op.ModelName, op.UniqueQuery)
}
