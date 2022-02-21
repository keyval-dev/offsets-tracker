package main

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/keyval-dev/offsets-tracker/binary"
	"github.com/keyval-dev/offsets-tracker/target"
	"github.com/keyval-dev/offsets-tracker/writer"
	"log"
)

func main() {
	minimunGoVersion, err := version.NewConstraint(">= 1.12")
	if err != nil {
		log.Fatalf("error in parsing version constraint: %v\n", err)
	}

	stdLibOffsets, err := target.New("go").
		FindVersionsBy(target.GoDevFileVersionsStrategy).
		DownloadBinaryBy(target.DownloadPreCompiledBinaryFetchStrategy).
		VersionConstraint(&minimunGoVersion).
		FindOffsets([]*binary.DataMember{
			{
				StructName: "runtime.g",
				Field:      "goid",
			},
			{
				StructName: "net/http.Request",
				Field:      "Method",
			},
			{
				StructName: "net/http.Request",
				Field:      "URL",
			},
			{
				StructName: "net/url.URL",
				Field:      "Path",
			},
		})

	if err != nil {
		log.Fatalf("error while fetching offsets: %v\n", err)
	}

	grpcOffsets, err := target.New("google.golang.org/grpc").
		FindOffsets([]*binary.DataMember{
			{
				StructName: "google.golang.org/grpc/internal/transport.Stream",
				Field:      "method",
			},
			{
				StructName: "google.golang.org/grpc.ClientConn",
				Field:      "target",
			},
		})

	if err != nil {
		log.Fatalf("error while fetching offsets: %v\n", err)
	}

	fmt.Println("Done collecting offsets, writing results to file ...")
	err = writer.WriteResults(stdLibOffsets, grpcOffsets)
	if err != nil {
		log.Fatalf("error while writing results to file: %v\n", err)
	}

	fmt.Println("Done!")
}
