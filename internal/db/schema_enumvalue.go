package db

type enumValue struct {
	rec Record
}

func (ev *enumValue) ID() ID {
	return ev.rec.ID()
}

func (ev *enumValue) Name() string {
	return ev.rec.MustGet("name").(string)
}
