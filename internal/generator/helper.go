package generator

import (
	"log"
	"strings"

	"github.com/eslamward/backforge/internal/parser"
)

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func toSingular(s string) string {
	if strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}
	return s
}

func mapType(t string) string {
	switch strings.ToLower(t) {
	case "integer":
		return "int64"
	case "bool":
		return "bool"
	case "text":
		return "string"
	case "datetime":
		return "time.Time"
	default:
		return "string"
	}
}

func getUniqeOrPrimary(model parser.Model) *parser.Field {

	var field *parser.Field

	for _, f := range model.Fields {
		if f.Primary {
			field = &f
		}
		if f.Unique {
			field = &f
			break
		}
	}

	return field

}
func uniqeField(model parser.Model) *parser.Field {

	var field *parser.Field
	for _, f := range model.Fields {

		if f.Unique {
			field = &f
			break
		}
	}

	return field

}

func primaryField(model parser.Model) *parser.Field {
	for i := range model.Fields {
		f := &model.Fields[i]

		if f.Primary {
			return f
		}
	}

	log.Fatal("model must have an PRIMARY KEY field")
	return nil
}

func isDefaultFieldTime(model parser.Model) bool {
	for i := range model.Fields {
		f := &model.Fields[i]

		if f.Default != "" && f.Type == "datetime" {
			return true
		}
	}
	return false
}
