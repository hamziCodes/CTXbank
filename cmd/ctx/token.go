package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hamziCodes/CTXbank/internal/core"
)

func runToken(args []string) {
	fs := flag.NewFlagSet("token", flag.ExitOnError)
	regen := fs.Bool("regenerate", false, "Replace the project token with a fresh value (old one stops working)")
	_ = fs.Parse(args)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	bankDir := filepath.Join(cwd, core.MemoryBankDir)
	if _, err := os.Stat(bankDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: memory-bank not found. Run 'ctx init' first.\n")
		os.Exit(1)
	}

	var token string
	if *regen {
		token, err = core.RegenerateProjectToken(bankDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error regenerating token: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Project token regenerated. The old token no longer works.")
	} else {
		token, err = core.GetProjectToken(bankDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading token: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\nProject dashboard token for %s:\n\n  %s\n\n", filepath.Base(cwd), token)
	fmt.Println("Paste it on the Connect page to open this project's dashboard.")
	fmt.Println("Keep it private — anyone with this token can read and edit")
	fmt.Println("this project's memory bank via the local dashboard.")
}
