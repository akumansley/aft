package db

//Enums
type Enum struct {
	ID   ID
	Name string
}

func (d Enum) FromJSON(arg interface{}) (interface{}, error) {
	//TODO if it's an enum, verify that it's a UUID.
	// I see two options. One we execute some code here that intelligently
	// checks if this enum id matches the enum values from the table.
	// Second option is that I plum through tx into FromJSON. Second seems worse.
	return UUIDFromJSON(arg)
}
