/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package simplepvtdata

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/msp"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

func marshalOrPanic(pb proto.Message) []byte {
	bytes, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}

	return bytes
}

// SimplePrivateDataCC example Chaincode implementation
type SimplePrivateDataCC struct {
}

// Init initializes chaincode
// ===========================
func (t *SimplePrivateDataCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimplePrivateDataCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	switch function {
	case "put":
		for i := 0; i < len(args); i = i + 3 {
			err := stub.PutPrivateData(args[i], args[i+1], []byte(args[i+2]))
			if err != nil {
				return shim.Error(err.Error())
			}
		}

		return shim.Success([]byte{})

	case "get":
		data, err := stub.GetPrivateData(args[0], args[1])
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(data)
	case "putandset":
		err := stub.PutPrivateData(args[0], args[1], []byte(args[2]))
		if err != nil {
			return shim.Error(err.Error())
		}

		noo := &common.SignaturePolicy_NOutOf{
			N: int32(len(args[3:])),
		}
		pt := &common.SignaturePolicyEnvelope{
			Rule: &common.SignaturePolicy{
				Type: &common.SignaturePolicy_NOutOf_{
					NOutOf: noo,
				},
			},
		}

		for i, org := range args[3:] {
			pt.Identities = append(pt.Identities, &msp.MSPPrincipal{
				PrincipalClassification: msp.MSPPrincipal_ROLE,
				Principal: marshalOrPanic(&msp.MSPRole{
					Role:          msp.MSPRole_MEMBER,
					MspIdentifier: org,
				}),
			})
			noo.Rules = append(noo.Rules, &common.SignaturePolicy{
				Type: &common.SignaturePolicy_SignedBy{
					SignedBy: int32(i),
				},
			})
		}

		err = stub.SetPrivateDataValidationParameter(args[0], args[1], marshalOrPanic(pt))
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)

	case "metaset":
		noo := &common.SignaturePolicy_NOutOf{
			N: int32(len(args[2:])),
		}
		pt := &common.SignaturePolicyEnvelope{
			Rule: &common.SignaturePolicy{
				Type: &common.SignaturePolicy_NOutOf_{
					NOutOf: noo,
				},
			},
		}

		for i, org := range args[2:] {
			pt.Identities = append(pt.Identities, &msp.MSPPrincipal{
				PrincipalClassification: msp.MSPPrincipal_ROLE,
				Principal: marshalOrPanic(&msp.MSPRole{
					Role:          msp.MSPRole_MEMBER,
					MspIdentifier: org,
				}),
			})
			noo.Rules = append(noo.Rules, &common.SignaturePolicy{
				Type: &common.SignaturePolicy_SignedBy{
					SignedBy: int32(i),
				},
			})
		}

		err := stub.SetPrivateDataValidationParameter(args[0], args[1], marshalOrPanic(pt))
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	case "metaget":
		meta, err := stub.GetPrivateDataValidationParameter(args[0], args[1])
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(meta)
	default:
		//error
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}
