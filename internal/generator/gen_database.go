package generator

import (
	"fmt"
	"log"
	"strings"

	"github.com/eslamward/backforge/internal/parser"
)

func InitDB(cfg *parser.DatabaseConfig, args ...string) string {

	dsn := buildDSN(cfg)
	im := `
	
	"fmt"
	"strings"
	`
	if cfg.Type == "sqlite" {
		im = ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
		package database

		import (
			"database/sql"
			"log"
			%s
			%s
		
		)

		func InitDB() *sql.DB {

		

		`, im, getDriverImport(cfg.Type)))
	if cfg.Type != "sqlite" {
		fmt.Println(cfg.Type, cfg.Name)
		sb.WriteString(fmt.Sprintf(`

			db,err := %s
			if err != nil{
				log.Fatal("error openinig database :",err)
			}
		`, openRootConnection(cfg)))
		sb.WriteString(fmt.Sprintf(`_,err = db.Exec("CREATE DATABASE %s")`, cfg.Name))
		sb.WriteString(fmt.Sprintf(
			`
			if err != nil{
				if strings.Contains(err.Error(),"exists"){
				fmt.Println("database already created")
				}else{		
				log.Fatal("error creating database",err)

			}
		}
			`))

		sb.WriteString(fmt.Sprintf(`
			db, err = sql.Open("%s", %s)
			if err != nil {
				log.Fatal("error opening db:", err)
			}
		`, getDriverName(cfg.Type), dsn))
	} else {

		sb.WriteString(fmt.Sprintf(`
			db, err := sql.Open("%s", %s)
			if err != nil {
				log.Fatal("error opening db:", err)
			}
		`, getDriverName(cfg.Type), dsn))
		sb.WriteString(`
		if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			log.Fatal("failed to enable foreign keys:", err)
	}
`)
	}

	for _, mod := range args {
		sb.WriteString(fmt.Sprintf("_,err= db.Exec(%sTable)\n", mod))
		sb.WriteString(`
		if err != nil{
			log.Fatal("error creating tables:",err)
		}
			`)
	}

	sb.WriteString(`
	return db
}

`)

	return sb.String()
}
func GenerateCreateTable(schema *parser.Schema, cfg *parser.DatabaseConfig) string {
	var sb strings.Builder
	sb.WriteString("package database\n\n")
	for _, model := range schema.Models {
		var cols []string
		var fks []string
		for _, f := range model.Fields {
			col := f.Name + " " + mapSQLType(f, cfg.Type)

			if f.AutoIncrement {
				err := validateAutoIncrement(f, cfg.Type)
				if err != nil {
					log.Fatal("error validation autoncrement: ", err)
				}
				col += " " + getAutoIncrementSQL(cfg.Type)
			}

			if f.Primary && cfg.Type != "sqlite" {

				col += " PRIMARY KEY"
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
