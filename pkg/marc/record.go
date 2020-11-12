package marc

import (
	"fmt"
	"strings"
)

// Record is a struct representing a MARC record. It has a Fields slice
// which contains both ControlFields and DataFields.
type Record struct {
	Data   []byte
	Fields []Field
	Leader Leader
}

// Contains returns true if Record contains the value passed.
func (r Record) Contains(searchValue string) bool {
	if searchValue == "" {
		return true
	}
	for _, field := range r.Fields {
		if field.Contains(searchValue) {
			return true
		}
	}
	return false
}

// HasFields returns true if the Record contains the fields indicated
func (r Record) HasFields(filters FieldFilters) bool {
	exclude := FieldFilters{}
	return len(r.Filter(filters, exclude)) > 0
}

// ControlNum returns the control number (tag 001) for the record.
func (r Record) ControlNum() string {
	for _, f := range r.Fields {
		if f.Tag == "001" {
			return f.Value
		}
	}
	return ""
}

func (r Record) Raw() []byte {
	// Include the record terminator.
	return append(r.Data, rt)
}

func (r Record) String() string {
	return fmt.Sprintf("Leader: %s", r.Leader)
}

// Filter returns the fields in the record that match
// the given filter.
func (r Record) Filter(include FieldFilters, exclude FieldFilters) []Field {
	if len(include.Fields) == 0 && len(exclude.Fields) == 0 {
		// Nothing to filter
		return r.Fields
	}

	if len(include.Fields) > 0 {
		return r.filterInclude(include)
	}

	return r.filterExclude(exclude)
}

func (r Record) filterInclude(filters FieldFilters) []Field {
	list := []Field{}
	for _, filter := range filters.Fields {
		// Get all the fields in the record that match the tag
		// (there could be more than one)
		for _, field := range r.FieldsByTag(filter.Tag) {
			if len(filter.Subfields) == 0 {
				// add the value as-is, no need to filter by subfield
				list = append(list, field)
			} else {
				// extract the indicated subfields from the field
				// before adding it to the list
				filteredField := Field{
					Tag:        field.Tag,
					Value:      field.Value,
					Indicator1: field.Indicator1,
					Indicator2: field.Indicator2,
					SubFields:  field.GetSubFields(filter.Subfields),
				}
				if len(filteredField.SubFields) > 0 {
					list = append(list, filteredField)
				}
			}
		}
	}
	return list
}

func (r Record) filterExclude(filters FieldFilters) []Field {
	list := []Field{}
	for _, field := range r.Fields {
		include := true
		for _, filter := range filters.Fields {
			if len(filter.Subfields) > 0 {
				panic("Exclude filter by subfields not supported yet")
			}
			if filter.Tag == field.Tag {
				include = false
				break
			}
		}
		if include {
			list = append(list, field)
		}
	}
	return list
}

// FieldsByTag returns an array with the fields in the record for the given tag
func (r Record) FieldsByTag(tag string) []Field {
	var fields []Field
	for _, field := range r.Fields {
		if field.Tag == tag {
			fields = append(fields, field)
		}
	}
	return fields
}

// GetValue returns the first value for a field tag/subfield combination.
func (r Record) GetValue(tag string, subfield string) string {
	for _, field := range r.FieldsByTag(tag) {
		if field.IsControlField() {
			return field.Value
		}
		if subfield == "" {
			// No subfield indicated, return the string version of the field
			// TODO: Return the values rather than "=NNN \\ $aAAA $bBBB"
			return field.String()
		}
		for _, sub := range field.SubFields {
			if sub.Code == subfield {
				// Return the first instance of the requested subfield
				return sub.Value
			}
		}
	}
	return ""
}

// GetValues returns the values that match the field tag/subfield combination.
func (r Record) GetValues(tag string, subfield string) []string {
	values := []string{}
	for _, field := range r.FieldsByTag(tag) {
		if strings.TrimSpace(subfield) == "" {
			// No subfield indicated, return the string version of the field
			values = append(values, field.String())
		} else {
			for _, sub := range field.SubFields {
				if sub.Code == subfield {
					// Return the first instance of the requested subfield
					values = append(values, sub.Value)
				}
			}
		}
	}
	return values
}
