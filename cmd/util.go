package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

func printResult(v interface{}, jqQuery *gojq.Query) error {
	if jqQuery == nil {
		return printJSON(v)
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var iface interface{}
	err = json.Unmarshal(b, &iface)
	if err != nil {
		return err
	}

	iter := jqQuery.Run(iface)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok && err != nil {
			return err
		}

		err := printJSON(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func printJSON(v interface{}) error {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
