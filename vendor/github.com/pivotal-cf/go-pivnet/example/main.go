package main

import (
	"fmt"
	"log"
	"os"

	pivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
)

func main() {
	config := pivnet.ClientConfig{
		Host:      pivnet.DefaultHost,
		Token:     "token-from-pivnet",
		UserAgent: "pivnet-cli-example",
	}

	stdoutLogger := log.New(os.Stdout, "", log.LstdFlags)
	stderrLogger := log.New(os.Stderr, "", log.LstdFlags)

	verbose := false
	logger := logshim.NewLogShim(stdoutLogger, stderrLogger, verbose)

	client := pivnet.NewClient(config, logger)

	products, err := client.Products.List()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("products: %v", products)
}
