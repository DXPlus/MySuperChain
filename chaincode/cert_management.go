package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Org struct {
	Name          string `json:"name"` // org name
	OrgCACert     string `json:"cert"` // org ca cert
}

// set chain org ca cert
func (sc *SuperChain) setRootCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	rootCertificateString := args[0]

	if err := stub.PutState(RootCertificate, []byte(rootCertificateString)); err != nil {
		return shim.Error(fmt.Errorf("save root certificate error: %w", err).Error())
	}

	return shim.Success(nil)
}

// set chain org ca cert
func (sc *SuperChain) setRootPrivateKey(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	rootPrivateKeyByteString := args[0]

	if err := stub.PutState(RootPrivateKey, []byte(rootPrivateKeyByteString)); err != nil {
		return shim.Error(fmt.Errorf("save root privateKey error: %w", err).Error())
	}

	return shim.Success(nil)
}

// set chain org ca cert
func (sc *SuperChain) setChainOrgCACert(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	chainID:= args[0]
	orgCACert:= args[1]

	// save chain's org ca cert
	if err := stub.PutState(ToChainOrgCertID(chainID), []byte(orgCACert)); err != nil {
		return shim.Error(fmt.Errorf("set chain error: %w", err).Error())
	}

	return shim.Success(nil)
}

// get chain org ca cert
func (sc *SuperChain) getChainOrgCACert(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	chainID:= args[0]

	chainOrgCACert, err := stub.GetState(ToChainOrgCertID(chainID))
	if err != nil {
		return shim.Error(fmt.Errorf("get chain org ca cert from chainID error: %w", err).Error())
	}

	return shim.Success(chainOrgCACert)
}

// [{"name":"ORG1MSP","cert":"====fawe8few8jfajef9aeffase9fkgkcae98f9643===="},{"name":"ORG2MSP","cert":"====45we8few8jfajef9aeffase9fkgkcae98f9666===="}]
// get chain's org ca cert
func (sc *SuperChain) updateOrgCACert(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	chainID:= args[0]
	orgName:= args[1]
	newCert:= args[2]

	// 获取链的org列表
	chainOrgCACert, err := stub.GetState(ToChainOrgCertID(chainID))
	if err != nil {
		return shim.Error(fmt.Errorf("get chain org ca cert from chainID error: %w", err).Error())
	}

	var orgs []Org
	err1 := json.Unmarshal([]byte(chainOrgCACert), &orgs)
	if err1 != nil {
		return shim.Error(fmt.Errorf("get org ca json error: %w", err1).Error())
	}

	// find org and modify cert
	for i := 0; i <= len(orgs); i++ {
		if orgs[i].Name == orgName{
			orgs[i].OrgCACert = newCert
		}
	}

	orgsJson, err := json.Marshal(orgs)
	if err != nil {
		return shim.Error(fmt.Errorf("marshal org json chain error: %w", err).Error())
	}
	// save chain's org ca cert
	if err := stub.PutState(ToChainOrgCertID(chainID), orgsJson); err != nil {
		return shim.Error(fmt.Errorf("save chain error: %w", err).Error())
	}

	return shim.Success(nil)
}

func ToChainOrgCertID(chainID string) string {
	return "CHAINORGS" + "-" + chainID
}