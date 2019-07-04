package payment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestChaincode_createWallet(t *testing.T) {
	logger := shim.NewLogger("payment-chaincode")
	logger.SetLevel(shim.LogDebug)
	stub := shim.NewMockStub("create_wallet", &Chaincode{
		logger: logger,
	})
	if stub == nil {
		t.Fatal("failed to create mock stub")
	}

	t.Run("Successfully create wallet", func(t *testing.T) {
		walletAddress := "GAVG7MMFCYAWZHISVBRLBFKIN5UP2EL42KQZLCBA3SHX2UWRCRJ4LCHC"
		result := stub.MockInvoke(createWallet, [][]byte{[]byte(createWallet), []byte(walletAddress)})
		if result.Status != shim.OK {
			t.Fail()
		}

		state, err := stub.GetState(walletAddress)
		if err != nil {
			t.Fail()
		}

		if string(state) != `{"balance":"100.00"}` {
			t.Fail()
		}
	})

	t.Run("Ivalid amount of args", func(t *testing.T) {
		walletAddress := "GAVG7MMFCYAWZHISVBRLBFKIN5UP2EL42KQZLCBA3SHX2UWRCRJ4LCHC"
		result := stub.MockInvoke(createWallet, [][]byte{[]byte(createWallet), []byte(walletAddress), []byte("invalid third arg")})
		if result.Status != shim.ERRORTHRESHOLD {
			t.Fail()
		}
	})
}

func TestChaincode_transferPayment(t *testing.T) {
	logger := shim.NewLogger("payment-chaincode")
	logger.SetLevel(shim.LogDebug)
	stub := shim.NewMockStub("transfer_payment", &Chaincode{
		logger: logger,
	})
	if stub == nil {
		t.Fatal("failed to create mock stub")
	}

	t.Run("Successfully transferred payment", func(t *testing.T) {
		const (
			senderWallet   = "GCQDW5IGMLVEUBJSG2AD2OROMH4IBTBQSYMZ6OWIBYTM5SXCKC2QP4E2"
			receiverWallet = "SC6NJXU7Z57VPY2H7OAFLPEYMZAKBMGSLPKZMCBVF7NP3NOCUDYJQSPH"
		)
		stub.MockTransactionStart("123")
		defer stub.MockTransactionEnd("123")

		assert.NoError(t, stub.PutState(senderWallet, []byte(`{"balance":"100.00"}`)))
		assert.NoError(t, stub.PutState(receiverWallet, []byte(`{"balance":"0"}`)))

		transferAmount := fmt.Sprintf("%d", 10)

		result := stub.MockInvoke(transferPayment,
			[][]byte{
				[]byte(transferPayment),
				[]byte(senderWallet),
				[]byte(receiverWallet),
				[]byte(transferAmount),
			})

		senderWalletState, err := stub.GetState(senderWallet)
		assert.NoError(t, err)
		assert.Equal(t, `{"balance":"90"}`, string(senderWalletState))

		receiverWalletState, err := stub.GetState(receiverWallet)
		assert.NoError(t, err)
		assert.Equal(t, `{"balance":"10"}`, string(receiverWalletState))

		assert.Equal(t, int(result.Status), shim.OK)
	})

	t.Run("Invalid amount of args", func(t *testing.T) {
		result := stub.MockInvoke(transferPayment,
			[][]byte{
				[]byte(transferPayment),
			})

		assert.Equal(t, int(result.Status), shim.ERRORTHRESHOLD)
		assert.Equal(t, result.Message, "invalid payment chaincode arguments")
	})
}
