package ec2fzf

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// New init sessions and prepare templates
func New() (*Ec2fzf, error) {
	options := ParseOptions()

	if options.Version {
		showVersion()
		os.Exit(0)
	}

	ec2resource := make([]*EC2Resource, 0)

	for _, region := range options.Regions {

		sess, err := session.NewSessionWithOptions(session.Options{

			Config: aws.Config{
				Region: aws.String(region),
			},
		})
		if err != nil {
			return nil, err
		}

		r := EC2Resource{
			Region: Region{
				Name: region,
			},
			Client: *sess,
		}
		ec2resource = append(ec2resource, &r)

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
		fzfInput:        new(bytes.Buffer),
		options:         options,
		listTemplate:    tmpl,
		previewTemplate: previewTemplate,
		EC2Resources:    ec2resource,
	}, nil
}

func updateInstances(e *Ec2fzf, res EC2Resource) {
	for _, val := range e.EC2Resources {
		if val.Region.Name == res.Region.Name {
			val.Instances = res.Instances
			// fmt.Println("znalazłem", key, val.Region.Name)
		}
	}
}

func (e *Ec2fzf) getEc2List(ctx context.Context) {

	wg := &sync.WaitGroup{}

	chResult := make(chan EC2Resource, len(e.EC2Resources))

	for _, res := range e.EC2Resources {
		wg.Add(1)

		go func(c context.Context, r EC2Resource) {

			// ec2Client := ec2.New(&r.Client)
			retrivedInstances, err := e.ListInstances(c, ec2.New(&r.Client))
			if err != nil {
				message := fmt.Sprintf("region: %s is not avaiable. err: %s", r.Region.Name, err.Error())
				fmt.Println(message)

			} else {
				resource := r.DeepCopy()
				resource.Instances = retrivedInstances
				resource.Region.Available = true
				chResult <- *resource
			}

			wg.Done()
		}(ctx, *res)
	}

	wg.Wait()
	close(chResult)

	for val := range chResult {
		updateInstances(e, val)
	}

	for _, val := range e.EC2Resources {
		fmt.Println("valll", val.Region.Name, len(val.Instances))
	}

}

// func

// Run create list of ec2
func (e *Ec2fzf) Run() {

	ctx := context.Background()
	e.getEc2List(ctx)

	os.Exit(0)

	indexes, err := finder.FindMulti(
	// 	instances,
	// 	func(i int) string {
	// 		str, _ := TemplateForInstance(instances[i], e.listTemplate)
	// 		return fmt.Sprintf("%s\n", str)
	// 	},
	// 	finder.WithPreviewWindow(func(i, w, h int) string {
	// 		if i == -1 {
	// 			return ""
	// 		}

	// 		str, _ := TemplateForInstance(instances[i], e.previewTemplate)

	// 		return str
	// 	}),
	)

	// if err != nil {
	// 	if errors.Is(err, finder.ErrAbort) {
	// 		os.Exit(1)
	// 	}
	// 	panic(err)
	// }

	// for _, idx := range indexes {
	// 	details := e.GetConnectionDetails(instances[idx])
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Printf("%s\n", details)
	// }
}
