package main

// 引用套件
import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"bytes"
)


// 定義合約的資產
type SampleChaincode struct {


}

type IamAsset struct {
    
}

// 定義合約初始化的行為方法
func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

// 定義合約的操作方法

// 新增資產
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// 取得資產

func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

// 批量取得資產
func getRange(stub shim.ChaincodeStubInterface, args []string) (string, error){


	//startKey := "0"
	//endKey := "999"

	resultsIterator, err := stub.GetStateByRange(args[0], args[1])
	if err != nil {
		return "", fmt.Errorf("something was wrong")
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf("result was wrong")
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer.String(),nil
}

// 移除資產

func deleteKey(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	err := stub.DelState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to delete asset: %s", args[0])
	}
	return args[0], nil

}

// 調閱歷史紀錄

func getHistory(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	valueIter, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if valueIter == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	defer  valueIter.Close()

	var value string
	for valueIter.HasNext() {
		result, err2 := valueIter.Next()
		if err2 != nil{
			return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
		}
		value +=  string(result.Value) + "||"
	}

	return string(value), nil

}

// 定義合約操作方法的調用列表
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error

	switch fn {
		case "set" :
			result, err = set(stub, args)
		case "get" :
			result, err = get(stub, args)
		case  "getHistory" :
			result, err = getHistory(stub, args)
		case  "deleteKey" :
			result, err = deleteKey(stub, args)
		case  "getRange":
			result, err = getRange(stub, args)
		default:
			return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// 啟用智能合約 
func main() {
	err := shim.Start(new(SampleChaincode))
	if err != nil {
		fmt.Println("Could not start SampleChaincode")
	} else {
		fmt.Println("SampleChaincode successfully started")
	}
}
