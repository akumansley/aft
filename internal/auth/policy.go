package auth

import (
	"fmt"

	"awans.org/aft/internal/db"
)

type PolicyType int

const (
	Create PolicyType = iota
	Update
	Read
)

func (pt PolicyType) String() string {
	switch pt {
	case Create:
		return "create"
	case Update:
		return "update"
	case Read:
		return "read"
	default:
		panic("unknown op")
	}
}

var PolicyModel = db.MakeModel(
	db.MakeID("ea5eda03-6780-4a31-8b9b-e5f16a98d8b3"),
	"policy",
	[]db.AttributeL{
		pAllowRead, pAllowCreate, pAllowUpdate,
		pReadWhere, pCreateWhere, pUpdateWhere,
	},
	[]db.RelationshipL{
		PolicyFor,
	},
	[]db.ConcreteInterfaceL{},
)

var pAllowRead = db.MakeConcreteAttribute(
	db.MakeID("5783649a-eb1f-4a96-ba00-219c5137641c"),
	"allowRead",
	db.Bool,
)

var pReadWhere = db.MakeConcreteAttribute(
	db.MakeID("55cfda72-c7f2-47aa-85ab-e54b98f1eda0"),
	"readWhere",
	db.String,
)

var pAllowCreate = db.MakeConcreteAttribute(
	db.MakeID("43592002-b914-4aed-9d7a-4292ba2b3467"),
	"allowCreate",
	db.Bool,
)

var pCreateWhere = db.MakeConcreteAttribute(
	db.MakeID("d4413ff7-391b-4dce-9fef-002c8cbc9441"),
	"createWhere",
	db.String,
)

var pAllowUpdate = db.MakeConcreteAttribute(
	db.MakeID("fe77d33a-7691-438d-8d61-c79d5fed2454"),
	"allowUpdate",
	db.Bool,
)

var pUpdateWhere = db.MakeConcreteAttribute(
	db.MakeID("c07a6822-9487-43a8-9b00-d3d87ff473d7"),
	"updateWhere",
	db.String,
)

var PolicyFor = db.MakeConcreteRelationship(
	db.MakeID("be24d5ca-48f4-4d6f-a550-5b969703f440"),
	"interface",
	false,
	db.InterfaceInterface,
)

var ConcreteInterfacePolicies = db.MakeReverseRelationshipWithSource(
	db.MakeID("79bb494e-0333-4778-9239-6ef3d5a9f720"),
	"policies",
	PolicyFor,
	db.InterfaceModel,
)

var ModelPolicies = db.MakeReverseRelationshipWithSource(
	db.MakeID("ebd7e315-78b4-4179-9f98-149ac7a257a0"),
	"policies",
	PolicyFor,
	db.ModelModel,
)

var InterfacePolicies = db.MakeInterfaceRelationshipWithSource(
	db.MakeID("a70fd021-8a7c-4cc1-ba56-b99a1ccc5393"),
	"policies",
	true,
	db.InterfaceInterface,
	PolicyModel,
)

var PolicyRole = db.MakeReverseRelationshipWithSource(
	db.MakeID("e7bb2583-ce26-4369-86dc-9a8f6952ad2e"),
	"role",
	RolePolicy,
	PolicyModel,
)

type PolicyL struct {
	ID_         db.ID  `record:"id"`
	AllowRead   bool   `record:"allowRead"`
	ReadWhere   string `record:"readWhere"`
	AllowCreate bool   `record:"allowCreate"`
	CreateWhere string `record:"createWhere"`
	AllowUpdate bool   `record:"allowUpdate"`
	UpdateWhere string `record:"updateWhere"`
	For_        db.ModelL
}

func (lit PolicyL) ID() db.ID {
	return lit.ID_
}

func (lit PolicyL) InterfaceID() db.ID {
	return PolicyModel.ID()
}

func (lit PolicyL) InterfaceName() string {
	return PolicyModel.Name_
}

func (lit PolicyL) MarshalDB(b *db.Builder) (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(b, lit)
	recs = append(recs, rec)
	links = append(links, db.Link{From: lit, To: lit.For_, Rel: PolicyFor})
	return
}

type policy struct {
	rec db.Record
	tx  db.Tx
}

func (p *policy) ReadWhere() string {
	return p.rec.MustGet("readWhere").(string)
}

func (p *policy) UpdateWhere() string {
	return p.rec.MustGet("updateWhere").(string)
}

func (p *policy) CreateWhere() string {
	return p.rec.MustGet("createWhere").(string)
}

func (p *policy) String() string {
	return fmt.Sprintf("policy{for: %v ..}", p.Interface().Name())
}

func (p *policy) Interface() db.Interface {
	tx := p.tx
	policies := tx.Ref(PolicyModel.ID())
	ifaces := tx.Ref(db.InterfaceInterface.ID())

	ifrec, err := tx.Query(ifaces,
		db.Join(policies, ifaces.Rel(InterfacePolicies.Load(tx))),
		db.Filter(policies, db.EqID(p.rec.ID())),
	).OneRecord()
	if err != nil {
		panic("No model")
	}
	// this is awkward and inefficient
	i, err := tx.Schema().GetInterfaceByID(ifrec.ID())
	if err != nil {
		panic("No model")
	}
	return i
}

func subJSON(data interface{}, subs map[string]interface{}) {
	switch data.(type) {
	case map[string]interface{}:
		subJSONObject(data.(map[string]interface{}), subs)
	case []interface{}:
		subJSONArray(data.([]interface{}), subs)
	default:
		return
	}
}
func subJSONArray(data []interface{}, subs map[string]interface{}) {
	for _, v := range data {
		subJSON(v, subs)
	}
}

func subJSONObject(data map[string]interface{}, subs map[string]interface{}) {
	for k, v := range data {
		if sv, ok := v.(string); ok {
			subVal, inSub := subs[sv]
			if inSub {
				data[k] = subVal
			}
		} else {
			subJSON(v, subs)
		}
	}
}
