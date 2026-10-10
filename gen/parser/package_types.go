package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/tools/go/packages"
)

func loadFileTypes(filename string) (*ast.File, *types.Info, error) {
	absolute, err := filepath.Abs(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("解析文件路径 %s: %w", filename, err)
	}
	if _, err := os.Stat(absolute); err != nil {
		return nil, nil, fmt.Errorf("读取源文件 %s: %w", filename, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	config := &packages.Config{
		Context: ctx,
		Dir:     filepath.Dir(absolute),
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Tests: strings.HasSuffix(absolute, "_test.go"),
	}
	pkgs, err := packages.Load(config, ".")
	if err != nil {
		return nil, nil, fmt.Errorf("加载源文件 %s 所属包: %w", filename, err)
	}
	file, info, loadErr := findFileTypes(pkgs, absolute)
	if file != nil || loadErr != nil {
		return file, info, loadErr
	}
	packageErr := missingFileError(pkgs, filename)
	if err := ctx.Err(); err != nil {
		return nil, nil, errors.Join(packageErr, err)
	}
	included, err := standaloneFileIncluded(ctx, config.Dir, absolute)
	if err != nil {
		return nil, nil, errors.Join(packageErr, err)
	}
	if !included {
		return nil, nil, packageErr
	}
	pkgs, err = packages.Load(config, absolute)
	if err != nil {
		return nil, nil, errors.Join(packageErr, fmt.Errorf("加载独立源文件 %s: %w", filename, err))
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, errors.Join(packageErr, err)
	}
	file, info, err = findFileTypes(pkgs, absolute)
	if err != nil {
		return nil, nil, errors.Join(packageErr, err)
	}
	if file != nil {
		return file, info, nil
	}
	return nil, nil, errors.Join(packageErr, fmt.Errorf("加载独立源文件: %w", missingFileError(pkgs, filename)))
}

func sameSourceFile(left, right string) bool {
	if filepath.Clean(left) == filepath.Clean(right) {
		return true
	}
	leftInfo, leftErr := os.Stat(left)
	rightInfo, rightErr := os.Stat(right)
	return leftErr == nil && rightErr == nil && os.SameFile(leftInfo, rightInfo)
}

func findFileTypes(pkgs []*packages.Package, filename string) (*ast.File, *types.Info, error) {
	for _, pkg := range pkgs {
		if pkg.Fset == nil {
			continue
		}
		for _, file := range pkg.Syntax {
			if !sameSourceFile(pkg.Fset.Position(file.Pos()).Filename, filename) {
				continue
			}
			for _, diagnostic := range pkg.Errors {
				if diagnostic.Kind == packages.ParseError {
					return nil, nil, fmt.Errorf("解析源文件 %s 所属包: %s", filename, diagnostic.Error())
				}
			}
			// 包含类型错误的包仍可提供部分有效信息；未知字段由消费方决定是否拒绝。
			return file, pkg.TypesInfo, nil
		}
	}
	return nil, nil, nil
}

func missingFileError(pkgs []*packages.Package, filename string) error {
	var diagnostics []string
	for _, pkg := range pkgs {
		for _, diagnostic := range pkg.Errors {
			diagnostics = append(diagnostics, diagnostic.Error())
		}
	}
	if len(diagnostics) > 0 {
		return fmt.Errorf("加载源文件 %s 所属包: %s", filename, strings.Join(diagnostics, "; "))
	}
	return fmt.Errorf("源文件 %s 未包含在当前构建配置的包中", filename)
}

func standaloneFileIncluded(ctx context.Context, dir, filename string) (bool, error) {
	cmd := exec.CommandContext(ctx, "go", "env", "GOMOD")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("检查源文件模块环境: %w: %s", err, output)
	}
	module := strings.TrimSpace(string(output))
	if module != "" && module != os.DevNull {
		return false, nil
	}
	cmd = exec.CommandContext(ctx, "go", "list", "-e", "-json", ".")
	cmd.Dir = dir
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(entry), "GO111MODULE=") {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	cmd.Env = append(cmd.Env, "GO111MODULE=off")
	output, err = cmd.Output()
	if err != nil {
		return false, fmt.Errorf("检查独立源文件构建约束: %w", err)
	}
	var pkg struct {
		GoFiles, CgoFiles, TestGoFiles, XTestGoFiles []string
	}
	if err := json.Unmarshal(output, &pkg); err != nil {
		return false, fmt.Errorf("解析独立源文件构建信息: %w", err)
	}
	files := append(pkg.GoFiles, pkg.CgoFiles...)
	if strings.HasSuffix(filename, "_test.go") {
		files = append(files, pkg.TestGoFiles...)
		files = append(files, pkg.XTestGoFiles...)
	}
	for _, file := range files {
		if sameSourceFile(filepath.Join(dir, file), filename) {
			return true, nil
		}
	}
	return false, nil
}
