package db

var FuncType = MakeEnum(
	MakeID("3c08b33f-526b-4408-926b-ad9dfc87bf04"),
	"funcType",
	[]EnumValueL{
		RPC,
		Validator,
		Internal,
	})

var RPC = MakeEnumValue(
	MakeID("4b8db42e-d084-4328-a758-a76939341ffa"),
	"rpc",
)

var Validator = MakeEnumValue(
	MakeID("75e1f16d-49f1-4f29-b69b-0b6d7eb1d2f8"),
	"validator",
)

var Internal = MakeEnumValue(
	MakeID("0751f10f-b1b1-4c32-a22e-4981be6d199d"),
	"internal",
)
