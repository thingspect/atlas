package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/atlas-api/crypto"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/dao/user"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

const usage = `Usage:
%[1]s uuid
%[1]s uniqid
%[1]s [options] org <org name> <admin email> <admin password>
%[1]s [options] user <org ID> <admin email> <admin password>
%[1]s [options] login <org name> <admin email> <admin password>
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	pgURI := flag.String("pgURI",
		"postgres://postgres:postgres@127.0.0.1/atlas_test", "PostgreSQL URI")
	grpcURI := flag.String("grpcURI", "127.0.0.1:50051", "gRPC URI")
	grpcTLS := flag.Bool("grpcTLS", false, "gRPC TLS")
	flag.Parse()

	if _, ok := map[string]struct{}{
		"uuid": {}, "uniqid": {}, "org": {}, "user": {}, "login": {},
	}[flag.Arg(0)]; !ok {
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
	// Log in user.
	case "login":
		opts := []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		}

		switch *grpcTLS {
		case false:
			opts = append(opts, grpc.WithInsecure())
		case true:
			opts = append(opts, grpc.WithTransportCredentials(
				credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})))
		}

		conn, err := grpc.Dial(*grpcURI, opts...)
		checkErr(err)

		cli := api.NewSessionServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		login, err := cli.Login(ctx, &api.LoginRequest{
			Email: flag.Arg(2), OrgName: flag.Arg(1), Password: flag.Arg(3),
		})
		cancel()
		checkErr(err)
		checkErr(conn.Close())

		fmt.Fprintf(os.Stdout, "Login: %+v\n", login)

		return
	}

	// Set up database connection.
	pg, err := dao.NewPgDB(*pgURI)
	checkErr(err)

	orgID := flag.Arg(1)

	switch flag.Arg(0) {
	// Create org and fall through to user.
	case "org":
		// Build org email.
		emailParts := strings.SplitN(flag.Arg(2), "@", 2)
		if len(emailParts) != 2 {
			fmt.Fprintln(os.Stderr, "invalid email address")
			os.Exit(1)
		}

		orgDAO := org.NewDAO(pg)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		createOrg, err := orgDAO.Create(ctx, &api.Org{
			Name:        flag.Arg(1),
			DisplayName: strings.Title(strings.ToLower(flag.Arg(1))),
			Email:       "noreply@" + emailParts[1],
		})
		cancel()
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

		u := &api.User{
			OrgId:  orgID,
			Email:  flag.Arg(2),
			Role:   common.Role_ADMIN,
			Status: common.Status_ACTIVE,
			Tags:   []string{strings.ToLower(common.Role_ADMIN.String())},
		}

		userDAO := user.NewDAO(pg)
		createUser, err := userDAO.Create(ctx, u)
		checkErr(err)

		checkErr(userDAO.UpdatePassword(ctx, createUser.Id, orgID, hash))
		fmt.Fprintf(os.Stdout, "User: %+v\n", createUser)
	}
}
