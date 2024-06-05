// Package main runs the Hash tool.
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/thingspect/atlas/internal/atlas-api/crypto"
	"github.com/thingspect/atlas/internal/atlas-api/session"
)

const usage = `Usage:
%[1]s pass <password>
%[1]s pwt <base64 key> <token>
`

func main() {
	checkErr := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	flag.Usage = func() {
		_, err := fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		checkErr(err)

		flag.PrintDefaults()
	}

	flag.Parse()

	switch {
	// Hash password.
	case flag.NArg() == 2 && flag.Arg(0) == "pass":
		hash, err := crypto.HashPass(flag.Arg(1))
		checkErr(err)

		_, err = fmt.Fprintf(os.Stdout, "%s\n", hash)
		checkErr(err)
	// Decrypt and validate web token.
	case flag.NArg() == 3 && flag.Arg(0) == "pwt":
		key, err := base64.StdEncoding.DecodeString(flag.Arg(1))
		checkErr(err)

		sess, err := session.ValidateWebToken(key, flag.Arg(2))
		_, err = fmt.Fprintf(os.Stdout, "Session: %+v\nError: %v\n", sess, err)
		checkErr(err)
	default:
		flag.Usage()
		os.Exit(2)
	}
}
