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
		response = c.createWallet(stub)
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
	args := stub.GetStringArgs()
	if len(args) != 2 {
		c.logger.Errorf("creat wallet action: invalid amount of arguments. want 2, got %d", len(args))
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: errInvalidArguments.Error(),
		}
	}

	bytes, err := json.Marshal(&wallet{
		Balance: defaultAmountOfPHRMTokens,
	})
	if err != nil {
		c.logger.Error("failed to serialize state to json format", err)
		return peer.Response{
			Status:  shim.ERROR,
			Message: "failed to serialize state to json format",
		}
	}

	if err := stub.PutState(string(args[1]), bytes); err != nil {
		c.logger.Error("failed to create new wallet", err)
		return peer.Response{
			Status:  shim.ERROR,
			Message: errors.Wrap(err, "failed to create wallet").Error(),
		}
	}

	return peer.Response{
		Status: shim.OK,
	}
}
