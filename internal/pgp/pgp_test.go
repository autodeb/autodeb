package pgp_test

import (
	"strings"
	"testing"

	"golang.org/x/crypto/openpgp"

	"salsa.debian.org/autodeb-team/autodeb/internal/pgp"
	"salsa.debian.org/autodeb-team/autodeb/internal/pgp/pgptest"

	"github.com/stretchr/testify/assert"
)

func TestVerifySignatureClearsigned(t *testing.T) {
	msg, entity, err := pgp.VerifySignatureClearsigned(
		strings.NewReader(signedMessage),
		strings.NewReader(pgptest.TestKeyPublic),
	)

	assert.NoError(t, err)
	assert.Equal(t, "this is a test\n", msg)
	assert.Equal(t, pgptest.TestKeyFingerprint, pgp.EntityFingerprint(entity))
}

func TestVerfySignatureClearsignedError(t *testing.T) {
	msg, entity, err := pgp.VerifySignatureClearsigned(
		strings.NewReader("POTATOOOO"),
		strings.NewReader(pgptest.TestKeyPublic),
	)

	assert.Error(t, err)
	assert.Equal(t, "", msg)
	assert.Nil(t, entity)
}

func TestClearsignMessage(t *testing.T) {
	msg, err := pgp.Clearsign(
		strings.NewReader("test message"),
		strings.NewReader(pgptest.TestKeyPrivate),
	)

	assert.NoError(t, err)
	assert.Contains(t, msg, "BEGIN PGP SIGNED MESSAGE", "the output should contaian the signed message")
	assert.Contains(t, msg, "BEGIN PGP SIGNATURE", "the output should contain the signature")

	msg, entity, err := pgp.VerifySignatureClearsigned(
		strings.NewReader(msg),
		strings.NewReader(pgptest.TestKeyPublic),
	)

	assert.NoError(t, err)
	assert.Equal(t, "test message", strings.TrimSpace(msg))
	assert.Equal(t, pgptest.TestKeyFingerprint, pgp.EntityFingerprint(entity))
}

func TestEntitySignatures(t *testing.T) {
	keyring, err := openpgp.ReadArmoredKeyRing(
		strings.NewReader(pgptest.TestKeyPublicSigned),
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(keyring))

	signatures := pgp.EntitySignatures(keyring[0])

	assert.Equal(t, 1, len(signatures), "the key should signatures on it")
}

func TestEntitySignaturesNoSignatures(t *testing.T) {
	keyring, err := openpgp.ReadArmoredKeyRing(
		strings.NewReader(pgptest.TestKeyPublic),
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(keyring))

	signatures := pgp.EntitySignatures(keyring[0])

	assert.Equal(t, 0, len(signatures), "the key should have no signatures on it")
}

const signedMessage = `
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

this is a test
-----BEGIN PGP SIGNATURE-----

iQEzBAEBCgAdFiEEi18odTLPQt8c9/fq0x67v8LkDPMFAlsDbBMACgkQ0x67v8Lk
DPMY5gf/XNzYK5CjFrUFMgfEujMffHBpU76LH9uc9iGi4W6wDE25wwM9o+ncqF9N
D50P5A1o5zhlWC0DgGwX/i2MKjdna05bjWTrgG6GPIRqsylPOsznFSjtuOOQymAa
+kqCTyqOByrrwYFChqdWbAXxzftsZMUA1H5M3P9hQFWnYMy8WUKTx/n+0DbebzYn
2iJsk2ZmkzwRRbx/y7oWv7Zl7DUjH8czdN6TZ7u2/kjJMAtMeLnO2BmdgmzGho76
O1Uk2WsCTL9skUyVgvgaBxwqcFkwTx+POX0hsx/14jebZrZPnJxWV0f3OlmtVtcP
XKFCklrFE+4eVrLTn0BoyomfNiTGQA==
=SmYp
-----END PGP SIGNATURE-----
`
