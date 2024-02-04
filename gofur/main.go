package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

type File struct {
	Path    string
	Content []byte
}

func writeFileWithCheck(file File) error {
	if err := os.WriteFile(file.Path, file.Content, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func main() {
	cmd := NewCommand()
	cmd.Register(
		runProject,
		installProject,
		generateProject,
		generateModel,
		generateView,
		generateHandler,
	)

	cmd.Execute()
}

func runProject() *cobra.Command {
	return &cobra.Command{
		Use:     "run",
		Example: "gofur run",
		Short:   "Run GoFur development server",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat("cmds/main.go"); err != nil {
				fmt.Println("not in GoFur app root: cmd/main.go not found")
				return
			}
			if err := exec.Command("templ", "generate").Run(); err != nil {
				fmt.Println(err)
				return
			}

			if err := exec.Command("go", "run", "cmds/main.go").Run(); err != nil {
				fmt.Println(err)
			}
		},
	}
}

func installProject() *cobra.Command {
	return &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Example: "gofur install",
		Short:   "Install project's dependency",
		Run: func(cmd *cobra.Command, args []string) {
			start := time.Now()
			fmt.Println("installing project...")
			if err := exec.Command("go", "get", "github.com/shtayeb/gofur@latest").Run(); err != nil {
				fmt.Println(err)
				return
			}

			if err := exec.Command("go", "get", "github.com/a-h/templ").Run(); err != nil {
				fmt.Println(err)
				return
			}
			if err := exec.Command("templ", "generate").Run(); err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("done installing project in %v\n", time.Since(start))
		},
	}
}

func generateProject() *cobra.Command {
	return &cobra.Command{
		Use:     "new",
		Example: "gofur new hello-world",
		Short:   "Create new GoFur project",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("invalid arguments")
				return
			}

			name := args[0]

			fmt.Println("creating new GoFur project:", name)
			if err := os.Mkdir(name, os.ModePerm); err != nil {
				fmt.Println(err)
				return
			}

			files := []File{
				// setup directory
				{Path: name + "/http", Content: nil},
				{Path: name + "/http/handlers", Content: nil},
				{Path: name + "/http/middlewares", Content: nil},
				{Path: name + "/http/session", Content: nil},
				{Path: name + "/http/types", Content: nil},

				{Path: name + "/internal", Content: nil},
				{Path: name + "/internal/models", Content: nil},
				{Path: name + "/internal/database", Content: nil},

				{Path: name + "/public", Content: nil},
				{Path: name + "/public/assets", Content: nil},
				{Path: name + "/public/css", Content: nil},
				{Path: name + "/public/js", Content: nil},

				{Path: name + "/sql", Content: nil},
				{Path: name + "/sql/queries", Content: nil},
				{Path: name + "/sql/schema", Content: nil},

				{Path: name + "/views", Content: nil},
				{Path: name + "/views/hello", Content: nil},
				{Path: name + "/views/layout", Content: nil},

				{Path: name + "/cmds", Content: nil},

				// setup files
				{Path: name + "/http/handlers/hello.go", Content: renderStub("./stubs/handler.go.stub", "hello.go", map[string]string{"mod": name})},
				{Path: name + "/http/middlewares/auth.go", Content: renderStub("./stubs/middlerwares/auth.go.stub", "auth.go", map[string]string{"mod": name})},
				{Path: name + "/http/types/types.go", Content: renderStub("./stubs/types/types.go.stub", "types.go", map[string]string{"mod": name})},

				{Path: name + "/public/css/app.css", Content: renderStub("./stubs/public/app_css.go.stub", "app.css", map[string]string{"mod": name})},
				{Path: name + "/public/js/app.js", Content: renderStub("./stubs/public/app_js.go.stub", "app.js", map[string]string{"mod": name})},

				{Path: name + "/cmds/main.go", Content: renderStub("./stubs/main.go.stub", "main.go", map[string]string{"mod": name})},

				{Path: name + "/views/layout/base.templ", Content: renderStub("./stubs/base_temp.go.stub", "base.templ", map[string]string{"mod": name})},
				{Path: name + "/views/hello/hello.templ", Content: renderStub("./stubs/base_temp.go.stub", "hello.templ", map[string]string{"mod": name})},

				{Path: name + "/internal/models/models.go", Content: renderStub("./stubs/internal/models.go.stub", "models.go", map[string]string{"mod": name})},
				{Path: name + "/internal/app.go", Content: renderStub("./stubs/internal/app.go.stub", "app.go", map[string]string{"mod": name})},

				{Path: name + "/sql/schema/001_users.sql", Content: renderStub("./stubs/sql/schema/users.go.stub", "users.sql", map[string]string{"mod": name})},
				{Path: name + "/sql/schema/002_sessions.sql", Content: renderStub("./stubs/sql/schema/sessions.go.stub", "sessions.sql", map[string]string{"mod": name})},
				{Path: name + "/sql/queries/sessions.sql", Content: renderStub("./stubs/sql/queries/sessions.go.stub", "sessions.sql", map[string]string{"mod": name})},

				{Path: name + "/go.mod", Content: renderStub("./stubs/go_mod.go.stub", "go.mod", map[string]string{"mod": name})},
				{Path: name + "/.air.toml", Content: renderStub("./stubs/airtoml.go.stub", "air.toml", map[string]string{"mod": name})},
				{Path: name + "/.env", Content: renderStub("./stubs/env.go.stub", "env.env", map[string]string{"mod": name})},
				{Path: name + "/.gitignore", Content: renderStub("./stubs/gitignore.go.stub", ".gitignore", map[string]string{"mod": name})},
			}

			errors := []error{}
			for _, file := range files {
				if file.Content == nil {
					if err := os.Mkdir(file.Path, os.ModePerm); err != nil {
						errors = append(errors, err)
					}
				} else {
					if err := writeFileWithCheck(file); err != nil {
						errors = append(errors, err)
					}
				}
			}

			if len(errors) != 0 {
				fmt.Println("GoFur encountered errors during file initialization:", errors)
				return
			}
		},
	}
}

type StubDetails struct {
	Name     string
	FileName string
	Values   map[string]string
}

func renderStub(name string, fileName string, values map[string]string) []byte {
	stub := StubDetails{
		Name:     name,
		FileName: fileName,
		Values:   values,
	}

	contentsBuff, err := os.ReadFile(stub.Name)
	if err != nil {
		log.Fatalf("RENDER STUB: Unable to read file: %s", stub.Name)
	}

	// Where to write the result
	// file, err := os.OpenFile(stub.Destination+stub.FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	// if err != nil {
	// 	log.Fatalf("Unable to open file: %s", stub.FileName)
	// }
	// defer file.Close()

	tem, err := template.New(stub.FileName).Parse(string(contentsBuff))
	if err != nil {
		log.Fatalf("RENDER STUB: Unable to parse template: %s", stub.Name)
	}

	var tem_buffer bytes.Buffer
	err = tem.Execute(&tem_buffer, stub.Values)
	if err != nil {
		log.Fatal("RENDER STUB: cant execute the template")
	}
	return tem_buffer.Bytes()
}

func generateModel() *cobra.Command {
	return &cobra.Command{
		Use:     "model",
		Example: "gofur model user",
		Short:   "Generate new model",
		Run: func(cmd *cobra.Command, args []string) {
			wdPath, _ := os.Getwd()
			// get the app name somehow
			file := File{
				// setup directory
				Path:    wdPath + "/models/" + args[0],
				Content: renderStub("./stubs/model.go.stub", fmt.Sprintf("%v.go", args[0]), map[string]string{"mod": args[0]}),
			}

			if err := os.Mkdir(file.Path, os.ModePerm); err != nil {
				log.Fatalln(err)
			}

			err := writeFileWithCheck(file)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
}

func generateView() *cobra.Command {
	return &cobra.Command{
		Use:     "view",
		Example: "GoFur view user",
		Short:   "Generate new view",
		Run: func(cmd *cobra.Command, args []string) {
			wdPath, _ := os.Getwd()
			file := File{
				// setup directory
				Path:    wdPath + "/views/" + args[0],
				Content: renderStub("./stubs/view.go.stub", fmt.Sprintf("%v.go", args[0]), map[string]string{"mod": args[0]}),
			}

			if err := os.Mkdir(file.Path, os.ModePerm); err != nil {
				log.Fatalln(err)
			}

			err := writeFileWithCheck(file)
			if err != nil {
				log.Fatalln(err)
			}

		},
	}
}

func generateHandler() *cobra.Command {
	return &cobra.Command{
		Use:     "handler",
		Example: "gofur handler home",
		Short:   "Generate new handler",
		Run: func(cmd *cobra.Command, args []string) {
			wdPath, _ := os.Getwd()
			file := File{
				// setup directory
				Path:    wdPath + "/handlers/" + args[0],
				Content: renderStub("./stubs/handler.go.stub", fmt.Sprintf("%v.go", args[0]), map[string]string{"mod": args[0]}),
			}

			if err := os.Mkdir(file.Path, os.ModePerm); err != nil {
				log.Fatalln(err)
			}

			err := writeFileWithCheck(file)
			if err != nil {
				log.Fatalln(err)
			}

		},
	}
}
