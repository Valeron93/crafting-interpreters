package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const imports = "import \"github.com/Valeron93/crafting-interpreters/scanner\""

type typ struct {
	typename string
	fields   []string
}

func generateStruct(w io.Writer, t typ, basetype string) {

	fmt.Fprintf(w, "type %v struct {\n", t.typename)

	for _, field := range t.fields {
		fmt.Fprintf(w, "\t%v\n", field)
	}
	fmt.Fprintf(w, "}\n\n")
	generateStructMethods(w, t, basetype)
}

func generateStructMethods(w io.Writer, t typ, basetype string) {
	firstSymbol := unicode.ToLower(rune(t.typename[0]))
	fmt.Fprintf(w, "func (%c *%v) Accept(visitor %vVisitor) any {\n", firstSymbol, t.typename, basetype)
	fmt.Fprintf(w, "\treturn visitor.Visit%v(%c)\n", t.typename, firstSymbol)
	fmt.Fprintf(w, "}\n\n")
}

func generateVisitor(w io.Writer, types []typ, basetype string) {

	fmt.Fprintf(w, "type %vVisitor interface {\n", basetype)

	for _, t := range types {
		fmt.Fprintf(w, "\tVisit%v(*%v) any\n", t.typename, t.typename)
	}
	fmt.Fprintf(w, "}\n\n")

	fmt.Fprintf(w, "type %v interface {\n", basetype)
	fmt.Fprintf(w, "\tAccept(%vVisitor) any\n", basetype)
	fmt.Fprintf(w, "}\n\n")

}

func parseType(s string) (typ, error) {
	typeAndFields := strings.Split(s, ":")
	if len(typeAndFields) != 2 {
		return typ{}, fmt.Errorf("failed on '%v'", s)
	}
	typename := strings.TrimSpace(strings.Split(s, ":")[0])
	fields := strings.Split(s, ":")[1]

	actualFields := []string{}

	for field := range strings.SplitSeq(fields, ",") {
		actualFields = append(actualFields, strings.TrimSpace(field))
	}

	return typ{
		typename: typename,
		fields:   actualFields,
	}, nil
}

func generateFromAst(w io.Writer, ast []string, basetype string) error {
	types := []typ{}
	for _, s := range ast {
		if len(s) > 0 {
			if s[0] == '#' {
				fmt.Fprintf(w, "%s\n", s[1:])
			} else {
				t, err := parseType(s)
				if err != nil {
					return err
				}
				types = append(types, t)
			}
		}
	}
	generateVisitor(w, types, basetype)
	for _, t := range types {
		generateStruct(w, t, basetype)
	}
	return nil
}
func readAstFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := make([]string, 0, 10)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret, nil
}

func main() {
	outputFile := flag.String("o", "", "output go file")
	inputFile := flag.String("i", "", "input ast file")
	typename := flag.String("type", "", "Go base type name")
	flag.Parse()

	if len(*outputFile) == 0 || len(*inputFile) == 0 || len(*typename) == 0 {
		flag.Usage()
		os.Exit(64)
	}

	file, err := os.Create(*outputFile)
	defer file.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(file, "// DO NOT MODIFY!!! This file is generated from %v\n\n", *inputFile)

	asts, err := readAstFromFile(*inputFile)
	generateFromAst(file, asts, *typename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate types: %v\n", err)
	}

}
