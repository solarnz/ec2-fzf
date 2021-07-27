package ec2fzf

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Options struct {
	Version         bool
	GetPrivateIp    bool
	Regions         []string
	Template        string
	PreviewTemplate string
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

	pflag.StringSliceP("regions", "r", []string{"eu-central-1"}, "The AWS region")
	pflag.BoolP("get-private-ip", "i", true, "Return the private ip of the instance selected")
	pflag.StringSlice("filters", []string{}, "Filters to apply with the ec2 api call")
	pflag.BoolP("version", "v", false, "Show version and exit")
	pflag.Parse()
	// pflag.PrintDefaults()

	viper.BindPFlags(pflag.CommandLine)

	// viper.RegisterAlias("UsePrivateIp", "use-private-ip")
	// viper.RegisterAlias("regions", "region")

	viper.SetDefault("Region", "eu-central-1")
	viper.SetDefault("GetPrivateIp", false)
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

	return Options{
		Regions:         viper.GetStringSlice("Regions"),
		GetPrivateIp:    viper.GetBool("GetPrivateIp"),
		Template:        viper.GetString("Template"),
		PreviewTemplate: viper.GetString("PreviewTemplate"),
		Filters:         viper.GetStringSlice("Filters"),
		Version:         viper.GetBool("Version"),
	}
}
