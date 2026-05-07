package generator

import (
	"fmt"
	"strings"

	"github.com/eslamward/backforge/internal/parser"
)

func GenerateServices(model parser.Model) string {
	var sb strings.Builder

	nameCapital := toPascalCase(toSingular(model.Name))
	nameLower := strings.ToLower(toSingular(model.Name))

	sb.WriteString(`
	package services

import (
	"backforge/internal/backerror"
	"backforge/internal/models"
	"backforge/internal/repository"
	"backforge/internal/validate"
	"errors"
)



	`)

	sb.WriteString(fmt.Sprintf(`
	
type %sService interface {
	Create%s(*models.%s) *backerror.BackForgeError
	Update%s(*models.Requested%s) (*models.%s, *backerror.BackForgeError)
	Get%s(%s) (*models.%s,*backerror.BackForgeError)
	Get%ss() ([]models.%s,*backerror.BackForgeError)
	`, nameCapital, nameCapital, nameCapital, nameCapital, nameCapital,
		nameCapital, nameCapital, mapType(primaryField(model).Type),
		nameCapital, nameCapital, nameCapital))

	if uniqeField(model) != nil {
		sb.WriteString(fmt.Sprintf("Get%sBy%s(%s) (*models.%s,*backerror.BackForgeError)\n",
			nameCapital, toPascalCase(uniqeField(model).Name), mapType(uniqeField(model).Type),
			nameCapital,
		))
	}

	sb.WriteString(fmt.Sprintf(`
	Delete%s(%s) *backerror.BackForgeError
}

type %sServices struct {
	%sRepo repository.%sRepository
}

func New%sServices(%sRepo repository.%sRepository) *%sServices {
	return &%sServices{
		%sRepo: %sRepo,
	}
}
	`,
		nameCapital, mapType(primaryField(model).Type),
		nameLower, nameLower, nameCapital, nameCapital,
		nameLower, nameCapital, nameLower,
		nameLower, nameLower, nameLower))

	/**********Creat***********/

	sb.WriteString(createSer(nameLower, nameCapital, model))

	/**********Update***********/

	sb.WriteString(updateSer(nameLower, nameCapital, model))

	/**********Get***********/

	sb.WriteString(getSer(nameLower, nameCapital, model))

	/**********Get By Unique***********/

	sb.WriteString(getByUniqeSer(nameLower, nameCapital, model))

	/**********Get All***********/

	sb.WriteString(getAllSer(nameLower, nameCapital, model))

	/**********Delete***********/

	sb.WriteString(deleteSer(nameLower, nameCapital, model))

	return sb.String()
}

func createSer(nameLower, nameCapital string, model parser.Model) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		`
	func (%s *%sServices) Create%s(%s *models.%s) *backerror.BackForgeError {

	validator := validate.New()
	%s.Validate%s(validator)
	if !validator.Valid() {
		return backerror.New(backerror.VALIDATION, validator.Errors, "services")
	}
	`, nameLower[0:1], nameLower, nameCapital, nameLower, nameCapital,
		nameLower, nameCapital))

	if uniqeField(model) != nil {
		sb.WriteString(fmt.Sprintf(
			`
	selected%s, err := %s.%sRepo.GetBy%s(%s.%s)//Todo check if uniqe field put this

	if selected%s != nil {
		return backerror.New(backerror.UNIQUE_CONSTRAINTS ,
		errors.New("this %s already exists"),
		 "services")
	}
	`, nameCapital, nameLower[0:1], nameLower,
			toPascalCase(uniqeField(model).Name),
			nameLower, toPascalCase(uniqeField(model).Name),
			nameCapital, uniqeField(model).Name))
	}

	if uniqeField(model) != nil {
		sb.WriteString(fmt.Sprintf("err = %s.%sRepo.Create(%s)", nameLower[:1], nameLower, nameLower))
	} else {
		sb.WriteString(fmt.Sprintf("err := %s.%sRepo.Create(%s)", nameLower[:1], nameLower, nameLower))

	}

	sb.WriteString(
		`

	if err != nil {
		return err
	}

	return nil
	}
		`)

	return sb.String()
}
func updateSer(nameLower, nameCapital string, model parser.Model) string {

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		`
		func (%s *%sServices) Update%s(req%s *models.Requested%s) (*models.%s,*backerror.BackForgeError) {
		%s, err := %s.%sRepo.Get(req%s.%s)
		//Todo handle errors
		if err != nil {
			return nil, err
		}
		%s.CheckUpdatedValue(req%s)
		validator := validate.New()
		%s.Validate%s(validator)

		if !validator.Valid() {
		return nil, backerror.New(backerror.VALIDATION, validator.Errors, "services")
		}
		err = %s.%sRepo.Update(%s)
		if err != nil {
			return nil, err
		}

		return %s,nil
		}
		`, nameLower[0:1], nameLower, nameCapital, nameCapital, nameCapital,
		nameCapital, toSingular(nameLower), nameLower[0:1], toSingular(nameLower),
		nameCapital,
		toPascalCase(primaryField(model).Name), nameLower,
		nameCapital, nameLower, nameCapital, nameLower[0:1],
		nameLower, nameLower, nameLower,
	))
	return sb.String()
}

func getSer(nameLower, nameCapital string, model parser.Model) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`
		func (%s *%sServices) Get%s(%s %s) (*models.%s,*backerror.BackForgeError){

		%s, err := %s.%sRepo.Get(%s)

		if err != nil {
			return nil, err
		}
		if %s != nil {

			return %s, nil
		}

		return nil,backerror.New(backerror.BAD_REQUEST,
		 errors.New("can't retrive the %s"), "services")

		}
		`, nameLower[0:1], nameLower, nameCapital, primaryField(model).Name,
		mapType(primaryField(model).Type), nameCapital,
		nameLower, nameLower[0:1], nameLower,
		primaryField(model).Name, nameLower, nameLower, nameLower,
	))
	return sb.String()
}

func getByUniqeSer(nameLower, nameCapital string, model parser.Model) string {
	var sb strings.Builder
	if uniqeField(model) != nil {

		sb.WriteString(fmt.Sprintf(
			`
		func (%s *%sServices) Get%sBy%s(%s %s) (*models.%s,*backerror.BackForgeError){

		
		req := models.Requested%s%s{
			
			%s : &%s,
		}

		validator := validate.New()
		req.Validate%s(validator)
		if !validator.Valid() {
		return nil,backerror.New(backerror.VALIDATION, validator.Errors, "services")
		}

		%s, err := %s.%sRepo.GetBy%s(%s)

		if err != nil {
			return nil, err
		}
		if %s != nil {


			return %s, nil
		}

		return nil,backerror.New(backerror.BAD_REQUEST,
		 errors.New("can't retrive the %s"), "services")

}
		`, nameLower[0:1], nameLower, nameCapital,
			toPascalCase(uniqeField(model).Name), uniqeField(model).Name,
			mapType(uniqeField(model).Type), nameCapital,
			nameCapital, toPascalCase(uniqeField(model).Name), toPascalCase(uniqeField(model).Name),
			uniqeField(model).Name, toPascalCase(uniqeField(model).Name),
			nameLower, nameLower[0:1], nameLower, toPascalCase(uniqeField(model).Name),
			uniqeField(model).Name, nameLower, nameLower, nameLower,
		))

	}
	return sb.String()
}

func getAllSer(nameLower, nameCapital string, model parser.Model) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		`
		func (%s *%sServices) Get%ss() ([]models.%s,*backerror.BackForgeError){

		%s, err := %s.%sRepo.GetAll()
		if err != nil{
		return nil,err
		}
		return %s, nil
		
}
		`, nameLower[0:1], nameLower, nameCapital,
		nameCapital,
		model.Name, nameLower[0:1], nameLower, model.Name,
	))
	return sb.String()
}

func deleteSer(nameLower, nameCapital string, model parser.Model) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`
		func (%s *%sServices) Delete%s(%s %s) *backerror.BackForgeError{

		err := %s.%sRepo.Delete(%s)

		if err != nil {
			return err
		}
	return nil
	}
		`, nameLower[0:1], nameLower, nameCapital, primaryField(model).Name,
		mapType(primaryField(model).Type), nameLower[0:1],
		nameLower, primaryField(model).Name,
	))

	return sb.String()
}
