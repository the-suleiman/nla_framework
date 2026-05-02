package types

type (
	OdataFld struct {
		Name         string
		Type         string
		CastToGoType string // optional custom cast from odata json to pg field
	}
)
