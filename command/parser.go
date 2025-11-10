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
		fmt.Println("Debug1: ", err.Error())
		return nil, err
	}

	if len(line) < 3 || line[len(line)-2] != '\r' {
		fmt.Println("Debug2: ", line)
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
		fmt.Println("Debug3: ", line)
		return nil, fmt.Errorf("invalid command")
	}
}

func parseInlineCommand(content string) (*Cmd, error) {
	fmt.Println("Debug4: ", content)
	return &Cmd{
		Name: content,
	}, nil
}

func parseMultipleCommand(reader *bufio.Reader, line []byte) (*Cmd, error) {
	fmt.Println("Debug5: ", line)
	count, err := strconv.Atoi(string(bytes.TrimPrefix(line, []byte{'*'})))
	if err != nil {
		fmt.Println("Debug6: ", err.Error())
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
			fmt.Println("Debug7: ", line)
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

func (cmd *Cmd) ToRespFormat() string {
	result := "*" + strconv.Itoa(len(cmd.Args)) + "\r\n"
	for _, arg := range cmd.Args {
		result += "$" + strconv.Itoa(len(arg)) + "\r\n" + arg + "\r\n"
	}
	return result
}
