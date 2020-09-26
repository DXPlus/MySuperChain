package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"io"
	"strings"
)

const (
	RootPrivateKey   = "superchain-root-private-key"
	RootCertificate  = "superchain-root-certificate"
)

type SuperChain struct {}

type Chain struct {
	ID            string `json:"id"`             //Chain ID
	INFO          string `json:"info"`           //Chain info
	IP            string `json:"ip"`             //Chain-PeerAPP-IP
	SERIAL        string `json:"serial"`         //Registration time
}

type ReturnToRegister struct {
	ID            string `json:"id"`             //Chain ID
	CERT          string `json:"cert"`           //Chain-PeerAPP-Cert
	ROOTCERT      string `json:"root_cert"`      //Superchain RootCert
}

func (sc *SuperChain) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return sc.initialize(stub)
}

func (sc *SuperChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("invoke: %s\n", function)
	switch function {
	case "chainRegister":
		return sc.chainRegister(stub, args)
	case "getChainInfo":
		return sc.getChainInfo(stub, args)
	case "deleteChain":
		return sc.deleteChain(stub, args)
	case "setChainOrgCACert":
		return sc.setChainOrgCACert(stub, args)
	case "getChainOrgCACert":
		return sc.getChainOrgCACert(stub, args)
	case "updateOrgCACert":
		return sc.updateOrgCACert(stub, args)

	default:
		return shim.Error("invalid function: " + function + ", args: " + strings.Join(args, ","))
	}
}

// init
func (sc *SuperChain) initialize(stub shim.ChaincodeStubInterface) pb.Response{

	//// Generate PrivateKey and RootCert
	//rootCertificateString,rootPrivateKeyByteString := sc.GeneratePrivateKeyAndRootCert()
	//// save rootCertificate and rootPrivateKey in superchain
	//if err := stub.PutState(RootCertificate, []byte(rootCertificateString)); err != nil {
	//	return shim.Error(fmt.Errorf("save rootCertificate error: %w", err).Error())
	//}
	//
	//if err := stub.PutState(RootPrivateKey, []byte(rootPrivateKeyByteString)); err != nil {
	//	return shim.Error(fmt.Errorf("save privateKey error: %w", err).Error())
	//}

	return shim.Success(nil)
}

// register chain to get cert and save chain info
func (sc *SuperChain) chainRegister(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	info := args[0]
	ip := args[1]
	serial := args[2]
	csr := args[3]
	orgCACert := args[4]

	// create cert
	cert,err := sc.CreateCertWithCsr(stub, []byte(csr))
	if err != nil{
		return shim.Error(fmt.Errorf("create cert with csr error: %w", err).Error())
	}

	// create ID
	tempString := info + ip + serial
	Sha1Inst := sha1.New()
	io.WriteString(Sha1Inst,tempString)
	chainID := fmt.Sprintf("%x",Sha1Inst.Sum(nil))

	// create chain struct
	chain := Chain{
		ID:            chainID,
		INFO:          info,
		IP:            ip,
		SERIAL:        serial,
	}
	chainJson, err := json.Marshal(chain)
	if err != nil {
		return shim.Error(fmt.Errorf("chain json marshal error: %w", err).Error())
	}

	// save chain
	if err := stub.PutState(chain.ID, chainJson); err != nil {
		return shim.Error(fmt.Errorf("save chain error: %w").Error())
	}

	// save chain's org ca cert
	if err := stub.PutState(ToChainOrgCertID(chain.ID), []byte(orgCACert)); err != nil {
		return shim.Error(fmt.Errorf("save chain's org ca cert error: %w", err).Error())
	}

	// get root certificate
	rootCertificateByte, err := stub.GetState(RootCertificate)
	if err != nil {
		return shim.Error(fmt.Errorf("get root certificate error: %w", err).Error())
	}

	// create return struct
	rtr := ReturnToRegister{
		ID:            chain.ID,
		CERT:          string(cert),
		ROOTCERT:      string(rootCertificateByte),
	}
	rtrJson, err := json.Marshal(rtr)
	if err != nil {
		return shim.Error(fmt.Errorf("rtr json marshal error: %w", err).Error())
	}

	return shim.Success(rtrJson)
}

// get ChainInfo from chainID
func (sc *SuperChain) getChainInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	chainID := args[0]

	chainInfo , err := stub.GetState(chainID)
	if err != nil {
		return shim.Error(fmt.Errorf("get chain from chainID error: %w", err).Error())
	}

	return shim.Success(chainInfo)
}

// delete chain from chainID
func (sc *SuperChain) deleteChain(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	chainID := args[0]

	// delete chain info
	err := stub.DelState(chainID)
	if err != nil {
		return shim.Error(fmt.Errorf("delete chain from chainID error: %w", err).Error())
	}

	// delete chain org ca cert info
	err1 := stub.DelState(ToChainOrgCertID(chainID))
	if err1 != nil {
		return shim.Error(fmt.Errorf("delete chain org ca cert from chainID error: %w", err1).Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SuperChain))
	if err != nil {
		fmt.Printf("Error starting StudentChainCode: %s", err)
	}
}



