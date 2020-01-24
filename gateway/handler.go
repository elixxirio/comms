////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for gateway functionality

package gateway

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"google.golang.org/grpc/reflection"
	"runtime/debug"
)

// Handler interface for the Gateway
type Handler interface {
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages(userID *id.User, messageID string, ipAddress string) ([]string, error)
	// Returns the message matching the given parameters to the client
	GetMessage(userID *id.User, msgID string, ipAddress string) (*pb.Slot, error)
	// Upload a message to the cMix Gateway
	PutMessage(message *pb.Slot, ipAddress string) error
	// Pass-through for Registration Nonce Communication
	RequestNonce(message *pb.NonceRequest, ipAddress string) (*pb.Nonce, error)
	// Pass-through for Registration Nonce Confirmation
	ConfirmNonce(message *pb.RequestRegistrationConfirmation, ipAddress string) (*pb.
		RegistrationConfirmation, error)
	// Ping gateway to ask for users to notify
	PollForNotifications(auth *connect.Auth) ([]string, error)
}

// Gateway object used to implement endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler
}

// Starts a new gateway on the address:port specified by localServer
// and a callback interface for gateway operations
// with given path to public and private key for TLS connection
func StartGateway(id, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {
	pc, lis, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock)
	if err != nil {
		jww.FATAL.Panicf("Unable to start comms server: %+v", err)
	}

	gatewayServer := Comms{
		ProtoComms: pc,
		handler:    handler,
	}

	go func() {
		pb.RegisterGatewayServer(gatewayServer.LocalServer, &gatewayServer)

		// Register reflection service on gRPC server.
		// This blocks for the lifetime of the listener.
		reflection.Register(gatewayServer.LocalServer)
		if err := gatewayServer.LocalServer.Serve(lis); err != nil {
			jww.FATAL.Panicf("Failed to serve: %+v",
				errors.New(err.Error()))
		}
		jww.INFO.Printf("Shutting down gateway server listener: %s",
			lis)

	}()

	return &gatewayServer
}

// Handler implementation for the Gateway
type implementationFunctions struct {
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages func(userID *id.User, messageID string, ipAddress string) ([]string, error)
	// Returns the message matching the given parameters to the client
	GetMessage func(userID *id.User, msgID string, ipAddress string) (*pb.Slot, error)
	// Upload a message to the cMix Gateway
	PutMessage func(message *pb.Slot, ipAddress string) error
	// Pass-through for Registration Nonce Communication
	RequestNonce func(message *pb.NonceRequest, ipAddress string) (*pb.Nonce, error)
	// Pass-through for Registration Nonce Confirmation
	ConfirmNonce func(message *pb.RequestRegistrationConfirmation, ipAddress string) (*pb.
			RegistrationConfirmation, error)
	// Ping gateway to ask for users to notify
	PollForNotifications func(auth *connect.Auth) ([]string, error)
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// Creates and returns a new Handler interface
func NewImplementation() *Implementation {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%s", debug.Stack())
	}
	return &Implementation{
		Functions: implementationFunctions{
			CheckMessages: func(userID *id.User, messageID string, ipAddress string) ([]string, error) {
				warn(um)
				return nil, nil
			},
			GetMessage: func(userID *id.User, msgID string, ipAddress string) (*pb.Slot, error) {
				warn(um)
				return &pb.Slot{}, nil
			},
			PutMessage: func(message *pb.Slot, ipAddress string) error {
				warn(um)
				return nil
			},
			RequestNonce: func(message *pb.NonceRequest, ipAddress string) (*pb.Nonce, error) {
				warn(um)
				return new(pb.Nonce), nil
			},
			ConfirmNonce: func(message *pb.RequestRegistrationConfirmation, ipAddress string) (*pb.RegistrationConfirmation, error) {
				warn(um)
				return new(pb.RegistrationConfirmation), nil
			},
			PollForNotifications: func(auth *connect.Auth) ([]string, error) {
				warn(um)
				return nil, nil
			},
		},
	}
}

// Return any MessageIDs in the buffer for this UserID
func (s *Implementation) CheckMessages(userID *id.User, messageID string, ipAddress string) (
	[]string, error) {
	return s.Functions.CheckMessages(userID, messageID, ipAddress)
}

// Returns the message matching the given parameters to the client
func (s *Implementation) GetMessage(userID *id.User, msgID string, ipAddress string) (
	*pb.Slot, error) {
	return s.Functions.GetMessage(userID, msgID, ipAddress)
}

// Upload a message to the cMix Gateway
func (s *Implementation) PutMessage(message *pb.Slot, ipAddress string) error {
	return s.Functions.PutMessage(message, ipAddress)
}

// Pass-through for Registration Nonce Communication
func (s *Implementation) RequestNonce(message *pb.NonceRequest, ipAddress string) (
	*pb.Nonce, error) {
	return s.Functions.RequestNonce(message, ipAddress)
}

// Pass-through for Registration Nonce Confirmation
func (s *Implementation) ConfirmNonce(message *pb.RequestRegistrationConfirmation,
	ipAddress string) (*pb.RegistrationConfirmation, error) {
	return s.Functions.ConfirmNonce(message, ipAddress)
}

// Ping gateway to ask for users to notify
func (s *Implementation) PollForNotifications(auth *connect.Auth) ([]string, error) {
	return s.Functions.PollForNotifications(auth)
}
