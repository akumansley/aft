package parsers

import (
	"fmt"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
)

func (p Parser) consumeCase(m db.Interface, keys api.Set, data map[string]interface{}) (operations.Case, error) {
	var c map[string]interface{}
	if v, ok := data["case"]; ok {
		c = v.(map[string]interface{})
		delete(keys, "case")
	}
	return p.parseCase(m, c)
}

func (p Parser) parseCase(iface db.Interface, data map[string]interface{}) (c operations.Case, err error) {
	var entries []operations.CaseEntry

	for k, val := range data {
		entryData, ok := val.(map[string]interface{})
		if !ok {
			return c, fmt.Errorf("expected an object\n")
		}
		var m db.Model
		m, err = p.Tx.Schema().GetModel(k)
		if err != nil {
			return
		}

		found := false
		var implements []db.Interface
		implements, err = m.Implements(p.Tx)
		if err != nil {
			return
		}
		for _, ifc := range implements {
			if ifc.ID() == iface.ID() {
				found = true
				break
			}
		}
		if !found {
			return c, fmt.Errorf("%v does not implement %v\n", m.Name(), iface.Name())
		}

		var entry operations.CaseEntry
		entry, err = p.parseCaseEntry(m, entryData)
		if err != nil {
			return
		}
		entries = append(entries, entry)
	}
	return operations.Case{Entries: entries}, nil
}

func (p Parser) parseCaseEntry(m db.Model, data map[string]interface{}) (c operations.CaseEntry, err error) {
	unusedKeys := make(api.Set)
	for k := range data {
		unusedKeys[k] = api.Void{}
	}

	i, s, err := p.consumeIncludeOrSelect(m, unusedKeys, data)

	// there are no (immediately) nested cases

	if len(unusedKeys) != 0 {
		return c, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return operations.CaseEntry{ModelID: m.ID(), Select: s, Include: i}, err
}
