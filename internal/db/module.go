package db

// Model
var ModuleModel = MakeModel(
	MakeID("f7c0fa36-3225-4996-bd10-116a257706d9"),
	"module",
	[]AttributeL{moduleName, moduleGoPackage},
	[]RelationshipL{},
	[]ConcreteInterfaceL{},
)

var moduleName = MakeConcreteAttribute(
	MakeID("37ff59a6-ecd1-49f6-a14b-2bc03a34e7b1"),
	"name",
	String,
)

var moduleGoPackage = MakeConcreteAttribute(
	MakeID("20144e78-c6e6-41fb-b548-d4964a6c0bbe"),
	"goPackage",
	String,
)

var ModuleInterfaces = MakeReverseRelationshipWithSource(
	MakeID("6409ef4a-6ceb-4ad8-891c-23b474b81ae9"),
	"interfaces",
	AbstractInterfaceModule,
	ModuleModel,
)

var ModuleFunctions = MakeReverseRelationshipWithSource(
	MakeID("26311d44-6fb1-407a-86df-a2710c538a33"),
	"functions",
	FunctionModule,
	ModuleModel,
)

var ModuleDatatypes = MakeReverseRelationshipWithSource(
	MakeID("28103c3d-4868-4aad-bf35-96ef9c2c3163"),
	"datatypes",
	DatatypeModule,
	ModuleModel,
)

func MakeModule(id ID, name, goPackage string, interfaces []InterfaceL,
	functions []FunctionL, datatypes []DatatypeL, modLits []ModuleLiteral) ModuleL {
	return ModuleL{
		id, name, goPackage, interfaces, functions, datatypes, modLits,
	}
}

type ModuleL struct {
	ID_        ID     `record:"id"`
	Name_      string `record:"name"`
	GoPackage  string `record:"goPackage"`
	interfaces []InterfaceL
	functions  []FunctionL
	datatypes  []DatatypeL
	modLits    []ModuleLiteral
}

func (lit ModuleL) MarshalDB(b *Builder) (recs []Record, links []Link) {
	rec := MarshalRecord(b, lit)
	recs = append(recs, rec)
	for _, iface := range lit.interfaces {
		links = append(links, Link{To: lit, From: iface, Rel: AbstractInterfaceModule})
	}
	for _, function := range lit.functions {
		links = append(links, Link{To: lit, From: function, Rel: FunctionModule})
	}
	for _, datatype := range lit.datatypes {
		links = append(links, Link{To: lit, From: datatype, Rel: DatatypeModule})
	}
	for _, modLit := range lit.modLits {
		links = append(links, Link{To: lit, From: modLit, Rel: modLit.ModuleRelationship()})
	}
	return
}

func (lit ModuleL) ID() ID {
	return lit.ID_
}

func (lit ModuleL) InterfaceID() ID {
	return ModuleModel.ID()
}

func (lit ModuleL) InterfaceName() string {
	return ModuleModel.Name_
}
