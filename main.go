package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/solarnz/fzf/src"
)

type Ec2fzf struct {
	ec2      *ec2.EC2
	fzfInput *bytes.Buffer
}

type Tags []*ec2.Tag

func (s Tags) Len() int {
	return len(s)
}

func (s Tags) Less(i, j int) bool {
	return *(s[i].Key) < *(s[j].Key)
}

func (s Tags) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var region string
var usePrivateIp bool

func main() {
	flag.StringVar(&region, "region", "us-east-1", "The AWS region")
	flag.BoolVar(&usePrivateIp, "private", false, "return the private IP address of the instance rather than the public dns")
	flag.Parse()

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ec2fzf := Ec2fzf{
		ec2:      ec2.New(sess),
		fzfInput: new(bytes.Buffer),
	}

	instances, err := ec2fzf.listInstances()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, i := range instances {
		ec2fzf.fzfInput.WriteString(StringFromInstance(i))
		ec2fzf.fzfInput.WriteString("\n")
	}

	options := fzf.DefaultOptions()
	fzf.PostProcessOptions(options)
	options.Header = []string{
		"AWS EC2 Instances",
	}
	options.Multi = false
	options.Input = ec2fzf.fzfInput
	options.Printer = func(str string) {
		i, err := InstanceIdFromString(str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		address, err := ec2fzf.GetConnectionDetails(i)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf(address)
	}
	fzf.Run(options)
}

func (e *Ec2fzf) listInstances() ([]*ec2.Instance, error) {
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

	if usePrivateIp {
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
