#!/usr/bin/env python3

from . import native

native.init_native()

# Function definitions

def get_key_fingerprints(ascii_armored_key):
  '''
    Returns all fingerprints contained in the specified key

    @param {string} ascii_armored_key - The private / public key you want the fingerprints in ASCII Armored Format
    @returns {string[]} the fingerprints in the specified keys
  '''
  return native.get_key_fingerprints(ascii_armored_key)

def load_key(ascii_armored_key):
  '''
    Loads the specified key into memory store for later use
    Both public and private keys can be loaded using load_key

    @param {string} ascii_armored_key - The private / public key you want the fingerprints in ASCII Armored Format
    @returns {string} the first fingerprint of the loaded key
  '''
  fps = get_key_fingerprints(ascii_armored_key)
  native.load_key(ascii_armored_key)
  return fps[0]

def unlock_key(fingerprint, password):
  '''
    Unlocks a pre-loaded private key with the specified password
    The private key should be pre-loaded with load_key function

    @param {string} fingerprint - Fingerprint of the private key that should be unlocked for use
    @param {string} password - Password of the private key
    @returns True
  '''
  return native.unlock_key(fingerprint, password)

def verify_signature(data, signature):
  '''
    Verifies the signature of a payload and returns true if it's valid.
    The signature field can be in ASCII Armored Format or base64 encoded binary PGP Signature

    @param {string} data - Data to be verified
    @param {string} signature - A ASCII Armored Format or Base64 Encoded Binary PGP Signature
    @returns {boolean}
  '''
  return native.verify_signature(data, signature)

def quanto_verify_signature(data, signature):
  '''
    Verifies the signature in Quanto Signature Format of a payload and returns true if it's valid.
    The signature field can be in ASCII Armored Format or base64 encoded binary PGP Signature

    @param {string} data - Data to be verified
    @param {string} signature - A ASCII Armored Format or Base64 Encoded Binary PGP Signature
    @returns {boolean}
  '''
  return native.quanto_verify_signature(data, signature)

def verify_base64_data_signature(b64data, signature):
  '''
    Verifies the signature of a payload and returns true if it's valid.
    The data should be encoded as base64
    The signature field can be in ASCII Armored Format or base64 encoded binary PGP Signature

    @param {string} data - Base64 Encoded Data to be verified
    @param {string} signature - A ASCII Armored Format or Base64 Encoded Binary PGP Signature
    @returns {boolean}
  '''
  return native.verify_base64_data_signature(b64data, signature)

def quanto_verify_base64_data_signature(b64data, signature):
  '''
    Verifies the signature in Quanto Signature Format of a payload and returns true if it's valid.
    The data should be encoded as base64
    The signature field can be in ASCII Armored Format or base64 encoded binary PGP Signature

    @param {string} data - Base64 Encoded Data to be verified
    @param {string} signature - A ASCII Armored Format or Base64 Encoded Binary PGP Signature
    @returns {boolean}
  '''
  return native.quanto_verify_base64_data_signature(b64data, signature)

def sign_data(data, fingerprint):
  '''
    Signs the specified data using a pre-loaded and pre-unlocked key specified by fingerprint
    The key should be previously loaded with load_key and unlocked with unlockKey

    @param {string} data - Data to be signed
    @param {string} fingerprint - Fingerprint of the key used to sign data
    @returns {string} - The ASCII Armored PGP Signature
  '''
  return native.sign_data(data, fingerprint)

def quanto_sign_data(data, fingerprint):
  '''
    Signs the specified data in Quanto Signature Format using a pre-loaded and pre-unlocked key specified by fingerprint
    The key should be previously loaded with load_key and unlocked with unlockKey

    @param {string} data - Data to be signed
    @param {string} fingerprint - Fingerprint of the key used to sign data
    @returns {string} - The ASCII Armored PGP Signature
  '''
  return native.quanto_sign_data(data, fingerprint)

def sign_base64_data(data, fingerprint):
  '''
    Signs the specified data using a pre-loaded and pre-unlocked key specified by fingerprint
    The key should be previously loaded with load_key and unlocked with unlockKey
    The data should be encoded as base64

    @param {string} data - Base64 Encoded Data to be signed
    @param {string} fingerprint - Fingerprint of the key used to sign data
    @returns {string} - The ASCII Armored PGP Signature
  '''
  return native.sign_base64_data(data, fingerprint)

def quanto_sign_base64_data(data, fingerprint):
  '''
    Signs the specified data in Quanto Signature Format using a pre-loaded and pre-unlocked key specified by fingerprint
    The key should be previously loaded with load_key and unlocked with unlockKey
    The data should be encoded as base64

    @param {string} data - Base64 Encoded Data to be signed
    @param {string} fingerprint - Fingerprint of the key used to sign data
    @returns {string} - The ASCII Armored PGP Signature
  '''
  return native.quanto_sign_base64_data(data, fingerprint)

def change_key_password(key_data, current_password, new_password):
  '''
    Changes a private key password

    @param {string} key_data - The private key
    @param {string} current_password - The current password of the key
    @param {string} new_password - The new password for the key
    @returns {string} the same private key encrypted with the new_password
  '''
  return native.change_key_password(key_data, current_password, new_password)

def generate_key(password, identifier, bits):
  '''
    Generates a new PGP Private Key by the specified password, identifier and bits (key-length)

    The Identifier should be in one of the following formats:
     - "Name"
     - "Name <email>"
    The key-length (bits parameter) should not be less than 2048
    The key is not automatically loaded into memory after generation

    @param {string} password - The password to encrypt the private key
    @param {string} identifier - The identifier of the key
    @param {number} bits - Number of bits of the RSA Key (recommended 3072)
    @returns {string} - The generated private key
  '''
  return native.generate_key(password, identifier, bits)

def get_public_key(fingerprint):
  '''
    Returns a ASCII Armored public key of a pre-loaded key

    The public/private key should be pre-loaded with loadKey

    @param {string} fingerprint - Fingerprint to fetch the public key
    @returns {string} - The public key
  '''
  return native.get_public_key(fingerprint)
