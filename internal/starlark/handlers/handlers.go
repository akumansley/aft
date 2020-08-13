package handlers

import (
	"awans.org/aft/internal/api/parsers"
	"go.starlark.net/starlark"
)

func (h Handler) findOne(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	modelName, body, err := unpack(thread, b, args, kwargs)
	if err != nil {
		return
	}

	p := parsers.Parser{Tx: h.tx}
	op, err := p.ParseFindOne(modelName, body)
	if err != nil {
		return
	}

	result, err := op.Apply(h.tx)
	if err != nil {
		return
	}
	return output(result)
}

func (h Handler) findMany(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	modelName, body, err := unpack(thread, b, args, kwargs)
	if err != nil {
		return
	}

	p := parsers.Parser{Tx: h.tx}
	op, err := p.ParseFindMany(modelName, body)
	if err != nil {
		return
	}

	result, err := op.Apply(h.tx)
	if err != nil {
		return
	}
	return output(result)
}

func (h Handler) count(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	modelName, body, err := unpack(thread, b, args, kwargs)
	if err != nil {
		return
	}

	p := parsers.Parser{Tx: h.tx}
	op, err := p.ParseCount(modelName, body)
	if err != nil {
		return
	}

	result, err := op.Apply(h.tx)
	if err != nil {
		return
	}
	return output(result)
}

func (h Handler) del(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	modelName, body, err := unpack(thread, b, args, kwargs)
	if err != nil {
		return
	}

	p := parsers.Parser{Tx: h.tx}
	op, err := p.ParseDelete(modelName, body)
	if err != nil {
		return
	}

	result, err := op.Apply(h.tx)
	if err != nil {
		return
	}
	return output(result)
}

func (h Handler) deleteMany(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	modelName, body, err := unpack(thread, b, args, kwargs)
	if err != nil {
		return
	}

	p := parsers.Parser{Tx: h.tx}
	op, err := p.ParseDeleteMany(modelName, body)
	if err != nil {
		return
	}

	result, err := op.Apply(h.tx)
	if err != nil {
		return
	}
	return output(result)
}

func (h Handler) update(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	modelName, body, err := unpack(thread, b, args, kwargs)
	if err != nil {
		return
	}

	p := parsers.Parser{Tx: h.tx}
	op, err := p.ParseUpdate(modelName, body)
	if err != nil {
		return
	}

	result, err := op.Apply(h.tx)
	if err != nil {
		return
	}
	return output(result)
}

func (h Handler) updateMany(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	modelName, body, err := unpack(thread, b, args, kwargs)
	if err != nil {
		return
	}

	p := parsers.Parser{Tx: h.tx}
	op, err := p.ParseUpdateMany(modelName, body)
	if err != nil {
		return
	}

	result, err := op.Apply(h.tx)
	if err != nil {
		return
	}
	return output(result)
}

func (h Handler) create(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	modelName, body, err := unpack(thread, b, args, kwargs)
	if err != nil {
		return
	}

	p := parsers.Parser{Tx: h.tx}
	op, err := p.ParseCreate(modelName, body)
	if err != nil {
		return
	}

	result, err := op.Apply(h.tx)
	if err != nil {
		return
	}
	return output(result)
}

func (h Handler) upsert(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	modelName, body, err := unpack(thread, b, args, kwargs)
	if err != nil {
		return
	}

	p := parsers.Parser{Tx: h.tx}
	op, err := p.ParseUpsert(modelName, body)
	if err != nil {
		return
	}

	result, err := op.Apply(h.tx)
	if err != nil {
		return
	}
	return output(result)
}
