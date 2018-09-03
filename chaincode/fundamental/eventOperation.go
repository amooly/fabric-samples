package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	EventTypeMedical   = "medical"
	EventTypeTransport = "transport"
)

func saveEvent(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	event := &Event{}
	if err := json.Unmarshal([]byte(args[0]), event); err != nil {
		return shim.Error("Failed to assemble event:" + err.Error())
	}

	if err := stub.PutState(event.EventNo, []byte(args[0])); err != nil {
		return shim.Error("Failed to add event:" + err.Error())
	}
	return shim.Success([]byte("Successfully add event:" + event.EventNo))
}

func queryEventByEventNo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	eventByte, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to query event:" + err.Error())
	}
	return shim.Success(eventByte)

}

func queryEventByDamaged(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect size of parameters")
	}

	eventType := args[0]
	damagedId := args[1]

	queryString := fmt.Sprintf("{\"selector\":{\"DataType\":\"%s\",\"Damaged\":\"%s\"}}", eventType, damagedId)

	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer iterator.Close()

	events := []Event{}
	for iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		var event Event
		if err := json.Unmarshal(responseRange.Value, &event); err != nil {
			return shim.Error(err.Error())
		}
		events = append(events, event)
	}

	eventsByte, err := json.Marshal(events)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(eventsByte)
}
