package render

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"
)

// Render 渲染模板并把结果按 Go 源码格式化，用于生成 Go 代码的生成器。
//
// 格式化失败时错误里会附上原始输出——生成的代码不合法时，只看
// "expected xxx" 这种语法错误几乎无法定位，必须能看到生成了什么。
func Render(name, tpl string, data any, opt ...Option) ([]byte, error) {
	opts := NewOptions(opt...)
	t, err := template.New(name).Funcs(opts.Funcs).Parse(tpl)
	if err != nil {
		return nil, fmt.Errorf("解析模板 %s: %w", name, err)
	}
	buf := bytes.Buffer{}
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("渲染模板 %s: %w", name, err)
	}
	out := buf.Bytes()
	if opts.Fmt {
		out, err = format.Source(out)
		if err != nil {
			return nil, fmt.Errorf("格式化 %s 的生成结果失败: %w\n--- 原始输出 ---\n%s", name, err, buf.Bytes())
		}
	}
	return out, nil
}
