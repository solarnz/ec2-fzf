package main

import (
	"fmt"
	"os"

	ec2fzf "github.com/solarnz/ec2-fzf"
)

func main() {
	fzf, err := ec2fzf.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fzf.Run()
}
