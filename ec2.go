package ec2fzf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (e *Ec2fzf) ListInstances() ([]*ec2.Instance, error) {
	instances := make([]*ec2.Instance, 0, 0)
	filters := make([]*ec2.Filter, 0, 0)
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
	params := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	err := e.ec2.DescribeInstancesPages(
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

func (e *Ec2fzf) GetConnectionDetails(instanceId string) (string, error) {
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceId)},
	}
	resp, err := e.ec2.DescribeInstances(params)
	if err != nil {
		return "", err
	}

	if !(len(resp.Reservations) == 1) || !(len(resp.Reservations[0].Instances) == 1) {
		return "", fmt.Errorf("No instance could be found for %s", instanceId)
	}

	if e.options.UsePrivateIp {
		return *resp.Reservations[0].Instances[0].PrivateIpAddress, nil
	}
	return *resp.Reservations[0].Instances[0].PublicDnsName, nil
}

func (e *Ec2fzf) StringFromInstance(i *ec2.Instance) (string, error) {
	tags := make(map[string]string)

	for _, t := range i.Tags {
		tags[*t.Key] = *t.Value
	}

	buffer := new(bytes.Buffer)
	err := e.template.Execute(
		buffer,
		struct {
			Tags map[string]string
		}{
			Tags: tags,
		},
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%19s: %s", *i.InstanceId, buffer.String()), nil
}

func InstanceIdFromString(s string) (string, error) {
	i := strings.Index(s, ":")

	if i < 0 {
		return "", fmt.Errorf("Unable to find instance id")
	}
	return strings.TrimSpace(s[0:i]), nil
}
