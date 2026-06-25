package lib

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// registerEncodingDeclarations 注册编码相关的CEL函数声明
func registerEncodingDeclarations() []*exprpb.Decl {
	return []*exprpb.Decl{
		// base64
		decls.NewFunction("base64",
			decls.NewOverload("base64_string",
				[]*exprpb.Type{decls.String},
				decls.String)),
		decls.NewFunction("base64",
			decls.NewOverload("base64_bytes",
				[]*exprpb.Type{decls.Bytes},
				decls.String)),

		// base64Decode
		decls.NewFunction("base64Decode",
			decls.NewOverload("base64Decode_string",
				[]*exprpb.Type{decls.String},
				decls.String)),
		decls.NewFunction("base64Decode",
			decls.NewOverload("base64Decode_bytes",
				[]*exprpb.Type{decls.Bytes},
				decls.String)),

		// urlencode
		decls.NewFunction("urlencode",
			decls.NewOverload("urlencode_string",
				[]*exprpb.Type{decls.String},
				decls.String)),
		decls.NewFunction("urlencode",
			decls.NewOverload("urlencode_bytes",
				[]*exprpb.Type{decls.Bytes},
				decls.String)),

		// urldecode
		decls.NewFunction("urldecode",
			decls.NewOverload("urldecode_string",
				[]*exprpb.Type{decls.String},
				decls.String)),
		decls.NewFunction("urldecode",
			decls.NewOverload("urldecode_bytes",
				[]*exprpb.Type{decls.Bytes},
				decls.String)),

		// hexdecode
		decls.NewFunction("hexdecode",
			decls.NewInstanceOverload("hexdecode",
				[]*exprpb.Type{decls.String},
				decls.Bytes)),
	}
}

// registerEncodingImplementations 注册编码相关的CEL函数实现
func registerEncodingImplementations() []*functions.Overload {
	return []*functions.Overload{
		// base64_string
		{
			Operator: "base64_string",
			Unary: func(value ref.Val) ref.Val {
				v, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(value, "unexpected type '%v' passed to base64_string", value.Type())
				}
				return types.String(base64.StdEncoding.EncodeToString([]byte(v)))
			},
		},

		// base64_bytes
		{
			Operator: "base64_bytes",
			Unary: func(value ref.Val) ref.Val {
				v, ok := value.(types.Bytes)
				if !ok {
					return types.ValOrErr(value, "unexpected type '%v' passed to base64_bytes", value.Type())
				}
				return types.String(base64.StdEncoding.EncodeToString(v))
			},
		},

		// base64Decode_string
		{
			Operator: "base64Decode_string",
			Unary: func(value ref.Val) ref.Val {
				v, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_string", value.Type())
				}
				decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
				if err != nil {
					return types.NewErr("%v", err)
				}
				return types.String(decodeBytes)
			},
		},

		// base64Decode_bytes
		{
			Operator: "base64Decode_bytes",
			Unary: func(value ref.Val) ref.Val {
				v, ok := value.(types.Bytes)
				if !ok {
					return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_bytes", value.Type())
				}
				decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
				if err != nil {
					return types.NewErr("%v", err)
				}
				return types.String(decodeBytes)
			},
		},

		// urlencode_string
		{
			Operator: "urlencode_string",
			Unary: func(value ref.Val) ref.Val {
				v, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_string", value.Type())
				}
				return types.String(url.QueryEscape(string(v)))
			},
		},

		// urlencode_bytes
		{
			Operator: "urlencode_bytes",
			Unary: func(value ref.Val) ref.Val {
				v, ok := value.(types.Bytes)
				if !ok {
					return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_bytes", value.Type())
				}
				return types.String(url.QueryEscape(string(v)))
			},
		},

		// urldecode_string
		{
			Operator: "urldecode_string",
			Unary: func(value ref.Val) ref.Val {
				v, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_string", value.Type())
				}
				decodeString, err := url.QueryUnescape(string(v))
				if err != nil {
					return types.NewErr("%v", err)
				}
				return types.String(decodeString)
			},
		},

		// urldecode_bytes
		{
			Operator: "urldecode_bytes",
			Unary: func(value ref.Val) ref.Val {
				v, ok := value.(types.Bytes)
				if !ok {
					return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_bytes", value.Type())
				}
				decodeString, err := url.QueryUnescape(string(v))
				if err != nil {
					return types.NewErr("%v", err)
				}
				return types.String(decodeString)
			},
		},

		// hexdecode
		{
			Operator: "hexdecode",
			Unary: func(lhs ref.Val) ref.Val {
				v1, ok := lhs.(types.String)
				if !ok {
					return types.ValOrErr(lhs, "unexpected type '%v' passed to hexdecode", lhs.Type())
				}
				out, err := hex.DecodeString(string(v1))
				if err != nil {
					return types.ValOrErr(lhs, "hexdecode error: %v", err)
				}
				return types.Bytes(out)
			},
		},
	}
}
