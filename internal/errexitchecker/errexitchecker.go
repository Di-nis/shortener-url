// package errexitchecker определяет анализатор, выполняющий проверку на использование os.Exit в main-функции main-пакета.

// Зачем это нужно:
// - `os.Exit` пропускает deferred-вызовы, что может приводить к утечкам и
// некорректному завершению программы.
// - Поведение усложняет тестирование, так как тесты прерываются немедленно.
//
// Рекомендуемые варианты:
// - Вынести основную логику в функцию `Run()`, возвращающую код/ошибку, а
// затем вызвать `os.Exit` только в обёртке main.
// - При необходимости подавить диагностику локальным комментарием.
package errexitchecker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ErrExitCheckAnalyzer = &analysis.Analyzer{
	Name: "errexitcheck",
	Doc:  "check for use of os.Exit in the main function of the main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if s, ok := call.Fun.(*ast.SelectorExpr); ok {
				if s.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "using a direct call to os.Exit in the main function of the main package")
				}
			}
		}
	}
	var isMainFunction, isMainPackage bool
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				isMainPackage = true
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					isMainFunction = true
				}
			case *ast.ExprStmt: // выражение
				if isMainFunction && isMainPackage {
					expr(x)
					isMainFunction, isMainPackage = false, false
				}
			}
			return true
		})
	}
	return nil, nil
}
