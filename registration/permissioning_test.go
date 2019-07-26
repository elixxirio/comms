////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package registration

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/comms/testkeys"
	"testing"
)

// Smoke test SendNodeTopology
func TestSendNodeTopology(t *testing.T) {
	ServerAddress := getNextServerAddress()
	RegAddress := getNextServerAddress()

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	reg := StartRegistrationServer(RegAddress,
		NewImplementation(), certData, keyData)
	defer server.Shutdown()
	defer reg.Shutdown()

	connID := MockID("permissioningToServer")
	regID := MockID("Permissioning")

	server.ConnectToRegistration(regID, RegAddress, certData)
	reg.ConnectToNode(connID, ServerAddress, nil)

	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(connID, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
	}
}

func TestSendNodeTopologyNilKeyError(t *testing.T) {
	ServerAddress := getNextServerAddress()
	RegAddress := getNextServerAddress()

	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	reg := StartRegistrationServer(RegAddress,
		NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()

	connID := MockID("permissioningToServer")
	regID := MockID("Permissioning")

	server.ConnectToRegistration(regID, RegAddress, nil)
	reg.ConnectToNode(connID, ServerAddress, nil)

	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(connID, msgs)
	if err == nil {
		t.Errorf("SendNodeTopology: did not receive missing private key error")
	}
}
