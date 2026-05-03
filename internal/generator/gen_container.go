package generator

import (
	"fmt"
	"strings"
)

func GenerateConatiner(args ...string) string {
	var sb strings.Builder

	sb.WriteString(`
package app

import (
	"backforge/internal/handler"
	"backforge/internal/repository"
	"backforge/internal/services"
	"database/sql"
)
	`)

	sb.WriteString("type Conatiner struct {\n\n")
	for _, i := range args {
		sb.WriteString(fmt.Sprintf("	%shand handler.%sHandler\n", toPascalCase(toSingular(i)), toPascalCase(toSingular(i))))
	}
	sb.WriteString("\n}\n")

	sb.WriteString("func NewConatiner(db *sql.DB) *Conatiner {\n\n")

	for _, i := range args {

		sb.WriteString(fmt.Sprintf(`
	%sRepo := repository.New%sRepository(db)
	%sServ := services.New%sServices(%sRepo)
	%shand := handler.New%sHandler(%sServ)
		`, toSingular(i), toPascalCase(toSingular(i)),
			toSingular(i), toPascalCase(toSingular(i)), toSingular(i),
			toSingular(i), toPascalCase(toSingular(i)), toSingular(i),
		))

	}

	sb.WriteString("\nreturn &Conatiner{\n\n")
	for _, i := range args {
		sb.WriteString(fmt.Sprintf("		%shand: %shand,\n", toPascalCase(toSingular(i)), toSingular(i)))
	}
	sb.WriteString("}\n")
	sb.WriteString("}")

	return sb.String()
}
