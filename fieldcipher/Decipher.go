package fieldcipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/openpgp"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var fieldMatchRegex = regexp.MustCompile(`\(([a-zA-Z0-9]*)\)\[([\-A-Za-z0-9+/\\=]*)\](.*)`)

type Decipher struct {
	privateKey openpgp.EntityList
}

func MakeDecipherWithASCIIPrivateKey(privateKey string) (*Decipher, error) {
	keys, err := remote_signer.ReadKey(privateKey)
	if err != nil {
		return nil, err
	}

	return MakeDecipher(keys)
}

func MakeDecipher(privateKey openpgp.EntityList) (*Decipher, error) {
	return &Decipher{
		privateKey: privateKey,
	}, nil
}

func (d *Decipher) Unlock(password string) bool {
	for _, v := range d.privateKey.DecryptionKeys() {
		err := v.PrivateKey.Decrypt([]byte(password))
		if err != nil {
			return false
		}
	}

	return true
}

func (d *Decipher) DecipherPacket(packet CipherPacket) (*DecipherPacket, error) {
	_, err := base64.StdEncoding.DecodeString(packet.EncryptedKey)

	if err != nil {
		return nil, fmt.Errorf("invalid encrypted key: %s", err)
	}

	decryptedKey, err := d.pgpDecrypt(packet.EncryptedKey)

	if err != nil {
		return nil, fmt.Errorf("error decrypting key: %s", err)
	}

	unmatchedFields, data, err := d.DecryptJsonFields(packet.EncryptedJSON, decryptedKey)

	if err != nil {
		return nil, err
	}

	for i := range unmatchedFields {
		unmatchedFields[i].Expected = CipherPathUnmangle(unmatchedFields[i].Expected)
		unmatchedFields[i].Got = CipherPathUnmangle(unmatchedFields[i].Got)
	}

	return &DecipherPacket{
		UnmatchedFields: unmatchedFields,
		DecryptedData:   data,
		JSONChanged:     len(unmatchedFields) > 0,
	}, nil
}

func (d *Decipher) DecryptJsonFields(data map[string]interface{}, baseKey []byte) ([]UnmatchedField, map[string]interface{}, error) {
	return d.decryptJsonObject(data, baseKey, "/", nil)
}

func (d *Decipher) decryptJsonObject(data map[string]interface{}, baseKey []byte, currentLevel string, unmatchedFields []UnmatchedField) ([]UnmatchedField, map[string]interface{}, error) {
	if unmatchedFields == nil {
		unmatchedFields = make([]UnmatchedField, 0)
	}
	var err error
	decData := make(map[string]interface{})

	for k, v := range data {
		nodePath := currentLevel + base64.StdEncoding.EncodeToString([]byte(k)) + "/"

		switch v2 := v.(type) {
		case map[string]interface{}:
			unmatchedFields, decData[k], err = d.decryptJsonObject(v2, baseKey, nodePath, unmatchedFields)
		case []interface{}:
			unmatchedFields, decData[k], err = d.decryptArray(v2, baseKey, nodePath, unmatchedFields)
		default:
			unmatchedFields, decData[k], err = d.decryptNode(v, baseKey, nodePath, unmatchedFields)
		}

		if err != nil {
			return nil, nil, err
		}
	}

	return unmatchedFields, decData, nil
}

func (d *Decipher) decryptArray(data []interface{}, baseKey []byte, currentLevel string, unmatchedFields []UnmatchedField) ([]UnmatchedField, interface{}, error) {
	var err error
	outArray := make([]interface{}, len(data))

	for i, v := range data {
		nodePath := currentLevel + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", i))) + "/"

		switch v2 := v.(type) {
		case map[string]interface{}:
			unmatchedFields, outArray[i], err = d.decryptJsonObject(v2, baseKey, nodePath, unmatchedFields)
		case []interface{}:
			unmatchedFields, outArray[i], err = d.decryptArray(v2, baseKey, nodePath, unmatchedFields)
		default:
			unmatchedFields, outArray[i], err = d.decryptNode(v, baseKey, nodePath, unmatchedFields)
		}

		if err != nil {
			return nil, nil, err
		}
	}

	return unmatchedFields, outArray, nil
}

func (d *Decipher) decryptNode(data interface{}, baseKey []byte, currentLevel string, unmatchedFields []UnmatchedField) ([]UnmatchedField, interface{}, error) {
	stringVal, ok := data.(string)

	if !ok { // If not string, not encrypted
		return unmatchedFields, data, nil
	}

	if len(stringVal) < len(MAGIC) || stringVal[:len(MAGIC)] != MAGIC { // Not Encrypted
		return unmatchedFields, data, nil
	}

	stringVal = stringVal[len(MAGIC):]
	encryptedData, err := base64.StdEncoding.DecodeString(stringVal)

	if err != nil {
		return nil, nil, fmt.Errorf("error decrypting field %s: %s", CipherPathUnmangle(currentLevel), err)
	}

	decryptedData, err := AESDecrypt(encryptedData, baseKey)

	if err != nil {
		return nil, nil, fmt.Errorf("error decrypting field %s: %s", CipherPathUnmangle(currentLevel), err)
	}

	if !fieldMatchRegex.MatchString(decryptedData) {
		return nil, nil, fmt.Errorf("invalid decrypted data: %s", decryptedData)
	}

	fields := fieldMatchRegex.FindStringSubmatch(decryptedData)

	if len(fields) != 4 {
		return nil, nil, fmt.Errorf("invalid decrypted data: %s", decryptedData)
	}

	dataType := fields[1]
	nodePath := fields[2]
	dataData := fields[3]

	dataDataS := strings.Split(dataData, "\x00") // Remove block padding
	dataData = dataDataS[0]

	if nodePath != currentLevel {
		unmatchedFields = append(unmatchedFields, UnmatchedField{
			Expected: nodePath,
			Got:      currentLevel,
		})
	}

	var objData interface{}

	switch dataType {
	case "string":
		objData = dataData
		err = nil
	case "float":
		objData, err = strconv.ParseFloat(dataData, 64)
	case "bool":
		objData, err = strconv.ParseBool(dataData)
	case "int":
		objData, err = strconv.ParseInt(dataData, 10, 64)
	case "datetime":
		objData, err = time.Parse(time.RFC3339, dataData)
	case "null":
		objData = nil
		err = nil
	default:
		return nil, nil, fmt.Errorf("unknown type %s at %s", dataType, CipherPathUnmangle(currentLevel))
	}

	if err != nil {
		return nil, nil, fmt.Errorf("error parsing data at %s: %s", CipherPathUnmangle(currentLevel), err)
	}

	return unmatchedFields, objData, nil
}

func (d *Decipher) pgpDecrypt(data string) ([]byte, error) {
	var err error
	var fps []string

	fps, err = remote_signer.GetFingerPrintsFromEncryptedMessageRaw(data)

	if err != nil {
		return nil, err
	}

	if len(fps) == 0 {
		return nil, fmt.Errorf("no encrypted payloads found")
	}

	var rd io.Reader

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding key base64: %s", err)
	}
	rd = bytes.NewReader(decoded)

	md, err := openpgp.ReadMessage(rd, d.privateKey, nil, nil)

	if err != nil {
		return nil, fmt.Errorf("error reading encrypted key: %s", err)
	}

	rawData, err := ioutil.ReadAll(md.LiteralData.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading key: %s", err)
	}

	return rawData, nil
}

func AESDecrypt(data, baseKey []byte) (string, error) {
	block, err := aes.NewCipher(baseKey)
	if err != nil {
		return "", err
	}

	iv := data[:16]
	data = data[16:]

	cbc := cipher.NewCBCDecrypter(block, iv)

	output := make([]byte, len(data))

	cbc.CryptBlocks(output, data)

	return string(output), nil
}
