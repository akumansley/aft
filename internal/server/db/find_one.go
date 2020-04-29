package db

func (op FindOneOperation) Apply(db DB) (st interface{}, err error) {
	return findOne(db, op.ModelName, op.UniqueQuery)
}
