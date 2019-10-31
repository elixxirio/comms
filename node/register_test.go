////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package node

import (
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"testing"
)

// Smoke test SendNodeRegistration
func TestSendNodeRegistration(t *testing.T) {
	RegAddress := getNextServerAddress()
	server := StartNode(getNextServerAddress(), NewImplementation(),
		nil, nil)
	reg := registration.StartRegistrationServer(RegAddress,
		registration.NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()

	msgs := &pb.NodeRegistration{}
	err := server.SendNodeRegistration(&connect.Host{
		Id:             "serverToPermissioning",
		Address:        RegAddress,
		Cert:           nil,
		DisableTimeout: false,
	}, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
	}
}
