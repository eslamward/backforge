package generator

import (
	"fmt"
	"strings"

	"github.com/eslamward/backforge/internal/parser"
)

func GenertateRepo(model parser.Model, cfg *parser.DatabaseConfig) string {
	var sb strings.Builder

	nameCapital := toPascalCase(toSingular(model.Name))
	nameLower := strings.ToLower(toSingular(model.Name))
	field := getUniqeOrPrimary(model)
	dTime := isDefaultFieldTime(model)
	sb.WriteString(`
	package repository

	import (
	"backforge/internal/models"
	"backforge/internal/backerror"

	"database/sql"
	"errors"
	"strings"
	"fmt"

	`)
	if dTime {
		sb.WriteString(`
		"time"
	)
		`)
	} else {
		sb.WriteString(`
	)
		`)
	}
	sb.WriteString(fmt.Sprintf(`
		type %sRepository interface {
		Create(*models.%s) *backerror.BackForgeError
		Update(*models.%s) *backerror.BackForgeError
	`, nameCapital, nameCapital, nameCapital))

	if uniqeField(model) != nil {
		sb.WriteString(fmt.Sprintf("GetBy%s(%s) (*models.%s ,*backerror.BackForgeError)",
			toPascalCase(field.Name), mapType(field.Type), nameCapital))
	}

	sb.WriteString(fmt.Sprintf(
		`	//Todo Field Name and Type
		Get(%s) (*models.%s ,*backerror.BackForgeError)//Todo Field Name and Type
		GetAll()([]models.%s,*backerror.BackForgeError)
		DeleteBy%s(%s) *backerror.BackForgeError
		Delete(%s) *backerror.BackForgeError
	}

	type %sRepository struct {
		DB *sql.DB
	}

	func New%sRepository(db *sql.DB) *%sRepository {
		return &%sRepository{
		DB: db,
	}
}
	`,
		mapType(primaryField(model).Type), nameCapital, nameCapital,
		toPascalCase(field.Name), mapType(field.Type), mapType(primaryField(model).Type),
		nameLower, nameCapital, nameLower, nameLower))

	sb.WriteString(createRepo(model, nameLower, nameCapital, cfg))
	/*******Update******/

	sb.WriteString(updateRepo(model, nameLower, nameCapital, cfg))
	/*******Get By******/
	sb.WriteString(getByRepoUnique(model, nameLower, nameCapital, cfg))
	/*******GET******/
	sb.WriteString(getRepo(model, nameLower, nameCapital, cfg))

	/*******GET ALL******/
	sb.WriteString(getAllRepo(field, model, nameLower, nameCapital))
	/*******Delete******/
	sb.WriteString(deleteRepo(model, nameLower, cfg))

	return sb.String()
}

func createRepo(model parser.Model, nameLower, nameCapital string, cfg *parser.DatabaseConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
	func (%s *%sRepository) Create(%s *models.%s) *backerror.BackForgeError {	
		%s

	%s.%s = %s
	//Todo Return all auto incremnet field and default field
	return nil
}
	`, nameLower[0:1], nameLower, nameLower, nameCapital,
		insertMatchField(model, cfg), nameLower,
		toPascalCase(primaryField(model).Name), primaryField(model).Name,
	))
	return sb.String()

}

func updateRepo(model parser.Model, nameLower, nameCapital string, cfg *parser.DatabaseConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
	func (%s *%sRepository) Update(%s *models.%s) *backerror.BackForgeError {
	%s
	if err != nil {
		fmt.Println("failed to update : ", err.Error())
		return backerror.New(
			backerror.DB_UPDATE,
			errors.New("failed to update %s"),
			"repository")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("failed to check affected rows: ", err.Error())
		return backerror.New(backerror.DB_UPDATE,
			errors.New("failed to check affected rows"),
			"repository")
	}

	if rows == 0 {
		return backerror.New(backerror.BAD_REQUEST,
			errors.New("%s not found"),
			"repository")
	}
	return nil
	}
	`, nameLower[0:1], nameLower, nameLower, nameCapital,
		updateMatchedField(model, cfg), nameLower, nameLower))

	return sb.String()

}
func getRepo(model parser.Model, nameLower, nameCapital string, cfg *parser.DatabaseConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
	func (%s *%sRepository) Get(%s %s) (*models.%s,*backerror.BackForgeError){ 
	%s
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, backerror.New(backerror.BAD_REQUEST,
				errors.New("the %s with this %s not found"),
				"repository")

		}
		fmt.Println("faild to query with error: ", err.Error())

		return nil, backerror.New(backerror.DB_GET,
			errors.New("faild to query the %s with error"),
			"repository")
	}

	return &%s, nil

	}
	`,
		nameLower[0:1], nameLower, primaryField(model).Name, mapType(primaryField(model).Type),
		nameCapital, getMatchedFieldForId(model, cfg),
		nameLower, primaryField(model).Name, nameLower, nameLower,
	))
	return sb.String()
}

func getAllRepo(field *parser.Field, model parser.Model, nameLower, nameCapital string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
	func (%s *%sRepository) GetAll() ([]models.%s,*backerror.BackForgeError){ 
		
	%s
	}
	`,
		nameLower[0:1], nameLower,
		nameCapital, getAllMatchedFieldForId(model),
	))

	/*******Delete By******/

	sb.WriteString(fmt.Sprintf(`
	func (%s *%sRepository) DeleteBy%s(%s %s) *backerror.BackForgeError {
	query := "DELETE FROM %s WHERE %s = ?"

	result, err := %s.DB.Exec(query, %s)
	if err != nil {
			fmt.Println("failed to delete: ", err.Error())

		return backerror.New(backerror.DB_DELETE,
			errors.New("failed to delete %s"),
			"repository")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("failed to check affected rows: ", err.Error())

		return backerror.New(backerror.DB_DELETE,
			errors.New("failed to check affected rows"),
			"repository")
	}

	if rows == 0 {
			return backerror.New(backerror.BAD_REQUEST,
			errors.New("the %s not found"),
			"repository")
	}

	return nil
}
	`, nameLower[0:1], nameLower, toPascalCase(field.Name),
		field.Name, mapType(field.Type), model.Name,
		field.Name, nameLower[0:1], field.Name, nameLower, nameLower,
	))

	return sb.String()

}

func getByRepoUnique(model parser.Model, nameLower, nameCapital string, cfg *parser.DatabaseConfig) string {
	var sb strings.Builder
	if uniqeField(model) != nil {
		sb.WriteString(fmt.Sprintf(`
	func (%s *%sRepository) GetBy%s(%s %s) (*models.%s,*backerror.BackForgeError){ 
	%s
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, backerror.New(backerror.BAD_REQUEST,
				errors.New("the %s with this %s not found"),
				"repository")

		}
		fmt.Println("faild to query with error: ", err.Error())

		return nil, backerror.New(backerror.DB_GET,
			errors.New("faild to query the %s with error"),
			"repository")
	}

	return &%s, nil

	}
	`,
			nameLower[0:1], nameLower, toPascalCase(uniqeField(model).Name), uniqeField(model).Name,
			mapType(uniqeField(model).Type), nameCapital, getMatchedField(model, cfg),
			nameLower, uniqeField(model).Name, nameLower, nameLower,
		))
	}

	return sb.String()
}
func deleteRepo(model parser.Model, nameLower string, cfg *parser.DatabaseConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
	func (%s *%sRepository) Delete(%s %s) *backerror.BackForgeError {
	query := "DELETE FROM %s WHERE %s = %s"

	result, err := %s.DB.Exec(query, %s)
	if err != nil {
			fmt.Println("failed to delete: ", err.Error())

		return backerror.New(backerror.DB_DELETE,
			errors.New("failed to delete %s"),
			"repository")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("failed to check affected rows: ", err.Error())

		return backerror.New(backerror.DB_DELETE,
			errors.New("failed to check affected rows"),
			"repository")
	}

	if rows == 0 {
			return backerror.New(backerror.BAD_REQUEST,
			errors.New("the %s not found"),
			"repository")
	}


	return nil
}
	`, nameLower[0:1], nameLower, primaryField(model).Name,
		mapType(primaryField(model).Type), model.Name,
		primaryField(model).Name, getSinglePlaceholder(cfg.Type), nameLower[0:1], primaryField(model).Name, nameLower, nameLower,
	))

	return sb.String()
}
