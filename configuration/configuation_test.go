package configuration

import (
	"testing"
)

func TestIsValidStreamMapFormat(t *testing.T) {
	// <Stream Name>@<Cluster Name>:<Logical Server Name>
	var testSet = []struct {
		input          string
		expectedResult bool
	}{
		{"VAS_V2_SR_V2@RYD2:Server", true},
		{"VAS SR @RYD2:Server", false}, // Whitespaces not allowed between characters of the stream name
		{"1234@RYD2:Server 5", false},  // No whitespaces allowed in logical server name
		{"1234@RYD2: Server5", true},   // Whitespaces are allowed after the colon
		{"1234@RYD2:Server5 ", true},   // Whitespaces are allowed at the end of the mapping string
		{" 1234@RYD2:Server5", true},   // Whitespaces are allowed at the beginning of the mapping string
		{"1234$RYD2:Server5", false},   // Only @ character is allowed as separator between stream name and logical server name
	}

	for _, test := range testSet {
		if result := isValidStreamMapFormat(test.input); result != test.expectedResult {
			t.Errorf("isValidStreamMapFormat(\"%s\") = %v, should be %v", test.input, result, test.expectedResult)
		}
	}

}
