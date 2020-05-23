package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/go/packages"
)

const (
	cliName    = "goor"
	cliVersion = "0.1"
)

// flags
var (
	typeName = flag.String("type", "", "[required] a struct type name.")
	output   = flag.String("output", "", `[option] output file name (default "srcdir/<type>_constructor_gen.go").`)
	pointer  = flag.Bool("pointer", true, `[option] set the return value to a pointer when creating the constructor.`)
	getter   = flag.Bool("getter", false, "[option] when you create a constructor, you also create a getter.")
	setter   = flag.Bool("setter", false, "[option] when you create a constructor, you also create a setter.")
	version  = flag.Bool("version", false, "outputs the current version.")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", cliName)
	fmt.Fprintf(os.Stderr, "\t%s [flags] -type=[struct type name]\n", cliName)
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttps://github.com/ttakuya50/goor\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetPrefix(fmt.Sprintf("%s:", cliName))
	flag.Usage = Usage
	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "%s\n", cliVersion)
		os.Exit(0)
	}

	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()
	if len(args) <= 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	pkg := parsePackage(args)

	astFiles, err := parseFiles(pkg.GoFiles)
	if err != nil {
		log.Fatalf("[error]parse files err: %v", err)
	}

	structType, err := searchStructType(astFiles)
	if err != nil {
		log.Fatalf("[error]: %v", err)
	}

	fileds, err := searchFiled(structType)
	if err != nil {
		log.Fatalf("[error]search filed err: %v", err)
	}

	src, err := createConstructor(pkg, fileds)
	if err != nil {
		log.Fatalf("[error]writing output: %s", err)
	}

	if err := ioutil.WriteFile(createFileName(args), src, 0644); err != nil {
		log.Fatalf("[error] write file:%s", err)
	}
}

// searchStructType search struct type.
func searchStructType(astFiles []*ast.File) (*ast.StructType, error) {
	for _, astFile := range astFiles {
		for _, decl := range astFile.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if *typeName != typeSpec.Name.Obj.Name {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				return structType, nil
			}
		}
	}

	return nil, fmt.Errorf("struct and type names do not match")
}

// createFileName create a file name.
func createFileName(args []string) string {
	if *output != "" {
		return *output
	}

	var dir string
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		dir = filepath.Dir(args[0])
	}
	return fmt.Sprintf("%s/%s_constructor_gen.go", dir, strcase.ToSnake(*typeName))
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// parsePackage parse package.
func parsePackage(patterns []string) *packages.Package {
	cfg := &packages.Config{
		Mode: packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("[error]: %d packages found", len(pkgs))
	}

	return pkgs[0]
}

// parseFiles parse files.
func parseFiles(files []string) ([]*ast.File, error) {
	fset := token.NewFileSet()

	astFiles := make([]*ast.File, 0, len(files))
	for _, file := range files {
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		astFiles = append(astFiles, f)
	}
	return astFiles, nil
}
