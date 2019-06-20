// Copyright 2019 Pharmeum
// Developed by Highchain Software https://highchain.io/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ------------------------------------------------------------------------------
package payment

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/hyperledger/fabric/protos/peer"
)

const defaultAmountOfPHRMTokens = "100.00"

//Chaincode payment chaincode implementation
type Chaincode struct {
	logger *shim.ChaincodeLogger
}

//Init do initialization of chaincode structure
func (c *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	c.logger = shim.NewLogger("payment-chaincode")
	return shim.Success([]byte("successfully initialized chaincode"))
}

//Invoke actions defined in chaincode
func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Retrieve the requested chaincode function and arguments
	function, _ := stub.GetFunctionAndParameters()

	c.logger.Info("invocation request started, method called = ", function)
	var response peer.Response

	startedTime := time.Now()
	defer func(funcName string) {
		c.logger.Infof("invocation request finished %s, status %d, duration %s", function, response.Status, time.Since(startedTime).String())
	}(function)

	switch function {
	case createWallet:
		c.createWallet(stub)
	default:
		c.logger.Debugf("no such function %s, invocation rejected", function)
		response = peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: fmt.Sprintf("no such function %s, invocation rejected", function),
		}

		return response
	}

	return response
}

func (c *Chaincode) createWallet(stub shim.ChaincodeStubInterface) peer.Response {
	x509, _ := cid.GetX509Certificate(stub)
	if x509.Subject.CommonName != adminIdentity {
		c.logger.Debug("failed to create wallet, only admin can create new wallet. Current identity = ", x509.Subject.CommonName)
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "only admin can create new wallet",
		}
	}

	// args[0] - public key
	args, err := stub.GetArgsSlice()
	if err != nil {
		c.logger.Debug("failed to read arguments from stub", err)
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: err.Error(),
		}
	}

	if len(args) != 1 {
		c.logger.Debug("invalid amount of arguments")
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: errInvalidArguments.Error(),
		}
	}

	ok, err := ValidPublicKey(args)
	if err != nil {
		c.logger.Debug("public key validation failed", err)
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: errors.Wrap(err, "public key validation failed").Error(),
		}
	}

	if !ok {
		c.logger.Debug(errInvalidPublicKey.Error())
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: errInvalidPublicKey.Error(),
		}
	}

	bytes, err := json.Marshal(&wallet{
		Balance: defaultAmountOfPHRMTokens,
	})

	if err := stub.PutState(string(args[1]), bytes); err != nil {
		return peer.Response{
			Status:  shim.ERROR,
			Message: errors.Wrap(err, "failed to create wallet").Error(),
		}
	}

	return peer.Response{
		Status: shim.OK,
	}
}
