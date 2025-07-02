package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

type Method struct {
	Name    string
	Params  []string
	Results []string
}

type Interface struct {
	Name    string
	Methods []Method
}

func main() {
	filePath := "internal/campaign/repository.go"

	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filePath, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("Erro ao analisar o arquivo:", err)
		return
	}

	fmt.Printf("Interfaces encontradas em %s:\n", filePath)

	var myInterface Interface
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		myInterface = findInterfaceAndMethods(genDecl, fs)
		break
	}

}

func findInterfaceAndMethods(genDecl *ast.GenDecl, fs *token.FileSet) Interface {
	myInterface := Interface{
		Methods: []Method{},
	}

	for _, spec := range genDecl.Specs {
		typeSpec := spec.(*ast.TypeSpec)
		if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
			fmt.Printf("- %s (na linha %d)\n", typeSpec.Name.Name, fs.Position(typeSpec.Pos()).Line)
			myInterface.Name = typeSpec.Name.Name

			// Listar métodos da interface
			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue // pode ser um mét\odo embutido, ignoramos por simplicidade
				}
				myMethods := Method{
					Name: method.Names[0].Name,
				}
				funcType, ok := method.Type.(*ast.FuncType)
				if !ok {
					continue
				}

				var params []string
				if funcType.Params != nil {
					for _, param := range funcType.Params.List {
						typeStr := fmt.Sprintf("%s", exprToString(param.Type))
						for range param.Names {
							params = append(params, typeStr)
						}
						// Caso o parâmetro não tenha nome, ainda adicionamos o tipo
						if len(param.Names) == 0 {
							params = append(params, typeStr)
						}
					}
				}

				// 4. Resultados
				var results []string
				if funcType.Results != nil {
					for _, result := range funcType.Results.List {
						typeStr := fmt.Sprintf("%s", exprToString(result.Type))
						for range result.Names {
							results = append(results, typeStr)
						}
						if len(result.Names) == 0 {
							results = append(results, typeStr)
						}
					}
				}

				myMethods.Params = params
				myMethods.Results = results
				myInterface.Methods = append(myInterface.Methods, myMethods)
			}

		}
	}
	return myInterface
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	case *ast.ArrayType:
		return "[]" + exprToString(e.Elt)
	case *ast.MapType:
		return "map[" + exprToString(e.Key) + "]" + exprToString(e.Value)
	case *ast.FuncType:
		return "func" // simplificado
	default:
		return fmt.Sprintf("%T", expr)
	}
}
