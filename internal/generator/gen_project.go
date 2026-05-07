package generator

import (
	"fmt"

	"github.com/eslamward/backforge/internal/filesystem"
	"github.com/eslamward/backforge/internal/parser"
)

func GenerateProject(schema *parser.Schema) {
	internal := "internal/"
	modelsName := []string{}
	for _, model := range schema.Models {
		/**/
		modelsName = append(modelsName, model.Name)

		/**/

		contentModel := GenerateModel(model)
		err := filesystem.WriteProjectFile(internal, "models/"+model.Name+"_model"+".go", contentModel)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Printf("%s model created successfully\n", model.Name)

		/**/
		contentHandler := GenerateHandler(model)
		err = filesystem.WriteProjectFile(internal, "handler/"+model.Name+"_handler"+".go", contentHandler)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Printf("%s handler created successfully\n", model.Name)

		/**/
		contentRepo := GenertateRepo(model, &schema.Configuration.DatabaseConfig)
		err = filesystem.WriteProjectFile(internal, "repository/"+model.Name+"_repository"+".go", contentRepo)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Printf("%s repository created successfully\n", model.Name)

		/**/
		contentServices := GenerateServices(model)
		err = filesystem.WriteProjectFile(internal, "services/"+model.Name+"_services"+".go", contentServices)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Printf("%s services created successfully\n", model.Name)

		/**/
		contentRoutes := GenerateRoutes(model)
		err = filesystem.WriteProjectFile(internal, "routes/"+model.Name+"_routes"+".go", contentRoutes)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Printf("%s Routes created successfully\n", model.Name)

	}

	/**/
	contentDB := InitDB(&schema.Configuration.DatabaseConfig, modelsName...)
	err := filesystem.WriteProjectFile(internal, "database/database.go", contentDB)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Database Init successfully\n")

	/**/
	contentDBTables := GenerateCreateTable(schema, &schema.Configuration.DatabaseConfig)
	err = filesystem.WriteProjectFile(internal, "database/database_tables.go", contentDBTables)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Database tables created successfully\n")

	/**/
	contentConatiner := GenerateConatiner(modelsName...)
	err = filesystem.WriteProjectFile(internal, "app/"+"container"+".go", contentConatiner)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("Container created successfully")

	/**/
	contentInjectRoutes := InjectRoutes(modelsName...)
	err = filesystem.WriteProjectFile(internal, "routes/"+"routes"+".go", contentInjectRoutes)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Routes Injected successfully\n")

	/**/
	contentValidator := GenerateValidator()
	err = filesystem.WriteProjectFile(internal, "validate/validator.go", contentValidator)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Validator created successfully\n")

	/**/
	contenthandHelper := GenerateHandlerHelper()
	err = filesystem.WriteProjectFile(internal, "handler/helper_handler.go", contenthandHelper)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Handler Helper created successfully\n")

	/**/
	contentHealth := GenerateHealthCheck()
	err = filesystem.WriteProjectFile(internal, "handler/health_check.go", contentHealth)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Handler Health Check created successfully\n")

	/**/
	contentError := GenerateError()
	err = filesystem.WriteProjectFile(internal, "backerror/back_error.go", contentError)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Error package created successfully\n")

	/**/
	contentMain := generateMain(&schema.Configuration.ServerConfig)
	err = filesystem.WriteProjectFile("cmd", "main.go", contentMain)
	if err != nil {
		fmt.Println("Error", err)
	}

	/***/
	err = filesystem.CreatFolder("bin")
	if err != nil {
		fmt.Println("Error", err)
	}
}
