package ec2fzf

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/solarnz/fzf/src"
)

type Ec2fzf struct {
	ec2      *ec2.EC2
	fzfInput *bytes.Buffer
	options  Options
}

func New() (*Ec2fzf, error) {
	options := ParseOptions()

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(options.Region),
		},
	})
	if err != nil {
		return nil, err
	}

	return &Ec2fzf{
		ec2:      ec2.New(sess),
		fzfInput: new(bytes.Buffer),
		options:  options,
	}, nil
}

func (e *Ec2fzf) Run() {
	instances, err := e.ListInstances()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, i := range instances {
		e.fzfInput.WriteString(StringFromInstance(i))
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
