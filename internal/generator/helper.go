package generator

import (
	"fmt"
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

func autoIncrementField(model parser.Model) *parser.Field {
	for i := range model.Fields {
		f := &model.Fields[i]

		if f.Primary && f.Type == "integer" {
			return f
		}
	}

	log.Fatal("model must have an INTEGER PRIMARY KEY field")
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

func insertMatchField(model parser.Model) string {
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
	var numOFQuestion []string
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
		numOFQuestion = append(numOFQuestion, "?")

	}

	sb.WriteString(fmt.Sprintf(
		`
		query := "INSERT INTO %s (%s) VALUES (%s)"
	result, err := %s.DB.Exec(query,%s)
	`,
		model.Name, strings.Join(fieldsName, ", "),
		strings.Join(numOFQuestion, ", "), model.Name[0:1],
		strings.Join(fieldNameModel, ", "),
	))
	sb.WriteString(fmt.Sprintf(`
		if err != nil {
		fmt.Println("failed to insert record: ", err.Error())
		return backerror.New(backerror.DB_INSERT,
			errors.New("failed to insert %s record"),
			"repository")
		}
			

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("failed to get last insert id: ", err.Error())

		return backerror.New(
			backerror.DB_INSERT,
			errors.New("failed to get last insert id"),
			"repository")
	}
	`, toSingular(model.Name)))

	for _, v := range timeList {
		sb.WriteString(fmt.Sprintf("\nvar %s %s\n", toPascalCase(v.name), v.typ))
		sb.WriteString(fmt.Sprintf(`
			err = s.DB.QueryRow(
			"SELECT %s FROM %s WHERE id = ?",
			id,
				).Scan(&%s)
			%s.%s = &%s
		`, v.name, model.Name, toPascalCase(v.name),
			toSingular(model.Name), toPascalCase(v.name), toPascalCase(v.name)))
	}
	return sb.String()
}

func updateMatchedField(model parser.Model) string {

	var modelName = toSingular(model.Name)
	var fields []parser.Field
	var fieldsName []string
	var fieldNameModel []string
	for _, f := range model.Fields {
		if f.Primary && f.Type == "integer" {
			continue
		}

		fields = append(fields, f)

	}

	for _, f := range fields {
		name := fmt.Sprintf("%s = ?", f.Name)
		modName := fmt.Sprintf("%s.%s", modelName, toPascalCase(f.Name))

		fieldsName = append(fieldsName, name)
		fieldNameModel = append(fieldNameModel, modName)

	}

	update := fmt.Sprintf(
		`query := "UPDATE %s SET %s WHERE %s = ?"

		result, err := %s.DB.Exec(query,
		%s,%s.Id,
	)
		`, model.Name, strings.Join(fieldsName, ", "),
		primaryField(model).Name, model.Name[0:1], strings.Join(fieldNameModel, ", "), modelName)

	return update
}

func getMatchedField(model parser.Model) string {

	var fieldNameModel []string

	for _, f := range model.Fields {
		modName := fmt.Sprintf("&%s.%s", toSingular(model.Name), toPascalCase(f.Name))
		fieldNameModel = append(fieldNameModel, modName)
	}

	get := fmt.Sprintf(`
	var %s models.%s
	query := "SELECT * FROM %s where %s == ?"

	row := %s.DB.QueryRow(query, %s)

	err := row.Scan(%s,)`,
		toSingular(model.Name), toPascalCase(toSingular(model.Name)),
		model.Name, uniqeField(model).Name, model.Name[0:1], uniqeField(model).Name,
		strings.Join(fieldNameModel, ", "),
	)
	return get
}

func getMatchedFieldForId(model parser.Model) string {

	var fieldNameModel []string

	for _, f := range model.Fields {
		modName := fmt.Sprintf("&%s.%s", toSingular(model.Name), toPascalCase(f.Name))
		fieldNameModel = append(fieldNameModel, modName)
	}

	get := fmt.Sprintf(`
	var %s models.%s
	query := "SELECT * FROM %s where %s == ?"

	row := %s.DB.QueryRow(query, %s)

	err := row.Scan(%s,)`,
		toSingular(model.Name), toPascalCase(toSingular(model.Name)),
		model.Name, primaryField(model).Name, model.Name[0:1], primaryField(model).Name,
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
