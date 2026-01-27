package crypto

import (
	"bytes"
	"fmt"
	"os"

	"filippo.io/age"
	"golang.org/x/term"
)

func PromptForPassphrase() (string, error) {
	fd := int(os.Stdin.Fd())

	fmt.Print("Set archive passphrase: ")
	bytePassphrase, err := term.ReadPassword(fd)
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("error reading password: %s", err)
	}

	fmt.Print("Confirm archive passphrase: ")
	bytePassphraseConfirmation, err := term.ReadPassword(fd)
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("error reading password: %s", err)
	}

	if !bytes.Equal(bytePassphrase, bytePassphraseConfirmation) {
		return "", fmt.Errorf("passwords do not match")
	}

	return string(bytePassphrase), nil
}

func EncryptWithPassphrase(data []byte, passphrase string) ([]byte, error) {
	recipient, err := age.NewScryptRecipient(passphrase)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	w, err := age.Encrypt(buf, recipient)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	for i := range data {
		data[i] = 0
	}

	return buf.Bytes(), nil
}
