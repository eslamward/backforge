package generator

import (
	"fmt"
	"strings"

	"github.com/eslamward/backforge/internal/parser"
)

type Routes struct {
	Method   string
	Url      string
	Response string
}

type RouteData struct {
	Model string
	Urls  []Routes
}

func AllRoutes(urls []Routes, model string) RouteData {

	rd := RouteData{
		Model: model,
		Urls:  urls,
	}
	return rd
}

var routesData []RouteData

func InjectRoutes(args ...string) string {
	var sb strings.Builder

	sb.WriteString(`
	package routes
	import (
	"backforge/internal/handler"
	"backforge/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	)

	`)

	sb.WriteString("func RegisterRoutes(router *chi.Mux, container *app.Conatiner){\n\n")

	sb.WriteString(`
	    router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: false,
		}))
		`)

	sb.WriteString(`router.Get("/api/health",handler.HealthCheck)


	`)

	for _, i := range args {

		sb.WriteString(fmt.Sprintf("		%sRoutes(container.%shand, router)\n",
			strings.ToLower(toSingular(i)), toPascalCase(toSingular(i))))
	}

	sb.WriteString("}")
	return sb.String()
}

func GenerateRoutes(model parser.Model) string {

	var sb strings.Builder
	var urls []Routes

	nameCapital := toPascalCase(toSingular(model.Name))
	nameLower := strings.ToLower(toSingular(model.Name))

	sb.WriteString("package routes\n\n")

	sb.WriteString(`import (
	"backforge/internal/handler"
	"github.com/go-chi/chi/v5"

)
`)

	sb.WriteString(fmt.Sprintf("\nfunc %sRoutes(%sHandler handler.%sHandler, router *chi.Mux) {\n",
		nameLower, nameLower, nameCapital))

	sb.WriteString(fmt.Sprintf(
		`    router.Post("/api/%s", %sHandler.Create%s)
	`, nameLower, nameLower, nameCapital))

	urls = append(
		urls,
		Routes{
			Method:   "POST",
			Url:      fmt.Sprintf("/api/%s", nameLower),
			Response: fmt.Sprintf("create %s", nameLower),
		})

	sb.WriteString(fmt.Sprintf(
		`    router.Get("/api/%s", %sHandler.Get%s)
	`, nameLower, nameLower, nameCapital))

	urls = append(
		urls,
		Routes{
			Method:   "GET",
			Url:      fmt.Sprintf("/api/%s", nameLower),
			Response: fmt.Sprintf("retrive %s", nameLower),
		})

	if uniqeField(model) != nil {
		sb.WriteString(fmt.Sprintf(
			`    router.Get("/api/%s/%s", %sHandler.Get%sBy%s)
	`, nameLower, uniqeField(model).Name, nameLower, nameCapital, toPascalCase(uniqeField(model).Name)))

		urls = append(
			urls,
			Routes{
				Method:   "GET",
				Url:      fmt.Sprintf("/api/%s/%s", nameLower, uniqeField(model).Name),
				Response: fmt.Sprintf("retrive %s by %s", nameLower, uniqeField(model).Name),
			})

	}

	sb.WriteString(fmt.Sprintf(
		`    router.Get("/api/%s/%s", %sHandler.Get%ss)
	`, nameLower, model.Name, nameLower, nameCapital))
	urls = append(
		urls,
		Routes{
			Method:   "GET",
			Url:      fmt.Sprintf("/api/%s/%s", nameLower, model.Name),
			Response: fmt.Sprintf("retrive all %s", model.Name),
		})

	sb.WriteString(fmt.Sprintf(
		`    router.Put("/api/%s", %sHandler.Update%s)
	`, nameLower, nameLower, nameCapital))

	urls = append(
		urls,
		Routes{
			Method:   "PUT",
			Url:      fmt.Sprintf("/api/%s", nameLower),
			Response: fmt.Sprintf("update %s", nameLower),
		})

	sb.WriteString(fmt.Sprintf(
		`    router.Delete("/api/%s", %sHandler.Delete%s)
	`, nameLower, nameLower, nameCapital))
	urls = append(
		urls,
		Routes{
			Method:   "DELETE",
			Url:      fmt.Sprintf("/api/%s", nameLower),
			Response: fmt.Sprintf("delete %s", nameLower),
		})

	sb.WriteString("}\n")

	routesData = append(routesData, RouteData{Urls: urls, Model: toPascalCase(model.Name)})

	return sb.String()

}

func GenerateHealthCheck() string {
	var sb strings.Builder

	sb.WriteString("package handler\n\n")
	sb.WriteString(`
	import 
	(
	"net/http"
	"fmt"
	)
	`)

	sb.WriteString("var healthCheck = `")

	sb.WriteString(`
<div style="font-family: Arial, sans-serif; padding: 20px;">

  <h2 style="color: #2c3e50;"> BackForge Status</h2>

  <div style="background: #ecf0f1; padding: 12px; border-radius: 8px; margin-bottom: 20px;">
    <p><b>Status:</b> Available</p>
    <p><b>Version:</b> 1.0.0</p>
    <p><b>Env:</b> Development</p>
  </div>

  <h3 style="color: #34495e;"> Available REST APIs</h3>
`)

	for _, rs := range routesData {

		sb.WriteString(fmt.Sprintf(`
  <div style="margin-top: 25px;">
    <h4 style="color: #e74c3c;"> %s</h4>
  </div>
	`, rs.Model))

		for _, r := range rs.Urls {

			methodColor := "#7f8c8d"

			switch r.Method {
			case "GET":
				methodColor = "#27ae60"
			case "POST":
				methodColor = "#2980b9"
			case "PUT":
				methodColor = "#f39c12"
			case "DELETE":
				methodColor = "#c0392b"
			}

			sb.WriteString(fmt.Sprintf(`
  <div style="
    border: 1px solid #ddd;
    border-radius: 10px;
    padding: 12px;
    margin: 12px 0;
    background: #fafafa;
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  ">
    <p>
      <b>Method:</b> 
      <span style="color: white; background: %s; padding: 3px 8px; border-radius: 5px;">
        %s
      </span>
    </p>

    <p>
      <b>URL:</b> 
      <code id="url-%s" style="background:#ecf0f1; padding:3px 6px; border-radius:5px;">
        %s
      </code>

      <button onclick="copyUrl('url-%s')" 
        style="
          margin-left:10px;
          padding:4px 8px;
          border:none;
          background:#2ecc71;
          color:white;
          border-radius:5px;
          cursor:pointer;
        ">
        Copy
      </button>
    </p>

    <p><b>Target:</b> %s</p>
  </div>
		`, methodColor, r.Method, r.Url, r.Url, r.Url, r.Response))
		}
	}

	sb.WriteString(`
<script>
function copyUrl(id) {
  const text = document.getElementById(id).innerText;
  navigator.clipboard.writeText(text).then(() => {
    alert("Copied: " + text);
  });
}
</script>
`)
	sb.WriteString(`
</div>
`)
	sb.WriteString("`")

	/**/

	sb.WriteString(
		`
		func HealthCheck(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html")

		w.WriteHeader(200)

		`)
	sb.WriteString(`w.Write([]byte(fmt.Sprintf("<div>%s</div>",healthCheck)))
	}
	`)

	return sb.String()
}
