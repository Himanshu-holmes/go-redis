package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	commandSet = "SET"
)
type Command interface {
}
type SetCommand struct {
	key, value string
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read %s\n", v.Type())
	
		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case commandSet:
					fmt.Println("length", len(v.Array()))
					if len(v.Array()) != 3 {
						return nil, fmt.Errorf("invalid number of variables for set command")
					}
					cmd := SetCommand{
						key:   v.Array()[1].String(),
						value: v.Array()[2].String(),
					}
					return cmd, nil

				default:
					return nil, fmt.Errorf("unknown command %s", v.String())
				}
				if v.String() == commandSet {
					panic("works")
				}
				fmt.Printf("%v\n", v)
			}
		}
	}
	return "boo",nil
}
