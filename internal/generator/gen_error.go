package generator

func GenerateError() string {

	return `
	package backerror

	const (
	DB_INSERT   = "db_insert"
	DB_UPDATE   = "db_update"
	DB_GET      = "db_get"
	DB_DELETE   = "db_delete"
	VALIDATION  = "validation"
	BAD_REQUEST = "bad_request"
	UNIQUE_CONSTRAINTS = "unique_constraints"
	FOREIGN_CONSTRAINTS = "foreign_constraints"

)

	type BackForgeError struct {
		Type    string
		Message any
		Layer   string // db - repo - service - handler

	}

	func New(typ string, message any, layer string) *BackForgeError {
			return &BackForgeError{
				Layer:   layer,
				Type:    typ,
				Message: message,
			}
	}

	
	`

}
