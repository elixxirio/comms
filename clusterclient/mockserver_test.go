////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// This sets up a dummy/mock server instance for testing purposes
package clusterclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"gitlab.com/privategrity/comms/mixserver"
	"os"
	"testing"
)

const SERVER_ADDRESS = "localhost:5555"

// Start server for testing
func TestMain(m *testing.M) {
	go mixserver.StartServer(SERVER_ADDRESS, TestInterface{})
	os.Exit(m.Run())
}

// Blank struct implementing ServerHandler interface for testing purposes (Passing to StartServer)
type TestInterface struct{}

func (m TestInterface) NewRound(roundId string) {}

func (m TestInterface) SetPublicKey(roundId string, pkey []byte) {}

func (m TestInterface) PrecompDecrypt(message *pb.PrecompDecryptMessage) {}

func (m TestInterface) PrecompEncrypt(message *pb.PrecompEncryptMessage) {}

func (m TestInterface) PrecompReveal(message *pb.PrecompRevealMessage) {}

func (m TestInterface) PrecompPermute(message *pb.PrecompPermuteMessage) {}

func (m TestInterface) PrecompShare(message *pb.PrecompShareMessage) {}

func (m TestInterface) PrecompShareInit(message *pb.PrecompShareInitMessage) {}

func (m TestInterface) PrecompShareCompare(message *pb.
	PrecompShareCompareMessage) {}

func (m TestInterface) PrecompShareConfirm(message *pb.
	PrecompShareConfirmMessage) {}

func (m TestInterface) RealtimeDecrypt(message *pb.RealtimeDecryptMessage) {}

func (m TestInterface) RealtimeEncrypt(message *pb.RealtimeEncryptMessage) {}

func (m TestInterface) RealtimePermute(message *pb.RealtimePermuteMessage) {}

func (m TestInterface) ClientPoll(message *pb.ClientPollMessage) *pb.CmixMessage {
	return &pb.CmixMessage{}
}

func (m TestInterface) RequestContactList(message *pb.ContactPoll) *pb.
ContactMessage {
	return &pb.ContactMessage{}
}

func (m TestInterface) SetNick(message *pb.Contact) {}

func (m TestInterface) ReceiveMessageFromClient(message *pb.CmixMessage) {}
