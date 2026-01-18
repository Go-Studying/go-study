package animals

import (
	"errors"
	"fmt"
	"regexp"
)

type Animal string

func (animal Animal) BackupAnimal() Animal {
	return "alternate"
}

var (
	seen = make(map[Animal]bool)

	//1. Error Structure를 잘 짜는 것이 중요
	//2. 캡슐화의 한 방식으로 canonical codes를 사용하는 방법도 있음
	//https://pkg.go.dev/google.golang.org/grpc/codes
	ErrDuplicate = errors.New("duplicate")

	ErrMarsupial = errors.New("marsupials are not supported")
)

func marsupial(animal Animal) bool {
	if animal == "marsupial" {
		return true
	}
	return false
}

func process(animal Animal) error {
	switch {
	case seen[animal]:
		return ErrDuplicate
	case marsupial(animal):
		return ErrMarsupial
	}
	seen[animal] = true
	return nil
}

// Good:
// The caller can simply compare the returned error value of the function with one of the known error values:
func handlePet(animal Animal) error {
	switch err := process(animal); err {
	case ErrDuplicate:
		return fmt.Errorf("feed %q: %v", animal, err)
	case ErrMarsupial:
		// try to recover with a friend instead.
		alternate := animal.BackupAnimal()
		return handlePet(alternate)
	}
	return nil
}

// Good:
// 위의 경우로도 대부분 충분하나,
// If process returns wrapped errors (discussed below), you can use errors.Is.
func handlePetWithWrappedErrors(animal Animal) error {
	switch err := process(animal); {
	case errors.Is(err, ErrDuplicate):
		return fmt.Errorf("feed %q: %v", animal, err)
	case errors.Is(err, ErrMarsupial):
		// try to recover with a friend instead.
		alternate := animal.BackupAnimal()
		return handlePet(alternate)
	}
	return nil
}

// Bad:
// Do not attempt to distinguish errors based on their string form.
func handlePetWithString(animal Animal) error {
	err := process(animal)
	var matched bool

	matched, err = regexp.MatchString(`duplicate`, err.Error())
	if matched {
		//
	}
	matched, err = regexp.MatchString(`marsupial`, err.Error())
	if matched {
		//
	}
	return nil
}
