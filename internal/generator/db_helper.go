package generator

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/eslamward/backforge/internal/parser"
)

func openRootConnection(cfg *parser.DatabaseConfig) string {
	switch cfg.Type {

	case "postgres":
		dsn := fmt.Sprintf(
			`"postgres://%s:%s@%s:%s/postgres?sslmode=disable"`,
			cfg.User, cfg.Password, cfg.Host, cfg.Port,
		)
		return fmt.Sprintf(`sql.Open("pgx", %s)`, dsn)

	case "mysql":
		dsn := fmt.Sprintf(
			`"%s:%s@tcp(%s:%s)/"`,
			cfg.User, cfg.Password, cfg.Host, cfg.Port,
		)
		return fmt.Sprintf(`sql.Open("mysql", %s)`, dsn)

	case "sqlserver":
		dsn := fmt.Sprintf(
			`"sqlserver://%s:%s@%s:%s"`,
			cfg.User, cfg.Password, cfg.Host, cfg.Port,
		)
		return fmt.Sprintf(`sql.Open("sqlserver", %s)`, dsn)

	default:
		return ""
	}
}
func buildDSN(cfg *parser.DatabaseConfig) string {

	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Name == "" {
		cfg.Type = "sqlite"
	}

	switch cfg.Type {
	case "postgres":
		return fmt.Sprintf(
			`"postgres://%s:%s@%s:%s/%s?sslmode=disable"`,
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		)

	case "mysql":
		return fmt.Sprintf(
			`"%s:%s@tcp(%s:%s)/%s?parseTime=true"`,
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		)

	case "sqlite":
		return `"file:./output/bin/app.db"`

	case "sqlserver":
		return fmt.Sprintf(
			`"sqlserver://%s:%s@%s:%s?database=%s"`,
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		)
	}

	return "file:./output/bin/app.db"
}
func getDriverName(dbType string) string {
	switch dbType {
	case "sqlite":
		return "sqlite"

	case "postgres":
		return "pgx"

	case "mysql":
		return "mysql"

	case "sqlserver":
		return "sqlserver"

	default:
		panic("unsupported database type: " + dbType)
	}
}

func getDriverImport(dbType string) string {
	switch dbType {
	case "sqlite":
		return `_ "modernc.org/sqlite"`
	case "postgres":
		return `_ "github.com/jackc/pgx/v5/stdlib"`
	case "mysql":
		return `_ "github.com/go-sql-driver/mysql"`
	case "sqlserver":
		return `_ "github.com/microsoft/go-mssqldb"`
	default:
		return `_ "modernc.org/sqlite"`
	}
}
func GetDriverGoGet(dbType string) string {
	switch dbType {
	case "sqlite":
		return "modernc.org/sqlite"

	case "postgres":
		return "github.com/jackc/pgx/v5/stdlib"

	case "mysql":
		return "github.com/go-sql-driver/mysql"

	case "sqlserver":
		return "github.com/microsoft/go-mssqldb"

	default:
		return "modernc.org/sqlite"
	}
}
func buildPlaceholders(dbType string, n int) string {
	placeholders := make([]string, n)

	switch dbType {
	case "postgres":
		for i := 0; i < n; i++ {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
	case "sqlserver":
		for i := 0; i < n; i++ {
			placeholders[i] = fmt.Sprintf("@p%d", i+1)
		}
	default:
		for i := 0; i < n; i++ {
			placeholders[i] = "?"
		}
	}

	return strings.Join(placeholders, ", ")
}
func getSinglePlaceholder(dbType string) string {
	switch dbType {
	case "postgres":
		return "$1"
	case "sqlserver":
		return "@p1"
	default:
		return "?"
	}
}
func getSinglePlaceholderWithoutNumber(dbType string) string {
	switch dbType {
	case "postgres":
		return "$"
	case "sqlserver":
		return "@p"
	default:
		return "?"
	}
}
func insertAndReturnFields(
	db *sql.DB,
	dbType string,
	query string,
	args ...any,
) (int64, error) {

	switch dbType {

	case "postgres":
		var id int64
		err := db.QueryRow(query+" RETURNING id", args...).Scan(&id)
		return id, err

	case "sqlserver":
		var id int64
		err := db.QueryRow(query+"; SELECT SCOPE_IDENTITY();", args...).Scan(&id)
		return id, err

	default:
		result, err := db.Exec(query, args...)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()
	}
}

func insertAndReturnID(modelName, db, dbType, query, sPr string, args ...string) string {
	var s string
	if len(args) == 0 {
		s = ""
	} else {
		s = strings.Join((args), ", ")
	}

	var sb strings.Builder

	switch dbType {

	case "postgres":
		sb.WriteString(fmt.Sprintf(`
		var %s int64
		err := %s.DB.QueryRow("%s RETURNING %s", %s).Scan(&%s)
		`, sPr, db, query, sPr, s, sPr))

	case "sqlserver":
		sb.WriteString(fmt.Sprintf(`var %s int64
		err := %s.DB.QueryRow("%s+; SELECT SCOPE_IDENTITY();", %s).Scan(&%s)
		`, sPr, db, query, s, sPr))

	default:
		sb.WriteString(fmt.Sprintf(`result, err := %s.DB.Exec("%s", %s)
		
		`, db, query, s))
		sb.WriteString(fmt.Sprintf(`
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()),"foreign"){
		return backerror.New(backerror.FOREIGN_CONSTRAINTS,
			errors.New("failed to insert %s record dueto foreign constraints check value of this field"),
			"repository")
		}
		fmt.Println("failed to insert record: ", err.Error())
		return backerror.New(backerror.DB_INSERT,
			errors.New("failed to insert %s record"),
			"repository")
		}
			
		`, modelName, modelName))
		sb.WriteString(fmt.Sprintf(`%s,err := result.LastInsertId()`, sPr))
		sb.WriteString(`
		
		if err != nil {
	
		return backerror.New(
			backerror.DB_INSERT,
			errors.New("failed to get last insert id"),
			"repository")
	 	}
			`)

	}

	if dbType == "postgres" || dbType == "sqlserver" {

		sb.WriteString(fmt.Sprintf(`
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()),"foreign"){
		return backerror.New(backerror.FOREIGN_CONSTRAINTS,
			errors.New("failed to insert %s record dueto foreign constraints check value of this field"),
			"repository")
		}
		fmt.Println("failed to insert record: ", err.Error())
		return backerror.New(backerror.DB_INSERT,
			errors.New("failed to insert %s record"),
			"repository")
		}
			
		`, modelName, modelName))
	}

	return sb.String()
}
func isNumericType(t string) bool {
	switch strings.ToLower(t) {

	case "integer", "int", "bigint", "smallint":
		return true

	default:
		return false
	}
}

func validateAutoIncrement(field parser.Field, dbType string) error {

	if field.AutoIncrement && !isNumericType(field.Type) {
		return fmt.Errorf(
			"field %s: auto_increment requires numeric type",
			field.Name,
		)
	}

	if dbType == "sqlite" &&
		field.AutoIncrement &&
		!field.Primary {

		return fmt.Errorf(
			"sqlite requires auto_increment field to be primary key: %s",
			field.Name,
		)
	}

	return nil
}

func mapSQLType(field parser.Field, dbType string) string {

	if dbType == "postgres" &&
		field.AutoIncrement {

		return "SERIAL"
	}

	switch strings.ToLower(field.Type) {

	case "integer", "int":
		switch dbType {

		case "mysql":
			return "INT"

		case "sqlserver":
			return "INT"

		default:
			return "INTEGER"
		}

	case "bigint":
		return "BIGINT"

	case "smallint":
		return "SMALLINT"

	case "text":
		switch dbType {

		case "mysql":
			return "TEXT"

		case "sqlserver":
			return "NVARCHAR(MAX)"

		default:
			return "TEXT"
		}

	case "boolean", "bool":
		switch dbType {

		case "mysql":
			return "BOOLEAN"

		case "sqlserver":
			return "BIT"

		default:
			return "BOOLEAN"
		}
	case "datetime":

		switch dbType {

		case "postgres":
			return "TIMESTAMP"

		default:
			return "DATETIME"
		}

	default:
		return strings.ToUpper(field.Type)
	}
}
func getAutoIncrementSQL(dbType string) string {

	// postgres SERIAL already handles auto increment
	if dbType == "postgres" {
		return ""
	}

	switch dbType {

	case "sqlite":
		return " PRIMARY KEY AUTOINCREMENT"

	case "mysql":
		return "AUTO_INCREMENT"

	case "sqlserver":
		return "IDENTITY(1,1)"

	default:
		return ""
	}
}
func isAutoIncrementField(f parser.Field) bool {

	if isNumericType(f.Type) && f.AutoIncrement {
		return true
	}

	return false
}
func autoIncrementField(model parser.Model) *parser.Field {
	for _, f := range model.Fields {
		if isNumericType(f.Type) && f.AutoIncrement {
			return &f
		}
	}

	return nil
}

func insertMatchField(model parser.Model, cfg *parser.DatabaseConfig) string {
	sPr := primaryField(model).Name
	var sb strings.Builder
	//Todo Insert into fields that not primary
	var timeList []struct {
		name string
		typ  string
	}
	var modelName = toSingular(model.Name)
	var fields []parser.Field
	var fieldsName []string
	var fieldNameModel []string
	numOfPlaceHolder := 0
	for _, f := range model.Fields {
		if f.Primary && f.Type == "integer" {
			continue
		}
		if f.Default != "" && f.Type == "datetime" {
			timeList = append(timeList, struct {
				name string
				typ  string
			}{name: f.Name, typ: mapType(f.Type)})
			continue
		}

		fields = append(fields, f)

	}

	for _, f := range fields {
		name := fmt.Sprintf("%s", f.Name)
		modName := fmt.Sprintf("%s.%s", modelName, toPascalCase(f.Name))

		fieldsName = append(fieldsName, name)
		fieldNameModel = append(fieldNameModel, modName)
		numOfPlaceHolder += 1

	}

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		model.Name, strings.Join(fieldsName, ", "),
		buildPlaceholders(cfg.Type, numOfPlaceHolder),
	)

	sb.WriteString(fmt.Sprintf(`%s`,
		insertAndReturnID(modelName, modelName[:1], cfg.Type, query, sPr,
			fieldNameModel...,
		)))

	for _, v := range timeList {
		sb.WriteString(fmt.Sprintf("\nvar %s %s\n", toPascalCase(v.name), v.typ))
		sb.WriteString(fmt.Sprintf(`
			err = s.DB.QueryRow(
			"SELECT %s FROM %s WHERE %s = %s",
			%s,
				).Scan(&%s)
			%s.%s = &%s
		`, v.name, model.Name, sPr,
			getSinglePlaceholder(cfg.Type), sPr, toPascalCase(v.name),
			toSingular(model.Name), toPascalCase(v.name), toPascalCase(v.name)))
	}
	return sb.String()
}

func updateMatchedField(model parser.Model, cfg *parser.DatabaseConfig) string {

	var modelName = toSingular(model.Name)
	var fields []parser.Field
	var fieldsName []string
	var fieldNameModel []string
	for _, f := range model.Fields {
		if f.Primary || isAutoIncrementField(f) {
			continue
		}

		fields = append(fields, f)

	}

	num := 0

	for i, f := range fields {
		name := ""
		if cfg.Type != "sqlite" {
			name = fmt.Sprintf("%s = %s%d", f.Name, getSinglePlaceholderWithoutNumber(cfg.Type), i+1)
		} else {
			name = fmt.Sprintf("%s = %s", f.Name, getSinglePlaceholder(cfg.Type))

		}

		modName := fmt.Sprintf("%s.%s", modelName, toPascalCase(f.Name))

		fieldsName = append(fieldsName, name)
		fieldNameModel = append(fieldNameModel, modName)

		num = i + 1
	}

	update := fmt.Sprintf(
		`query := "UPDATE %s SET %s WHERE %s = %s%d"

		result, err := %s.DB.Exec(query,
		%s,%s.%s,
	)
		`, model.Name, strings.Join(fieldsName, ", "),
		primaryField(model).Name, getSinglePlaceholderWithoutNumber(cfg.Type), num+1,
		model.Name[0:1], strings.Join(fieldNameModel, ", "), modelName,
		toPascalCase(primaryField(model).Name))

	return update
}

func getMatchedField(model parser.Model, cfg *parser.DatabaseConfig) string {

	var fieldNameModel []string

	for _, f := range model.Fields {
		modName := fmt.Sprintf("&%s.%s", toSingular(model.Name), toPascalCase(f.Name))
		fieldNameModel = append(fieldNameModel, modName)
	}

	get := fmt.Sprintf(`
	var %s models.%s
	query := "SELECT * FROM %s where %s = %s"

	row := %s.DB.QueryRow(query, %s)

	err := row.Scan(%s,)`,
		toSingular(model.Name), toPascalCase(toSingular(model.Name)),
		model.Name, uniqeField(model).Name, getSinglePlaceholder(cfg.Type), model.Name[0:1], uniqeField(model).Name,
		strings.Join(fieldNameModel, ", "),
	)
	return get
}

func getMatchedFieldForId(model parser.Model, cfg *parser.DatabaseConfig) string {

	var fieldNameModel []string

	for _, f := range model.Fields {
		modName := fmt.Sprintf("&%s.%s", toSingular(model.Name), toPascalCase(f.Name))
		fieldNameModel = append(fieldNameModel, modName)
	}

	get := fmt.Sprintf(`
	var %s models.%s
	query := "SELECT * FROM %s where %s = %s"

	row := %s.DB.QueryRow(query, %s)

	err := row.Scan(%s,)`,
		toSingular(model.Name), toPascalCase(toSingular(model.Name)),
		model.Name, primaryField(model).Name, getSinglePlaceholder(cfg.Type), model.Name[0:1], primaryField(model).Name,
		strings.Join(fieldNameModel, ", "),
	)
	return get
}
func getAllMatchedFieldForId(model parser.Model) string {

	var fieldNameModel []string

	for _, f := range model.Fields {
		modName := fmt.Sprintf("&%s.%s", toSingular(model.Name), toPascalCase(f.Name))
		fieldNameModel = append(fieldNameModel, modName)
	}

	getAll := fmt.Sprintf(`
	var %s []models.%s
	query := "SELECT * FROM %s"

	rows,err := %s.DB.Query(query)
	if err != nil {

		fmt.Println("faild to query with error: ", err.Error())

		return nil, backerror.New(backerror.DB_GET,
			errors.New("faild to query the %s with error"),
			"repository")
	}

	
	for rows.Next() {
		var %s models.%s
		err := rows.Scan(%s)
		if err != nil {

			fmt.Println("faild to query with error: ", err.Error())
			return nil, backerror.New(backerror.DB_GET,
				errors.New("faild to query the %s with error"),
				"repository")
		}

		%s = append(%s, %s)

	}
	if err := rows.Err(); err != nil {

		fmt.Println("faild to query with error: ", err.Error())

		return nil, backerror.New(backerror.DB_GET,
			errors.New("faild to query the %s with error"),
			"repository")
	}
		
	return %s,nil

	`,
		model.Name, toPascalCase(toSingular(model.Name)),
		model.Name, model.Name[0:1], model.Name,
		toSingular(model.Name), toPascalCase(toSingular(model.Name)),
		strings.Join(fieldNameModel, ", "), model.Name, model.Name, model.Name,
		toSingular(model.Name), model.Name, model.Name,
	)
	return getAll
}
