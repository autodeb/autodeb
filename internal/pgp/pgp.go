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
	"golang.org/x/crypto/openpgp/packet"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

// Entity represents a PGP entity
type Entity = openpgp.Entity

// EntityList is a list of PGP entities
type EntityList = openpgp.EntityList

//Clearsign will sign a message with a cleartext signature and return
//it as a string
func Clearsign(msg io.Reader, key io.Reader) (string, error) {
	var buf bytes.Buffer
	if err := clearsignWriter(msg, key, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

//clearsignWriter will sign a message with a cleartext signature and write it
//to Writer.
func clearsignWriter(msg io.Reader, key io.Reader, w io.Writer) error {
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

//VerifySignatureClearsignedKeyRing verifies the signature of a clearsigned gpg
//message, returning the message contents and the signer's entity.
func VerifySignatureClearsignedKeyRing(msg io.Reader, keyring openpgp.KeyRing) (string, *Entity, error) {
	msgBytes, err := ioutil.ReadAll(msg)
	if err != nil {
		return "", nil, err
	}

	block, _ := clearsign.Decode(msgBytes)
	if block == nil {
		return "", nil, errors.New("could not decode clearsigned message")
	}

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

//VerifySignatureClearsigned verifies the signature of a clearsigned gpg
//message, returning the message contents and the signer's entity.
func VerifySignatureClearsigned(msg io.Reader, keys io.Reader) (string, *Entity, error) {
	keyring, err := openpgp.ReadArmoredKeyRing(keys)
	if err != nil {
		return "", nil, errors.WithMessage(err, "could not read keyring")
	}
	return VerifySignatureClearsignedKeyRing(msg, keyring)
}

//ReadArmoredKeyRing retrieves a keyring from a reader
func ReadArmoredKeyRing(r io.Reader) (openpgp.EntityList, error) {
	return openpgp.ReadArmoredKeyRing(r)
}

//ReadKeyRing reads one or more public/private keys.
func ReadKeyRing(r io.Reader) (openpgp.EntityList, error) {
	return openpgp.ReadKeyRing(r)
}

//EntitySignatures returns all signatures of an entity
func EntitySignatures(entity *openpgp.Entity) []*packet.Signature {
	var signatures []*packet.Signature
	for _, identity := range entity.Identities {
		signatures = append(signatures, identity.Signatures...)
	}
	return signatures
}

//EntityFingerprint returns the hex representation of the
//entity's PrimaryKey fingerprint.
func EntityFingerprint(entity *openpgp.Entity) string {
	hexFingerprint := hex.EncodeToString(entity.PrimaryKey.Fingerprint[:])
	return fmt.Sprintf("0x%s", strings.ToUpper(hexFingerprint))
}
