package baseserver

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
)

func (bs *BaseServer) checkIDFile(ctx context.Context) {
	if bs.BaseServerConfiguration.IDFileLocation == "" {
		return
	}

	// Try reading this file.
	_, err := os.Stat(bs.BaseServerConfiguration.IDFileLocation)
	if err != nil { // Create the file in this case.
		if os.IsNotExist(err) {
			bs.createIDFile()
			return
		}
	}

	// Get it and convert to pub/private keypair.
	fileC, err := os.ReadFile(bs.BaseServerConfiguration.IDFileLocation)
	if err != nil {
		bs.Fail(fmt.Sprintf("service error with id from %s: %v", bs.BaseServerConfiguration.IDFileLocation, err))
	}

	sDec := make([]byte, base64.StdEncoding.DecodedLen(len(fileC)))
	n, err := base64.StdEncoding.Decode(sDec, fileC)
	if err != nil {
		bs.Fail(fmt.Sprintf("decode error with id: %v", err))
	}
	sDec = sDec[:n]

	pub := sDec[0:ed25519.PublicKeySize]
	priv := sDec[ed25519.PublicKeySize:ed25519.PrivateKeySize]
	sig := sDec[ed25519.PublicKeySize+ed25519.PrivateKeySize : ed25519.PublicKeySize+ed25519.PrivateKeySize+ed25519.SignatureSize]
	msg := sDec[ed25519.PublicKeySize+ed25519.PrivateKeySize+ed25519.SignatureSize:]

	// Verify signature of file
	if !ed25519.Verify(pub, msg, sig) {
		bs.Fail(fmt.Sprintf("Invalid ID file detected at %s", bs.BaseServerConfiguration.IDFileLocation))
	}

	bs.config.SetKeys(pub, priv)
	bs.Logger.Infof(bs.LogPrefix, "Loaded existing id %x from file %s", pub, bs.BaseServerConfiguration.IDFileLocation)
}

func (bs *BaseServer) createIDFile() {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		bs.Fail(fmt.Sprintf("service create key err %s", err))
	}

	msg := []byte("kentik/ktranslate")
	sig := ed25519.Sign(priv, msg)

	full := make([]byte, len(pub)+len(priv)+len(sig)+len(msg))
	copy(full, pub)
	copy(full[len(pub):], priv)
	copy(full[len(pub)+len(priv):], sig)
	copy(full[len(pub)+len(priv)+len(sig):], msg)

	sEnc := make([]byte, base64.StdEncoding.EncodedLen(len(full)))
	base64.StdEncoding.Encode(sEnc, full)

	err = os.WriteFile(bs.BaseServerConfiguration.IDFileLocation, sEnc, 0600)
	if err != nil {
		bs.Fail(fmt.Sprintf("encode error write to %s with id: %v", bs.BaseServerConfiguration.IDFileLocation, err))
	}

	bs.config.SetKeys(pub, priv)
	bs.Logger.Infof(bs.LogPrefix, "Created new id %x from file %s", pub, bs.BaseServerConfiguration.IDFileLocation)
}
