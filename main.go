package main

import (
	"fmt"
	"log"
	"os"

	"github.com/t1anchen/gogmlib/sm3"
	"github.com/urfave/cli"
)

// VERSION 主程序的版本
const VERSION = "1.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "gogmlib"
	app.Version = VERSION
	app.Commands = []cli.Command{
		{
			Name: "sm2",
			Subcommands: []cli.Command{
				{
					Name:  "encrypt",
					Usage: "使用 sm2 加密",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name: "sm3",
			Subcommands: []cli.Command{
				{
					Name:  "compute",
					Usage: "使用 sm3 哈希值计算",
					Action: func(c *cli.Context) error {
						for _, v := range c.Args() {
							hashStr := sm3.NewContext().ComputeFromString(v).ToHexString()
							fmt.Printf("%s* %s\n", hashStr, v)
						}
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
