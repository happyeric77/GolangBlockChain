package main
import (
	"project4/cli"
	"os"
)

func main(){
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()	
}