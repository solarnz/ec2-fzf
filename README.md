# ec2-fzf

ec2-fzf is a tool that utilised the [fzf](https://github.com/junegunn/fzf)
fuzzy matcher in order to retrieve the public or private address of an ec2
instance.

![GIF](https://raw.githubusercontent.com/solarnz/ec2-fzf/master/img/ec2-fzf.gif)

## Installation

```
go get github.com/solarnz/ec2-fzf/cmd/ec2-fzf
```

## Usage

You can pass `--private` to `ec2-fzf`, and it will return the private ip address
of the instance, rather than the public dns record. This is useful for
instances within a VPC.

You can also set `--region` and pass the ec2 region you would like to list
instances in.

You can use `ec2-fzf` with ssh with `ssh $(ec2-fzf --region ap-southeast-2)`

You can add filters on to the instances that are listed with the `--filters`
flag. You can define `--filters` multiple times to filter by more than one
value. Valid values are those used in the [aws-sdk-go
sdk](http://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeInstancesInput)
```
ec2-fzf --filters tag:realm=test --filters tag:Name=Test
```

## Configuration

You can set the default configuration options in `~/.config/ec2-fzf`, example
```
Region = "us-east-1"
Template = "{{index .Tags \"Name\"}}"
```
