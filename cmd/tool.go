package main

import (
	"fmt"
	"github.com/sunzhenyucn/lets-block-chain"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
)

const (
	builtInDBPath = "./db/lbc.db"
)

var (
	app         *cli.App
	initCommand = cli.Command{
		Name:   "init",
		Usage:  "Init LBC chain",
		Action: initCommandAction,
	}
	addCommand = cli.Command{
		Name:  "add",
		Usage: "Add block to block-chain",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "data, d",
				Usage: "Spec block data",
			},
		},
		Action: addCommandAction,
	}
	printCommand = cli.Command{
		Name:   "print",
		Usage:  "Print all blocks in block-chain",
		Action: printCommandAction,
	}
)

// init cli app meta data
func init() {
	app = cli.NewApp()
	app.Name = "LBC Chain Command Line Tool"
	app.Usage = "Used for access LBC Chain"
	app.Version = "0.1"
	app.UsageText = "tool [init|add|print|help]"
	app.Authors = []cli.Author{
		{
			Name:  "Gavin Sun",
			Email: "sunzhenyucn@gmail.com",
		},
	}
	app.Commands = []cli.Command{
		initCommand,
		addCommand,
		printCommand,
	}
}

// initCommandAction will use user spec db path to
// create block-chain
func initCommandAction(ctx *cli.Context) {
	createDirIfNotExist()
	if _, err := os.Stat(builtInDBPath); os.IsNotExist(err) != true {
		panic("[ERROR] Already initialized!")
	}
	lbc.NewBlockChain(builtInDBPath)
	fmt.Printf("[INFO] Successful initialized block chain on `%s`\n", builtInDBPath)
}

// addCommandAction will use spec data
// to create block on user spec db path's block-chain
func addCommandAction(ctx *cli.Context) {
	isInitialized()
	if ctx.String("data") == "" {
		panic("[ERROR] Block data cannot be empty!\n")
	}
	_chain := lbc.NewBlockChain(builtInDBPath)
	_chain.AddBlock(ctx.String("data"))
	fmt.Println("[INFO] Successful add block to block-chain!")
}

// printCommandAction will print all blocks in block-chain
func printCommandAction(ctx *cli.Context) {
	isInitialized()
	_chain := lbc.NewBlockChain(builtInDBPath)
	_iterator := _chain.Iterator()
	for {
		block := _iterator.Next()

		pow := lbc.NewProofWork(block)
		fmt.Printf("\nPrev. Hash: %x\nData: %s\nHash: %x\nPoW State: %s\n",
			block.Prev,
			block.Data,
			block.Hash,
			strconv.FormatBool(pow.Validate()))

		if len(block.Prev) == 0 {
			break
		}
	}
}

func isInitialized() {
	if _, err := os.Stat(builtInDBPath); os.IsNotExist(err) != false {
		panic("[ERROR] Please initialization first!")
	}
}

func createDirIfNotExist() {
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		err = os.MkdirAll("./db", 0755)
		if err != nil {
			panic(err)
		}
	}
}

// main function, cli app entry-point
func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
