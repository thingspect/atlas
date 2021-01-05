package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/crypto"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/dao/user"
	"github.com/thingspect/atlas/pkg/postgres"
)

const usage = `Usage:
%[1]s org <org name> <admin email> <admin password>
%[1]s user <org ID> <user email> <user password>
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	pgURI := flag.String("pgURI",
		"postgres://postgres:postgres@127.0.0.1/atlas_test", "PostgreSQL URI")
	flag.Parse()

	if flag.NArg() != 4 || (flag.Arg(0) != "org" && flag.Arg(0) != "user") {
		flag.Usage()
		os.Exit(2)
	}

	checkErr := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Set up database connection.
	pg, err := postgres.New(*pgURI)
	checkErr(err)
	orgDAO := org.NewDAO(pg)
	userDAO := user.NewDAO(pg)
	orgID := flag.Arg(1)

	switch {
	// Create org and fall through to user.
	case flag.Arg(0) == "org":
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		org := org.Org{Name: flag.Arg(1)}
		createOrg, err := orgDAO.Create(ctx, org)
		checkErr(err)
		orgID = createOrg.ID
		fmt.Fprintf(os.Stdout, "Org: %+v\n", createOrg)
		fallthrough
	// Create user.
	case flag.Arg(0) == "user":
		hash, err := crypto.HashPass(flag.Arg(3))
		checkErr(err)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		user := &api.User{OrgId: orgID, Email: flag.Arg(2)}
		createUser, err := userDAO.Create(ctx, user, hash)
		checkErr(err)
		fmt.Fprintf(os.Stdout, "User: %+v\n", createUser)
	}
}
