package fieldcipher

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/openpgp"
	"github.com/quan-to/remote-signer/openpgp/packet"
	"reflect"
	"strconv"
	"time"
)

const MAGIC = "FCMN"

type Cipher struct {
	publicKeys []*openpgp.Entity
}

func MakeCipherFromASCIIArmoredKeys(publicKeys []string) *Cipher {
	entities := make([]*openpgp.Entity, 0)

	for i, v := range publicKeys {
		keys, err := remote_signer.ReadKey(v)
		if err != nil {
			SLog.Error("Error parsing key %d: %s", i, err)
			continue
		}

		for _, k := range keys {
			entities = append(entities, k)
		}
	}

	return MakeCipher(entities)
}

func MakeCipher(publicKeys []*openpgp.Entity) *Cipher {
	return &Cipher{
		publicKeys: publicKeys,
	}
}

func (c *Cipher) GenerateEncryptedPacket(data map[string]interface{}, skipFields []string) (*CipherPacket, error) {
	key, err := GenerateKey()

	if err != nil {
		return nil, fmt.Errorf("error generating key: %s", err)
	}

	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return nil, fmt.Errorf("error serializing data: %s", err)
	}

	jsonData := string(jsonBytes)

	encJson, err := c.EncryptJSONFields(jsonData, key, skipFields)

	if err != nil {
		return nil, fmt.Errorf("error ciphering packet: %s", err)
	}

	encKey, err := c.PGPEncryptToBase64(key, "field-cipher-key.gpg")

	if err != nil {
		return nil, fmt.Errorf("error ciphering packet: %s", err)
	}

	// Clear the memory to let GC run whenever it can and we dont keep the key in the ram memory
	for i := 0; i < len(key); i++ {
		key[i] = 0x00
	}

	return &CipherPacket{
		EncryptedJSON: encJson,
		EncryptedKey:  encKey,
	}, nil
}

func (c *Cipher) EncryptJSONFields(jsonData string, key []byte, skipFields []string) (map[string]interface{}, error) {
	// We want to receive a json string so we can constrain the types for the cipher
	if skipFields == nil {
		skipFields = make([]string, 0)
	}

	var realData map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &realData)
	if err != nil {
		return nil, err
	}

	encData, err := c.encryptJsonField(realData, key, "/", skipFields)
	if err != nil {
		return nil, err
	}

	return encData, nil
}

func (c *Cipher) encryptJsonField(data map[string]interface{}, baseKey []byte, currentLevel string, skipFields []string) (map[string]interface{}, error) {
	var err error
	encData := map[string]interface{}{}

	for k, v := range data {
		nodePath := currentLevel + base64.StdEncoding.EncodeToString([]byte(k)) + "/"

		// Check if its in skip
		if remote_signer.StringIndexOf(nodePath, skipFields) > -1 {
			encData[k] = v
			continue
		}

		// Check if its an object
		v2, ok := v.(map[string]interface{})
		if ok {
			encData[k], err = c.encryptJsonField(v2, baseKey, nodePath, skipFields)
			if err != nil {
				return nil, fmt.Errorf("error serializing field %s: %s", nodePath, err)
			}
			continue
		}

		// Otherwise, its an node
		encData[k], err = c.encryptNode(v, baseKey, nodePath, skipFields)
		if err != nil {
			return nil, fmt.Errorf("error serializing field %s: %s", nodePath, err)
		}
	}

	return encData, nil
}

func (c *Cipher) encryptNode(obj interface{}, baseKey []byte, currentLevel string, skipFields []string) (interface{}, error) {
	// The Golang Unmarshal have these output types:
	// bool, for JSON booleans
	// float64, for JSON numbers
	// string, for JSON strings
	// []interface{}, for JSON arrays
	// map[string]interface{}, for JSON objects
	// nil for JSON null

	switch v := obj.(type) {
	case []interface{}:
		return c.encryptArray(v, baseKey, currentLevel, skipFields)
	case bool:
		return c.encryptBool(v, baseKey, currentLevel)
	case float64:
		return c.encryptFloat64(v, baseKey, currentLevel)
	case string:
		return c.encryptString(v, baseKey, currentLevel)
	case map[string]interface{}:
		return c.encryptJsonField(v, baseKey, currentLevel, skipFields)
	}

	return nil, fmt.Errorf("unknown type %s", reflect.TypeOf(obj))
}

func (c *Cipher) encryptArray(obj []interface{}, baseKey []byte, currentLevel string, skipFields []string) ([]interface{}, error) {
	var err error
	out := make([]interface{}, len(obj))
	for i, v := range obj {
		nodePath := currentLevel + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", i))) + "/"
		err = nil
		if remote_signer.StringIndexOf(nodePath, skipFields) > -1 {
			out[i] = v
			continue
		}
		switch v2 := v.(type) {
		case bool:
			out[i], err = c.encryptBool(v2, baseKey, nodePath)
		case string:
			out[i], err = c.encryptString(v2, baseKey, nodePath)
		case float64:
			out[i], err = c.encryptFloat64(v2, baseKey, nodePath)
		case map[string]interface{}:
			out[i], err = c.encryptJsonField(v2, baseKey, nodePath, skipFields)
		default:
			return nil, fmt.Errorf("unknown type %s", reflect.TypeOf(v))
		}
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (c *Cipher) encryptBool(v bool, baseKey []byte, currentLevel string) (string, error) {
	payload := c.genDataPayload("bool", strconv.FormatBool(v), currentLevel)
	return AESEncrypt(payload, baseKey)
}

func (c *Cipher) encryptString(v string, baseKey []byte, currentLevel string) (string, error) {
	payload := c.genDataPayload("string", v, currentLevel)
	return AESEncrypt(payload, baseKey)
}

func (c *Cipher) encryptFloat64(v float64, baseKey []byte, currentLevel string) (string, error) {
	payload := c.genDataPayload("float", strconv.FormatFloat(v, 'f', -1, 64), currentLevel)
	return AESEncrypt(payload, baseKey)
}

func (c *Cipher) genDataPayload(dataType, data, currentLevel string) []byte {
	return []byte(fmt.Sprintf("(%s)[%s]%s", dataType, currentLevel, data))
}

func (c *Cipher) PGPEncryptToBase64(data []byte, filename string) (string, error) {
	hints := &openpgp.FileHints{
		FileName: filename,
		IsBinary: true,
		ModTime:  time.Now(),
	}

	config := &packet.Config{
		DefaultHash:            crypto.SHA512,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.NoCompression,
		CompressionConfig:      &packet.CompressionConfig{},
	}

	buf := bytes.NewBuffer(nil)

	closer, err := openpgp.Encrypt(buf, c.publicKeys, nil, hints, config)
	if err != nil {
		return "", err
	}

	_, err = closer.Write(data)

	if err != nil {
		return "", err
	}

	err = closer.Close()
	if err != nil {
		return "", err
	}

	encData := buf.Bytes()

	return base64.StdEncoding.EncodeToString(encData), nil
}

func GenerateKey() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func AESEncrypt(data, baseKey []byte) (string, error) {
	block, err := aes.NewCipher(baseKey)
	if err != nil {
		return "", err
	}

	if len(data)%block.BlockSize() != 0 {
		// Make the zero pad
		pad := block.BlockSize() - len(data)%block.BlockSize()
		newData := make([]byte, len(data)+pad)
		copy(newData, data)
		for i := len(data); i < len(newData); i++ {
			newData[i] = 0x00
		}
		data = newData
	}

	iv, err := GenerateKey()

	if err != nil {
		return "", err
	}

	iv = iv[:16]

	cbc := cipher.NewCBCEncrypter(block, iv)

	output := make([]byte, len(data)+16)

	copy(output, iv)
	cbc.CryptBlocks(output[16:], data)

	return MAGIC + base64.StdEncoding.EncodeToString(output), nil
}
