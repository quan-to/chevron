// +build js,wasm

package main

import (
	"encoding/json"
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/fieldcipher"
	"github.com/quan-to/chevron/keyBackend"
	"github.com/quan-to/chevron/keymagic"
	"syscall/js"
	"time"
)

const testPrivateKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: OpenPGP.js v2.5.4
Comment: http://openpgpjs.org

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
avi/TFCd0pdY7CA=
=1sh5
-----END PGP PRIVATE KEY BLOCK-----`

const testEncryptedData = `{
    "data": {
      "Bank_GetAccountStatement": [{
        "category": "FCMNaEb9cOfPyHCiUxLdSrapV/N4bC6lsNbCW6lGl6Gymw2pCbtGO5ewOwwh3RlbjM0N+RHWK82rCar1ZoTashopNAUcLx3B5udxv9G2CAxHkpJAYjQkvKWcc3XwvJZ5gYUU",
        "historyCode": "FCMNfLg4J33qOpF+5rTWPdIinU4Ib2gu6JZsk296fb7oA/seW35YKnO6ZEnDeOafIiW2GeFxQJGGn5L1xWIqs2/vUM1djofmzDRV6MTG5BpqCvzOVGHO3zPrQd8v8WVnW3GgZr0cx8L9abSetB4gl/9tqQ==",
        "historyComplement": "FCMNAd1YrFqd+/vg5edyOZHBo5SSylakfyNv4Y+LQAsFuzOm0+9WkzRf9QctewAZNCZchh9DUVylva8A7xJQ2obeBdskkRRWfcHnHy9smPlYIMi8I/u4JdQXrbScl4Yw5Qlc5gqyHfVQfHKh/XVtZAmFaw==",
        "historyDescription": "FCMNlpgHmWyK6wKeEyQHgOsf44BSlh1stP9Qo2PAG3mIVqjJ/UkYnJgdFa9V3IwPE1lMFwIZLn82rGe+qchrdvJLyU4bof3xeomqERiwtdNcxEFLvR8kxt/vMrvY0OiRQNw8GkJI3VetjSJBGJP+T/qOykhsSdRCYO1O+HntsKdpM9w=",
        "isoDateTime": "FCMNI+KphZVRPD46iuY8ZMVn3OedsBi3kd8NYNGfW680fTDyjTlhA5fF++kU7zTbMZIfhjcH2luSjEddSFFm6FGkXZuq1GA2GcBhdM7T48/wQ4yySdx3F9UPzGeZiZChxDHbgl2qgDgbgu3MUWrvfDUVQelTVOwUEGgtqfeGoHaQe+g=",
        "name": "FCMNnyUQJFkVap6RFRPE3nPMRfGiS93/g+2WIf3J4QDvxtF8CyDD7Ui9/28NaY+Ur4RPsuqXeRQoa8CG8uMp+OyFyYRLo4GKfwfJUyy+PpwTk1LnKuz5TKULGcsyKEE0Opk1",
        "timestamp": "FCMN+EWKdRrPL57QsscVhTucax/IIkh7MfCx3jLy5iWIu1AJez1FiNMfbKzP/Fw3POZ0qC9gOYi1vfqduSqWIMePxeDkPYnVLXMYrpL4xqlNo35pbiKqtqwO6K7Bf6dH3sZURFItfSBsTpRTtBNkgaoGbA==",
        "transactionDetails": {
          "destinationAccountNumber": "FCMNL3DMkhw6/37rLQi3TAT1Z1CM9ZbAA9CKir/DKQV+NtOTOjKLb67/AegctNGzThM8qnZcGd2PNPcidbugOH1RDiY8zGJj9O+nqFF3AycM0We7xVFS0z4mkJNtlDHyspmK8aKngQUx+1aQY23cPNfJoJpRRdek5fnNrToB8CUO5VpEUqptkN+YBtmNK31BSuzR",
          "destinationBranchNumber": "FCMN/V7B63cC2R3JnAYoM2IGaFp/09fPr0K8HSpzHu89KNU3DstNGGGKv86Yo1eG4zgL1/OILleIqxK3zA1+w5+3btSySVzppVCdvLZrUVVU3WY+emOnjKOeLWsMfisMNO8cPdgFPUhTTP8w0bGXoGPrZKia/wJV2+QQa7Exd4+91ReYwFwLFOVyTn+FbhId2ExC",
          "destinationName": "FCMN0UpDtRolfx4U1yT24MIN9PQpx5KxjDXthRypxxG9y0FDWW37dZcDZ3ilLbGGl9UhSI/HimE/sFNeW6srz6BaGUYEev5xWYjfkK9XXtUIL5u0KAI2p8TlOOfeMwTbbqOkKAAdY5t9sLbtlCOdp5zIx5B57chGGcXxByE7Wvly7Jk=",
          "dstRoutingNumber": "FCMN+Kv03FYgz2mftFpKh7HGZa0vdIqlfTUXAauXQRC4wVHQ4C69eWv7t4/yziJyHyFiJbKRjE0WZYTaZSlC3no6RwdCfGxCCTCInjszAjGbTk3tGFTkmlGg43F7fLz/4eQBVr7/aBDFuQwWtRAd2ZZf4KKQYLof3V612r3lT2iXvT8=",
          "dstRoutingType": "FCMN7Yyr1408UFwZR3SiRDnO5fDsL2AucT4kUw+DMWKM8t7n65wZ4g88ZEpPwxDB1mCwnbE26KPKch9+Ic0MTF9dqkSKAmg9CG6emByrUyF2UxR77pjujhf6XXHyeTRt7DvpAC+XuIiWcId877gpt4oH7FjeMrjtLYgVosx5EvWia+w=",
          "metadata": "FCMNTLdZ2ouvyOwyFKAkjmXikfVI8f47eKUY5skoXJetJ9QjVFiDiftbEOXrkJhMVsbI6Uxo3L32OqdhYV1LqWIU1/VybiDMwlOt03DLXSTQ2jD4jS1qD8VTHFe8g50P3G5BGfxaKNkHu8z3P/s+eeE4KHWNENQyydUuw5Hqd7jm8Qo=",
          "sourceAccountNumber": "FCMNsck+WU5qFsCJ2T5tJPvY3gMgvtsAwSCqXhQvIIouGYtKpu7lebc+K2P0tk9osyk2Gxjow+rvYj9idJBDivfsCTUTE4ehYFc0dX2Fcpj3X4mujFcPulGtFE6JxFILgIB47Oa4lSPptIVhKBPKEG44bxrQ/ElyODxiMS5nscCXCcs8rjlPmJNqMzF5ibIHEy59",
          "sourceBranchNumber": "FCMNjAmRLd5qs/OCVNERRKfrOnqg8w9x9UafMEF8OpGxRV+sx9LD3mORO5JCwBvEoLZEpGaT4F4a/qgTSXtdV8f73u4ihNP5CogQPyd9jqSdpX5uLfRay2o2LILAurl6wCjLTYtNnoyzrYN5Fw2uv6Nj8LKB8yxVsxfEXbhnNq2okXg=",
          "sourceName": "FCMNOTviN+9C4Qtfo1KT9HuSJqQCBe1syE9X4uKDiMTgKfDa9xR0YpyxHrLzQYehtEmeeRQx7d6IsLsiRMeFDxWZbqrjwLyITltsjuuNjjate9IIhg0zKOEMCgdeWEAl8BtKy4pQ5lK+nDM967MDayL9pKs7iY+5ZHT19uW8YIfIBMk=",
          "srcRoutingNumber": "FCMN+uMo5hknn86je0anNbETUGRWE5na8WU8QF/bkK+8EV58OTMwxckJiqhyklbWSL4n/X9CF0OC2olPIhqnCKNpDq+Gd1YYgYhaOftAMu78yzugZEARr4p6u1XLjpBNyWSLHaVJnvUzboKtLCldsbCXyPs7GTWa5ls+V06mxRwLLOg=",
          "srcRoutingType": "FCMNqS1lfzQPux/oQGeURnKw8ZMmLfbVS5bfTPtac+FukLpN0cQcO1HA/95+A/GYqJ840xKi78PR5wtQEoAHcn3GG2QXfGyElxkzFNAyMMJo3hFkoKEnH8B6UazTCz/uX9DP+w2Hh5mKqSNB1i5eiQpl6WSU++IcfV3k5UHueFPa5qo=",
          "timestamp": "FCMN8EEP7s+gtL7rN6n05onUHYel+xB7ByQaqq+f5P2MCenko+5bXU8TRxiqBiQA9xYJKd94WemOxj0Wfc4FNtNueSog66uIDrrSdBGqqn4n0V0u6QdjPFz8d65htFzMBHD2TtNbY2p8aBA8feyLw2QCSQKzdMJw9nOGVKlTyN71REc=",
          "type": "FCMNGxnv9AGgiP134rs+XhilSQd3y7A1PLCH2EDSPuidBfMW8zBj3aFi+Gwp2MgBlmvlqogt0PUwDYX+CA7/p9zjYz34HL0h1lPReTFObX0Dfwb68SaFGETQlcmQK9hCG5NILUx0nk6mf6jCBnv5N1FBSA=="
        },
        "type": "FCMNhscy0uBNW8IMh6OgCGvxpV2+oFebDRT5fbVUvqyOMDC/6Rlkec47Jdz2GpFKXMgyeuVS4NN9fFi0K03o4Pj1keDDDcAM93+B/TZiPTOuTg+B9qnWUKW3XozA6Hknz/K5"
      }]
    }
  }`

var log = slog.Scope("GOLANG")

func main() {
	krm := keymagic.MakeKeyRingManager()
	kb := keyBackend.MakeSaveToDiskBackend("a", "")
	pgp := keymagic.MakePGPManagerWithKRM(kb, krm)

	log.Info("Loading Private Key")
	err, n := pgp.LoadKey(testPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Loaded %d private keys", n)
	log.Info("Unlocking Key")
	err = pgp.UnlockKey(rstest.TestKeyFingerprint, rstest.TestKeyFingerprint)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Loading public key")
	pubKey, err := pgp.GetPublicKeyAscii(rstest.TestKeyFingerprint)
	//if err != nil {
	//	panic(err)
	//}

	log.Info("Public Key loaded. Creating Field Cipher")
	cipher := fieldcipher.MakeCipherFromASCIIArmoredKeys([]string{pubKey})
	decipher, err := fieldcipher.MakeDecipher(pgp.GetPrivate(rstest.TestKeyFingerprint))

	if err != nil {
		log.Fatal(err)
	}

	var encData map[string]interface{}

	err = json.Unmarshal([]byte(testEncryptedData), &encData)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Decrypting Data")
	t := time.Now()
	d, err := decipher.DecipherPacket(fieldcipher.CipherPacket{
		EncryptedKey: "wcFMA6uJF6HKi88OARAAGKMZRQ+UwIQUTnBkfpJQqTCryzmgDXnpNgxnPvjbOXgVZmG7XK28mb5W4IHMQWwpZeu0dZ18MmQhc5w5rVVqr8xYGfcYUA2dUiSJD68Kg8ZrGe8pKrev69rfFIRS5PLSpq9T7w/dNNhS1cDyULZCFni/1kRRtkqd/QAyxkLrKgbBNR56cBmP4beCMkTXgzMGlAilR8zJ2+mYSQAMBmg8/V8QnIDwPjuKcU8Kwsa/T209g9oMuEmqAifaCS9dn/3+ALJtbI5o1ajX5S6NJGDmZFfDa1QQf2yK3n8l3s/FBoZ/zfIodC6WYsTv505lIpIEABDIW5pt5B65ZGG0/Nnr6LpG1shsQZEb9eHrpcQQEdq1wI7c/qxNuY/PpSSnnWCdMQnPclOo1koPY45cysK5DaD7YyAXj9lNxyKjznmVXpiqv0hJ/tr/TKrVm8vMuzG9iPFF3aaPiZacA5gS+JB7/uZrGgfVzLw5BA9aBJNCnGj0SNNn2ULjctpkknLQ35fURgayRWc8YBf5EVnxZBYCeE+8OxvfkKA/rES8uKiPrfIXSd6P7NMMDea33ROOlEMVv95gMviO75cN8Xi0QsgkSX8KhGfM2P6Gj8LwIKTcf0oV7zzRbNwY3ao9zRCJ370WcLW/vMaFy8cORe65p4S1lJvw1RRiKi5OjxuAIthNBj7S4AHkKaKjH53yJw4pwnG3GqyN4uHPwOAH4HXh6xjg+eQ5zMU0yLgKnUCPL9urojou4FDiRjKjo+CC4kQ80uDgKeXq+2oP1MH0szgZpEYjBGHajtd3kfUqLpb6CXczf1dLWODI5PHfHxjMnPI6EyR0/6c2tjviWLMQdOE9DgA=",
		EncryptedJSON: encData,
	})
	delta := time.Since(t)
	if err != nil {
		log.Fatal(err)
	}
	sd, _ := json.MarshalIndent(d, "", "   ")

	log.Info("Took %s", delta)
	log.Info("%s", sd)

	cb := js.NewCallback(func(args []js.Value) {
		t = time.Now()
		cp, _ := cipher.GenerateEncryptedPacket(d.DecryptedData, nil)
		delta = time.Since(t)

		log.Info("Took %s", delta)
		sd, _ = json.MarshalIndent(cp, "", "   ")

		log.Info("%s", sd)

		//fmt.Printf("Generating key at %v\n", t)
		//key, err := pgp.GeneratePGPKey("HUEBR", "123456", 2048)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(time.Since(t))
		//fmt.Println(key)
	})

	df := js.NewCallback(func(args []js.Value) {
		log.Info("Decrypting Data")
		t := time.Now()
		d, err := decipher.DecipherPacket(fieldcipher.CipherPacket{
			EncryptedKey: "wcFMA6uJF6HKi88OARAAGKMZRQ+UwIQUTnBkfpJQqTCryzmgDXnpNgxnPvjbOXgVZmG7XK28mb5W4IHMQWwpZeu0dZ18MmQhc5w5rVVqr8xYGfcYUA2dUiSJD68Kg8ZrGe8pKrev69rfFIRS5PLSpq9T7w/dNNhS1cDyULZCFni/1kRRtkqd/QAyxkLrKgbBNR56cBmP4beCMkTXgzMGlAilR8zJ2+mYSQAMBmg8/V8QnIDwPjuKcU8Kwsa/T209g9oMuEmqAifaCS9dn/3+ALJtbI5o1ajX5S6NJGDmZFfDa1QQf2yK3n8l3s/FBoZ/zfIodC6WYsTv505lIpIEABDIW5pt5B65ZGG0/Nnr6LpG1shsQZEb9eHrpcQQEdq1wI7c/qxNuY/PpSSnnWCdMQnPclOo1koPY45cysK5DaD7YyAXj9lNxyKjznmVXpiqv0hJ/tr/TKrVm8vMuzG9iPFF3aaPiZacA5gS+JB7/uZrGgfVzLw5BA9aBJNCnGj0SNNn2ULjctpkknLQ35fURgayRWc8YBf5EVnxZBYCeE+8OxvfkKA/rES8uKiPrfIXSd6P7NMMDea33ROOlEMVv95gMviO75cN8Xi0QsgkSX8KhGfM2P6Gj8LwIKTcf0oV7zzRbNwY3ao9zRCJ370WcLW/vMaFy8cORe65p4S1lJvw1RRiKi5OjxuAIthNBj7S4AHkKaKjH53yJw4pwnG3GqyN4uHPwOAH4HXh6xjg+eQ5zMU0yLgKnUCPL9urojou4FDiRjKjo+CC4kQ80uDgKeXq+2oP1MH0szgZpEYjBGHajtd3kfUqLpb6CXczf1dLWODI5PHfHxjMnPI6EyR0/6c2tjviWLMQdOE9DgA=",
			EncryptedJSON: encData,
		})
		delta := time.Since(t)
		if err != nil {
			log.Fatal(err)
		}
		sd, _ := json.MarshalIndent(d, "", "   ")

		log.Info("Took %s", delta)
		log.Info("%s", sd)
	})

	js.Global().Set("myFunc", cb)
	js.Global().Set("decryptTest", df)

	select {}
}
