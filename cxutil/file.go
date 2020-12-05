package cxutil

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	cxFileSuffix = ".cx"
)

// ExtractedResults contains the extracted results from parsing args for
// cx-related parts.
type ExtractedResults struct {
	CXFlags   []string // CX flags. key: flag key, value: flag value.
	CXSources []*os.File // CX source files.
}

// ListSourceNames returns a list of source filenames from an array of files.
func ListSourceNames(srcs []*os.File, stripSuffix bool) []string {
	names := make([]string, 0, len(srcs))

	for _, src := range srcs {
		name := src.Name()
		if stripSuffix {
			name = strings.TrimSuffix(name, cxFileSuffix)
		}
		names = append(names, name)
	}

	return names
}

// ExtractCXArgs extracts CX args and CX source files from CLI args.
func ExtractCXArgs(cmd *flag.FlagSet, nested bool) (*ExtractedResults, error) {
	args := cmd.Args()

	out := &ExtractedResults{
		CXFlags:   make([]string, 0, len(args)),
		CXSources: make([]*os.File, 0, len(args)),
	}

	for i, arg := range args {
		if len(arg) < 1 {
			// Just in case.
			continue
		}

		// Check for CX flag.
		if arg[0] == '+' {
			if err := parseCXFlag(arg, &out.CXFlags); err != nil {
				return nil, fmt.Errorf("failed at cli arg[%d]: %w", i, err)
			}
			continue
		}

		// Check whether arg is file/folder and exists.
		info, err := os.Stat(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file defined at arg[%d]: %w", i, err)
		}

		if info.IsDir() {
			if err := parseDir(&out.CXSources, arg, nested); err != nil {
				return nil, fmt.Errorf("failed to parse cx package at arg[%d]: %w", i, err)
			}
			continue
		}

		/* At this point, path is of a file. */

		// Skip if incorrect file suffix.
		if !strings.HasSuffix(info.Name(), cxFileSuffix) {
			continue
		}

		if err := parseFile(&out.CXSources, arg); err != nil {
			return nil, fmt.Errorf("failed to parse cx source file at arg[%d]: %w", i, err)
		}
	}

	return out, nil
}

func parseDir(files *[]*os.File, dir string, nested bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if !nested {
				return nil
			}
			return parseDir(files, path, nested)
		}

		/* At this point, path is of a file. */

		// Skip if incorrect file suffix.
		if !strings.HasSuffix(info.Name(), cxFileSuffix) {
			return nil
		}

		fmt.Println(info.Name())
		return parseFile(files, path)

	})
}

func parseFile(files *[]*os.File, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	*files = append(*files, f)
	return nil
}

// parseCXFlag checks whether the cx flag is valid and appends it to cxFlags.
func parseCXFlag(arg string, cxFlags *[]string) error {
	if len(arg) < 3 || arg[:2] != "++" {
		return fmt.Errorf("flag '%s' appears to be a cx flag, but is invalid", arg)
	}

	*cxFlags = append(*cxFlags, arg)
	return nil
}
