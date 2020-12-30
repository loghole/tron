package generate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/loghole/tron/cmd/tron/internal/helpers"
	"github.com/loghole/tron/cmd/tron/internal/models"
	"github.com/loghole/tron/cmd/tron/internal/stdout"
	"github.com/loghole/tron/cmd/tron/internal/templates"
)

func Implement(project *models.Project, printer stdout.Printer) error {
	printer.VerbosePrintln(color.FgMagenta, "Generate implementations")

	for _, proto := range project.Protos {
		if !proto.WithImpl() {
			continue
		}

		printer.VerbosePrintf(color.Reset, "\tgenerate from file '%s': ", color.YellowString(proto.RelPath))

		if err := generateHandler(project, proto); err != nil {
			printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

			return err
		}

		if err := generateMethods(project, proto); err != nil {
			printer.VerbosePrintln(color.FgRed, "FAIL: %v", err)

			return err
		}

		printer.VerbosePrintln(color.FgGreen, "OK")
	}

	printer.VerbosePrintln(color.FgBlue, "\tSuccess")

	return nil
}

func generateHandler(project *models.Project, proto *models.Proto) error {
	const handlerName = "handler.go"

	path := filepath.Join(project.AbsPath, models.ProjectPathImplementation,
		strings.ReplaceAll(proto.Package, ".", string(filepath.Separator)),
		handlerName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		handlerData := templates.NewHandlerData(proto.GoPackage, proto.Service.Name)
		handlerData.Desc = helpers.CamelCase(helpers.GoName(proto.Package))

		handlerData.AddImport(handlerData.Desc, strings.Join([]string{
			project.Module,
			models.ProjectPathPkgClients,
			strings.ReplaceAll(proto.Package, ".", "/"),
		}, "/"))

		template, err := helpers.ExecTemplate(templates.HandlerTemplate, handlerData)
		if err != nil {
			return fmt.Errorf("exec template: %w", err)
		}

		if err := helpers.WriteGoFile(path, template); err != nil {
			return fmt.Errorf("write file '%s': %w", path, err)
		}
	}

	return nil
}

func generateMethods(project *models.Project, proto *models.Proto) error {
	for _, method := range proto.Service.Methods {
		if err := generateMethod(project, proto, method); err != nil {
			return err
		}
	}

	return nil
}

func generateMethod(project *models.Project, proto *models.Proto, method *models.Method) error {
	path := filepath.Join(project.AbsPath, models.ProjectPathImplementation,
		strings.ReplaceAll(proto.Package, ".", string(filepath.Separator)),
		fmt.Sprintf("%s.go", helpers.SnakeCase(method.Name)))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		methodData := templates.NewMethodData()
		methodData.GoPackage = proto.GoPackage
		methodData.Name = method.Name

		methodData.AddInput(parseMethodType(project.Module, method.Input))
		methodData.AddOutput(parseMethodType(project.Module, method.Output))

		template, err := helpers.ExecTemplate(templates.MethodTemplate, methodData)
		if err != nil {
			return fmt.Errorf("exec template: %w", err)
		}

		if err := helpers.WriteGoFile(path, template); err != nil {
			return fmt.Errorf("write file '%s': %w", path, err)
		}
	}

	return nil
}

func parseMethodType(module, value string) (alias, importPath, typeName string) {
	parts := strings.Split(value, ".")
	if len(parts) == 0 {
		panic("can't be 0")
	}

	buf := bytes.NewBufferString(module)
	buf.WriteString("/")
	buf.WriteString(models.ProjectPathPkgClients)

	for _, part := range parts[:len(parts)-1] {
		if part == "" {
			continue
		}

		buf.WriteString("/")
		buf.WriteString(part)
	}

	typeName = parts[len(parts)-1]
	importPath = buf.String()
	alias = importAlias(importPath)

	return alias, importPath, typeName
}

func importAlias(s string) string {
	parts := strings.Split(s, "/")

	if len(parts) <= 1 {
		return s
	}

	return helpers.CamelCase(helpers.GoName(strings.Join(parts[len(parts)-2:], ".")))
}
