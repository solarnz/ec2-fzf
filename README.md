**This is fork of [solarnz/ec2-fzf](https://github.com/solarnz/ec2-fzf) repository, who inspired me to add some features.**
<br>
<br>


# IN DEVELOP

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

- Arguments
  - `-h,--help` show help
  - `-v,--version` Show version
  - `-r,--regions` List of regions split by comma. Each AWS region from the list will be scanned result of EC2 instances will be merged in one final result.  
  - `-i,--get-private-ip` Return private IP address from selectec instance.
  - `-f,--filters` Use filters to limit results, Valid values are those used in the [aws-sdk-go
sdk](http://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeInstancesInput)

- Keyboard shortcuts
  - `<ctrl-f>` - Refresh list
  - `<ctrl-l>` - Show instances list only from a specific region

- Command line
  ```sh
  ssh $(ec2-fzf --regions ap-southeast-2)
  ```

- Bash function
  ```sh
  function sshe(){
    local ip=$(ec2-fzf --regions ap-southeast-2 --get-private-ip)
    ssh root@$ip
  }
  ```

## Configuration

You can set the default configuration options in `~/.config/ec2-fzf/config.toml`, example
```
Region = "us-east-1"
Template = "{{index .Tags \"Name\"}}"
```
