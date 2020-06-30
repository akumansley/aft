package db

// // remove
// GetModel(string) (Model, error)
// GetRelationships(Model) ([]Relationship, error)
// GetRelationship(ID) (Relationship, error)
// GetModelByID(ModelID) (Model, error)
// SaveModel(Model) error
// SaveRelationship(Relationship) error

type Schema struct {
	tx *holdTx
}

func (s *Schema) GetModel(mid ModelID) Model {
	mrec, err := s.tx.FindOne(ModelModel.ID(), ID(mid))
	if err != nil {
		panic("GetModel failed")
	}
	return &model{mrec, tx}
}
