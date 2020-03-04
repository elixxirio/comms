////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains notificationBot -> all servers functionality

package notificationBot

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
)

// NotificationBot -> Permissioning
func (nb *Comms) RequestNdf(host *connect.Host, message *pb.NDFHash) (*pb.NDF, error) {

	// Call the ProtoComms RequestNdf call
	return nb.ProtoComms.RequestNdf(host, message)
}

// Notification Bot -> Gateway
func (nb *Comms) RequestNotifications(host *connect.Host) (*pb.IDList, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		authMsg, err := nb.PackAuthenticatedMessage(&pb.Ping{}, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).PollForNotifications(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.DEBUG.Printf("Sending Request Notification message")
	resultMsg, err := nb.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.IDList{}
	return result, ptypes.UnmarshalAny(resultMsg, result)

}