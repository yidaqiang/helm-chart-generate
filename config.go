package helm_chart_generate

import "embed"

// 嵌入整个目录（递归）
//
//go:embed assets/chart-template
var templatesFS embed.FS

func GetTemplatesFS() embed.FS {
	return templatesFS
}
