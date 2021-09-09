#!/usr/bin/env python3
#
import os
from ctypes import *
from sys import platform

TRUE = 1
FALSE = 0
OK = TRUE
BUFFER_SIZE = 4096

native_handler = None

# region Structures

class ChevronError(Exception):
  '''
  Exception raised when some error occur calling chevronlib
  '''
  def __init__(self, message):
    self.message = message
    super().__init__(self.message)

class loadkey_return(Structure):
  '''
  LoadKey_return struct
  '''
  _fields_ = [('err', c_int), ('loadedPrivateKeys', c_int)]

# endregion

# region Initialization

def init_native():
  global native_handler
  '''
    Initializes the native handlers for chevronlib
  '''

  shared_library_path = "./chevron.so" # Linux
  if platform.startswith('win32'):
    # Windows has two different DLL, one for 32 bit and one for 64 bit
    # No one should be using that in 32 bit windows, but who knows
    if platform.architecture()[0] == "32bit":
      shared_library_path = "./chevron32.dll"
    else:
      shared_library_path = "./chevron.dll"
  elif platform.startswith('darwin'): # MacOSX
    # TODO: Apple M1
    shared_library_path = "./chevron.dylib"

  # Load library

  try:
    libbasepath = os.path.dirname(__file__)
    native_handler = CDLL(os.path.join(libbasepath, shared_library_path))
  except Exception as e:
    raise ChevronError("cannot load shared library %s: %s" % (shared_library_path, e))

  native_handler.GetKeyFingerprints.restype=c_int
  native_handler.GetKeyFingerprints.argtypes=[c_char_p,c_char_p,c_int]

  native_handler.LoadKey.restype=loadkey_return
  native_handler.LoadKey.argtypes=[c_char_p,c_char_p,c_int]

  native_handler.UnlockKey.restype=c_int
  native_handler.UnlockKey.argtypes=[c_char_p,c_char_p,c_char_p,c_int]

  native_handler.VerifySignature.restype=c_int
  native_handler.VerifySignature.argtypes=[c_char_p,c_int,c_char_p,c_char_p,c_int]

  native_handler.QuantoVerifySignature.restype=c_int
  native_handler.QuantoVerifySignature.argtypes=[c_char_p,c_int,c_char_p,c_char_p,c_int]

  native_handler.VerifyBase64DataSignature.restype=c_int
  native_handler.VerifyBase64DataSignature.argtypes=[c_char_p,c_char_p,c_char_p,c_int]

  native_handler.QuantoVerifyBase64DataSignature.restype=c_int
  native_handler.QuantoVerifyBase64DataSignature.argtypes=[c_char_p,c_char_p,c_char_p,c_int]

  native_handler.SignData.restype=c_int
  native_handler.SignData.argtypes=[c_char_p,c_int,c_char_p,c_char_p,c_int]

  native_handler.QuantoSignData.restype=c_int
  native_handler.QuantoSignData.argtypes=[c_char_p,c_int,c_char_p,c_char_p,c_int]

  native_handler.SignBase64Data.restype=c_int
  native_handler.SignBase64Data.argtypes=[c_char_p,c_char_p,c_char_p,c_int]

  native_handler.QuantoSignBase64Data.restype=c_int
  native_handler.QuantoSignBase64Data.argtypes=[c_char_p,c_char_p,c_char_p,c_int]

  native_handler.ChangeKeyPassword.restype=c_int
  native_handler.ChangeKeyPassword.argtypes=[c_char_p,c_char_p,c_char_p,c_char_p,c_int]

  native_handler.GetPublicKey.restype=c_int
  native_handler.GetPublicKey.argtypes=[c_char_p,c_char_p,c_int]

  native_handler.GenerateKey.restype=c_int
  native_handler.GenerateKey.argtypes=[c_char_p,c_char_p,c_int,c_char_p,c_int]

# endregion

# region Calls

def get_key_fingerprints(ascii_armored_key):
  '''
  get_key_fingerprints returns all fingerprints in CSV format from a ASCII Armored PGP Keychain

  extern int GetKeyFingerprints(char* key_data, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  result = int(native_handler.GetKeyFingerprints(c_char_p(ascii_armored_key.encode('utf-8')), buff, len(buff)))
  buff = str(buff.value.decode("utf-8"))
  if result != OK:
    raise ChevronError(buff)

  return buff.split(",")


def load_key(key_data):
  '''
  load_key loads a private or public key into the memory keyring

  extern struct loadkey_return LoadKey(char* key_data, char* result, int resultLen);
  '''
  global native_handler

  buff = create_string_buffer(BUFFER_SIZE)

  result = native_handler.LoadKey(c_char_p(key_data.encode('utf-8')), buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  if result.err != OK:
    raise ChevronError(buff)

  return result.loadedPrivateKeys

def unlock_key(fingerprint, password):
  '''
  unlock_key unlocks a private key to be used

  extern int UnlockKey(char* fingerprint, char* password, char* result, int resultLen);
  '''
  global native_handler

  buff = create_string_buffer(BUFFER_SIZE)

  result = native_handler.UnlockKey(c_char_p(fingerprint.encode('utf-8')), c_char_p(password.encode('utf-8')), buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  if result != OK:
    raise ChevronError(buff)

  return True

def verify_signature(data, signature):
  '''
  VerifySignature verifies a signature using a already loaded public key

  extern int VerifySignature(char* data, int dataLen, char* signature, char* result, int resultLen);
  '''
  global native_handler


  buff = create_string_buffer(BUFFER_SIZE)

  data_p = c_char_p(data.encode('utf-8'))
  signature = c_char_p(signature.encode('utf-8'))

  result = native_handler.VerifySignature(data_p, len(data), signature, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  return result == 1, buff

def quanto_verify_signature(data, signature):
  '''
  QuantoVerifySignature verifies a signature in Quanto Signature Format using a already loaded public key

  extern int QuantoVerifySignature(char* data, int dataLen, char* signature, char* result, int resultLen);
  '''
  global native_handler

  buff = create_string_buffer(BUFFER_SIZE)

  data_p = c_char_p(data.encode('utf-8'))
  signature = c_char_p(signature.encode('utf-8'))

  result = native_handler.QuantoVerifySignature(data_p, len(data), signature, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  return result == 1, buff

def verify_base64_data_signature(b64data, signature):
  '''
  verify_base64_data_signature verifies a signature using a already loaded public key. The b64data is a raw binary data encoded in base64 string

  extern int VerifyBase64DataSignature(char* b64data, char* signature, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  data_p = c_char_p(b64data)
  signature = c_char_p(signature.encode('utf-8'))

  result = native_handler.VerifyBase64DataSignature(data_p, signature, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  return result == 1, buff

def quanto_verify_base64_data_signature(b64data, signature):
  '''
  quanto_verify_base64_data_signature verifies a signature in Quanto Signature Format using a already loaded public key.
  The b64data is a raw binary data encoded in base64 string

  extern int QuantoVerifyBase64DataSignature(char* b64data, char* signature, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  data_p = c_char_p(b64data)
  signature = c_char_p(signature.encode('utf-8'))

  result = native_handler.QuantoVerifyBase64DataSignature(data_p, signature, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  return result == 1, buff

def sign_data(data, fingerprint):
  '''
  sign_data signs data using a already loaded and unlocked private key

  extern int SignData(char* data, int dataLen, char* fingerprint, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  data_p = c_char_p(data.encode('utf-8'))
  fingerprint = c_char_p(fingerprint.encode('utf-8'))

  result = native_handler.SignData(data_p, len(data), fingerprint, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  if result != OK:
    raise ChevronError(buff)

  return buff

def quanto_sign_data(data, fingerprint):
  '''
  quanto_sign_data signs data using a already loaded and unlocked private key and returns in Quanto Signature Format

  extern int QuantoSignData(char* data, int dataLen, char* fingerprint, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  data_p = c_char_p(data.encode('utf-8'))
  fingerprint = c_char_p(fingerprint.encode('utf-8'))

  result = native_handler.QuantoSignData(data_p, len(data), fingerprint, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  if result != OK:
    raise ChevronError(buff)

  return buff

def sign_base64_data(b64data, fingerprint):
  '''
  SignBase64Data signs data using a already loaded and unlocked private key.
  The b64data is a raw binary data encoded in base64 string

  extern int SignBase64Data(char* b64data, char* fingerprint, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  data_p = c_char_p(b64data)
  fingerprint = c_char_p(fingerprint.encode('utf-8'))

  result = native_handler.SignBase64Data(data_p, fingerprint, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  if result != OK:
    raise ChevronError(buff)

  return buff

def quanto_sign_base64_data(b64data, fingerprint):
  '''
  quanto_sign_base64_data signs data using a already loaded and unlocked private key. Returns in Quanto Signature Format
  The b64data is a raw binary data encoded in base64 string

  extern int QuantoSignBase64Data(char* b64data, char* fingerprint, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  data_p = c_char_p(b64data)
  fingerprint = c_char_p(fingerprint.encode('utf-8'))

  result = native_handler.QuantoSignBase64Data(data_p, fingerprint, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  if result != OK:
    raise ChevronError(buff)

  return buff

def change_key_password(key_data, current_password, new_password):
  '''
  change_key_password re-encrypts the input key using new_password

  extern int ChangeKeyPassword(char* key_data, char* current_password, char* new_password, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  key_data = c_char_p(key_data.encode('utf-8'))
  current_password = c_char_p(current_password.encode('utf-8'))
  new_password = c_char_p(new_password.encode('utf-8'))

  result = native_handler.ChangeKeyPassword(key_data, current_password, new_password, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  if result != OK:
    raise ChevronError(buff)

  return buff

def get_public_key(fingerprint):
  '''
  get_public_key returns the cached public key from the specified fingerprint

  extern int GetPublicKey(char* fingerprint, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  fingerprint = c_char_p(fingerprint.encode('utf-8'))

  result = native_handler.GetPublicKey(fingerprint, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  if result != OK:
    raise ChevronError(buff)

  return buff

def generate_key(password, identifier, bits):
  '''
  generate_key generates a new key using specified bits and identifier and encrypts it using the specified password

  extern int GenerateKey(char* password, char* identifier, int bits, char* result, int resultLen);
  '''
  global native_handler
  buff = create_string_buffer(BUFFER_SIZE)

  password = c_char_p(password.encode('utf-8'))
  identifier = c_char_p(identifier.encode('utf-8'))

  result = native_handler.GenerateKey(password, identifier, bits, buff, len(buff))
  buff = str(buff.value.decode("utf-8"))
  if result != OK:
    raise ChevronError(buff)

  return buff
# endregion