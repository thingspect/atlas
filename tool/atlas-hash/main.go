// Package main runs the Hash tool.
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thingspect/atlas/internal/atlas-api/auth"
	"github.com/thingspect/atlas/internal/atlas-api/session"
)

const usage = `Usage of %[1]s:
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
		p, err := os.Executable()
		checkErr(err)

		_, err = fmt.Fprintf(flag.CommandLine.Output(), usage, filepath.Base(p))
		checkErr(err)

		flag.PrintDefaults()
	}

	flag.Parse()

	switch {
	// Hash password.
	case flag.NArg() == 2 && flag.Arg(0) == "pass":
		hash, err := auth.HashPass(flag.Arg(1))
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
