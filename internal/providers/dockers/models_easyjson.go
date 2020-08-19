// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dockers

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeArwosServerInternalProvidersDockers(in *jlexer.Lexer, out *DockerMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "stream":
			out.Stream = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "error":
			out.Error = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeArwosServerInternalProvidersDockers(out *jwriter.Writer, in DockerMessage) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Stream != "" {
		const prefix string = ",\"stream\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Stream))
	}
	if in.Status != "" {
		const prefix string = ",\"status\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Status))
	}
	if in.Error != "" {
		const prefix string = ",\"error\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Error))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v DockerMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeArwosServerInternalProvidersDockers(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v DockerMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeArwosServerInternalProvidersDockers(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *DockerMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeArwosServerInternalProvidersDockers(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *DockerMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeArwosServerInternalProvidersDockers(l, v)
}
