package ec2fzf

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	finder "github.com/ktr0731/go-fuzzyfinder"
)

type Ec2fzf struct {
	ec2             *ec2.EC2
	fzfInput        *bytes.Buffer
	options         Options
	listTemplate    *template.Template
	previewTemplate *template.Template
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

	tmpl, err := template.New("Instance").Funcs(sprig.TxtFuncMap()).Parse(options.Template)
	if err != nil {
		panic(err)
	}

	previewTemplate, err := template.New("Preview").Funcs(sprig.TxtFuncMap()).Parse(options.PreviewTemplate)
	if err != nil {
		panic(err)
	}

	return &Ec2fzf{
		ec2:             ec2.New(sess),
		fzfInput:        new(bytes.Buffer),
		options:         options,
		listTemplate:    tmpl,
		previewTemplate: previewTemplate,
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
			str, _ := TemplateForInstance(instances[i], e.listTemplate)
			return fmt.Sprintf("%s\n", str)
		},
		finder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			str, _ := TemplateForInstance(instances[i], e.previewTemplate)

			return str
		}),
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
