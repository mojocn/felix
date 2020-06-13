package cmd

//
//import (
//	"compress/gzip"
//	"crypto"
//	"crypto/rand"
//	"crypto/rsa"
//	"errors"
//	"github.com/spf13/cobra"
//	"golang.org/x/crypto/openpgp"
//	"golang.org/x/crypto/openpgp/armor"
//	"golang.org/x/crypto/openpgp/packet"
//	"gopkg.in/alecthomas/kingpin.v2"
//	"io"
//	"os"
//	"path/filepath"
//	"time"
//)
//
//// cronCmd represents the cron command
//var pgpCmd = &cobra.Command{
//	Use:   "pgp",
//	Short: "",
//	Long:  ``,
//	Run: func(cmd *cobra.Command, args []string) {
//
//
//	},
//}
//var (
//	// Goencrypt app
//	app  = kingpin.New("goencrypt", "A command line tool for encrypting files")
//	bits = app.Flag("bits", "Bits for keys").Default("4096").Int()
//	privateKey = app.Flag("private", "Private key").String()
//	publicKey = app.Flag("public", "Public key").String()
//	signatureFile = app.Flag("sig", "Signature File").String()
//
//	// Generates new public and private keys
//	keyGenCmd       = app.Command("keygen", "Generates a new public/private key pair")
//	keyOutputPrefix = keyGenCmd.Arg("prefix", "Prefix of key files").Required().String()
//	keyOutputDir    = keyGenCmd.Flag("d", "Output directory of key files").Default(".").String()
//
//	// Encrypts a file with a public key
//	encryptionCmd = app.Command("encrypt", "Encrypt from stdin")
//
//	// Signs a file with a private key
//	signCmd = app.Command("sign", "Sign stdin")
//
//	// Verifies a file was signed with the public key
//	verifyCmd = app.Command("verify", "Verify a signature of stdin")
//
//	// Decrypts a file with a private key
//	decryptionCmd = app.Command("decrypt", "Decrypt from stdin")
//)
//
//func init() {
//	rootCmd.AddCommand(cronCmd)
//}
//
//
//func encodePublicKey(out io.Writer, key *rsa.PrivateKey) {
//	w, err := armor.Encode(out, openpgp.PublicKeyType, make(map[string]string))
//	kingpin.FatalIfError(err, "Error creating OpenPGP Armor: %s", err)
//
//	pgpKey := packet.NewRSAPublicKey(time.Now(), &key.PublicKey)
//	kingpin.FatalIfError(pgpKey.Serialize(w), "Error serializing public key: %s", err)
//	kingpin.FatalIfError(w.Close(), "Error serializing public key: %s", err)
//}
//
//func decodePublicKey(filename string) *packet.PublicKey {
//
//	// open ascii armored public key
//	in, err := os.Open(filename)
//	kingpin.FatalIfError(err, "Error opening public key: %s", err)
//	defer in.Close()
//
//	block, err := armor.Decode(in)
//	kingpin.FatalIfError(err, "Error decoding OpenPGP Armor: %s", err)
//
//	if block.Type != openpgp.PublicKeyType {
//		kingpin.FatalIfError(errors.New("Invalid private key file"), "Error decoding private key")
//	}
//
//	reader := packet.NewReader(block.Body)
//	pkt, err := reader.Next()
//	kingpin.FatalIfError(err, "Error reading private key")
//
//	key, ok := pkt.(*packet.PublicKey)
//	if !ok {
//		kingpin.FatalIfError(errors.New("Invalid public key"), "Error parsing public key")
//	}
//	return key
//}
//
//func decodeSignature(filename string) *packet.Signature {
//
//	// open ascii armored public key
//	in, err := os.Open(filename)
//	kingpin.FatalIfError(err, "Error opening public key: %s", err)
//	defer in.Close()
//
//	block, err := armor.Decode(in)
//	kingpin.FatalIfError(err, "Error decoding OpenPGP Armor: %s", err)
//
//	if block.Type != openpgp.SignatureType {
//		kingpin.FatalIfError(errors.New("Invalid signature file"), "Error decoding signature")
//	}
//
//	reader := packet.NewReader(block.Body)
//	pkt, err := reader.Next()
//	kingpin.FatalIfError(err, "Error reading signature")
//
//	sig, ok := pkt.(*packet.Signature)
//	if !ok {
//		kingpin.FatalIfError(errors.New("Invalid signature"), "Error parsing signature")
//	}
//	return sig
//}
//
//func encryptFile() {
//	pubKey := decodePublicKey(*publicKey)
//	privKey := decodePrivateKey(*privateKey)
//
//	to := createEntityFromKeys(pubKey, privKey)
//
//	w, err := armor.Encode(os.Stdout, "Message", make(map[string]string))
//	kingpin.FatalIfError(err, "Error creating OpenPGP Armor: %s", err)
//	defer w.Close()
//
//	plain, err := openpgp.Encrypt(w, []*openpgp.Entity{to}, nil, nil, nil)
//	kingpin.FatalIfError(err, "Error creating entity for encryption")
//	defer plain.Close()
//
//	compressed, err := gzip.NewWriterLevel(plain, gzip.BestCompression)
//	kingpin.FatalIfError(err, "Invalid compression level")
//
//	n, err := io.Copy(compressed, os.Stdin)
//	kingpin.FatalIfError(err, "Error writing encrypted file")
//	kingpin.Errorf("Encrypted %d bytes", n)
//
//	compressed.Close()
//}
//
//func decryptFile() {
//	pubKey := decodePublicKey(*publicKey)
//	privKey := decodePrivateKey(*privateKey)
//
//	entity := createEntityFromKeys(pubKey, privKey)
//
//	block, err := armor.Decode(os.Stdin)
//	kingpin.FatalIfError(err, "Error reading OpenPGP Armor: %s", err)
//
//	if block.Type != "Message" {
//		kingpin.FatalIfError(err, "Invalid message type")
//	}
//
//	var entityList openpgp.EntityList
//	entityList = append(entityList, entity)
//
//	md, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
//	kingpin.FatalIfError(err, "Error reading message")
//
//	compressed, err := gzip.NewReader(md.UnverifiedBody)
//	kingpin.FatalIfError(err, "Invalid compression level")
//	defer compressed.Close()
//
//	n, err := io.Copy(os.Stdout, compressed)
//	kingpin.FatalIfError(err, "Error reading encrypted file")
//	kingpin.Errorf("Decrypted %d bytes", n)
//}
//
//func signFile() {
//	pubKey := decodePublicKey(*publicKey)
//	privKey := decodePrivateKey(*privateKey)
//
//	signer := createEntityFromKeys(pubKey, privKey)
//
//	err := openpgp.ArmoredDetachSign(os.Stdout, signer, os.Stdin, nil)
//	kingpin.FatalIfError(err, "Error signing input")
//}
//
//func verifyFile() {
//	pubKey := decodePublicKey(*publicKey)
//	sig := decodeSignature(*signatureFile)
//
//	hash := sig.Hash.New()
//	io.Copy(hash, os.Stdin)
//
//	err := pubKey.VerifySignature(hash, sig)
//	kingpin.FatalIfError(err, "Error signing input")
//	kingpin.Errorf("Verified signature")
//}
//
//func createEntityFromKeys(pubKey *packet.PublicKey, privKey *packet.PrivateKey) *openpgp.Entity {
//	config := packet.Config{
//		DefaultHash:            crypto.SHA256,
//		DefaultCipher:          packet.CipherAES256,
//		DefaultCompressionAlgo: packet.CompressionZLIB,
//		CompressionConfig: &packet.CompressionConfig{
//			Level: 9,
//		},
//		RSABits: *bits,
//	}
//	currentTime := config.Now()
//	uid := packet.NewUserId("", "", "")
//
//	e := openpgp.Entity{
//		PrimaryKey: pubKey,
//		PrivateKey: privKey,
//		Identities: make(map[string]*openpgp.Identity),
//	}
//	isPrimaryId := false
//
//	e.Identities[uid.Id] = &openpgp.Identity{
//		Name:   uid.Name,
//		UserId: uid,
//		SelfSignature: &packet.Signature{
//			CreationTime: currentTime,
//			SigType:      packet.SigTypePositiveCert,
//			PubKeyAlgo:   packet.PubKeyAlgoRSA,
//			Hash:         config.Hash(),
//			IsPrimaryId:  &isPrimaryId,
//			FlagsValid:   true,
//			FlagSign:     true,
//			FlagCertify:  true,
//			IssuerKeyId:  &e.PrimaryKey.KeyId,
//		},
//	}
//
//	keyLifetimeSecs := uint32(86400 * 365)
//
//	e.Subkeys = make([]openpgp.Subkey, 1)
//	e.Subkeys[0] = openpgp.Subkey{
//		PublicKey: pubKey,
//		PrivateKey: privKey,
//		Sig: &packet.Signature{
//			CreationTime:              currentTime,
//			SigType:                   packet.SigTypeSubkeyBinding,
//			PubKeyAlgo:                packet.PubKeyAlgoRSA,
//			Hash:                      config.Hash(),
//			PreferredHash:             []uint8{8}, // SHA-256
//			FlagsValid:                true,
//			FlagEncryptStorage:        true,
//			FlagEncryptCommunications: true,
//			IssuerKeyId:               &e.PrimaryKey.KeyId,
//			KeyLifetimeSecs:           &keyLifetimeSecs,
//		},
//	}
//	return &e
//}
//
//func generateKeys() {
//	key, err := rsa.GenerateKey(rand.Reader, *bits)
//	kingpin.FatalIfError(err, "Error generating RSA key: %s", err)
//
//	priv, err := os.Create(filepath.Join(*keyOutputDir, *keyOutputPrefix+".privkey"))
//	kingpin.FatalIfError(err, "Error writing private key to file: %s", err)
//	defer priv.Close()
//
//	pub, err := os.Create(filepath.Join(*keyOutputDir, *keyOutputPrefix+".pubkey"))
//	kingpin.FatalIfError(err, "Error writing public key to file: %s", err)
//	defer pub.Close()
//
//	encodePrivateKey(priv, key)
//	encodePublicKey(pub, key)
//}
