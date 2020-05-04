package db

func (op FindOneOperation) Apply(tx Tx) (st interface{}, err error) {
	return tx.FindOne(op.ModelName, op.UniqueQuery)
}
