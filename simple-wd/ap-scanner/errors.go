package scanner

import "fmt"

type CommandNotFoundError struct {
	command string
}

func (e *CommandNotFoundError) Error() string {
	return fmt.Sprintf("Command '%s' not found in $PATH", e.command)
}
