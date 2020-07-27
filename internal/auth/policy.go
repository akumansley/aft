package auth

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/db"
	"encoding/json"
	"fmt"
)

var PolicyModel = db.MakeModel(
	db.MakeID("ea5eda03-6780-4a31-8b9b-e5f16a98d8b3"),
	"policy",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("7ebfbce0-3280-4067-8cce-c00efa89bb43"),
			"name",
			db.String,
		),
		pText,
	},
	// set in init
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{},
)

var pText = db.MakeConcreteAttribute(
	db.MakeID("55cfda72-c7f2-47aa-85ab-e54b98f1eda0"),
	"text",
	db.String,
)

var PolicyFor = db.MakeConcreteRelationship(
	db.MakeID("be24d5ca-48f4-4d6f-a550-5b969703f440"),
	"model",
	false,
	db.ModelModel,
)

var ModelPolicies = db.MakeReverseRelationship(
	db.MakeID("09579552-6982-4732-9d69-585f2e6a74b1"),
	"policies",
	PolicyFor,
)

var PolicyRoles = db.MakeReverseRelationship(
	db.MakeID("e7bb2583-ce26-4369-86dc-9a8f6952ad2e"),
	"roles",
	RolePolicies,
)

type PolicyL struct {
	ID_   db.ID  `record:"id"`
	Name_ string `record:"name"`
	Text_ string `record:"text"`
	For_  db.ModelL
}

func (lit PolicyL) ID() db.ID {
	return lit.ID_
}

func (lit PolicyL) MarshalDB() (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(lit, PolicyModel)
	recs = append(recs, rec)
	links = append(links, db.Link{rec.ID(), lit.For_.ID(), PolicyFor})
	return
}

type policy struct {
	rec db.Record
	tx  db.Tx
}

func (p *policy) String() string {
	return fmt.Sprintf("policy{\"%v\"}", p.Text())
}

func (p *policy) Text() string {
	return pText.MustGet(p.rec).(string)
}

func (p *policy) Model() db.Model {
	tx := p.tx
	policies := tx.Ref(PolicyModel.ID())
	models := tx.Ref(db.ModelModel.ID())
	mrec, err := tx.Query(models, db.Join(policies, models.Rel(ModelPolicies)), db.Filter(policies, db.EqID(p.rec.ID()))).OneRecord()
	if err != nil {
		panic("No model")
	}
	// this is awkward and inefficient
	m, err := tx.Schema().GetModelByID(mrec.ID())
	if err != nil {
		panic("No model")
	}
	return m
}

func (p *policy) Apply(tx db.Tx, ref db.ModelRef) []db.QueryClause {
	iface, err := tx.Schema().GetInterfaceByID(ref.InterfaceID)
	if err != nil {
		panic("bad")
	}
	var data map[string]interface{}

	json.Unmarshal([]byte(p.Text()), &data)
	w, err := parsers.Parser{tx}.ParseWhere(iface.Name(), data)
	if err != nil {
		panic("bad")
	}
	clauses := operations.HandleWhere(tx, ref, w)
	return clauses
}
