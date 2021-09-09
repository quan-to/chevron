#!/usr/bin/env python3

from base64 import b64encode
from chevron import *
from .testdata import *

def test_get_key_fingerprints():
  assert get_key_fingerprints(TestKey) == [TestKeyFingerprint]

def test_load_key():
  assert load_key(TestKey) == TestKeyFingerprint

def test_unlock_key():
  load_key(TestKey)
  assert unlock_key(TestKeyFingerprint, TestKeyPassword) == True

def test_verify_signature():
  load_key(TestKey)
  result, err = verify_signature(PayloadToSign, TestSignature)
  assert result == True
  assert err == ""
  result, err = verify_signature("ABC" + PayloadToSign, TestSignature)
  assert result == False
  assert err == "openpgp: invalid signature: hash tag doesn't match"

def test_quanto_verify_signature():
  load_key(TestKey)
  result, err = quanto_verify_signature(PayloadToSign, TestQuantoSignature)
  assert result == True
  assert err == ""
  result, err = quanto_verify_signature("ABC" + PayloadToSign, TestQuantoSignature)
  assert result == False
  assert err == "openpgp: invalid signature: hash tag doesn't match"

def test_verify_base64_data_signature():
  load_key(TestKey)
  result, err = verify_base64_data_signature(b64encode(PayloadToSign.encode("utf-8")), TestSignature)
  assert result == True
  assert err == ""
  result, err = verify_base64_data_signature(b64encode(("ABC" + PayloadToSign).encode("utf-8")), TestSignature)
  assert result == False
  assert err == "openpgp: invalid signature: hash tag doesn't match"

def test_quanto_verify_base64_data_signature():
  load_key(TestKey)
  result, err = quanto_verify_base64_data_signature(b64encode(PayloadToSign.encode("utf-8")), TestQuantoSignature)
  assert result == True
  assert err == ""
  result, err = quanto_verify_base64_data_signature(b64encode(("ABC" + PayloadToSign).encode("utf-8")), TestQuantoSignature)
  assert result == False
  assert err == "openpgp: invalid signature: hash tag doesn't match"

def test_sign_data():
  load_key(TestKey)
  unlock_key(TestKeyFingerprint, TestKeyPassword)
  signature = sign_data(PayloadToSign, TestKeyFingerprint)
  assert "PGP SIGNATURE" in signature
  assert verify_signature(PayloadToSign, signature)[0] == True

def test_quanto_sign_data():
  load_key(TestKey)
  unlock_key(TestKeyFingerprint, TestKeyPassword)
  signature = quanto_sign_data(PayloadToSign, TestKeyFingerprint)
  assert not "PGP SIGNATURE" in signature
  assert quanto_verify_signature(PayloadToSign, signature)[0] == True

def test_sign_base64_data():
  load_key(TestKey)
  unlock_key(TestKeyFingerprint, TestKeyPassword)
  signature = sign_base64_data(b64encode(PayloadToSign.encode("utf-8")), TestKeyFingerprint)
  assert "PGP SIGNATURE" in signature
  assert verify_base64_data_signature(b64encode(PayloadToSign.encode("utf-8")), signature)[0] == True

def test_quanto_sign_base64_data():
  load_key(TestKey)
  unlock_key(TestKeyFingerprint, TestKeyPassword)
  signature = quanto_sign_base64_data(b64encode(PayloadToSign.encode("utf-8")), TestKeyFingerprint)
  assert not "PGP SIGNATURE" in signature
  assert quanto_verify_base64_data_signature(b64encode(PayloadToSign.encode("utf-8")), signature)[0] == True

def test_change_key_password():
  tmp_pass = "0912345aseuahse"
  key = generate_key(tmp_pass, tmp_pass, 2048)
  new_key = change_key_password(key, tmp_pass, TestKeyPassword)
  fingerprint = load_key(new_key)

  assert unlock_key(fingerprint, TestKeyPassword) == True

def test_get_public_key():
  load_key(TestKey)
  public_key = get_public_key(TestKeyFingerprint)
  fps =  get_key_fingerprints(public_key)
  assert len(public_key) > 0
  assert len(fps) > 0
  assert fps[0] == TestKeyFingerprint
