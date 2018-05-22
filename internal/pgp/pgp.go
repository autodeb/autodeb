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
