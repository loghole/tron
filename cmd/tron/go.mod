module github.com/loghole/tron/cmd/tron

go 1.15

require (
	github.com/Masterminds/semver v1.5.0
	github.com/fatih/color v1.10.0
	github.com/json-iterator/go v1.1.10
	github.com/lissteron/simplerr v0.8.0
	github.com/loghole/tron v0.15.2
	github.com/manifoldco/promptui v0.8.0
	github.com/spf13/cobra v1.1.1
	golang.org/x/tools v0.1.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace github.com/loghole/tron => ../..
