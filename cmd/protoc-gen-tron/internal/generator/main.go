package generator

import (
	"os"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
)

// nolint:funlen // generation can be big
func (gen *Generator) generateMain(p *protogen.Plugin) {
	mainPath := filepath.Join("cmd", gen.moduleName, "main.go")

	if st, err := os.Stat(mainPath); err == nil && !st.IsDir() {
		return
	}

	g := p.NewGeneratedFile(mainPath, "main")

	g.P("package main")
	g.P()

	g.P("import (")
	g.P(`"log"`)
	g.P()
	g.P(`"` + gen.module + `/config"`)
	g.P()
	g.P(`"github.com/loghole/tron"`)
	g.P(")")
	g.P()

	g.P("func main() {")
	g.P("app, err := tron.New(tron.AddLogCaller())")
	g.P("if err != nil {")
	g.P(`log.Fatalf("can't create app: %s", err)`)
	g.P("}")
	g.P()

	g.P("defer app.Close()")
	g.P()

	g.P("app.Logger().Info(config.GetExampleValue())")
	g.P()

	g.P("// Init handlers")

	g.P("var (")
	g.P(")")
	g.P()

	g.P("if err := app.WithRunOptions().Run(")
	g.P("); err != nil {")
	g.P(`app.Logger().Fatalf("can't run app: %v", err)`)
	g.P("}")

	g.P("}")
	g.P()
}
