////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> gateway functionality

package node

import (
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// SendReceiveBatch sends a batch to the gateway
func (s *Server) SendReceiveBatch(id fmt.Stringer, message []*pb.Batch) error {
	// Attempt to connect to addr
	c := s.manager.ConnectToGateway(id, nil)
	ctx, cancel := connect.DefaultContext()

	outputMessages := pb.Output{Messages: message}
	_, err := c.ReceiveBatch(ctx, &outputMessages)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("ReceiveBatch(): Error received: %+v", err)
	}

	cancel()
	return err
}
