package implement

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/loghole/tron/cmd/protoc-gen-tron/internal/helpers"
)

const handlerName = "handler.go"

func Generate(f *protogen.File, implPath, pkgPath, module string) error {
	if len(f.Services) == 0 {
		return nil
	}

	service := f.Services[0]

	if len(service.Methods) == 0 {
		return nil
	}

	if err := generateImplementation(f, service, implPath, pkgPath, module); err != nil {
		return fmt.Errorf("can't generate implementation: %w", err)
	}

	if err := generateMethods(f, service, implPath, pkgPath, module); err != nil {
		return fmt.Errorf("can't generate method: %w", err)
	}

	return nil
}

func generateImplementation(f *protogen.File, s *protogen.Service, implPath, pkgPath, module string) error {
	path := filepath.Join(
		implPath,
		strings.ReplaceAll(*f.Proto.Package, ".", string(filepath.Separator)),
		handlerName,
	)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		handlerData := NewHandlerData(string(f.GoPackageName), s.GoName)
		handlerData.Desc = helpers.CamelCase(helpers.GoName(*f.Proto.Package))

		desc := strings.Join([]string{module, pkgPath, strings.ReplaceAll(*f.Proto.Package, ".", "/")}, "/")
		handlerData.AddImport(handlerData.Desc, desc)

		template, err := helpers.ExecTemplate(HandlerTemplate, handlerData)
		if err != nil {
			return fmt.Errorf("can't exec template: %w", err)
		}

		result, err := helpers.FormatGoFile(template)
		if err != nil {
			return fmt.Errorf("format failed: %s, err: %w", template, err)
		}

		if err := helpers.WriteToFile(path, result); err != nil {
			return fmt.Errorf("can't write file: %w", err)
		}
	}

	return nil
}

func generateMethods(f *protogen.File, s *protogen.Service, implPath, pkgPath, module string) error {
	for _, method := range s.Methods {
		path := filepath.Join(
			implPath,
			strings.ReplaceAll(*f.Proto.Package, ".", string(filepath.Separator)),
			fmt.Sprintf("%s.go", helpers.SnakeCase(method.GoName)),
		)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			methodData := NewMethodData()
			methodData.GoPackage = string(f.GoPackageName)
			methodData.Name = method.GoName

			var (
				input  = getImportPath(f, method.Input, pkgPath, module)
				output = getImportPath(f, method.Output, pkgPath, module)
			)

			methodData.AddInput(importAlias(input), input, method.Input.GoIdent.GoName)
			methodData.AddOutput(importAlias(output), output, method.Output.GoIdent.GoName)

			template, err := helpers.ExecTemplate(MethodTemplate, methodData)
			if err != nil {
				return fmt.Errorf("can't exec template: %w", err)
			}

			result, err := helpers.FormatGoFile(template)
			if err != nil {
				return fmt.Errorf("format failed: %s, err: %w", template, err)
			}

			if err := helpers.WriteToFile(path, result); err != nil {
				return fmt.Errorf("can't write file: %w", err)
			}
		}
	}

	return nil
}

func getImportPath(f *protogen.File, message *protogen.Message, pkgPath, module string) string {
	if output := strings.Trim(message.GoIdent.GoImportPath.String(), `"`); output != `.` {
		return output
	}

	return strings.Join([]string{module, pkgPath, strings.ReplaceAll(*f.Proto.Package, ".", "/")}, "/")
}

func importAlias(s string) string {
	parts := strings.Split(s, "/")

	if len(parts) <= 1 {
		return s
	}

	return helpers.CamelCase(helpers.GoName(strings.Join(parts[len(parts)-2:], ".")))
}
