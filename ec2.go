package ec2fzf

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (e *Ec2fzf) ListInstances() ([]*ec2.Instance, error) {
	instances := make([]*ec2.Instance, 0, 0)
	params := &ec2.DescribeInstancesInput{}

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

func StringFromInstance(i *ec2.Instance) string {
	sortedTags := make(Tags, len(i.Tags))
	copy(sortedTags, i.Tags)
	sort.Sort(sortedTags)

	tagStrings := make([]string, 0, 0)
	for _, t := range sortedTags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", *t.Key, *t.Value))
	}
	return fmt.Sprintf("%s: Tags=(%s)", *i.InstanceId, strings.Join(tagStrings, " "))
}

func InstanceIdFromString(s string) (string, error) {
	i := strings.Index(s, ":")
	if i < 0 {
		return "", fmt.Errorf("Unable to find instance id")
	}
	return s[0:i], nil
}
