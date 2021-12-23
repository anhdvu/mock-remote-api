package main

import "sync"

type application struct {
	mu            sync.Mutex
	compTerminal  string
	compPassword  string
	otherTerminal string
	otherPassword string
}

func (app *application) updateCompTerminal(terminal string) error {
	return nil
}

func (app *application) updateCompPassword(terminal string) error {
	return nil
}

func (app *application) updateOtherTerminal(terminal string) error {
	return nil
}

func (app *application) updateOtherPassword(terminal string) error {
	return nil
}
