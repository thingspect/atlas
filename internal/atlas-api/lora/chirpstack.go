package lora

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Chirpstack holds references to the ChirpStack server and implements the
// Loraer interface.
type Chirpstack struct {
	tenantID  string
	appID     string
	devProfID string

	conn *grpc.ClientConn
}

// Verify Chirpstack implements Loraer.
var _ Loraer = &Chirpstack{}

// NewChirpstack builds a new Loraer and returns it and an error value.
func NewChirpstack(addr, apiKey, tenantID, appID, devProfID string) (
	Loraer, error,
) {
	// Build Chirpstack gRPC connection.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&credential{token: apiKey}),
	}
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}

	return &Chirpstack{
		tenantID:  tenantID,
		appID:     appID,
		devProfID: devProfID,

		conn: conn,
	}, nil
}
