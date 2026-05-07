package cli

import (
	"fmt"

	"github.com/eslamward/backforge/internal/generator"
	"github.com/eslamward/backforge/internal/parser"
	"github.com/eslamward/backforge/internal/toolchain"
)

const PREFEX = "************"

func Build() {

	//Reading and parsing the schema
	fmt.Println(PREFEX, "Parsing schema...", PREFEX)
	schema, err := parser.Parse("app.yaml")
	if err != nil {
		fmt.Println("error in parsing the schema please follow the roles :", err)
	}

	//Generate the project structure and all the files
	generator.GenerateProject(schema)

	//ceate mod file enshure go and download it and get some packege
	err = toolchain.RunGoBuild("output", generator.GetDriverGoGet(schema.Configuration.DatabaseConfig.Type))
	if err != nil {
		fmt.Println("Build error:", err)
		return
	}

	fmt.Println(PREFEX, "App built successfully", PREFEX)

}
