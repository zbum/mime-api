package main

import (
	"bytes"
	"context"
	"io"
	"net/mail"

	pb "mime_api/proto"

	"google.golang.org/grpc"
)

type displayPartServer struct {
	pb.UnimplementedDisplayPartServiceServer
}

func NewDisplayPartServer() *displayPartServer {
	return &displayPartServer{}
}

func (s *displayPartServer) ExtractDisplayPart(ctx context.Context, req *pb.ExtractRequest) (*pb.ExtractResponse, error) {
	reader := bytes.NewReader(req.EmlContent)
	message, err := mail.ReadMessage(reader)
	if err != nil {
		return &pb.ExtractResponse{
			Error: err.Error(),
		}, nil
	}

	content, contentType, err := ExtractDisplayPart(message)
	if err != nil {
		return &pb.ExtractResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.ExtractResponse{
		Content:     content,
		ContentType: contentType,
	}, nil
}

func (s *displayPartServer) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status: "OK",
	}, nil
}

func (s *displayPartServer) ExtractDisplayPartStream(stream grpc.ClientStreamingServer[pb.FileChunk, pb.ExtractResponse]) error {
	var buffer bytes.Buffer

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		buffer.Write(chunk.Data)
	}

	message, err := mail.ReadMessage(&buffer)
	if err != nil {
		return stream.SendAndClose(&pb.ExtractResponse{
			Error: err.Error(),
		})
	}

	content, contentType, err := ExtractDisplayPart(message)
	if err != nil {
		return stream.SendAndClose(&pb.ExtractResponse{
			Error: err.Error(),
		})
	}

	return stream.SendAndClose(&pb.ExtractResponse{
		Content:     content,
		ContentType: contentType,
	})
}