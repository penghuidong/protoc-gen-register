package plugin

import register "github.com/penghuidong/protoc-gen-register/extensions"

type RegisterBuilder struct {
	MsgOptions map[string]*register.RegisterMessageOptions
}
