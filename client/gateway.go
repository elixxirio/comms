///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains client -> gateway functionality

package client

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"google.golang.org/grpc"
	"io"
	"strconv"
)

// Client -> Gateway Send Function
func (c *Comms) SendPutMessage(host *connect.Host, message *pb.GatewaySlot) (*pb.GatewaySlotResponse, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).PutMessage(ctx, message)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Put message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.GatewaySlotResponse{}

	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> Gateway Send Function
func (c *Comms) SendPutManyMessages(host *connect.Host, messages *pb.GatewaySlots) (*pb.GatewaySlotResponse, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).PutManyMessages(ctx, messages)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending PutManyMessages: %+v", messages)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.GatewaySlotResponse{}

	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> Gateway Send Function
func (c *Comms) SendRequestNonceMessage(host *connect.Host,
	message *pb.NonceRequest) (*pb.Nonce, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).RequestNonce(ctx, message)

		// Make sure there are no errors with sending the message
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Request Nonce message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.Nonce{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> Gateway Send Function
func (c *Comms) SendConfirmNonceMessage(host *connect.Host,
	message *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).ConfirmNonce(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Confirm Nonce message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RegistrationConfirmation{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> Gateway Send Function
func (c *Comms) SendPoll(host *connect.Host,
	message *pb.GatewayPoll) (*pb.GatewayPollResponse, error) {
	// Set up the context
	ctx, cancel := connect.StreamingContext()
	defer cancel()

	// Create the Stream Function
	f := func(conn *grpc.ClientConn) (interface{}, error) {

		// Send the message
		clientStream, err := pb.NewGatewayClient(conn).Poll(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return clientStream, nil
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Poll message: %+v", message)
	resultClient, err := c.Stream(host, f)
	if err != nil {
		return nil, err
	}

	stream := resultClient.(pb.Gateway_PollClient)
	jww.DEBUG.Printf("Receiving chunks for gateway poll")

	// Get the total number of chunks from the header
	md, err := stream.Header()
	if err != nil {
		return nil, errors.Errorf("Could not receive streaming header: %v", err)
	}

	// Return with an error if there is no header being received via the stream
	if md.Len() == 0 {
		return nil, errors.New(pb.NoStreamingHeaderErr)
	}

	totalChunks, err := strconv.Atoi(md.Get(pb.ChunkHeader)[0])
	if err != nil {
		return nil, errors.Errorf("Invalid header received: %v", err)
	}

	// Receive the chunks
	chunks := make([]*pb.StreamChunk, 0, totalChunks)
	chunk, err := stream.Recv()
	receivedChunks := 0
	for ; err == nil; chunk, err = stream.Recv() {
		chunks = append(chunks, chunk)
		receivedChunks++
	}
	if err != io.EOF { // EOF is an expected error after server-side has completed streaming
		return nil, errors.Errorf("Failed to complete streaming, received %d of %d messages: %v", receivedChunks, totalChunks, err)
	}

	// Assemble the result
	result := &pb.GatewayPollResponse{}
	return result, pb.AssembleChunksIntoResponse(chunks, result)
}

// Client -> Gateway Send Function
func (c *Comms) RequestHistoricalRounds(host *connect.Host,
	message *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).RequestHistoricalRounds(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Requesting Historical Rounds: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.HistoricalRoundsResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> Gateway Send Function
func (c *Comms) RequestMessages(host *connect.Host,
	message *pb.GetMessages) (*pb.GetMessagesResponse, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).RequestMessages(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Requesing Messages: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.GetMessagesResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
