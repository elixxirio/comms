////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains registration server comms initialization functionality

package registration

import (
	"crypto/tls"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"math"
	"net"
)

// Registration object used to implement
// endpoints and top-level comms functionality
type Comms struct {
	connect.ProtoComms
	handler Handler
}

// Starts a new server on the address:port specified by localServer
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartRegistrationServer(localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {
	var grpcServer *grpc.Server

	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
	if err != nil {
		err = errors.New(err.Error())
		jww.FATAL.Panicf("Failed to listen: %+v", err)
	}

	// If TLS was specified
	if certPEMblock != nil && keyPEMblock != nil {
		// Create the TLS certificate
		x509cert, err2 := tls.X509KeyPair(certPEMblock, keyPEMblock)
		if err2 != nil {
			err = errors.New(err2.Error())
			jww.FATAL.Panicf("Could not load TLS keys: %+v", err)
		}

		creds := credentials.NewServerTLSFromCert(&x509cert)

		// Create the gRPC server with TLS
		jww.INFO.Printf("Starting server with TLS...")
		grpcServer = grpc.NewServer(grpc.Creds(creds),
			grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(math.MaxInt32))
	} else {
		// Create the gRPC server without TLS
		jww.WARN.Printf("Starting server with TLS disabled...")
		grpcServer = grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(math.MaxInt32))
	}

	registrationServer := Comms{
		ProtoComms: connect.ProtoComms{
			LocalServer:   grpcServer,
			ListeningAddr: localServer,
		},
		handler: handler,
	}

	if keyPEMblock != nil {
		err = registrationServer.SetPrivateKey(keyPEMblock)
		if err != nil {
			jww.ERROR.Printf("Error setting RSA private key: %+v", err)
		}
	} else {
		jww.WARN.Println("Starting registration server with no private key...")
	}

	go func() {
		pb.RegisterRegistrationServer(registrationServer.LocalServer, &registrationServer)
		pb.RegisterGenericServer(registrationServer.LocalServer, &registrationServer)

		// Register reflection service on gRPC server.
		reflection.Register(registrationServer.LocalServer)
		if err = registrationServer.LocalServer.Serve(lis); err != nil {
			err = errors.New(err.Error())
			jww.FATAL.Panicf("Failed to serve: %+v", err)
		}
		jww.INFO.Printf("Shutting down registration server listener:"+
			" %s", lis)
	}()

	return &registrationServer
}
