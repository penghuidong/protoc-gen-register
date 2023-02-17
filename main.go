package main

import (
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var request pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(input, &request)
	if err != nil {
		panic(err)
	}
	opts := protogen.Options{}
	plugin, err := opts.New(&request)
	if err != nil {
		panic(err)
	}

	for _, file := range plugin.Files {
		for _, message := range file.Messages {
			msgOpts := message.Desc.Options()

			log.Printf("extension %v")
		}
	}
}
