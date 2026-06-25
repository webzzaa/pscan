package lib

import (
	"crypto/md5" //nolint:gosec // G501: MD5用于POC检测逻辑，非加密用途
	"fmt"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// registerCryptoDeclarations 注册加密相关的CEL函数声明
func registerCryptoDeclarations() []*exprpb.Decl {
	return []*exprpb.Decl{
		decls.NewFunction("md5",
			decls.NewOverload("md5_string",
				[]*exprpb.Type{decls.String},
				decls.String)),
		decls.NewFunction("shirokey",
			decls.NewOverload("shiro_key",
				[]*exprpb.Type{decls.String, decls.String},
				decls.String)),
	}
}

// registerCryptoImplementations 注册加密相关的CEL函数实现
func registerCryptoImplementations() []*functions.Overload {
	return []*functions.Overload{
		{
			Operator: "md5_string",
			Unary: func(value ref.Val) ref.Val {
				v, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(value, "unexpected type '%v' passed to md5", value.Type())
				}
				//nolint:gosec // G401: MD5用于POC检测，非加密用途
				return types.String(fmt.Sprintf("%x", md5.Sum([]byte(v))))
			},
		},
		{
			Operator: "shiro_key",
			Binary: func(key ref.Val, mode ref.Val) ref.Val {
				v1, ok := key.(types.String)
				if !ok {
					return types.ValOrErr(key, "unexpected type '%v' passed to shiro_key", key.Type())
				}
				v2, ok := mode.(types.String)
				if !ok {
					return types.ValOrErr(mode, "unexpected type '%v' passed to shiro_mode", mode.Type())
				}
				cookie := GetShrioCookie(string(v1), string(v2))
				if cookie == "" {
					return types.NewErr("%v", "key b64decode failed")
				}
				return types.String(cookie)
			},
		},
	}
}
