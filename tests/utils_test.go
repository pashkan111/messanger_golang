package tests

import (
	"messanger/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStructToMap(t *testing.T) {
	type test_struct struct {
		Param1 string
		Param2 int
		Param3 bool
	}

	test := test_struct{
		Param1: "test",
		Param2: 1,
		Param3: true,
	}

	mapping := utils.ConvertStructToMap(test)
	assert.Equal(t, mapping, map[string]interface{}{
		"Param1": "test",
		"Param2": 1,
		"Param3": true,
	})

}
