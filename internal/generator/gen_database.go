package generator

import (
	"fmt"
	"strings"

	"github.com/eslamward/backforge/internal/parser"
)

func InitDB(args ...string) string {

	var sb strings.Builder

	sb.WriteString(`
	package database

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "file:./output/bin/app.db")
	if err != nil {
		log.Fatal("error opening db:", err)
	}
	db.Exec("PRAGMA foreign_keys = ON;")

`)
	for _, mod := range args {
		sb.WriteString(fmt.Sprintf("db.Exec(%sTable)\n", mod))
	}

	sb.WriteString(`
	return db
}

`)

	return sb.String()
}
func GenerateCreateTable(schema *parser.Schema) string {
	var sb strings.Builder
	sb.WriteString("package database\n\n")
	for _, model := range schema.Models {
		var cols []string
		var fks []string
		for _, f := range model.Fields {
			col := f.Name + " " + strings.ToUpper(f.Type)

			if f.Primary {
				col += " PRIMARY KEY"
			}
			if f.AutoIncrement {
				col += " AUTOINCREMENT"
			}
			if f.NotNull {
				col += " NOT NULL"
			}

			if f.Unique {
				col += " UNIQUE"
			}
			if f.Default != "" {
				col += " DEFAULT " + f.Default
			}
			if f.Check != "" {
				col += " CHECK(" + f.Check + ")"
			} else {
				if f.NotNull {
					if strings.ToLower(f.Type) == strings.ToLower("integer") {
						col += " CHECK(" + f.Name + " > 0  )"

					}
					if strings.ToLower(f.Type) == strings.ToLower("text") {
						col += " CHECK(length(" + f.Name + ") > 0  )"

					}

				}
			}
			col = "\n" + col
			cols = append(cols, col)

			if f.ForeignKey != nil {
				fk := fmt.Sprintf(
					"FOREIGN KEY (%s) REFERENCES %s(%s)",
					f.Name,
					f.ForeignKey.Model,
					f.ForeignKey.Field,
				)

				if f.ForeignKey.OnDelete != "" {
					fk += " ON DELETE " + f.ForeignKey.OnDelete
				}
				if f.ForeignKey.OnUpdate != "" {
					fk += " ON UPDATE " + f.ForeignKey.OnUpdate
				}

				fks = append(fks, fk)
			}
		}

		all := append(cols, fks...)

		sb.WriteString(fmt.Sprintf(
			"var %sTable = `CREATE TABLE IF NOT EXISTS %s (%s);`\n\n",
			model.Name, model.Name,
			strings.Join(all, ", "),
		))
	}
	return sb.String()
}
