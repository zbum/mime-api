package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "mime_api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const chunkSize = 64 * 1024 // 64KB chunks

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: streamtest <eml_file>")
	}

	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDisplayPartServiceClient(conn)

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	stream, err := client.ExtractDisplayPartStream(context.Background())
	if err != nil {
		log.Fatalf("failed to create stream: %v", err)
	}

	buf := make([]byte, chunkSize)
	chunkCount := 0

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read file: %v", err)
		}

		chunk := &pb.FileChunk{Data: buf[:n]}
		if err := stream.Send(chunk); err != nil {
			log.Fatalf("failed to send chunk: %v", err)
		}
		chunkCount++
		fmt.Printf("Sent chunk %d (%d bytes)\n", chunkCount, n)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
	}

	fmt.Println("\n--- Response ---")
	if resp.Error != "" {
		fmt.Printf("Error: %s\n", resp.Error)
	} else {
		fmt.Printf("Content-Type: %s\n", resp.ContentType)
		fmt.Printf("Content:\n%s\n", resp.Content)
	}
}