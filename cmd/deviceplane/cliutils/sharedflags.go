package cliutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

type HasCommand interface {
	Command(name string, help string) *kingpin.CmdClause
}

func GlobalAndCategorizedCmd(globalApp *kingpin.Application, categoryCmd *kingpin.CmdClause, do func(HasCommand)) {
	for _, attachmentPoint := range []HasCommand{globalApp, categoryCmd} {
		do(attachmentPoint)
	}
}

const (
	FormatTable      string = "table"
	FormatJSON       string = "json"
	FormatJSONStream string = "json-stream"
	FormatYAML       string = "yaml"
)

func AddFormatFlag(formatVar *string, categoryCmd *kingpin.CmdClause, allowedFormats ...string) {
	fFlag := categoryCmd.Flag("output", fmt.Sprintf("Output format to use. (%s)", strings.Join(allowedFormats, ", ")))
	fFlag.Short('o')
	fFlag.Default(allowedFormats[0])
	fFlag.EnumVar(formatVar, allowedFormats...)
}

func PrintWithFormat(obj interface{}, format string) error {
	switch format {
	case FormatJSONStream:
		if reflect.TypeOf(obj).Kind() != reflect.Slice {
			return errors.New("obj type is not an array")
		}

		s := reflect.ValueOf(obj)

		for i := 0; i < s.Len(); i++ {
			bytes, err := json.Marshal(s.Index(i).Interface())
			if err != nil {
				return err
			}
			fmt.Println(string(bytes))
		}
		return nil

	case FormatJSON:
		bytes, err := json.MarshalIndent(obj, "", strings.Repeat(" ", 4))
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
		return nil

	case FormatYAML:
		bytes, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
		return nil
	}
	return fmt.Errorf("format (%s) not supported", format)
}
