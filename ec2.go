package ec2fzf

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//test
//test
//test
//test
//test
//test

func (e *Ec2fzf) ListInstances(ctx context.Context, ec2Client *ec2.EC2) ([]*ec2.Instance, error) {
	instances := make([]*ec2.Instance, 0)
	filters := make([]*ec2.Filter, 0, 0)

	filters = append(filters, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []*string{aws.String("pending"), aws.String("running"), aws.String("shutting-down")},
	})

	for _, filter := range e.options.Filters {
		split := strings.SplitN(filter, "=", 2)
		if len(split) < 2 {
			return nil, fmt.Errorf("Filters can only contain one '='. Filter \"%s\" has %d", filter, len(split))
		}

		filters = append(filters, &ec2.Filter{
			Name:   aws.String(split[0]),
			Values: []*string{aws.String(split[1])},
		})
	}
	params := &ec2.DescribeInstancesInput{}

	if len(filters) > 0 {
		params.Filters = filters
	}

	contextWithTimeout, cancelTimeout := context.WithTimeout(ctx, timeContext*time.Second)
	defer cancelTimeout()

	err := ec2Client.DescribeInstancesPagesWithContext(
		contextWithTimeout,
		params,
		func(p *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, r := range p.Reservations {
				for _, i := range r.Instances {
					instances = append(instances, i)
				}
			}
			return !lastPage
		},
	)

	return instances, err
}

func (e *Ec2fzf) GetConnectionDetails(instance *ec2.Instance) string {
	if e.options.GetPrivateIP {
		return *instance.PrivateIpAddress
	}
	return *instance.PublicDnsName
}

func TemplateForInstance(i *ec2.Instance, t *template.Template) (output string, err error) {
	tags := make(map[string]string)

	for _, t := range i.Tags {
		tags[*t.Key] = *t.Value
	}

	buffer := new(bytes.Buffer)
	err = t.Execute(
		buffer,
		struct {
			Tags map[string]string
			*ec2.Instance
		}{
			tags,
			i,
		},
	)

	output = buffer.String()
	return
}

func InstanceIdFromString(s string) (string, error) {
	i := strings.Index(s, ":")

	if i < 0 {
		return "", fmt.Errorf("Unable to find instance id")
	}
	return strings.TrimSpace(s[0:i]), nil
}
