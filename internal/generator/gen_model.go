package generator

import (
	"fmt"
	"strings"

	"github.com/eslamward/backforge/internal/parser"
)

func GenerateModel(model parser.Model) string {
	var sb strings.Builder
	importTime := false

	structName := toPascalCase(toSingular(model.Name))
	reqName := "Requested" + structName

	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	for _, field := range model.Fields {
		goType := mapType(field.Type)

		if goType == "time.Time" {
			importTime = true
		}

		fieldName := toPascalCase(field.Name)

		if !field.NotNull && !field.Primary && !field.Unique {
			goType = "*" + goType
		}

		tag := fmt.Sprintf("`json:\"%s\" db:\"%s\"`", field.Name, field.Name)

		line := fmt.Sprintf("\t%s %s %s\n", fieldName, goType, tag)
		sb.WriteString(line)
	}

	sb.WriteString("}\n\n")

	sb.WriteString(fmt.Sprintf("type %s struct {\n\n", reqName))
	for _, field := range model.Fields {
		goType := mapType(field.Type)

		if goType == "time.Time" {
			importTime = true
		}
		if field.Primary {
		} else {
			goType = "*" + goType

		}

		fieldName := toPascalCase(field.Name)

		tag := fmt.Sprintf("`json:\"%s\" db:\"%s\"`", field.Name, field.Name)

		line := fmt.Sprintf("\t%s %s %s\n", fieldName, goType, tag)
		sb.WriteString(line)
	}

	sb.WriteString("}\n\n")
	sb.WriteString(
		fmt.Sprintf("type Requested%s%s struct  {\n%s %s `json:\"%s\"`\n}\n\n",
			structName, toPascalCase(primaryField(model).Name),
			toPascalCase(primaryField(model).Name),
			mapType(primaryField(model).Type),
			primaryField(model).Name))

	if uniqeField(model) != nil {
		uLower := uniqeField(model).Name
		uUpper := toPascalCase(uniqeField(model).Name)
		uType := mapType(uniqeField(model).Type)

		sb.WriteString(fmt.Sprintf("type Requested%s%s struct  {\n%s *%s `json:\"%s\"`\n}\n\n",
			structName, uUpper, uUpper, uType, uLower))
		sb.WriteString(fmt.Sprintf("func (%s Requested%s%s) Validate%s(validator *validate.Validator){\n\n",
			model.Name[:1], structName, uUpper, uUpper))
		sb.WriteString(validationField(*uniqeField(model), model, "pointer"))
		sb.WriteString("}\n\n")

	}
	//Todo Create Validation
	sb.WriteString(fmt.Sprintf("func (%s *%s) Validate%s(validator *validate.Validator){\n\n",
		model.Name[0:1], structName, structName))

	// Validation
	for _, f := range model.Fields {
		sb.WriteString(validationField(f, model, ""))

	}

	sb.WriteString("}\n\n")

	sb.WriteString(fmt.Sprintf(
		"func (%s *%s) CheckUpdatedValue(requested%s *Requested%s) {\n\n",
		model.Name[0:1], structName, structName, structName))

	for _, f := range model.Fields {

		if f.Primary {
			continue
		}
		if f.Type == "datetime" {
			sb.WriteString(fmt.Sprintf(`
			if requested%s.%s != nil {
			*%s.%s = *requested%s.%s
			}
		`, structName, toPascalCase(f.Name), model.Name[:1],
				toPascalCase(f.Name), structName, toPascalCase(f.Name)))
			continue
		}

		sb.WriteString(fmt.Sprintf(`
			if requested%s.%s != nil {
			%s.%s = *requested%s.%s
			}
		`, structName, toPascalCase(f.Name), model.Name[:1],
			toPascalCase(f.Name), structName, toPascalCase(f.Name)))

	}

	sb.WriteString("}\n\n")

	if importTime {
		return "package models\n\n" +
			`import(
			 	"time"
				"backforge/internal/validate"
				)		
				` + sb.String()

	}

	return "package models\n\n" +
		`import "backforge/internal/validate"
					
		
				` + sb.String()
}

func validationField(f parser.Field, model parser.Model, point string) string {

	var sb strings.Builder

	if point == "pointer" {
		if f.Unique || f.NotNull || f.ForeignKey != nil || f.MaxLength != 0 || f.MinLength != 0 {
			if f.Type == "integer" {
				sb.WriteString(fmt.Sprintf(`
				validator.Check(*%s.%s == 0 ,"%s","the value for %s not provided")

				`,
					model.Name[0:1], toPascalCase(f.Name), f.Name, f.Name))
				if f.MinValue != 0 {
					sb.WriteString(fmt.Sprintf(`
				validator.Check(*%s.%s < %d ,"%s","the  %s must be larger or equal to %d")

				`,
						model.Name[0:1], toPascalCase(f.Name), f.MinValue, f.Name, f.Name, f.MinValue))
				}
				if f.MaxValue != 0 {
					sb.WriteString(fmt.Sprintf(`
				validator.Check(*%s.%s > %d ,"%s","the  %s must be smaller or equal to %d")

				`,
						model.Name[0:1], toPascalCase(f.Name), f.MaxValue, f.Name, f.Name, f.MaxValue))
				}

			}

			if f.Type == "text" {
				sb.WriteString(fmt.Sprintf(`validator.Check(len(*%s.%s) == 0 ,"%s","the value for %s not provided")
				
				`,
					model.Name[0:1], toPascalCase(f.Name), f.Name, f.Name))
				if f.MinLength != 0 {
					sb.WriteString(fmt.Sprintf(`
				validator.Check(len(*%s.%s) < %d ,"%s","the length of %s must be larger or equal to %d")

				`,
						model.Name[0:1], toPascalCase(f.Name), f.MinLength, f.Name, f.Name, f.MinLength))
				}
				if f.MaxLength != 0 {
					sb.WriteString(fmt.Sprintf(`
				validator.Check(len(*%s.%s) > %d ,"%s","the length of %s must be smaller or equal to %d")

				`,
						model.Name[0:1], toPascalCase(f.Name), f.MaxLength, f.Name, f.Name, f.MaxLength))
				}

			}

			if strings.ToLower(f.Name) == "email" {
				sb.WriteString(fmt.Sprintf(`validator.Check(!validate.Match(*%s.%s, validate.EmailRxp),"%s","the value for %s not valid")
				
				`,
					model.Name[0:1], toPascalCase(f.Name), f.Name, f.Name))
			}

		}
	} else {
		if f.Unique || f.NotNull || f.ForeignKey != nil || f.MaxLength != 0 || f.MinLength != 0 {
			if f.Type == "integer" {
				sb.WriteString(fmt.Sprintf(`
				validator.Check(%s.%s == 0 ,"%s","the value for %s not provided")

				`,
					model.Name[0:1], toPascalCase(f.Name), f.Name, f.Name))
				if f.MinValue != 0 {
					sb.WriteString(fmt.Sprintf(`
				validator.Check(%s.%s < %d ,"%s","the  %s must be larger or equal to %d")

				`,
						model.Name[0:1], toPascalCase(f.Name), f.MinValue, f.Name, f.Name, f.MinValue))
				}
				if f.MaxValue != 0 {
					sb.WriteString(fmt.Sprintf(`
				validator.Check(%s.%s > %d ,"%s","the  %s must be smaller or equal to %d")

				`,
						model.Name[0:1], toPascalCase(f.Name), f.MaxValue, f.Name, f.Name, f.MaxValue))
				}

			}

			if f.Type == "text" {
				sb.WriteString(fmt.Sprintf(`validator.Check(len(%s.%s) == 0 ,"%s","the value for %s not provided")
				
				`,
					model.Name[0:1], toPascalCase(f.Name), f.Name, f.Name))
				if f.MinLength != 0 {
					sb.WriteString(fmt.Sprintf(`
				validator.Check(len(%s.%s) < %d ,"%s","the length of %s must be larger or equal to %d")

				`,
						model.Name[0:1], toPascalCase(f.Name), f.MinLength, f.Name, f.Name, f.MinLength))
				}
				if f.MaxLength != 0 {
					sb.WriteString(fmt.Sprintf(`
				validator.Check(len(%s.%s) > %d ,"%s","the length of %s must be smaller or equal to %d")

				`,
						model.Name[0:1], toPascalCase(f.Name), f.MaxLength, f.Name, f.Name, f.MaxLength))
				}

			}

			if strings.ToLower(f.Name) == "email" {
				sb.WriteString(fmt.Sprintf(`validator.Check(!validate.Match(%s.%s, validate.EmailRxp),"%s","the value for %s not valid")
				
				`,
					model.Name[0:1], toPascalCase(f.Name), f.Name, f.Name))
			}

		}
	}
	return sb.String()
}
