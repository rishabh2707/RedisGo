package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func ReadCommand(reader *bufio.Reader) (*Cmd, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	if len(line) < 3 || line[len(line)-2] != '\r' {
		return nil, fmt.Errorf("invalid command")
	}

	line = bytes.TrimSuffix(line, []byte{'\r', '\n'})

	switch line[0] {
	case '+':
		content := string(line[1:])
		return parseInlineCommand(content)
	case '*':
		return parseMultipleCommand(reader, line)
	default:
		return nil, fmt.Errorf("invalid command")
	}
}

func parseInlineCommand(content string) (*Cmd, error) {
	return &Cmd{
		Name: content,
	}, nil
}

func parseMultipleCommand(reader *bufio.Reader, line []byte) (*Cmd, error) {

	count, err := strconv.Atoi(string(bytes.TrimPrefix(line, []byte{'*'})))
	if err != nil {
		return nil, err
	}

	cmd := &Cmd{}
	cmd.Args = make([]string, 0, count)

	for i := 0; i < count; i++ {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		if len(line) < 2 {
			return nil, fmt.Errorf("invalid command")
		}

		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		charCount, err := strconv.Atoi(string(bytes.TrimPrefix(line, []byte{'$'})))
		if err != nil {
			return nil, err
		}
		data := make([]byte, charCount+2)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			return nil, err
		}
		cmd.Args = append(cmd.Args, string(data[:charCount]))
	}

	cmd.Name = cmd.Args[0]

	return cmd, nil
}
