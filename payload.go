package ec2fzf

import (
	"bytes"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const timeContext = 5

// Ec2fzf default struct used for communication
type Ec2fzf struct {
	fzfInput        *bytes.Buffer
	options         Options
	listTemplate    *template.Template
	previewTemplate *template.Template
	EC2Resources    []*EC2Resource
}

// Options struct allow to transport the options to other part of program
type Options struct {
	Version         bool
	GetPrivateIP    bool
	Regions         []string
	Template        string
	PreviewTemplate string
	Filters         []string
}

// EC2Resource struct contain session created from RegionID & AWS Credentials
type EC2Resource struct {
	Region    Region
	Client    session.Session
	Instances []*ec2.Instance
}

// Region struct
type Region struct {
	Name      string
	Available bool
}

func (r *EC2Resource) DeepCopy() *EC2Resource {

	// cc := resource.Client.Copy()
	// resource.Client = nil

	// copyByte, err := json.Marshal(resource)
	// if err != nil {
	// 	fmt.Println("error to ", err.Error())
	// 	panic(err)
	// 	// return EC2Resource{}, nil
	// }

	// var copyResource EC2Resource

	// err = json.Unmarshal(copyByte, &copyResource)
	// if err != nil {
	// 	panic(err)
	// }

	// // copyResource.Client = cc

	// return copyResource, err

	newResource := &EC2Resource{
		Region: r.Region.Copy(),
		Client: *r.Client.Copy(),
		// Instances: r.I,
	}

	// result := copy(newResource.Instances, r.Instances)
	copy(newResource.Instances, r.Instances)

	// fmt.Println("result copy,", result)

	return newResource
}

func (r *Region) Copy() Region {
	newRegion := Region{
		Name:      r.Name,
		Available: r.Available,
	}

	return newRegion
}

// func getEC2ResourceIndex() int {
// 	for key, val := range
// }
