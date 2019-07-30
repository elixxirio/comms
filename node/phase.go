////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> server functionality for precomputation operations

package node

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc/metadata"
)

func (s *NodeComms) SendPostPhase(id fmt.Stringer,
	message *pb.Batch) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := s.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.PostPhase(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("PostPhase: Error received: %+v", err)
	}

	cancel()
	return result, err
}

// GetPostPhaseStreamClient gets the streaming client
// using a header and returns the stream and the cancel context
// if there are no connection errors
func (s *NodeComms) GetPostPhaseStreamClient(id fmt.Stringer,
	header pb.BatchInfo) (pb.Node_StreamPostPhaseClient, context.CancelFunc, error) {

	ctx, cancel := s.getPostPhaseStreamContext(header)

	streamClient, err := s.getPostPhaseStream(id, ctx)

	if err != nil {
		return nil, nil, err
	}

	return streamClient, cancel, nil

}

// getPostPhaseStreamContext is given batchInfo PostPhase header
// and creates a streaming context, adds the header to the context
// and returns the context with the header and a cancel func
func (s *NodeComms) getPostPhaseStreamContext(batchInfo pb.BatchInfo) (context.Context, context.CancelFunc) {

	// Create streaming context so you can close stream later
	ctx, cancel := connect.StreamingContext()

	// Create a new context with some metadata
	// using the batch info batchInfo
	ctx = metadata.AppendToOutgoingContext(ctx, "batchinfo", batchInfo.String())

	return ctx, cancel
}

// getPostPhaseStream uses an id and streaming context to retrieve
// a Node_StreamPostPhaseClient object otherwise it returns
// an error if the connection is unavailable
func (s *NodeComms) getPostPhaseStream(id fmt.Stringer, ctx context.Context) (
	pb.Node_StreamPostPhaseClient, error) {

	// Attempt to connect to addr
	c := s.GetNodeConnection(id)

	// Get the stream client using streaming context
	streamClient, err := c.StreamPostPhase(ctx,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with getting the stream client
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("getPostPhaseStream: Error received: %+v", err)
		return nil, err
	}

	return streamClient, nil
}

// GetPostPhaseStreamHeader gets the header
// in the metadata from the server stream
// and returns it with an error if it fails.
func GetPostPhaseStreamHeader(stream pb.Node_StreamPostPhaseServer) (*pb.BatchInfo, error) {

	// Unmarshal header into batch info
	batchInfo := pb.BatchInfo{}

	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok {
		return nil, errors.New("unable to retrieve meta data / header %v")
	}

	err := proto.UnmarshalText(md.Get("batchinfo")[0], &batchInfo)
	if err != nil {
		return nil, err
	}

	return &batchInfo, nil

}
