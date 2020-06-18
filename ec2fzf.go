package ec2fzf

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	finder "github.com/ktr0731/go-fuzzyfinder"
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

	idx, err := finder.Find(
		instances,
		func(i int) string {
			str, _ := e.StringFromInstance(instances[i])
			return fmt.Sprintf("%s\n", str)
		},
	)

	if err != nil {
		panic(err)
	}

	details := e.GetConnectionDetails(instances[idx])
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", details)
}
