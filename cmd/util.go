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

	// Convert to interface{}
	// https://github.com/itchyny/gojq#usage-as-a-library
	iface := v
	if _, ok := iface.([]interface{}); ok {
		// nop
	} else if _, ok := iface.(map[string]interface{}); ok {
		// nop
	} else {
		b, err := json.Marshal(iface)
		if err != nil {
			return err
		}

		var v interface{}
		err = json.Unmarshal(b, &v)
		if err != nil {
			return err
		}

		iface = v
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
