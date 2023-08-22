package main

import "testing"

func Test_relayManagerTestProcessRelayMessageSanityCheck(t *testing.T) {
	// Test case 1: Method is blank
	rmsg1 := &baseresponse{
		Method: "",
		Data:   map[string]string{"key": "test data"},
	}
	err := processRelayMessageSanityCheck(rmsg1)
	if err.Error() != "processRelayMessage: method cannot be blank" {
		t.Errorf("Expected error: 'processRelayMessage: method cannot be blank', got: '%s'", err.Error())
	}

	// Test case 2: Data is empty
	rmsg2 := &baseresponse{
		Method: "test method",
		Data:   nil,
	}
	err = processRelayMessageSanityCheck(rmsg2)
	if err.Error() != "processRelayMessage: No data provided" {
		t.Errorf("Expected error: 'processRelayMessage: No data provided.', got: '%s'", err.Error())
	}

	// Test case 3: Method and data are valid
	rmsg3 := &baseresponse{
		Method: "test method",
		Data:   map[string]string{"key": "test data"},
	}
	err = processRelayMessageSanityCheck(rmsg3)
	if err != nil {
		t.Errorf("Expected no error, got: '%s'", err.Error())
	}

	// Test case 4: Method is valid, data is empty map
	rmsg4 := &baseresponse{
		Method: "test method",
		Data:   make(map[string]string),
	}
	err = processRelayMessageSanityCheck(rmsg4)
	if err.Error() != "processRelayMessage: No data provided" {
		t.Errorf("Expected error: 'processRelayMessage: No data provided.', got: '%s'", err.Error())
	}

	// Test case 5: Method is valid, data has multiple key-value pairs
	rmsg5 := &baseresponse{
		Method: "test method",
		Data:   map[string]string{"key1": "value1", "key2": "value2"},
	}
	err = processRelayMessageSanityCheck(rmsg5)
	if err != nil {
		t.Errorf("Expected no error, got: '%s'", err.Error())
	}
}
