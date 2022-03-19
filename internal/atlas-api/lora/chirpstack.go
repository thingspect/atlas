package lora

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Chirpstack holds references to the ChirpStack Application Server and
// implements the Loraer interface.
type Chirpstack struct {
	orgID     int64
	nsID      int64
	appID     int64
	devProfID string

	conn *grpc.ClientConn
}

// Verify Chirpstack implements Loraer.
var _ Loraer = &Chirpstack{}

// NewChirpstack builds a new Loraer and returns it and an error value.
func NewChirpstack(
	addr, apiKey string, orgID, nsID, appID int, devProfID string,
) (Loraer, error) {
	// Build Chirpstack AS gRPC connection.
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&credential{token: apiKey}),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return &Chirpstack{
		orgID:     int64(orgID),
		nsID:      int64(nsID),
		appID:     int64(appID),
		devProfID: devProfID,

		conn: conn,
	}, nil
}
