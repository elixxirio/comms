////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"golang.org/x/net/context"
)

// Handle a GetMessage event
func (s *gateway) GetMessage(ctx context.Context, msg *pb.ClientPollMessage) (
	*pb.CmixMessage, error) {
	returnMsg, _ := gatewayHandler.GetMessage(msg.UserID, msg.MessageID)
	return returnMsg, nil
}