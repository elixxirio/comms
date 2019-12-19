////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains registration server gRPC endpoints

package registration

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
)

// Handles validation of reverse-authentication tokens
func (s *Comms) AuthenticateToken(ctx context.Context,
	msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	return &pb.Ack{}, s.ValidateToken(msg)
}

// Handles reception of reverse-authentication token requests
func (s *Comms) RequestToken(context.Context, *pb.Ping) (*pb.AssignToken, error) {
	token, err := s.GenerateToken()
	return &pb.AssignToken{
		Token: token,
	}, err
}

//Reviewer: Do I add an auth here? flim flam
// RegisterUser event handler which registers a user with the platform
func (r *Comms) RegisterUser(ctx context.Context, msg *pb.UserRegistration) (
	*pb.UserRegistrationConfirmation, error) {
	// Obtain the signed key by passing to registration server
	pubKey := msg.GetClientRSAPubKey()
	signature, err := r.handler.RegisterUser(msg.GetRegistrationCode(), pubKey)
	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
		err = errors.New(err.Error())
	}

	// Return the confirmation message
	return &pb.UserRegistrationConfirmation{
		ClientSignedByServer: &pb.RSASignature{
			Signature: signature,
		},
		Error: errMsg,
	}, err
}

//Reviewer: Do I add an auth here? flim flam
// CheckClientVersion event handler which checks whether the client library
// version is compatible with the network
func (r *Comms) GetCurrentClientVersion(ctx context.Context, ping *pb.Ping) (*pb.ClientVersion, error) {

	version, err := r.handler.GetCurrentClientVersion()

	// Return the confirmation message
	return &pb.ClientVersion{
		Version: version,
	}, err
}

// Handle a node registration event
func (r *Comms) RegisterNode(ctx context.Context, msg *pb.NodeRegistration) (
	*pb.Ack, error) {

	// Obtain peer IP address
	ip, port, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return &pb.Ack{}, err
	}
	address := fmt.Sprintf("%s:%s", ip, port)

	// Pass information for Node registration
	err = r.handler.RegisterNode(msg.GetID(), address, msg.GetServerTlsCert(),
		msg.GetGatewayAddress(), msg.GetGatewayTlsCert(),
		msg.GetRegistrationCode())
	return &pb.Ack{}, err
}

// Handles incoming requests for the NDF
func (r *Comms) PollNdf(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.NDF, error) {
	//Marshall the any message to the message type needed
	ndfHash := &pb.NDFHash{}
	err := ptypes.UnmarshalAny(msg.Message, ndfHash)
	if err != nil {
		return nil, err
	}
	authMsg := r.AuthenticatedReceiver(msg)

	newNDF, err := r.handler.PollNdf(ndfHash.Hash, authMsg)
	//Return the new ndf
	return &pb.NDF{Ndf: newNDF}, err
}
