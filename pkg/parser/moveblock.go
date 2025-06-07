package parser

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/msarfaty/tuf/internal/logging"
	filestats "github.com/msarfaty/tuf/pkg/file"
	"go.uber.org/zap"
)

// the character that preceeds the moved block if file is not empty
const (
	BUFFER_CHAR     = "\n"
	BUFFER_CHAR_MAX = 2
)

var logger *zap.SugaredLogger

type MoveOptions struct {
	// the address to move
	Address string
	// the BlockDescription of the block to move
	BlockDescription *BlockDescription
	// directory to move from
	FromDirectory string
	// file to move from
	FromFile string
	// directory to move to
	ToDirectory string
	// file to move to
	ToFile string

	// the terraform files to where the source block may live
	sourceWorkspaceFiles []string
}

// validate MoveOptions and set defaults
func (mo *MoveOptions) validate() error {

	if mo.FromDirectory == "" && mo.FromFile == "" {
		return fmt.Errorf("must set source directory or file")
	}
	if mo.FromDirectory != "" && mo.FromFile != "" {
		return fmt.Errorf("cannot choose both a file and directory to move from")
	}

	if mo.ToDirectory == "" && mo.ToFile == "" {
		return fmt.Errorf("must set destination file or directory")
	}
	if mo.ToDirectory != "" && mo.ToFile != "" {
		return fmt.Errorf("cannot choose both a file and directory to move to")
	}

	if mo.Address == "" && mo.BlockDescription == nil {
		return fmt.Errorf("must include a blockdescription or resource address to filter for")
	}
	if mo.Address != "" && mo.BlockDescription != nil {
		return fmt.Errorf("cannot supply both a blockdescription and address for moving")
	}

	if mo.BlockDescription == nil {
		bd, err := New(mo.Address)
		if err != nil {
			return err
		}
		mo.BlockDescription = &bd
	}

	mo.sourceWorkspaceFiles = []string{}
	if mo.FromFile != "" {
		mo.sourceWorkspaceFiles = append(mo.sourceWorkspaceFiles, mo.FromFile)
	} else {
		dNodes, err := os.ReadDir(mo.FromDirectory)
		if err != nil {
			return fmt.Errorf("failed to open source workspace directory %s: %w", mo.FromDirectory, err)
		}
		for _, dNode := range dNodes {
			if dNode.Type().IsRegular() {
				mo.sourceWorkspaceFiles = append(mo.sourceWorkspaceFiles, dNode.Name())
			}
		}
	}

	return nil
}

// delete the HCL selection from the source file
func deleteRange(blockrange *hcl.Range) error {
	contents, err := os.ReadFile(blockrange.Filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s for cleaning up old HCL block: %w", blockrange.Filename, err)
	}
	contents = append(contents[:blockrange.Start.Byte], contents[blockrange.End.Byte:]...)

	// we want to prettify the contents that we delete somewhat
	// after we delete, we want to remove newlines so that we only preserve 2
	contents = filestats.DeleteOverNOccurrences(contents, []byte(BUFFER_CHAR), blockrange.Start.Byte, 2)

	// lastly, we want to trim the file to only have one trailing newline if we have to
	contents = filestats.DeleteOverNOccurrences(contents, []byte(BUFFER_CHAR), len(contents)-1, 1)

	// finally, we write the removed + prettified contents
	err = os.WriteFile(blockrange.Filename, contents, 0644)
	if err != nil {
		return fmt.Errorf("failed to write prettified from-file %s: %w", blockrange.Filename, err)
	}

	return nil
}

// prettify the copied HCL block selection
// for example, adding newlines where needed if not present
func prettifyCopySelection(copyBytes []byte, dest string) ([]byte, error) {
	mt, err := filestats.FileIsEmpty(dest)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file (%s): %w", dest, err)
	}
	endsWithBufferChar, err := filestats.FileEndsWith(dest, BUFFER_CHAR)
	if err != nil {
		return nil, fmt.Errorf("failed to check contents of file (%s): %w", dest, err)
	}
	if !mt && !endsWithBufferChar {
		copyBytes = append([]byte(BUFFER_CHAR), copyBytes...)
	}

	return append(copyBytes, []byte(BUFFER_CHAR)...), nil
}

// copy an HCL range to a new file
func copyRange(blockRange *hcl.Range, dest string) error {
	// read source file from range
	content, err := os.ReadFile(blockRange.Filename)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", blockRange.Filename, err)
	}
	copyBytes := content[blockRange.Start.Byte:blockRange.End.Byte]

	// open destination
	file, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open destination file %s: %w", dest, err)
	}
	defer file.Close()

	copyBytes, err = prettifyCopySelection(copyBytes, dest)
	if err != nil {
		return err
	}
	_, err = io.Writer.Write(file, copyBytes)
	if err != nil {
		return fmt.Errorf("failed to write hcl block to dest file: %w", err)
	}

	return nil
}

// move an HCL block according to the given options
func MoveHclBlock(mo *MoveOptions) error {
	mo.validate()

	p := hclparse.NewParser()

	for _, fname := range mo.sourceWorkspaceFiles {
		hclFile, diags := p.ParseHCLFile(fname)
		if diags.HasErrors() {
			return fmt.Errorf("failed to move file while parsing file %s: %s", fname, diags.Error())
		}
		body, ok := hclFile.Body.(*hclsyntax.Body)
		if !ok {
			return fmt.Errorf("error casting hcl in file=(%s) to hclsyntax", fname)
		}
		for _, block := range body.Blocks {
			if (*mo.BlockDescription).Matches(*block.AsHCLBlock()) {
				blockRange := block.Range()
				logger.Debugf("found match for address=[%s] in file %s[%d:%d]", (*mo.BlockDescription).address(), fname, blockRange.Start.Line, blockRange.Start.Column)
				err := copyRange(&blockRange, mo.ToFile)
				if err != nil {
					return fmt.Errorf("failed to copy range (%s[%d:%d]) to (%s): %w", blockRange.Filename, blockRange.Start.Byte, blockRange.End.Byte, mo.ToFile, err)
				}
				err = deleteRange(&blockRange)
				if err != nil {
					return fmt.Errorf("failed to delete range (%s[%d:%d]): %w", blockRange.Filename, blockRange.Start.Byte, blockRange.End.Byte, err)
				}
				return nil
			}
		}
	}

	return errors.New("no block was found in any file matching the block description or address")
}

func init() {
	var err error
	logger, err = logging.GetLogger()
	if err != nil {
		panic(fmt.Errorf("failed to setup logger: %w", err))
	}
}
