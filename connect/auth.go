////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles authentication logic for the top-level comms object

package connect

import (
	"bytes"
	"crypto/rand"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/nonce"
	"gitlab.com/elixxir/crypto/signature/rsa"
)

// Auth represents an authorization state for a message or host
type Auth struct {
	IsAuthenticated bool
	Sender          Host
}

// Perform the client handshake to establish reverse-authentication
func (c *ProtoComms) clientHandshake(host *Host) (err error) {

	// Set up the context
	client := pb.NewGenericClient(host.connection)
	ctx, cancel := MessagingContext()
	defer cancel()

	// Send the token request message
	result, err := client.RequestToken(ctx,
		&pb.Ping{})
	if err != nil {
		return errors.New(err.Error())
	}

	// Assign the host token
	host.token = result.Token

	// Pack the authenticated message with signature enabled
	msg, err := c.PackAuthenticatedMessage(&pb.AssignToken{
		Token: host.token,
	}, host, true)

	// Set up the context
	ctx, cancel = MessagingContext()
	defer cancel()

	// Send the authenticate token message
	_, err = client.AuthenticateToken(ctx, msg)
	if err != nil {
		err = errors.New(err.Error())
	}

	return
}

// Convert any message type into a authenticated message
func (c *ProtoComms) PackAuthenticatedMessage(msg proto.Message, host *Host,
	enableSignature bool) (*pb.AuthenticatedMessage, error) {

	// Marshall the provided message into an Any type
	anyMsg, err := ptypes.MarshalAny(msg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Build the authenticated message
	authMsg := &pb.AuthenticatedMessage{
		ID:        host.id,
		Signature: nil,
		Token:     host.token,
		Message:   anyMsg,
	}

	// If signature is enabled, sign the message and add to payload
	if enableSignature {
		authMsg.Signature, err = c.signMessage(anyMsg)
		if err != nil {
			return nil, err
		}
	}

	return authMsg, nil
}

// Generates a new token and adds it to internal state
func (c *ProtoComms) GenerateToken() ([]byte, error) {
	token, err := nonce.NewNonce(nonce.RegistrationTTL)
	if err != nil {
		return nil, err
	}

	c.tokens.Store(string(token.Bytes()), &token)
	return token.Bytes(), nil
}

// Validates a signed token using internal state
func (c *ProtoComms) ValidateToken(msg *pb.AuthenticatedMessage) error {

	// Verify the Host exists for the provided ID
	host, ok := c.GetHost(msg.ID)
	if !ok {
		return errors.Errorf("Invalid token for host ID: %+v", msg.ID)
	}

	// Verify the token signature
	if err := c.verifyMessage(msg, host); err != nil {
		return errors.Errorf("Invalid token signature: %+v", err)
	}

	// Get the signed token
	tokenMsg := &pb.AssignToken{}
	err := ptypes.UnmarshalAny(msg.Message, tokenMsg)
	if err != nil {
		return errors.Errorf("Unable to unmarshal token: %+v", err)
	}

	// Verify the signed token was actually assigned
	token, ok := c.tokens.Load(string(tokenMsg.Token))
	if !ok {
		return errors.Errorf("Unable to locate token: %+v", msg.Token)
	}

	// Verify the signed token is not expired
	if !token.(*nonce.Nonce).IsValid() {
		return errors.Errorf("Invalid or expired token: %+v", msg.Token)
	}

	// Token has been validated and can be safely stored
	host.token = msg.Token
	return nil
}

// AuthenticatedReceiver handles reception of an AuthenticatedMessage,
// checking if the host is authenticated & returning an Auth state
func (c *ProtoComms) AuthenticatedReceiver(msg *pb.AuthenticatedMessage) *Auth {
	res := &Auth{
		IsAuthenticated: false,
		Sender:          Host{},
	}

	// Check if the sender is authenticated, and if the token is valid
	host, ok := c.GetHost(msg.ID)
	if ok && bytes.Compare(host.token, msg.Token) == 0 {
		res.Sender = *host
		res.IsAuthenticated = true
	}
	return res
}

// Takes a generic-type message, returns the signature
// The message is signed with the ProtoComms RSA PrivateKey
func (c *ProtoComms) signMessage(anyMessage *any.Any) ([]byte, error) {
	// Hash the message data
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	hash.Write([]byte(anyMessage.String()))
	hashed := hash.Sum(nil)

	// Obtain the private key
	key := c.GetPrivateKey()
	if key == nil {
		return nil, errors.Errorf("Cannot sign message: No private key")
	}

	// Sign the message and return the signature
	signature, err := rsa.Sign(rand.Reader, key, options.Hash, hashed, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return signature, nil
}

// Takes an AuthenticatedMessage and a Host, verifies the signature
// using Host public key, returning an error if invalid
func (c *ProtoComms) verifyMessage(msg *pb.AuthenticatedMessage, host *Host) error {

	// Get hashed data of the message
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	hash.Write([]byte(msg.Message.String()))
	hashed := hash.Sum(nil)

	// Verify signature of message using host public key
	err := rsa.Verify(host.rsaPublicKey, options.Hash, hashed, msg.Signature, nil)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
