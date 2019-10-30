////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway gRPC endpoints

package gateway

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

// Sends new MessageIDs in the buffer to a client
func (g *GatewayComms) CheckMessages(ctx context.Context, msg *pb.ClientRequest) (
	*pb.IDList, error) {
	userID := id.NewUserFromBytes(msg.UserID)
	msgIds, ok := g.handler.CheckMessages(userID, msg.LastMessageID)
	returnMsg := &pb.IDList{}
	if ok {
		returnMsg.IDs = msgIds
	}
	return returnMsg, nil
}

// Sends a message matching the given parameters to a client
func (g *GatewayComms) GetMessage(ctx context.Context, msg *pb.ClientRequest) (
	*pb.Slot, error) {
	userID := id.NewUserFromBytes(msg.UserID)
	returnMsg, ok := g.handler.GetMessage(userID, msg.LastMessageID)
	if !ok {
		// Return an empty message if no results
		returnMsg = &pb.Slot{}
	}
	return returnMsg, nil
}

// Receives a single message from a client
func (g *GatewayComms) PutMessage(ctx context.Context, msg *pb.Slot) (*pb.Ack,
	error) {

	p, ok := peer.FromContext(ctx)
	if ok {
		gatewayHandler.PutMessage(msg, p.Addr.String())
	}

	return &pb.Ack{}, nil
}

// Receives a batch of messages from a server
func (s *gateway) ReceiveBatch(ctx context.Context, msg *pb.OutputMessages) (*pb.Ack,
	error) {
	g.handler.PutMessage(msg)
	return &pb.Ack{}, nil
}

// Pass-through for Registration Nonce Communication
func (g *GatewayComms) RequestNonce(ctx context.Context,
	msg *pb.NonceRequest) (*pb.Nonce, error) {
	return g.handler.RequestNonce(msg)
}

// Pass-through for Registration Nonce Confirmation
func (g *GatewayComms) ConfirmNonce(ctx context.Context,
	msg *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation, error) {
	return g.handler.ConfirmNonce(msg)
}
