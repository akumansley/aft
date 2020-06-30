package db

type attr struct {
	rec Record
	tx  Tx
}

func (a *attr) ID() ID {
	return a.rec.ID()
}

func (a *attr) Name() string {
	return a.rec.MustGet("name").(string)
}
