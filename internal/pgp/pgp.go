package pgp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/clearsign"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

//Clearsign will sign a message with a cleartext signature
func Clearsign(msg io.Reader, key io.Reader, w io.Writer) error {
	keyring, err := openpgp.ReadArmoredKeyRing(key)
	if err != nil {
		return errors.WithMessage(err, "could not read key")
	}
	if numKeys := len(keyring); numKeys > 1 {
		return errors.Errorf("more than one key was provided (%d)", numKeys)
	}

	privateKey := keyring[0].PrivateKey

	writeCloser, err := clearsign.Encode(w, privateKey, nil)
	if err != nil {
		return errors.WithMessage(err, "couldn't setup clearsign encoder")
	}

	defer writeCloser.Close()

	if _, err := io.Copy(writeCloser, msg); err != nil {
		return errors.WithMessage(err, "could not copy the message")
	}

	return nil
}

//VerifySignatureClearsigned verifies the signature of a clearsigned gpg
//message, returning the message contents and the signer's entity.
func VerifySignatureClearsigned(msg io.Reader, keys io.Reader) (string, *openpgp.Entity, error) {
	keyring, err := openpgp.ReadArmoredKeyRing(keys)
	if err != nil {
		return "", nil, errors.WithMessage(err, "could not read keyring")
	}

	msgBytes, err := ioutil.ReadAll(msg)
	if err != nil {
		return "", nil, err
	}

	block, _ := clearsign.Decode(msgBytes)

	signer, err := openpgp.CheckDetachedSignature(
		keyring,
		bytes.NewReader(block.Bytes),
		block.ArmoredSignature.Body,
	)
	if err != nil {
		return "", nil, errors.WithMessage(err, "could not check signature")
	}

	messageText := string(block.Plaintext)

	return messageText, signer, nil
}

//EntityFingerprint returns the hex representation of the
//entity's PrimaryKey fingerprint.
func EntityFingerprint(entity *openpgp.Entity) string {
	hexFingerprint := hex.EncodeToString(entity.PrimaryKey.Fingerprint[:])
	return fmt.Sprintf("0x%s", strings.ToUpper(hexFingerprint))
}
