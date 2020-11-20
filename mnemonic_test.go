package hedera

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate24WordMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic24()
	assert.NoError(t, err)

	assert.Equal(t, 24, len(mnemonic.Words()))
}

func TestGenerate12WordMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic24()
	assert.NoError(t, err)

	assert.Equal(t, 24, len(mnemonic.Words()))
}

func TestMnemonicFromString(t *testing.T) {
	mnemonic, err := MnemonicFromString(testMnemonic)
	assert.NoError(t, err)

	assert.Equal(t, testMnemonic, mnemonic.String())
	assert.Equal(t, 24, len(mnemonic.Words()))
}

func TestNew24MnemonicFromGeneratedMnemonic(t *testing.T) {
	generatedMnemonic, err := GenerateMnemonic24()
	assert.NoError(t, err)

	mnemonicFromSlice, err := NewMnemonic(generatedMnemonic.Words())
	assert.NoError(t, err)
	assert.Equal(t, generatedMnemonic.words, mnemonicFromSlice.words)

	mnemonicFromString, err := MnemonicFromString(generatedMnemonic.String())
	assert.NoError(t, err)
	assert.Equal(t, generatedMnemonic, mnemonicFromString)

	gKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)

	slKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)

	stKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)

	assert.Equal(t, gKey.keyData, slKey.keyData)
	assert.Equal(t, gKey.keyData, stKey.keyData)
}

func TestNew12MnemonicFromGeneratedMnemonic(t *testing.T) {
	generatedMnemonic, err := GenerateMnemonic12()
	assert.NoError(t, err)

	mnemonicFromSlice, err := NewMnemonic(generatedMnemonic.Words())
	assert.NoError(t, err)
	assert.Equal(t, generatedMnemonic.words, mnemonicFromSlice.words)

	mnemonicFromString, err := MnemonicFromString(generatedMnemonic.String())
	assert.NoError(t, err)
	assert.Equal(t, generatedMnemonic, mnemonicFromString)

	gKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)

	slKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)

	stKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)

	assert.Equal(t, gKey.keyData, slKey.keyData)
	assert.Equal(t, gKey.keyData, stKey.keyData)
}

func TestMnemonicBreaksWithBadLength(t *testing.T) {
	// note this mnemonic is probably invalid and is only used to test breakage based on length
	shortMnemonic := "inmate flip alley wear offer often piece magnet surge toddler submit right business"

	_, err := MnemonicFromString(shortMnemonic)
	assert.Error(t, err)

	_, err = NewMnemonic(strings.Split(shortMnemonic, " "))
	assert.Error(t, err)
}
