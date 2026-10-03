package commentgap

import (
	"go/build/constraint"
	"go/parser"
	"go/scanner"
	"go/token"
	"slices"
	"strconv"
	"strings"
)

type lexeme struct {
	kind    token.Token
	literal string
}

func scan(source []byte, mode scanner.Mode) ([]lexeme, error) {
	var s scanner.Scanner
	s.Init(token.NewFileSet().AddFile("", -1, len(source)), source, nil, mode)
	var result []lexeme
	for {
		_, kind, literal := s.Scan()
		if kind == token.EOF {
			break
		}
		result = append(result, lexeme{kind, literal})
	}
	if s.ErrorCount != 0 {
		return nil, ErrScan
	}
	return result, nil
}

func proveGo(path string, before, after []byte) error {
	left, leftErr := scan(before, 0)
	right, rightErr := scan(after, 0)
	if leftErr != nil || rightErr != nil {
		return ErrScan
	}
	if !slices.Equal(left, right) {
		return ErrTokens
	}
	var directives [2][]string
	var output bool
	for i, source := range [][]byte{before, after} {
		comments, err := scan(source, scanner.ScanComments)
		if err != nil {
			return err
		}
		for _, comment := range comments {
			if comment.kind != token.COMMENT {
				continue
			}
			if directive(comment.literal) {
				directives[i] = append(directives[i], comment.literal)
			}
			output = output || exampleOutput(comment.literal)
		}
	}
	if !slices.Equal(directives[0], directives[1]) {
		return ErrDirective
	}
	for _, source := range [][]byte{before, after} {
		file, err := parser.ParseFile(token.NewFileSet(), path, source, parser.ImportsOnly)
		if err != nil {
			return ErrScan
		}
		for _, item := range file.Imports {
			imported, err := strconv.Unquote(item.Path.Value)
			if err != nil {
				return ErrScan
			}
			if imported == "C" {
				return ErrCgo
			}
		}
	}
	if strings.HasSuffix(path, "_test.go") && output {
		return ErrExampleOutput
	}
	return nil
}

func directive(comment string) bool {
	return strings.HasPrefix(comment, "//") && len(comment) > 2 && comment[2] != ' ' && comment[2] != '\t' ||
		strings.HasPrefix(comment, "/*line ") || constraint.IsPlusBuild(comment)
}

func exampleOutput(comment string) bool {
	text := strings.TrimPrefix(comment, "//")
	if strings.HasPrefix(comment, "/*") {
		text = strings.TrimSuffix(comment[2:], "*/")
	}
	text = strings.ToLower(strings.TrimLeft(text, " \t"))
	return strings.HasPrefix(text, "output:") || strings.HasPrefix(text, "unordered output:")
}
