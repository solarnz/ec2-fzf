package ec2fzf

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/solarnz/fzf/src"
)

type Ec2fzf struct {
	ec2      *ec2.EC2
	fzfInput *bytes.Buffer
	options  Options
	template *template.Template
}

func New() (*Ec2fzf, error) {
	options, err := ParseOptions()
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(options.Region),
		},
	})
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("Instance").Parse(options.Template)

	return &Ec2fzf{
		ec2:      ec2.New(sess),
		fzfInput: new(bytes.Buffer),
		options:  options,
		template: tmpl,
	}, nil
}

func (e *Ec2fzf) Run() {
	instances, err := e.ListInstances()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, i := range instances {
		str, err := e.StringFromInstance(i)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		e.fzfInput.WriteString(str)
		e.fzfInput.WriteString("\n")
	}

	fzfOptions := fzf.DefaultOptions()
	fzf.PostProcessOptions(fzfOptions)
	fzfOptions.Header = []string{
		"AWS EC2 Instances",
	}
	fzfOptions.Multi = false
	fzfOptions.Input = e.fzfInput
	fzfOptions.Printer = func(str string) {
		i, err := InstanceIdFromString(str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		address, err := e.GetConnectionDetails(i)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf(address)
	}
	fzf.Run(fzfOptions)
}
