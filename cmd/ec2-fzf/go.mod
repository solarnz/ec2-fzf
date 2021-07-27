module main

go 1.16

require (
	github.com/richie-tt/ec2-fzf v1.0.0
	// github.com/richie-tt/ec2-fzf/pkg/version v0.0.0

)

replace (
	github.com/richie-tt/ec2-fzf => ../../
	// github.com/richie-tt/ec2-fzf/pkg/version => ../../pkg/version
)