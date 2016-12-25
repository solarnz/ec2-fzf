package ec2fzf

import (
	"flag"
)

type Options struct {
	Region       string
	UsePrivateIp bool
}

func ParseOptions() Options {
	region := flag.String("region", "us-east-1", "The AWS region")
	usePrivateIp := flag.Bool("private", false, "return the private IP address of the instance rather than the public dns")

	flag.Parse()

	return Options{
		Region:       (*region),
		UsePrivateIp: *usePrivateIp,
	}
}
