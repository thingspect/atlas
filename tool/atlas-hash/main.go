package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/crypto"
)

const usage = `Usage:
%[1]s pass <password>
%[1]s pwt <key, base64> <token>
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	checkErr := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	switch {
	// Hash password.
	case flag.NArg() == 2 && flag.Arg(0) == "pass":
		hash, err := crypto.HashPass(flag.Arg(1))
		checkErr(err)
		fmt.Fprintf(os.Stdout, "%s\n", hash)
	// Decrypt and validate web token.
	case flag.NArg() == 3 && flag.Arg(0) == "pwt":
		key, err := base64.StdEncoding.DecodeString(flag.Arg(1))
		checkErr(err)
		sess, err := session.ValidateWebToken(key, flag.Arg(2))
		fmt.Fprintf(os.Stdout, "Session: %+v\nError: %v\n", sess, err)
	default:
		flag.Usage()
		os.Exit(2)
	}
}
