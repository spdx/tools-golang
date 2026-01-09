package shaclgen

import (
	"reflect"
	"slices"
	"strings"
	"time"

	. "github.com/dave/jennifer/jen"

	"github.com/spdx/tools-golang/spdx/v3/internal/ld"
)

func (g *generator) appendValidations(f *File) {
	for _, name := range keys(g.nameToIRI) {
		iri := g.nameToIRI[name]
		c := g.iriToType[iri]

		if g.flatStruct && c.Abstract {
			continue
		}

		validationFunc := Func().Params(Id("o").Op("*").Id(g.className(iri))).Id("Validate").Params().Id("error").Block(
			Return(Qual(ldImport, "JoinErrors").ParamsFunc(func(f *Group) {
				g.appendValidationFuncs(f, c)
			})),
		)
		f.Add(validationFunc)
	}
}

func (g *generator) appendValidationFuncs(f *Group, c *Class) {
	if c.ParentIRI != "" {
		if g.flatStruct {
			g.appendValidationFuncs(f, g.iriToType[c.ParentIRI])
		} else {
			f.Line().Qual(ldImport, "ValidateProperty").Params(Id("o"), Op("&").Id("o").Dot(g.className(c.ParentIRI)))
		}
	}
	for _, p := range c.Properties {
		fieldType := g.fieldType(c, p)

		validatePropParams := []Code{Id("o"), Op("&").Id("o").Dot(g.propName(c, p))}

		// some validations are unnecessary such as time.Time, since we store a time representation instead of arbitrary string
		skipTypeValidations := slices.Contains(validatingTypes, ld.TypeForIRI(p.TypeIRI))
		if g.isList(c, p) && p.MinCount > 0 {
			validatePropParams = append(validatePropParams, Line().Qual(ldImport, "ValidateMinCount").Index(fieldType).Params(Lit(p.MinCount)))
		}

		if g.isList(c, p) && p.MaxCount > 1 {
			validatePropParams = append(validatePropParams, Line().Qual(ldImport, "ValidateMaxCount").Index(fieldType).Params(Lit(p.MaxCount)))
		}

		var allowedIRIs []string
		for _, validation := range p.Validations {
			switch v := validation.(type) {
			case AllowedIRIValidation:
				allowedIRIs = append(allowedIRIs, string(v))
			case MatchPatternValidation:
				if skipTypeValidations {
					continue
				}
				expr := strings.ReplaceAll(string(v), "\\\\", "\\")
				validatePropParams = append(validatePropParams, Line().Qual(ldImport, "ValidateExpression").Params(Lit(expr)))
			}
		}

		if len(allowedIRIs) > 0 {
			var validateValuesParams []Code
			for _, allowedIRI := range allowedIRIs {
				validateValuesParams = append(validateValuesParams, Line().Id(g.namedIndividualName(&Individual{
					IRI:     cleanIRI(allowedIRI),
					TypeIRI: p.TypeIRI,
				})))
			}
			idCheck := Qual(ldImport, "ValidateIRI").Params(append(validateValuesParams, Line())...)
			if g.isList(c, p) {
				idCheck = Qual(ldImport, "ValidateAll").Params(idCheck)
			}
			validatePropParams = append(validatePropParams, Line().Add(idCheck))
		}

		// first 2 params are initialized as object, property --
		// only append a property validation if we added any validations or if the property is required
		if len(validatePropParams) > 2 || p.MinCount > 0 {
			f.Line().Qual(ldImport, "ValidateProperty").Params(validatePropParams...)
		}
	}
}

// these types do not need further pattern validation or otherwise implement their own Validate function
var validatingTypes = []reflect.Type{
	reflect.TypeOf(ld.URI("")),
	reflect.TypeOf(time.Time{}),
	reflect.TypeOf(ld.DateTime{}),
	reflect.TypeOf(ld.PositiveInt(0)),
	reflect.TypeOf(ld.NonNegativeInt(0)),
}
