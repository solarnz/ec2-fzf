package ec2fzf

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	finder "github.com/ktr0731/go-fuzzyfinder"
)

type Ec2fzf struct {
	fzfInput        *bytes.Buffer
	options         Options
	listTemplate    *template.Template
	previewTemplate *template.Template
	detailTemplate  *template.Template
	ec2Sessions     []*session.Session
}

func New() (*Ec2fzf, error) {
	options := ParseOptions()

	sessions := make([]*session.Session, 0)
	for _, region := range options.Regions {
		sess, err := session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(region),
			},
		})
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, sess)
	}

	tmpl, err := template.New("Instance").Funcs(sprig.TxtFuncMap()).Parse(options.Template)
	if err != nil {
		panic(err)
	}

	previewTemplate, err := template.New("Preview").Funcs(sprig.TxtFuncMap()).Parse(options.PreviewTemplate)
	if err != nil {
		panic(err)
	}

	detailTemplate, err := template.New("Details").Funcs(sprig.TxtFuncMap()).Parse(options.DetailTemplate)
	if err != nil {
		panic(err)
	}

	return &Ec2fzf{
		fzfInput:        new(bytes.Buffer),
		options:         options,
		listTemplate:    tmpl,
		detailTemplate:  detailTemplate,
		previewTemplate: previewTemplate,
		ec2Sessions:     sessions,
	}, nil
}

func (e *Ec2fzf) Run() {
	instances := make([]*ec2.Instance, 0)
	instancesLock := &sync.Mutex{}

	wg := &sync.WaitGroup{}
	for _, sess := range e.ec2Sessions {
		wg.Add(1)
		go func(s *session.Session) {
			retrivedInstances, err := e.ListInstances(ec2.New(s))
			if err != nil {
				panic(err)
			}

			instancesLock.Lock()
			instances = append(instances, retrivedInstances...)
			instancesLock.Unlock()
			wg.Done()
		}(sess)
	}

	wg.Wait()

	indexes, err := finder.FindMulti(
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
		if errors.Is(err, finder.ErrAbort) {
			os.Exit(1)
		}
		panic(err)
	}

	for _, idx := range indexes {
		details, err := e.GetConnectionDetails(instances[idx])
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", details)
	}
}
