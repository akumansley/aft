package db

type InterfaceL struct {
	ID         ID     `record:"id"`
	Name       string `record:"name"`
	Attributes []InterfaceAttributeL
}

func (lit InterfaceL) GetID() ID {
	return lit.ID
}

func (lit InterfaceL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, InterfaceModel)
	recs = append(recs, rec)
	for _, a := range lit.Attributes {
		ars, al := a.MarshalDB()
		recs = append(recs, ars...)
		links = append(links, al...)

		links = append(links, Link{rec.ID(), a.GetID(), ModelAttributes})
	}
	return
}
