package ec2fzf

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Options struct {
	Regions         []string
	UsePrivateIp    bool
	Template        string
	PreviewTemplate string
	DetailTemplate  string
	Filters         []string
}

func ParseOptions() Options {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.config/ec2-fzf")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(err)
		}
	}

	pflag.StringSlice("region", []string{"us-east-1"}, "The AWS region")
	pflag.Bool("use-private-ip", true, "Return the private ip of the instance selected")
	pflag.StringSlice("filters", []string{}, "Filters to apply with the ec2 api call")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.RegisterAlias("UsePrivateIp", "use-private-ip")
	viper.RegisterAlias("regions", "region")

	viper.SetDefault("Region", "us-east-1")
	viper.SetDefault("UsePrivateIp", false)
	viper.SetDefault("Template", `{{ .InstanceId }}: {{index .Tags "Name"}}`)
	viper.SetDefault("PreviewTemplate", `
			Instance Id: {{.InstanceId}}
			Name:        {{index .Tags "Name"}}
			Private IP:  {{.PrivateIpAddress}}
			Public IP:   {{.PublicIpAddress}}

			Tags:
			{{ range $key, $value := .Tags }}
				{{ indent 2 $key }}: {{ $value }}
			{{- end -}}
		`,
	)
	viper.SetDefault("DetailTemplate", `{{if .UsePrivateIp}}{{.PrivateIpAddress}}{{else}}{{.PublicDnsName}}{{end}}`)

	return Options{
		Regions:         viper.GetStringSlice("Regions"),
		UsePrivateIp:    viper.GetBool("UsePrivateIp"),
		Template:        viper.GetString("Template"),
		PreviewTemplate: viper.GetString("PreviewTemplate"),
		DetailTemplate:  viper.GetString("DetailTemplate"),
		Filters:         viper.GetStringSlice("Filters"),
	}
}
