package args

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// argSet 是解析后的命令行参数集合，Args 是它的全局实例。
//
// Typ/FileName/StructName/Output/Inplace 是框架级的内置参数；
// Extra 会自动收集命令行上所有"多余参数"（即未被内置声明的参数），
// 这样新增一个 gen_impl 时不必回来改 args，直接在命令行传参即可。
//
// 例：
//
//	z_gen -f a.go -s Foo -tp validate -validate_tag check -validate_tag_sep ,
//
// 则在 validate 的 getParams 里：
//
//	args.Args.GetOr("validate_tag", "validate")   // => "check"
//	args.Args.GetOr("validate_tag_sep", "|")      // => ","
type argSet struct {
	Typ        string
	FileName   string
	StructName string
	Output     string
	Inplace    bool

	// Extra 存放所有未被内置声明的参数，key 已去掉前导横线。
	// -k=v / --k=v / -k v 三种写法都会被收集；单独的 -k 记为 "true"。
	Extra map[string]string
}

var Args argSet

var registerOnce sync.Once

// Parse 解析命令行。可重复调用，只有第一次真正解析。
func (c *argSet) Parse() *argSet {
	if flag.Parsed() {
		return c
	}
	registerOnce.Do(func() {
		c.registerFlags()
		flag.Usage = c.usage
	})
	// 先把多余参数从 os.Args 中摘出去，剩下的再交给标准 flag 包，
	// 否则 flag 包遇到未声明的参数会直接报错退出。
	rest := splitArgs(os.Args[1:], knownFlags(flag.CommandLine), &c.Extra)
	// 提前把 tp 取出来：flag 包是按顺序解析的，-h 若写在 -tp 之前，
	// 触发 usage 时 Typ 还是空的，会漏掉该生成器对应的扩展参数。
	c.Typ = peekFlagValue(os.Args[1:], "tp")
	flag.CommandLine.Parse(rest)
	return c
}

// peekFlagValue 在正式解析前先从 args 里取出指定参数的值。
// 只看 -k v / -k=v 两种写法，取第一个命中的。
func peekFlagValue(args []string, name string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if len(a) < 2 || a[0] != '-' {
			continue
		}
		k, v, ok := strings.Cut(strings.TrimLeft(a, "-"), "=")
		if ok {
			if k == name {
				return v
			}
			continue
		}
		if k == name && i+1 < len(args) && !isFlagLike(args[i+1]) {
			return args[i+1]
		}
	}
	return ""
}

// typeLister 返回所有可用的生成器类型，由主程序通过 SetTypeLister 注入。
// 用注入而非直接引用生成器注册表，是为了让 args 包不依赖具体业务包。
var typeLister func() []string

// SetTypeLister 注入可用生成器类型的提供者，必须在 Parse 之前调用。
// 注入后 -tp 的帮助信息会列出所有可选类型；不注入则只显示"生成类型"。
//
//	args.SetTypeLister(gen.ListTypes)
func SetTypeLister(f func() []string) {
	typeLister = f
}

func (c *argSet) registerFlags() {
	fs := flag.CommandLine
	fs.StringVar(&c.FileName, "filename", "", "源文件路径")
	fs.StringVar(&c.FileName, "f", "", "源文件路径(简写)")
	fs.StringVar(&c.StructName, "structname", "", "结构体名")
	fs.StringVar(&c.StructName, "s", "", "结构体名(简写)")
	fs.StringVar(&c.Output, "output", "", "输出文件路径")
	fs.StringVar(&c.Output, "o", "", "输出文件路径(简写)")
	fs.BoolVar(&c.Inplace, "inplace", false, "追加写入源文件")

	fs.StringVar(&c.Typ, "tp", "", tpUsage())
}

// tpUsage 生成 -tp 的帮助文本，把可用类型一并列出。
// flag 包会把 usage 里的 \n 自动替换成带缩进的换行，所以这里可以放心折行。
func tpUsage() string {
	if typeLister == nil {
		return "生成类型"
	}
	types := typeLister()
	if len(types) == 0 {
		return "生成类型"
	}
	return "生成类型，可选:\n" + wrapList(types, 60)
}

// wrapList 把列表按最大行宽折成多行，用 ", " 连接。
func wrapList(items []string, width int) string {
	var b strings.Builder
	line := ""
	for _, it := range items {
		switch {
		case line == "":
			line = it
		case len(line)+2+len(it) > width:
			b.WriteString(line)
			b.WriteString(",\n")
			line = it
		default:
			line += ", " + it
		}
	}
	return b.String() + line
}

// knownFlags 返回已声明的参数名集合。
func knownFlags(fs *flag.FlagSet) map[string]bool {
	// h/help 留给 flag 包自己处理，不要被收进 Extra
	m := map[string]bool{"h": true, "help": true}
	fs.VisitAll(func(f *flag.Flag) {
		m[f.Name] = true
	})
	return m
}

// splitArgs 按 flag 包的规则扫描 args，把未声明的参数抽到 extra，
// 返回清洗后的参数列表供标准 flag 包解析。
func splitArgs(args []string, known map[string]bool, extra *map[string]string) []string {
	if *extra == nil {
		*extra = map[string]string{}
	}
	rest := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		if len(a) < 2 || a[0] != '-' || a == "--" {
			// 从这里开始 flag 包也会停止解析，后面原样保留
			return append(rest, args[i:]...)
		}

		name := strings.TrimLeft(a, "-")
		if k, v, ok := strings.Cut(name, "="); ok {
			// -k=v / --k=v 形式，值一定存在
			if !known[k] {
				(*extra)[k] = v
				continue
			}
			rest = append(rest, a)
			continue
		}

		if known[name] {
			// 内置参数：连同它的值原样保留，由 flag 包自己判断类型，
			// 避免我们在这边误判 bool / string。
			rest = append(rest, a)
			if i+1 < len(args) && !isFlagLike(args[i+1]) {
				i++
				rest = append(rest, args[i])
			}
			continue
		}

		// 未知参数：-k v 形式取下一个作为值，否则当作 bool 开关
		if i+1 < len(args) && !isFlagLike(args[i+1]) {
			i++
			(*extra)[name] = args[i]
			continue
		}
		(*extra)[name] = "true"
	}
	return rest
}

// isFlagLike 判断 s 是否看起来像 flag 而非某个参数的值。
// 负数（如 -3）不算 flag，避免 -min -3 这类写法被拆错。
func isFlagLike(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}
	return s[1] < '0' || s[1] > '9'
}

// Get 取多余参数，不存在返回空串。
func (c *argSet) Get(key string) string {
	return c.Extra[key]
}

// GetOr 取多余参数，不存在或为空时返回默认值。
func (c *argSet) GetOr(key, def string) string {
	if v := c.Extra[key]; v != "" {
		return v
	}
	return def
}

// GetBoolOr 取多余参数并按 bool 解析，不存在或解析失败时返回默认值。
// 用于"默认开启、可显式关闭"的开关：-get=false。
func (c *argSet) GetBoolOr(key string, def bool) bool {
	v, ok := c.Extra[key]
	if !ok || v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

// GetBool 取多余参数并按 bool 解析，解析失败或不存在返回 false。
func (c *argSet) GetBool(key string) bool {
	return c.GetBoolOr(key, false)
}

// GetInt 取多余参数并按 int 解析，解析失败或不存在返回默认值。
func (c *argSet) GetInt(key string, def int) int {
	v, ok := c.Extra[key]
	if !ok || v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// Sub 返回带指定前缀的参数子集（key 去掉前缀），
// 便于各 gen_impl 用前缀隔离命名空间，避免参数撞名。
//
//	z_gen ... -php.pkg model -php.strict
//	args.Args.Sub("php.") // => {"pkg": "model", "strict": "true"}
func (c *argSet) Sub(prefix string) map[string]string {
	ret := map[string]string{}
	for k, v := range c.Extra {
		if strings.HasPrefix(k, prefix) {
			ret[strings.TrimPrefix(k, prefix)] = v
		}
	}
	return ret
}

// 以下是各 gen_impl 声明自己支持的扩展参数，仅用于 -h 展示。
// 声明与否都不影响 Extra 的收集，纯粹是让帮助信息更好读。

type declaredFlag struct {
	GenTyp string
	Key    string
	Usage  string
}

var (
	declareMu sync.Mutex
	declared  []declaredFlag
)

// Declare 声明某个生成器支持的扩展参数，仅用于 -h 展示，声明与否都不影响传参。
// genTyp 传该 gen_impl 的 GenType 常量，这样只有 -tp 指定它时才展示这些参数。
// 建议在 gen_impl 的 init 里调用。
func Declare(genTyp, key, usage string) {
	declareMu.Lock()
	defer declareMu.Unlock()
	declared = append(declared, declaredFlag{GenTyp: genTyp, Key: key, Usage: usage})
}

// filterDeclared 按生成器类型过滤扩展参数，genTyp 为空时返回全部。
func filterDeclared(list []declaredFlag, genTyp string) []declaredFlag {
	if genTyp == "" {
		return list
	}
	ret := make([]declaredFlag, 0, len(list))
	for _, d := range list {
		if d.GenTyp == genTyp {
			ret = append(ret, d)
		}
	}
	return ret
}

func (c *argSet) usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
	flag.CommandLine.PrintDefaults()

	declareMu.Lock()
	list := make([]declaredFlag, len(declared))
	copy(list, declared)
	declareMu.Unlock()
	if len(list) == 0 {
		return
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].GenTyp != list[j].GenTyp {
			return list[i].GenTyp < list[j].GenTyp
		}
		return list[i].Key < list[j].Key
	})

	if c.Typ == "" {
		// 没指定生成器就不展开，只提示哪些生成器有扩展参数，避免默认帮助太长
		types := declaredGenTypes(list)
		if len(types) == 0 {
			return
		}
		fmt.Fprintf(out, "\n有扩展参数的生成器: %s\n", strings.Join(types, ", "))
		fmt.Fprintf(out, "用 -tp <类型> -h 查看对应生成器的扩展参数\n")
		return
	}

	// 指定了生成器，只展示它支持的参数
	sub := filterDeclared(list, c.Typ)
	if len(sub) == 0 {
		return
	}
	fmt.Fprintf(out, "\n-tp %s 支持的扩展参数:\n", c.Typ)
	for _, d := range sub {
		fmt.Fprintf(out, "  -%s string\n    \t%s\n", d.Key, d.Usage)
	}
}

// declaredGenTypes 提取所有声明过扩展参数的生成器类型，去重且按字典序。
func declaredGenTypes(list []declaredFlag) []string {
	seen := make(map[string]bool, len(list))
	ret := make([]string, 0, len(list))
	for _, d := range list {
		if d.GenTyp == "" || seen[d.GenTyp] {
			continue
		}
		seen[d.GenTyp] = true
		ret = append(ret, d.GenTyp)
	}
	sort.Strings(ret)
	return ret
}
