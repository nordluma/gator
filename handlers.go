package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the `login` handler expects a single argument, the username")
	}

	user := cmd.args[0]
	s.config.SetUser(user)

	fmt.Printf("%s has been set as user\n", user)

	return nil
}
