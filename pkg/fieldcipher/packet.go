package fieldcipher

type CipherPacket struct {
	EncryptedKey  string `example:"wcDMA8HPMfuMKotZAQwAo62NR4snfqbT3S3EBd3xKAJjJuRx42hJU/f+p0eiQvXNuitRuLe0rF0U8YB3ArAhMX1OZ27t/QE7LKDd1T1oY28kUnHzkKzaIBoted7YXveXLRgr5WI1L6impgxlv+88C81Q6h7RqVWG2Vo6+rXdtg7GdK/VEOtJezIlRJ9Od/gBxmGFjtbSzeoQUTXyzN+xPY60PjpX1FXx+gmM1wHGvZjNLUSsMoKE01JtJJQj1kD4MX9nusp0CONzY4oCNptxgFgcSI/AFj7MZJAW9nH4yR+lQrjw+2KeAhWsWebGK4WiZFdxbEkVJ26GSawCTUqvqJJVt3R7N8vEmgNmM5u+QugM9inFQVa8SUTfqdHmpxq/QO+HtOqbsEiBZWHfNIC1muqjEshwpGhvqfajinSkyR2PbzwUgxPneTrGHiV/cG2LdriAy2zUjNSyoXsYqB9sp3gs9KdKg6nh+f0YE4fAwnb91+2B7xJz0wJFm25iAT4VkJCZjWOULVOzAEJv/38C0uAB5JjPc2wU364MKjwj+/iNutXhTHLgbeD44dmp4GfkCmlO9Hh8TM9JZXOLSgd7uODg4nkTH3ngvOIoe2BD4NTlA/sAp+tXZWRSPinNVIvt1MQgi4QcnHl3SNajPDCqZT7gSOTi7B1gegeF7uHuCeA7IaxK4qmlohvhOFAA"`
	EncryptedJSON map[string]interface{}
}

type UnmatchedField struct {
	Expected string `example:"/data/Test/0/test/"`
	Got      string `example:"/data/Test/0/name/"`
}

type DecipherPacket struct {
	DecryptedData   map[string]interface{}
	JSONChanged     bool `example:"false"`
	UnmatchedFields []UnmatchedField
}
