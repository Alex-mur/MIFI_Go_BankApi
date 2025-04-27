package crypto

import (
	"bytes"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
)

func EncryptPGP(data []byte, pubKey *openpgp.Entity) ([]byte, error) {
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, []*openpgp.Entity{pubKey}, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecryptPGP(data []byte, privKey *openpgp.Entity) ([]byte, error) {
	md, err := openpgp.ReadMessage(bytes.NewReader(data), openpgp.EntityList{privKey}, nil, nil)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(md.UnverifiedBody)
}

func LoadPublicKey(path string) (*openpgp.Entity, error) {
	keyFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()

	entityList, err := openpgp.ReadArmoredKeyRing(keyFile)
	if err != nil {
		return nil, err
	}

	return entityList[0], nil
}
