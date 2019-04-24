package main

import (
	"fmt"
	"github.com/sunzhenyucn/lets-block-chain"
	"github.com/urfave/cli"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	builtInDBPath = "./db/lbc.db"
)

var (
	app         *cli.App
	initCommand = cli.Command{
		Name:  "init",
		Usage: "Init LBC chain",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "address, a",
				Usage: "init address",
			},
		},
		Action: initCommandAction,
	}
	sendCommand = cli.Command{
		Name:  "send",
		Usage: "Send funds",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "from, f",
				Usage: "Sender name",
			},
			cli.StringFlag{
				Name:  "to, t",
				Usage: "Receiver name",
			},
			cli.IntFlag{
				Name:  "amount, n",
				Usage: "Funds amount",
			},
		},
		Action: sendCommandAction,
	}
	getBalanceCommand = cli.Command{
		Name:  "balance",
		Usage: "Get spec address's balance",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "address",
				Usage: "Address",
			},
		},
		Action: getBalanceCommandAction,
	}
	printCommand = cli.Command{
		Name:  "print",
		Usage: "Print all blocks in block-chain",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "address",
				Usage: "Init address",
			},
		},
		Action: printCommandAction,
	}
)

// init cli app meta data
func init() {
	app = cli.NewApp()
	app.Name = "LBC Chain Command Line Tool"
	app.Usage = "Used for access LBC Chain"
	app.Version = "0.1"
	app.UsageText = "tool [init|send|balance|print|help]"
	app.Authors = []cli.Author{
		{
			Name:  "Gavin Sun",
			Email: "sunzhenyucn@gmail.com",
		},
	}
	app.Commands = []cli.Command{
		initCommand,
		sendCommand,
		getBalanceCommand,
		printCommand,
	}
}

// initCommandAction will use user spec db path and
// address to create block-chain
func initCommandAction(ctx *cli.Context) {
	if ctx.String("a") == "" {
		panic("[ERROR] Init address cannot be empty!\n")
	}

	createDirIfNotExist()
	if _, err := os.Stat(builtInDBPath); os.IsNotExist(err) != true {
		panic("[ERROR] Already initialized!")
	}
	_bc := lbc.NewBlockChain(ctx.String("a"))
	defer _bc.Close()
	fmt.Printf("[INFO] Successful initialized block chain on `%s` and init address is `%s`\n", builtInDBPath, ctx.String("a"))
}

// sendCommandAction will send specified amount num funds from sender to
// receiver, must ensure sender's balance enough to close the deal
func sendCommandAction(ctx *cli.Context) {
	isInitialized()
	if ctx.String("f") == "" {
		panic("[ERROR] Sender address cannot be empty!\n")
	} else if ctx.String("t") == "" {
		panic("[ERROR] Receiver address cannot be empty!\n")
	} else if math.Signbit(float64(ctx.Int("n"))) && ctx.Int("n") != 0 {
		panic("[ERROR] Invalid amount number!\n")
	}

	_chain := lbc.NewBlockChain(ctx.String("f"))
	defer _chain.Close()
	_tx := lbc.NewUTXOTransaction(ctx.String("f"), ctx.String("t"), ctx.Int("n"), _chain)
	_chain.AddBlock([]*lbc.Transaction{_tx})
	fmt.Printf("[INFO] Successful send `%d` funds from `%s` to `%s`!\n", ctx.Int("n"), ctx.String("f"), ctx.String("t"))
}

// printCommandAction will print all blocks in block-chain
func printCommandAction(ctx *cli.Context) {
	isInitialized()
	if ctx.String("address") == "" {
		panic("[ERROR] Init address cannot be empty!\n")
	}
	_chain := lbc.NewBlockChain(ctx.String("address"))
	defer _chain.Close()
	_iterator := _chain.Iterator()
	for {
		block := _iterator.Next()

		pow := lbc.NewProofWork(block)
		fmt.Printf("\nPrev. Hash: %x\nTransaction Num: %d\nHash: %x\nPoW State: %s\n",
			block.Prev,
			len(block.Transactions),
			block.Hash,
			strconv.FormatBool(pow.Validate()))

		if len(block.Prev) == 0 {
			break
		}
	}
}

// getBalanceCommandAction will print spec address's UTXO value sum
func getBalanceCommandAction(ctx *cli.Context) {
	isInitialized()
	if ctx.String("address") == "" {
		panic("[ERROR] Address cannot be empty!\n")
	}

	chain := lbc.NewBlockChain(ctx.String("address"))
	defer chain.Close()

	balance := 0
	utxos := chain.FindUTXO(ctx.String("address"))
	for _, utxo := range utxos {
		balance += utxo.Value
	}

	fmt.Printf("[INFO] Balance of `%s`: %d", ctx.String("address"), balance)
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
