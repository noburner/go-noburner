package noburner

import (

	// standard

	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestVerify(t *testing.T) {

	fmt.Println()

	defer fmt.Println()

	nb := New(Config{

		Secret: "abc",
	})

	result, err := nb.Verify(context.Background(), "info@noburner.com")

	if err != nil {

		fmt.Println(err)

		return
	}

	data, err := json.MarshalIndent(result, "", "  ")

	if err != nil {

		fmt.Println(err)

		return
	}

	fmt.Println(string(data))
}
