package main

import (
	"testing"

	protocols "github.com/hotarublaze/gotes3mp/src/protocols"
)

func Test_relayManagerTestProcessRelayMessageSanityCheck(t *testing.T) {

	// Test case 1: Method is blank
	rmsg1 := &protocols.BaseResponse{
		Method: "",
		JobId:  "test jobid",
		Data:   map[string]string{"key": "test data"},
	}
	err := processRelayMessageSanityCheck(rmsg1)
	if err.Error() != "method cannot be blank" {
		t.Errorf("Expected error: 'processRelayMessage: method cannot be blank', got: '%s'", err.Error())
	}

	// Test case 2: Data is empty
	rmsg2 := &protocols.BaseResponse{
		Method: "test method",
		JobId:  "test jobid",
		Data:   nil,
	}
	err = processRelayMessageSanityCheck(rmsg2)
	if err.Error() != "no data provided" {
		t.Errorf("Expected error: 'no data provided', got: '%s'", err.Error())
	}

	// Test case 3: Method and data are valid
	rmsg3 := &protocols.BaseResponse{
		Method: "test method",
		JobId:  "test jobid",
		Data:   map[string]string{"key": "test data"},
	}
	err = processRelayMessageSanityCheck(rmsg3)
	if err != nil {
		t.Errorf("Expected no error, got: '%s'", err.Error())
	}

	// Test case 4: Method is valid, data is empty map
	rmsg4 := &protocols.BaseResponse{
		Method: "test method",
		JobId:  "test jobid",
		Data:   make(map[string]string),
	}
	err = processRelayMessageSanityCheck(rmsg4)
	if err.Error() != "no data provided" {
		t.Errorf("Expected error: 'no data provided', got: '%s'", err.Error())
	}

	// Test case 5: Method is valid, data has multiple key-value pairs
	rmsg5 := &protocols.BaseResponse{
		Method: "test method",
		JobId:  "test jobid",
		Data:   map[string]string{"key1": "value1", "key2": "value2"},
	}
	err = processRelayMessageSanityCheck(rmsg5)
	if err != nil {
		t.Errorf("Expected no error, got: '%s'", err.Error())
	}

	// Test case 6: Method is blank
	rmsg6 := &protocols.BaseResponse{
		Method: "Testing",
		JobId:  "",
		Data:   map[string]string{"key": "test data"},
	}
	err = processRelayMessageSanityCheck(rmsg6)
	if err.Error() != "jobid cannot be blank" {
		t.Errorf("Expected error: 'jobid cannot be blank', got: '%s'", err.Error())
	}
}
