///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains user discovery server gRPC endpoint wrappers
// When you add the udb server to mixmessages/mixmessages.proto and add the
// first function, a version of that goes here which calls the "handler"
// version of the function, with any mappings/wrappings necessary.

package udb

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/messages"
)

// Handles validation of reverse-authentication tokens
func (u *Comms) AuthenticateToken(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	err := u.ValidateToken(msg)
	if err != nil {
		jww.ERROR.Printf("Unable to authenticate token: %+v", err)
	}
	return &messages.Ack{}, err
}

// Handles reception of reverse-authentication token requests
func (u *Comms) RequestToken(context.Context, *messages.Ping) (*messages.AssignToken, error) {
	token, err := u.GenerateToken()
	return &messages.AssignToken{
		Token: token,
	}, err
}

func (u *Comms) RegisterUser(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Create an auth object
	authState, err := u.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	registration := &pb.UDBUserRegistration{}
	err = ptypes.UnmarshalAny(msg.Message, registration)
	if err != nil {
		return nil, err
	}

	return u.handler.RegisterUser(registration, authState)
}

func (u *Comms) RegisterFact(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.FactRegisterResponse, error) {
	// Create an auth object
	authState, err := u.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	request := &pb.FactRegisterRequest{}
	err = ptypes.UnmarshalAny(msg.Message, request)
	if err != nil {
		return nil, err
	}

	return u.handler.RegisterFact(request, authState)
}

func (u *Comms) ConfirmFact(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Create an auth object
	authState, err := u.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	request := &pb.FactConfirmRequest{}
	err = ptypes.UnmarshalAny(msg.Message, request)
	if err != nil {
		return nil, err
	}

	return u.handler.ConfirmFact(request, authState)
}

func (u *Comms) RemoveFact(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Create an auth object
	authState, err := u.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	request := &pb.FactRemovalRequest{}
	err = ptypes.UnmarshalAny(msg.Message, request)
	if err != nil {
		return nil, err
	}

	return u.handler.RemoveFact(request, authState)
}
