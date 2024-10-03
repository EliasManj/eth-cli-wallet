// other_test.go
package wallet

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestCreateNewAccount(t *testing.T) {
	account, err := CreateNewAccount(db, "testnewacc")
	require.NoError(t, err)
	require.Equal(t, "testnewacc", account.Label)
	require.NotEmpty(t, account.Privatey)
	require.NotEmpty(t, account.Publicy)
	acc, err := GetAccount(db, account.Label)
	require.NoError(t, err)
	require.Equal(t, account.Label, acc.Label)
	require.Equal(t, account.Privatey, acc.Privatey)
	require.Equal(t, account.Publicy, acc.Publicy)
	require.True(t, acc.Selected)
}

func TestGenerateKeyPair(t *testing.T) {
	private, public, err := GenerateKeyPair()
	require.NoError(t, err)
	require.NotEmpty(t, private)
	require.NotEmpty(t, public)
	fmt.Println(public, private)
	// Ensure the private key is a 64-character hexadecimal string
	require.Equal(t, 64, len(private))
	_, err = crypto.HexToECDSA(private)
	require.NoError(t, err, "Invalid private key format")

	// Ensure the public address is a 40-character hexadecimal string prefixed with "0x"
	require.Equal(t, 42, len(public))
	require.Equal(t, "0x", public[:2])

	// Verify that the public address is correctly derived from the private key
	privateKey, err := crypto.HexToECDSA(private)
	require.NoError(t, err)
	derivedAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	require.Equal(t, public, derivedAddress)
}

func TestImportAccount(t *testing.T) {
	private, public, err := GenerateKeyPair()
	require.NoError(t, err)
	account := Account{Label: "testimport", Publicy: public, Privatey: private}
	err = ImportAccount(db, account)
	require.NoError(t, err)
	getAccount, err := GetAccount(db, account.Label)
	require.NoError(t, err)
	require.Equal(t, account.Label, getAccount.Label)
	require.Equal(t, account.Privatey, getAccount.Privatey)
	require.Equal(t, account.Publicy, getAccount.Publicy)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 5; i++ {
		account, err := CreateNewAccount(db, fmt.Sprintf("testlistacc%d", i))
		require.NoError(t, err)
		require.NotEmpty(t, account)
	}
	for i := 0; i < 5; i++ {
		label := fmt.Sprintf("testlistacc%d", i)
		account, err := GetAccount(db, label)
		require.NoError(t, err)
		require.Equal(t, label, account.Label)
	}
}

func TestRemoveAccount(t *testing.T) {
	account, err := CreateNewAccount(db, "testremoveacc")
	require.NoError(t, err)
	err = RemoveAccount(db, account.Label)
	require.NoError(t, err)
	allAccounts, err := ListAccounts(db)
	require.NoError(t, err)
	require.NotContains(t, allAccounts, account.Label)
}

func TestSelectAccount(t *testing.T) {
	account1, err := CreateNewAccount(db, "testselectacc1")
	require.NoError(t, err)
	_, err = CreateNewAccount(db, "testselectacc2")
	require.NoError(t, err)
	err = SelectAccount(db, account1.Label)
	require.NoError(t, err)
	activeAccount, err := GetSelectedAccount(db)
	require.NoError(t, err)
	require.Equal(t, account1.Label, activeAccount.Label)
	require.True(t, activeAccount.Selected)
	require.Equal(t, account1.Privatey, activeAccount.Privatey)
	require.Equal(t, account1.Publicy, activeAccount.Publicy)
}

func TestOnlyOneSelectedAccount(t *testing.T) {
	account1, err := CreateNewAccount(db, "testonlyoneacc1")
	require.NoError(t, err)
	account2, err := CreateNewAccount(db, "testonlyoneacc2")
	require.NoError(t, err)
	account3, err := CreateNewAccount(db, "testonlyoneacc3")
	require.NoError(t, err)

	err = SelectAccount(db, account1.Label)
	require.NoError(t, err)
	err = SelectAccount(db, account2.Label)
	require.NoError(t, err)
	err = SelectAccount(db, account3.Label)
	require.NoError(t, err)

	activeAccount, err := GetSelectedAccount(db)
	require.NoError(t, err)
	require.Equal(t, account3.Label, activeAccount.Label)
	require.True(t, activeAccount.Selected)
	require.Equal(t, account3.Privatey, activeAccount.Privatey)
	require.Equal(t, account3.Publicy, activeAccount.Publicy)

	acc1, err := GetAccount(db, account1.Label)
	require.NoError(t, err)
	require.False(t, acc1.Selected)

	acc2, err := GetAccount(db, account2.Label)
	require.NoError(t, err)
	require.False(t, acc2.Selected)
}

func TestCreateTwoAccountsSelect(t *testing.T) {
	account1, err := CreateNewAccount(db, "testcreatetwoaccountsselect1")
	require.NoError(t, err)
	account2, err := CreateNewAccount(db, "testcreatetwoaccountsselect2")
	require.NoError(t, err)

	activeAccount, err := GetSelectedAccount(db)
	require.NoError(t, err)
	require.Equal(t, account2.Label, activeAccount.Label)

	acc1, err := GetAccount(db, account1.Label)
	require.NoError(t, err)
	require.False(t, acc1.Selected)

	acc2, err := GetAccount(db, account2.Label)
	require.NoError(t, err)
	require.True(t, acc2.Selected)
}
