package db

import (
	"fmt"
)

type migrateStep interface {
	Migrate(*Builder, RWTx) error
}

func MakeAutomigrateHandler(b *Builder) func(BeforeCommit) {
	handler := func(event BeforeCommit) {
		tx := event.Tx
		ops := tx.Operations()
		if len(ops) == 0 {
			return
		}
		rwtx := tx.(RWTx)
		var steps []migrateStep
		for _, op := range ops {
			step := stepForOp(op)
			if step != nil {
				steps = append(steps, step)
			}
		}
		for _, step := range steps {
			err := step.Migrate(b, rwtx)
			if err != nil {
				rwtx.Abort(err)
				return
			}
		}
	}
	return handler
}
func stepForOp(op Operation) migrateStep {
	switch op.(type) {
	case *CreateOp:
		return nil
	case *UpdateOp:
		update := op.(*UpdateOp)
		switch update.NewRecord.InterfaceID() {
		case ConcreteRelationshipModel.ID():
			// don't need to handle renames b/c they're not
			// stored on the records themselves
			if changed(update, "multi") {
				return emptyRelationship{update.NewRecord.ID()}
			}
		case ConcreteAttributeModel.ID():
			if changed(update, "name") {
				return renameAttribute{update}
			}
		case CoreDatatypeModel.ID():
			err := fmt.Errorf("unsupported migration %v", update)
			panic(err)
		}
	case *DeleteOp:
		delete := op.(*DeleteOp)
		switch delete.Record.InterfaceID() {
		case ModelModel.ID():
			return dropModel{delete.Record.ID()}
		case ConcreteRelationshipModel.ID():
			return dropRelationship{delete.Record.ID()}
		case ConcreteAttributeModel.ID():
			return dropAttribute{delete.Record.ID()}
		case CoreDatatypeModel.ID():
			err := fmt.Errorf("unsupported migration %v", delete)
			panic(err)
		}
	case *ConnectOp:
		connect := op.(*ConnectOp)
		switch connect.RelID {
		case ModelAttributes.ID():
			return addAttribute{connect.Target}
		case ModelRelationships.ID():
			return addRelationship{connect}
		case ConcreteAttributeDatatype.ID():
			return retypeAttribute{connect.Source}
		}
	case *DisconnectOp:
		disconnect := op.(*DisconnectOp)
		switch disconnect.RelID {
		case ModelRelationships.ID():
			return dropRelationship{disconnect.Target}
		case ModelAttributes.ID():
			return dropAttribute{disconnect.Target}
		case ConcreteAttributeDatatype.ID():
			// this leaves you in a weird/invalid state
			// but should be fixed up by another op
		}
	default:
		return nil
	}
	return nil
}

func changed(op *UpdateOp, fieldName string) bool {
	return op.OldRecord.MustGet(fieldName) != op.NewRecord.MustGet(fieldName)
}

type dropModel struct {
	modelID ID
}

func (d dropModel) String() string {
	return fmt.Sprintf("dropModel{%v}", d.modelID)
}

// need to cascade to rels pointing to this model
func (d dropModel) Migrate(b *Builder, tx RWTx) error {
	model, err := tx.AsOfStart().Schema().GetModelByID(d.modelID)
	if err != nil {
		return err
	}

	// clean up records
	mref := tx.Ref(d.modelID)
	recs := tx.Query(mref).Records()
	for _, r := range recs {
		tx.unloggedDelete(r)
	}

	// clean up implements records
	ifaces, err := model.Implements()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		tx.dropImplements(model.ID(), iface.ID())
	}

	// clean up rel records
	rels, err := model.Relationships()
	if err != nil {
		return err
	}
	for _, rel := range rels {
		tx.dropRel(model.ID(), rel.Target().ID(), rel.ID())
	}

	// TODO: clean up relationships pointing at this model
	panic("not implemented")

	// b.InterfaceUpdated(model)
	return nil
}

func dropLinks(tx RWTx, rel Relationship) {
	source := tx.Ref(rel.Source().ID())
	target := tx.Ref(rel.Target().ID())
	var q Q
	if rel.Multi() {
		q = tx.Query(source,
			Join(target, source.Rel(rel)),
			Aggregate(source, Some),
		)
	} else {
		q = tx.Query(source,
			Join(target, source.Rel(rel)),
		)
	}
	qrs := q.All()
	for _, qr := range qrs {
		if rel.Multi() {
			childQRs := qr.GetChildRelMany(rel)
			for _, childQR := range childQRs {
				tx.Disconnect(qr.Record.ID(), childQR.Record.ID(), rel.ID())
			}
		} else {
			childQR := qr.GetChildRelOne(rel)
			tx.Disconnect(qr.Record.ID(), childQR.Record.ID(), rel.ID())
		}
	}
}

type addRelationship struct {
	connect *ConnectOp
}

func (a addRelationship) String() string {
	return fmt.Sprintf("addRelationship{%v}", a.connect.Target)
}

func (a addRelationship) Migrate(b *Builder, tx RWTx) error {
	rel, err := tx.Schema().GetRelationshipByID(a.connect.Target)
	if err != nil {
		return err
	}

	return tx.addRel(rel.ID(), rel.Source().ID(), rel.Target().ID())
}

type dropRelationship struct {
	relID ID
}

func (d dropRelationship) String() string {
	return fmt.Sprintf("dropRelationship{%v}", d.relID)
}

func (d dropRelationship) Migrate(b *Builder, tx RWTx) error {
	rel, err := tx.AsOfStart().Schema().GetRelationshipByID(d.relID)
	if err != nil {
		return err
	}
	referencing, _ := tx.Schema().GetRelationshipByID(ReverseRelationshipReferencing.ID())

	revrels := tx.Ref(ReverseRelationshipModel.ID())
	rels := tx.Ref(ConcreteRelationshipModel.ID())
	q := tx.Query(revrels,
		Join(rels, revrels.Rel(referencing)),
		Filter(rels, EqID(rel.ID())),
	)
	revrelRecs := q.Records()

	for _, reverseRel := range revrelRecs {
		tx.unloggedDelete(reverseRel)
	}

	return tx.dropRel(rel.ID(), rel.Source().ID(), rel.Target().ID())
}

type emptyRelationship struct {
	relID ID
}

func (d emptyRelationship) String() string {
	return fmt.Sprintf("emptyRelationship{%v}", d.relID)
}

func (e emptyRelationship) Migrate(b *Builder, tx RWTx) error {
	rel, err := tx.Schema().GetRelationshipByID(e.relID)
	if err != nil {
		return err
	}
	dropLinks(tx, rel)
	return nil
}

func mapModel(tx RWTx, modelID ID, f func(Record) Record) {
	mref := tx.Ref(modelID)
	recs := tx.Query(mref).Records()
	for _, r := range recs {
		newR := f(r)
		tx.unloggedUpdate(r, newR)
	}
}

type renameAttribute struct {
	op *UpdateOp
}

func (r renameAttribute) String() string {
	return fmt.Sprintf("renameAttribute{%v}", r.op.NewRecord.ID())
}

func modelForAttr(tx Tx, attrID ID) Model {
	models := tx.Ref(ModelModel.ID())
	attrs := tx.Ref(ConcreteAttributeModel.ID())

	modelAttributes, _ := tx.Schema().GetRelationshipByID(ModelAttributes.ID())
	q := tx.Query(models,
		Join(attrs, models.Rel(modelAttributes)),
		Aggregate(attrs, Some),
		Filter(attrs, EqID(attrID)))
	rec, err := q.OneRecord()
	if err != nil {
		panic(err)
	}
	model := tx.Schema().LoadModel(rec)
	return model
}

func (r renameAttribute) Migrate(b *Builder, tx RWTx) error {
	model := modelForAttr(tx, r.op.NewRecord.ID())
	attrs, err := model.Attributes()
	if err != nil {
		return err
	}

	b.InterfaceUpdated(model)
	oldName := r.op.OldRecord.MustGet("name").(string)
	newName := r.op.NewRecord.MustGet("name").(string)

	rename := func(old Record) Record {
		newRec, err := b.RecordForInterface(model)
		if err != nil {
			panic(err)
		}
		for _, a := range attrs {
			if a.Name() == "type" {
				continue
			}
			if a.Name() == newName {
				oldVal := old.MustGet(oldName)
				newRec.Set(newName, oldVal)
				continue
			}
			newRec.Set(a.Name(), old.MustGet(a.Name()))
		}

		return newRec
	}
	mapModel(tx, model.ID(), rename)
	return nil
}

type addAttribute struct {
	attrID ID
}

func (a addAttribute) Migrate(b *Builder, tx RWTx) error {
	model := modelForAttr(tx, a.attrID)
	b.InterfaceUpdated(model)

	add := func(old Record) Record {
		newRec, err := b.RecordForInterface(model)
		if err != nil {
			panic(err)
		}

		attrs, err := model.Attributes()
		if err != nil {
			panic(err)
		}
		for _, attr := range attrs {
			if attr.ID() == a.attrID || attr.Name() == "type" {
				continue
			}
			newRec.Set(attr.Name(), old.MustGet(attr.Name()))
		}
		return newRec
	}
	mapModel(tx, model.ID(), add)
	return nil
}

func (a addAttribute) String() string {
	return fmt.Sprintf("addAttribute{%v}", a.attrID)
}

type dropAttribute struct {
	attrID ID
}

func (a dropAttribute) String() string {
	return fmt.Sprintf("dropAttribute{%v}", a.attrID)
}

func (d dropAttribute) Migrate(b *Builder, tx RWTx) error {
	attr, err := tx.AsOfStart().Schema().GetAttributeByID(d.attrID)
	if err != nil {
		return err
	}
	model := modelForAttr(tx, d.attrID)
	b.InterfaceUpdated(model)

	drop := func(old Record) Record {
		newRec, err := b.RecordForInterface(model)
		if err != nil {
			panic(err)
		}
		attrs, err := model.Attributes()
		if err != nil {
			panic(err)
		}
		for _, a := range attrs {
			if a.ID() == attr.ID() {
				continue
			}
			newRec.Set(attr.Name(), old.MustGet(attr.Name()))
		}
		return newRec
	}
	mapModel(tx, model.ID(), drop)
	return nil
}

type retypeAttribute struct {
	attrID ID
}

func (a retypeAttribute) String() string {
	return fmt.Sprintf("retypeAttribute{%v}", a.attrID)
}

func (r retypeAttribute) Migrate(b *Builder, tx RWTx) error {
	attr, err := tx.Schema().GetAttributeByID(r.attrID)
	if err != nil {
		return err
	}
	model := modelForAttr(tx, r.attrID)
	b.InterfaceUpdated(model)

	retype := func(old Record) Record {
		newRec, err := b.RecordForInterface(model)
		if err != nil {
			panic(err)
		}
		attrs, err := model.Attributes()
		if err != nil {
			panic(err)
		}
		for _, a := range attrs {
			if a.ID() == attr.ID() || a.Name() == "type" {
				continue
			}
			newRec.Set(a.Name(), old.MustGet(a.Name()))
		}
		return newRec
	}
	mapModel(tx, model.ID(), retype)
	return nil
}

type addInterface struct {
	op *ConnectOp
}

func (a addInterface) String() string {
	return fmt.Sprintf("addInterface{%v}", a.op.Target)
}

func (a addInterface) Migrate(b *Builder, tx RWTx) error {
	model, err := tx.Schema().GetModelByID(a.op.Source)
	if err != nil {
		return err
	}
	b.InterfaceUpdated(model)
	tx.addImplements(a.op.Source, a.op.Target)
	return nil
}

type removeInterface struct {
	op DisconnectOp
}

func (r removeInterface) String() string {
	return fmt.Sprintf("removeInterface{%v}", r.op.Target)
}

func (r removeInterface) Migrate(b *Builder, tx RWTx) error {
	model, err := tx.Schema().GetModelByID(r.op.Source)
	if err != nil {
		return err
	}
	b.InterfaceUpdated(model)
	tx.dropImplements(model.ID(), r.op.Target)
	return nil
}
