package generator

import (
	"fmt"
	"strings"

	"github.com/eslamward/backforge/internal/parser"
)

func GenerateHandler(model parser.Model) string {
	var sb strings.Builder

	nameCapital := toPascalCase(toSingular(model.Name))
	nameLower := strings.ToLower(toSingular(model.Name))

	sb.WriteString("package handler\n\n")
	sb.WriteString(`import (
	"backforge/internal/models"
	"backforge/internal/services"
	"net/http"
	"errors"
	
	)
	`)

	sb.WriteString(fmt.Sprintf(`
	type %sHandler interface {
	Create%s(http.ResponseWriter, *http.Request)
	Update%s(http.ResponseWriter, *http.Request)
	Get%s(http.ResponseWriter, *http.Request)
	Get%ss(http.ResponseWriter, *http.Request)
		`, nameCapital, nameCapital, nameCapital, nameCapital, nameCapital))

	if uniqeField(model) != nil {
		sb.WriteString(fmt.Sprintf("Get%sBy%s(http.ResponseWriter, *http.Request)\n",
			nameCapital, toPascalCase(uniqeField(model).Name),
		))
	}

	sb.WriteString(fmt.Sprintf(`
	Delete%s(http.ResponseWriter, *http.Request)
}

type %sHandler struct {
	%sServices services.%sService
}

func New%sHandler(%sServices services.%sService) *%sHandler {
	return &%sHandler{
		%sServices: %sServices,
	}
}
	`, nameCapital, nameLower, nameLower, nameCapital, nameCapital,
		nameLower, nameCapital, nameLower, nameLower, nameLower, nameLower,
	))
	//create
	sb.WriteString(createHand(nameLower, nameCapital))
	// GET

	sb.WriteString(getHand(nameLower, nameCapital, model))

	// GET By UniqField
	sb.WriteString(getByUniqueHand(nameLower, nameCapital, model))

	// GET All

	sb.WriteString(getAllHand(nameLower, nameCapital, model))

	// UPDATE
	sb.WriteString(updateHand(nameLower, nameCapital, model))

	// DELETE
	sb.WriteString(deleteHand(nameLower, nameCapital, model))

	return sb.String()
}

func createHand(nameLower, nameCapital string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
	func (%s *%sHandler)Create%s(w http.ResponseWriter ,r *http.Request){
		var %s models.%s

		err := readJson(w, r, &%s)
		if err != nil {
			badErrorRequest(w, r, err)
			return
		}

		serErr := %s.%sServices.Create%s(&%s)
		if serErr != nil {
			mapResponseToError(serErr,w,r)
			return
		}

		err = writeJSon(w, http.StatusOK, Encapsulate{"%s": %s}, nil)
		if err != nil {
			serverErrorResonse(w, r, err)
			return
	}

		


	}
	`, nameLower[0:1], nameLower, nameCapital, nameLower, nameCapital,
		nameLower, nameLower[0:1], nameLower, nameCapital, nameLower,
		nameLower, nameLower,
	))

	return sb.String()
}

func getHand(nameLower, nameCapital string, model parser.Model) string {
	var sb strings.Builder
	pr := toPascalCase(primaryField(model).Name)
	sPr := primaryField(model).Name
	sb.WriteString(fmt.Sprintf(`
	func (%s *%sHandler)Get%s(w http.ResponseWriter, r *http.Request) {
		var req%s models.Requested%s%s

	err := readJson(w, r, &req%s)
	if err != nil {
		badErrorRequest(w, r, err)
		return
	}
	if req%s.%s == 0 {
		badErrorRequest(w, r, errors.New("please put the %s field in the json body"))
		return
	}
	if req%s.%s < 0 {
		badErrorRequest(w, r, errors.New("please use valid %s"))
		return
	}

	%s, serErr := %s.%sServices.Get%s(req%s.%s)
	if serErr != nil {
			mapResponseToError(serErr,w,r)
		return
	}

	err = writeJSon(w, http.StatusOK, Encapsulate{"%s": %s}, nil)
	if err != nil {
		serverErrorResonse(w, r, err)
	}
	}
`, nameLower[0:1], nameLower, nameCapital,
		nameCapital, nameCapital, pr, nameCapital,
		nameCapital, pr, sPr, nameCapital, pr, sPr, nameLower,
		nameLower[0:1], nameLower, nameCapital, nameCapital,
		pr, nameLower, nameLower))
	return sb.String()
}

func getByUniqueHand(nameLower, nameCapital string, model parser.Model) string {
	var sb strings.Builder
	if uniqeField(model) != nil {
		sb.WriteString(fmt.Sprintf(`
	func (%s *%sHandler)Get%sBy%s(w http.ResponseWriter, r *http.Request) {
		var req%s models.Requested%s%s

	err := readJson(w, r, &req%s)
	if err != nil {
		badErrorRequest(w, r, err)
		return
	}
	if req%s.%s == nil {
		badErrorRequest(w, r, errors.New("please put the %s field in the json body"))
		return
	}
	

	%s, serErr := %s.%sServices.Get%sBy%s(*req%s.%s)
	if serErr != nil {
			mapResponseToError(serErr,w,r)
		return
	}

	err = writeJSon(w, http.StatusOK, Encapsulate{"%s": %s}, nil)
	if err != nil {
		serverErrorResonse(w, r, err)
	}
	}
`, nameLower[0:1], nameLower, nameCapital, toPascalCase(uniqeField(model).Name),
			nameCapital, nameCapital, toPascalCase(uniqeField(model).Name), nameCapital,
			nameCapital, toPascalCase(uniqeField(model).Name),
			uniqeField(model).Name,
			// nameCapital, toPascalCase(uniqeField(model).Name),uniqeField(model).Name,
			nameLower, nameLower[0:1], nameLower, nameCapital,
			toPascalCase(uniqeField(model).Name), nameCapital,
			toPascalCase(uniqeField(model).Name), nameLower, nameLower))
	}
	return sb.String()
}

func getAllHand(nameLower, nameCapital string, model parser.Model) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
	func (%s *%sHandler)Get%ss(w http.ResponseWriter, r *http.Request) {


	%s, serErr := %s.%sServices.Get%ss()
	if serErr != nil {
			mapResponseToError(serErr,w,r)
		return
	}

	err := writeJSon(w, http.StatusOK, Encapsulate{"%s": %s}, nil)
	if err != nil {
		serverErrorResonse(w, r, err)
	}
	}
	`, nameLower[0:1], nameLower, nameCapital,
		model.Name, nameLower[0:1], nameLower, nameCapital, model.Name, model.Name))
	return sb.String()
}

func updateHand(nameLower, nameCapital string, model parser.Model) string {
	var sb strings.Builder
	pr := toPascalCase(primaryField(model).Name)
	sPr := primaryField(model).Name
	sb.WriteString(fmt.Sprintf(`
	func (%s *%sHandler)Update%s(w http.ResponseWriter, r *http.Request) {
	var req%s models.Requested%s

	err := readJson(w, r, &req%s)
	if err != nil {
			badErrorRequest(w,r,err)
		return
	}

	if req%s.%s == 0 {
		badErrorRequest(w, r, errors.New("please put the %s in the json body"))
		return
	}
	if req%s.%s < 0 {
		badErrorRequest(w, r, errors.New("please use valid %s"))
		return
	}

	%s, serErr := %s.%sServices.Update%s(&req%s)
	if serErr != nil {
			mapResponseToError(serErr,w,r)
		return
	}

	err = writeJSon(w, http.StatusOK, Encapsulate{"%s": %s}, nil)
	if err != nil {
		serverErrorResonse(w, r, err)
	}
	}
`, nameLower[0:1], nameLower, nameCapital,
		nameCapital, nameCapital, nameCapital, nameCapital,
		pr, sPr, nameCapital, pr, sPr, nameLower,
		nameLower[0:1], nameLower, nameCapital, nameCapital, nameLower, nameLower,
	))

	return sb.String()
}

func deleteHand(nameLower, nameCapital string, model parser.Model) string {
	pr := toPascalCase(primaryField(model).Name)

	sPr := primaryField(model).Name
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
	func(%s *%sHandler) Delete%s(w http.ResponseWriter, r *http.Request) {
		
	var req%s models.Requested%s%s

	err := readJson(w, r, &req%s)
	if err != nil {
		badErrorRequest(w, r, err)
		return
	}

	
	if req%s.%s == 0 {
		badErrorRequest(w, r, errors.New("please put the %s field in the json body"))
		return
	}
	if req%s.%s < 0 {
		badErrorRequest(w, r, errors.New("please use valid %s"))
		return
	}
	serErr := %s.%sServices.Delete%s(req%s.%s)
	if serErr != nil {
			mapResponseToError(serErr,w,r)
		return
	}

	err = writeJSon(w, http.StatusOK, Encapsulate{"delete": "success"}, nil)
	if err != nil {
		serverErrorResonse(w, r, err)
	}
	}
`, nameLower[0:1], nameLower, nameCapital, pr, nameCapital, pr, pr,
		pr, pr, sPr, pr, pr, sPr,
		nameLower[0:1], nameLower, nameCapital, pr, pr,
	))
	return sb.String()
}

/**************************Gen helpers function*********************************/

func GenerateHandlerHelper() string {

	return `
	
	package handler

import (
	"backforge/internal/backerror"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Encapsulate map[string]any

func writeJSon(w http.ResponseWriter, status int, data Encapsulate, headers http.Header) error {

	jsonedData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonedData)
	return err

}
func readJson(w http.ResponseWriter, r *http.Request, data any) error {

	var maxByteSize int64 = 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxByteSize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)

	if err != nil {
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("json body constain invalid json at character: %d", syntaxError.Offset)

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("json body has incorrect json type for field: %s", unmarshalTypeError.Field)
			}
			return fmt.Errorf("json body is invalid type at charcter: %d", unmarshalTypeError.Offset)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("json body can't be larger than: %d", maxByteSize)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("json body conatin invalid json value")

		case errors.Is(err, io.EOF):
			return errors.New("json body can't be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body conatin unknow field: %s", field)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}

	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must conatin one json value")

	}

	return nil
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {

	errData := Encapsulate{"error": message}

	err := writeJSon(w, status, errData, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func serverErrorResonse(w http.ResponseWriter, r *http.Request, err error) {

	message := "the serve encounter a problem"
	errorResponse(w, r, http.StatusInternalServerError, message)

}

func badErrorRequest(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "can't find resources with this request"
	errorResponse(w, r, http.StatusNotFound, message)
}
func notAllowedResponse(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	message := fmt.Sprintf("the %s mthod not allowed", method)
	errorResponse(w, r, http.StatusNotFound, message)
}

func badErrorValueValidator(w http.ResponseWriter, r *http.Request, err map[string]string) {
	errorResponse(w, r, http.StatusUnprocessableEntity, err)
}

func errorConflict(w http.ResponseWriter, r *http.Request,err error) {
	errorResponse(w, r, http.StatusConflict, err.Error())
}



	func mapResponseToError(bk *backerror.BackForgeError, w http.ResponseWriter, r *http.Request) {
		var err error
	switch bk.Type {
	case backerror.DB_DELETE:
		err = bk.Message.(error)
		serverErrorResonse(w, r, err)
	case backerror.DB_INSERT:
		err = bk.Message.(error)
		serverErrorResonse(w, r, err)
	case backerror.DB_UPDATE:
		err = bk.Message.(error)
		serverErrorResonse(w, r, err)
	case backerror.DB_GET:
		err = bk.Message.(error)
		serverErrorResonse(w, r, err)
	case backerror.BAD_REQUEST:
		err = bk.Message.(error)
		badErrorRequest(w, r,err)
	case backerror.VALIDATION:
		vErr := bk.Message.(map[string]string)
		badErrorValueValidator(w, r, vErr)
	case backerror.UNIQUE_CONSTRAINTS:
		err = bk.Message.(error)
		fmt.Println(err,bk.Message,bk)

		errorConflict(w, r, err)
	case backerror.FOREIGN_CONSTRAINTS:
		err = bk.Message.(error)
		fmt.Println(err,bk.Message,bk)
		badErrorRequest(w, r,err)

	default:
		serverErrorResonse(w, r, bk.Message.(error))
	}
}

	
	`
}
