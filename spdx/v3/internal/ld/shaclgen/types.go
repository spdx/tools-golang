package shaclgen

type Individual struct {
	IRI     string
	TypeIRI string
	Label   string
	Comment string
}

type Class struct {
	IRI        string
	Abstract   bool
	GoName     string
	Kind       string
	Comment    string
	ParentIRI  string
	Properties []*Property
}

type Property struct {
	IRI         string
	GoName      string
	Comment     string
	TypeIRI     string
	MinCount    int
	MaxCount    int
	Validations []any
}

type AllowedIRIValidation string

type MatchPatternValidation string
