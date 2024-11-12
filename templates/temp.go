package templates

import (
	"embed"
)

//go:embed *.tmpl

// FS文件引用
var Fs embed.FS
