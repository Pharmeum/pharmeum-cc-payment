package payment

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestChaincode_createWallet(t *testing.T) {
	stub := shim.NewMockStub("create_wallet", &Chaincode{
		logger: shim.NewLogger("payment-chaincode"),
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
