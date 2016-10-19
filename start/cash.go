package main

import (
	"errors"
	"fmt"
	"strconv"
	//"reflect"
	//"unsafe"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Chaincode example simple Chaincode implementation
type Chaincode struct {
}

// args[0]	salerParty
// args[1]	salerPartycashAmount
// args[2]	buyinParty
// args[3]	buyinPartycashAmount

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	if len(args) != 4 {
		return nil, errors.New("Cash Init Expecting 2 number of arguments.")
	}

	salerParty := args[0]
	salerPartycashAmount := args[1]

	buyinParty := args[2]
	buyinPartycashAmount := args[3]

	err = stub.PutState(salerParty, []byte(salerPartycashAmount))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(buyinParty, []byte(buyinPartycashAmount))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// args[0]	partySrc
// args[1]	partyDst
// args[2]	X, transferring cash amount
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	var partySrc, partyDst string

	if len(args) != 3 {
		return nil, errors.New("Cash Invoke Expecting 3 number of arguments.")
	}

	partySrc = args[0]
	partyDst = args[1]
	X, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		return nil, err
	}

	// Get the state from the ledger
	amountSrc_Byte, err := stub.GetState(partySrc)
	if err != nil {
		return nil, err
	}
	if amountSrc_Byte == nil {
		return nil, errors.New("amountSrc not found.\n")
	}
	amountSrc_Str := string(amountSrc_Byte)
	amountSrc_Uint64, err := strconv.ParseUint(amountSrc_Str, 10, 64)

	amountDst_Byte, err := stub.GetState(partyDst)
	if err != nil {
		return nil, err
	}
	if amountDst_Byte == nil {
		return nil, errors.New("amountDst not found.\n")
	}
	amountDst_Str := string(amountDst_Byte)
	if err != nil {
		return nil, err
	}
	amountDst_Uint64, err := strconv.ParseUint(amountDst_Str, 10, 64)

	fmt.Printf("after cash GetState, amountSrc = %s, amountDst = %s \n", amountSrc_Str, amountDst_Str)

	// Perform the execution
	amountSrc_Uint64 = amountSrc_Uint64 - X
	amountDst_Uint64 = amountDst_Uint64 + X

	// Write the state back to the ledger
	err = stub.PutState(partySrc, []byte(strconv.FormatUint(amountSrc_Uint64, 10)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(partyDst, []byte(strconv.FormatUint(amountDst_Uint64, 10)))
	if err != nil {
		return nil, err
	}

	fmt.Printf("after cash transfer, amountSrc = %d, amountDst = %d \n", amountSrc_Uint64, amountDst_Uint64)

	return nil, nil

}

// args[0]	party
// return	{"party":"","cashAmount":""}
func (t *Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	var party string

	if len(args) != 1 {
		return nil, errors.New("Cash Query Expecting 1 number of arguments.")
	}
	party = args[0]
	cashAmount, err := stub.GetState(party)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + party + "\"}"
		return nil, errors.New(jsonResp)
	}
	if cashAmount == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + party + "\"}"
		return nil, errors.New(jsonResp)
	}
	// {"party":"","cashAmount":""}
	jsonResp := "{\"party\":\"" + party + "\",\"cashAmount\":\"" + string(cashAmount) + "\"}"

	fmt.Printf("Query Response:%s\n", jsonResp)
	return []byte(jsonResp), nil

}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting billTransfer chaincode: %s", err)
	}
}

//for test
//#register:
//CORE_CHAINCODE_ID_NAME=billTransfer CORE_PEER_ADDRESS=0.0.0.0:30303 /opt/gopath/src/github.com/hyperledger/fabric/BC_NXY/chaincode/billTransfer/billTransfer
//
//#deploy
//./peer chaincode deploy -n billTransfer -c '{"Function":"Init","Args": ["SalerParty", "100000000"]}'
//./peer chaincode deploy -n billTransfer -c '{"Function":"Init","Args": ["BuyinParty", "100000000"]}'
//
//#query
//./peer chaincode query -n billTransfer -c '{"Function":"Query","Args": ["SalerParty"]}'
//./peer chaincode query -n billTransfer -c '{"Function":"Query","Args": ["BuyinParty"]}'
//
//#invoke
//./peer chaincode invoke -n billTransfer -c '{"Function":"Invoke","Args": ["SalerParty","BuyinParty","1000"]}'
//
//#deploy 只能运行一次在 register之后，在Invoke 之前。
