package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/crypto"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/dao/user"
	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/test/random"
)

const usage = `Usage:
%[1]s uuid
%[1]s uniqid
%[1]s [options] org <org name> <admin email> <admin password>
%[1]s [options] user <org ID> <admin email> <admin password>
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	pgURI := flag.String("pgURI",
		"postgres://postgres:postgres@127.0.0.1/atlas_test", "PostgreSQL URI")
	flag.Parse()

	if _, ok := map[string]struct{}{"uuid": {}, "uniqid": {}, "org": {},
		"user": {}}[flag.Arg(0)]; !ok {
		flag.Usage()
		os.Exit(2)
	}

	checkErr := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	switch flag.Arg(0) {
	// Generate UUID and return.
	case "uuid":
		fmt.Fprintln(os.Stdout, uuid.NewString())
		return
	// Generate UniqID and return.
	case "uniqid":
		fmt.Fprintln(os.Stdout, random.String(16))
		return
	}

	// Set up database connection.
	pg, err := postgres.New(*pgURI)
	checkErr(err)
	orgDAO := org.NewDAO(pg)
	userDAO := user.NewDAO(pg)
	orgID := flag.Arg(1)

	switch flag.Arg(0) {
	// Create org and fall through to user.
	case "org":
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		createOrg, err := orgDAO.Create(ctx, &api.Org{Name: flag.Arg(1)})
		checkErr(err)
		orgID = createOrg.Id
		fmt.Fprintf(os.Stdout, "Org: %+v\n", createOrg)
		fallthrough
	// Create user.
	case "user":
		hash, err := crypto.HashPass(flag.Arg(3))
		checkErr(err)

		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()

		user := &api.User{
			OrgId:  orgID,
			Email:  flag.Arg(2),
			Role:   common.Role_ADMIN,
			Status: api.Status_ACTIVE,
		}
		createUser, err := userDAO.Create(ctx, user)
		checkErr(err)

		checkErr(userDAO.UpdatePassword(ctx, user.Id, orgID, hash))
		fmt.Fprintf(os.Stdout, "User: %+v\n", createUser)
	}
}
