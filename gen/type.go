package gen

import (
	"sort"

	"github.com/SVz777/z_generator/gen/parser"
)

// IGen 是生成器接口。
//
// Gen 返回生成的代码；出错时**必须返回 error**，不要把错误信息塞进返回值里
// —— 返回值会被直接写进目标源文件，把错误文本当代码写进去会生成不可编译的文件，
// 且排查时看不到任何有效上下文。
type IGen interface {
	Gen(*parser.Struct) ([]byte, error)
}

var generatorMap = map[string]IGen{}

func Register(name string, gen IGen) {
	generatorMap[name] = gen
}

func GetGenerator(name string) IGen {
	return generatorMap[name]
}

// ListTypes 返回所有已注册的生成器类型名，按字典序排序。
// 各 gen_impl 在 init 里调用 Register 完成注册，因此在 main 里
// 调用本函数时列表已就绪。
func ListTypes() []string {
	ret := make([]string, 0, len(generatorMap))
	for name := range generatorMap {
		ret = append(ret, name)
	}
	sort.Strings(ret)
	return ret
}
