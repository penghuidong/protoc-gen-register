package main

import (
	"fmt"
	"io/ioutil"
	"os"

	register "github.com/penghuidong/protoc-gen-register/extensions"
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
		fileName := file.GeneratedFilenamePrefix + ".pb.register.go"
		g := plugin.NewGeneratedFile(fileName, ".")
		if g == nil {
			continue
		}

		if !file.Generate {
			g.Skip()
			continue
		}

		fileOpts := proto.GetExtension(file.Desc.Options(), register.E_RegFileOpts)
		if fileOpts == nil {
			g.Skip()
			continue
		}

		regFileOpts, ok := fileOpts.(*register.RegisterFileOptions)
		if !ok || regFileOpts == nil {
			continue
		}

		g.P("package ", file.GoPackageName)

		g.P("type MsgRegister struct {")
		g.P("MsgIdMap map[uint32]interface{}")
		g.P("MsgNameMap map[string]interface{}")
		g.P("}")

		msgRegister := fmt.Sprintf("%sMsgRegister", regFileOpts.MsgPrefix)

		g.P("var ", msgRegister, "* MsgRegister")

		g.P("func RegistMessage(registerId uint32, registerName string, msg interface{}) {")
		g.P("if _, ok := ", msgRegister, ".MsgIdMap[registerId]; !ok {")
		g.P("msgRegister.MsgIdMap[registerId] = msg\n}")
		g.P("}")

		g.P("func init() {")
		g.P(msgRegister, " = new(MsgRegister)")
		g.P(msgRegister, ".MsgIdMap = make(map[uint32]interface{})")
		g.P(msgRegister, ".MsgNameMap = make(map[string]interface{})")

		for _, message := range file.Messages {
			msgOpts := message.Desc.Options()
			v := proto.GetExtension(msgOpts, register.E_RegMsgOpts)
			if v == nil {
				continue
			}

			regOpts, ok := v.(*register.RegisterMessageOptions)
			if !ok || regOpts == nil {
				continue
			}

			if regOpts.RegistName != "" {
				g.P("RegistMessage(", regOpts.Regist_Id, fmt.Sprintf(", \"%s\", &", regOpts.RegistName), message.GoIdent.GoName, "{})")
			} else {
				g.P("RegistMessage(", regOpts.Regist_Id, fmt.Sprintf(", \"%s\", &", message.GoIdent.GoName), message.GoIdent.GoName, "{})")
			}
		}
		g.P("}")

		g.P("\n")

		for _, message := range file.Messages {
			msgOpts := message.Desc.Options()
			v := proto.GetExtension(msgOpts, register.E_RegMsgOpts)
			if v == nil {
				continue
			}

			regOpts, ok := v.(*register.RegisterMessageOptions)
			if !ok || regOpts == nil {
				continue
			}
			g.P("func new", regFileOpts.MsgPrefix, message.GoIdent.GoName, "() interface{} {")
			g.P("return &", message.GoIdent.GoName, "{}")
			g.P("}\n")
		}

		g.P("func New", regFileOpts.MsgPrefix, "Message(msgId uint32) interface{} {")
		g.P("if msg, ok := ", msgRegister, ".MsgIdMap[mgId]; ok {")
		g.P("return msg}")
		g.P("return nil}\n")
	}

	response := plugin.Response()
	out, err := proto.Marshal(response)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, string(out))
}
