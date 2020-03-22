#include <napi.h>
#include "chevronwrap.h"


class GenerateKeyAsyncWorker : public Napi::AsyncWorker {
 public:
  GenerateKeyAsyncWorker(Napi::Function& callback, const std::string& password, const std::string& identifier, int bits) :
    Napi::AsyncWorker(callback),
    password(password),
    identifier(identifier),
    bits(bits),
    resultVal(ERROR) {}
  ~GenerateKeyAsyncWorker() {}

  // Executed inside the worker-thread.
  // It is not safe to access JS engine data structure
  // here, so everything we need for input and output
  // should go on `this`.
  void Execute() {
    char resultArray[_CHEVRON_BUFFER_SIZE];
    memset(resultArray, 0x00, _CHEVRON_BUFFER_SIZE);

    resultVal = chevronlib_generatekey((char *)password.c_str(), (char *)identifier.c_str(), bits, resultArray, _CHEVRON_BUFFER_SIZE);
    result = std::string(resultArray);
  }

  // Executed when the async work is complete
  // this function will be run inside the main event loop
  // so it is safe to use JS engine data again
  void OnOK() {
    Napi::HandleScope scope(Env());
    if (resultVal == ERROR) { // Error
        Callback().Call({Napi::String::New(Env(), result), Env().Undefined()});
    } else {
        Callback().Call({Env().Undefined(), Napi::String::New(Env(), result)});
    }
  }

 private:
    std::string password;
    std::string identifier;
    int bits;
    int resultVal;
    std::string result;
};

class LoadKeyAsyncWorker : public Napi::AsyncWorker {
 public:
  LoadKeyAsyncWorker(Napi::Function& callback, const std::string& keyData) :
    Napi::AsyncWorker(callback),
    keyData(keyData),
    resultVal(ERROR) {}
  ~LoadKeyAsyncWorker() {}

  // Executed inside the worker-thread.
  // It is not safe to access JS engine data structure
  // here, so everything we need for input and output
  // should go on `this`.
  void Execute() {
    char resultArray[_CHEVRON_BUFFER_SIZE];
    memset(resultArray, 0x00, _CHEVRON_BUFFER_SIZE);

    LoadKey_return res = chevronlib_loadkey((char *)keyData.c_str(), resultArray, _CHEVRON_BUFFER_SIZE);
    resultVal = res.err;
    result = std::string(resultArray);
  }

  // Executed when the async work is complete
  // this function will be run inside the main event loop
  // so it is safe to use JS engine data again
  void OnOK() {
    Napi::HandleScope scope(Env());
    if (resultVal == ERROR) { // Error
        Callback().Call({Napi::String::New(Env(), result), Env().Undefined()});
    } else {
        Callback().Call({Env().Undefined(), Napi::String::New(Env(), result)});
    }
  }

 private:
    std::string keyData;
    int resultVal;
    std::string result;
};

class UnlockKeyAsyncWorker : public Napi::AsyncWorker {
 public:
  UnlockKeyAsyncWorker(Napi::Function& callback, const std::string& fingerprint, const std::string& password) :
    Napi::AsyncWorker(callback),
    fingerprint(fingerprint),
    password(password),
    resultVal(ERROR) {}
  ~UnlockKeyAsyncWorker() {}

  // Executed inside the worker-thread.
  // It is not safe to access JS engine data structure
  // here, so everything we need for input and output
  // should go on `this`.
  void Execute() {
    char resultArray[_CHEVRON_BUFFER_SIZE];
    memset(resultArray, 0x00, _CHEVRON_BUFFER_SIZE);

    resultVal = chevronlib_unlockkey((char *)fingerprint.c_str(), (char *)password.c_str(), resultArray, _CHEVRON_BUFFER_SIZE);
    result = std::string(resultArray);
  }

  // Executed when the async work is complete
  // this function will be run inside the main event loop
  // so it is safe to use JS engine data again
  void OnOK() {
    Napi::HandleScope scope(Env());
    if (resultVal == ERROR) { // Error
        Callback().Call({Napi::String::New(Env(), result), Env().Undefined()});
    } else {
        Callback().Call({Env().Undefined(), Napi::String::New(Env(), result)});
    }
  }

 private:
    std::string fingerprint;
    std::string password;
    int resultVal;
    std::string result;
};


class VerifySignatureAsyncWorker : public Napi::AsyncWorker {
 public:
  VerifySignatureAsyncWorker(Napi::Function& callback, const std::string& b64data, const std::string& signature) :
    Napi::AsyncWorker(callback),
    b64data(b64data),
    signature(signature),
    resultVal(ERROR) {}
  ~VerifySignatureAsyncWorker() {}

  // Executed inside the worker-thread.
  // It is not safe to access JS engine data structure
  // here, so everything we need for input and output
  // should go on `this`.
  void Execute() {
    char resultArray[_CHEVRON_BUFFER_SIZE];
    memset(resultArray, 0x00, _CHEVRON_BUFFER_SIZE);

    resultVal = chevronlib_verifybase64datasignature((char *)b64data.c_str(), (char *)signature.c_str(), resultArray, _CHEVRON_BUFFER_SIZE);
    result = std::string(resultArray);
  }

  // Executed when the async work is complete
  // this function will be run inside the main event loop
  // so it is safe to use JS engine data again
  void OnOK() {
    Napi::HandleScope scope(Env());
    if (resultVal == ERROR) { // Error
        Callback().Call({Napi::String::New(Env(), result), Env().Undefined()});
    } else {
        Callback().Call({Env().Undefined(), Napi::Boolean::New(Env(), resultVal == TRUE)});
      }
  }

 private:
    std::string b64data;
    std::string signature;
    int resultVal;
    std::string result;
};

class SignDataAsyncWorker : public Napi::AsyncWorker {
 public:
  SignDataAsyncWorker(Napi::Function& callback, const std::string& b64data, const std::string& fingerprint) :
    Napi::AsyncWorker(callback),
    b64data(b64data),
    fingerprint(fingerprint),
    resultVal(ERROR) {}
  ~SignDataAsyncWorker() {}

  // Executed inside the worker-thread.
  // It is not safe to access JS engine data structure
  // here, so everything we need for input and output
  // should go on `this`.
  void Execute() {
    char resultArray[_CHEVRON_BUFFER_SIZE];
    memset(resultArray, 0x00, _CHEVRON_BUFFER_SIZE);

    resultVal = chevronlib_signbase64data((char *)b64data.c_str(), (char *)fingerprint.c_str(), resultArray, _CHEVRON_BUFFER_SIZE);
    result = std::string(resultArray);
  }

  // Executed when the async work is complete
  // this function will be run inside the main event loop
  // so it is safe to use JS engine data again
  void OnOK() {
    Napi::HandleScope scope(Env());
    if (resultVal == ERROR) { // Error
        Callback().Call({Napi::String::New(Env(), result), Env().Undefined()});
    } else {
        Callback().Call({Env().Undefined(), Napi::String::New(Env(), result)});
      }
  }

 private:
    std::string b64data;
    std::string fingerprint;
  int resultVal;
    std::string result;
};

class ChangeKeyPasswordAsyncWorker : public Napi::AsyncWorker {
 public:
  ChangeKeyPasswordAsyncWorker(Napi::Function& callback, const std::string& keyData, const std::string& currentPassword, const std::string& newPassword) :
    Napi::AsyncWorker(callback),
    keyData(keyData),
    currentPassword(currentPassword),
    newPassword(newPassword),
    resultVal(ERROR) {}
  ~ChangeKeyPasswordAsyncWorker() {}

  // Executed inside the worker-thread.
  // It is not safe to access JS engine data structure
  // here, so everything we need for input and output
  // should go on `this`.
  void Execute() {
    char resultArray[_CHEVRON_BUFFER_SIZE];
    memset(resultArray, 0x00, _CHEVRON_BUFFER_SIZE);

    resultVal = chevronlib_changekeypassword((char *)keyData.c_str(), (char *)currentPassword.c_str(), (char *)newPassword.c_str(), resultArray, _CHEVRON_BUFFER_SIZE);
    result = std::string(resultArray);
  }

  // Executed when the async work is complete
  // this function will be run inside the main event loop
  // so it is safe to use JS engine data again
  void OnOK() {
    Napi::HandleScope scope(Env());
    if (resultVal == ERROR) { // Error
        Callback().Call({Napi::String::New(Env(), result), Env().Undefined()});
    } else {
        Callback().Call({Env().Undefined(), Napi::String::New(Env(), result)});
      }
  }

 private:
    std::string keyData;
    std::string currentPassword;
    std::string newPassword;
    int resultVal;
    std::string result;
};


////////////////////

Napi::Value GenerateKeyAsync(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();

    if (info.Length() < 4) {
        Napi::TypeError::New(env, "Wrong number of arguments").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[0].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"password\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[1].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"identifier\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }
    if (!info[2].IsNumber()) {
        Napi::TypeError::New(env, "Expected first argument \"bits\" to be a number.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    std::string password = info[0].As<Napi::String>().Utf8Value();
    std::string identifier = info[1].As<Napi::String>().Utf8Value();
    int bits = info[2].As<Napi::Number>().Uint32Value();
    Napi::Function callback = info[3].As<Napi::Function>();

    GenerateKeyAsyncWorker* asyncWorker = new GenerateKeyAsyncWorker(callback, password, identifier, bits);
    asyncWorker->Queue();
    return info.Env().Undefined();
}


Napi::Value LoadKeyAsync(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();

    if (info.Length() < 2) {
        Napi::TypeError::New(env, "Wrong number of arguments").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[0].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"keyData\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    std::string keyData = info[0].As<Napi::String>().Utf8Value();
    Napi::Function callback = info[1].As<Napi::Function>();

    LoadKeyAsyncWorker* asyncWorker = new LoadKeyAsyncWorker(callback, keyData);
    asyncWorker->Queue();
    return info.Env().Undefined();
}

Napi::Value UnlockKeyAsync(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();

    if (info.Length() < 3) {
        Napi::TypeError::New(env, "Wrong number of arguments").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[0].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"fingerprint\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[1].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"password\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    std::string fingerprint = info[0].As<Napi::String>().Utf8Value();
    std::string password = info[1].As<Napi::String>().Utf8Value();
    Napi::Function callback = info[2].As<Napi::Function>();

    UnlockKeyAsyncWorker* asyncWorker = new UnlockKeyAsyncWorker(callback, fingerprint, password);
    asyncWorker->Queue();
    return info.Env().Undefined();
}

Napi::Value VerifySignatureAsync(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();

    if (info.Length() < 3) {
        Napi::TypeError::New(env, "Wrong number of arguments").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[0].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"b64data\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[1].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"signature\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    std::string b64data = info[0].As<Napi::String>().Utf8Value();
    std::string signature = info[1].As<Napi::String>().Utf8Value();
    Napi::Function callback = info[2].As<Napi::Function>();

    VerifySignatureAsyncWorker* asyncWorker = new VerifySignatureAsyncWorker(callback, b64data, signature);
    asyncWorker->Queue();
    return info.Env().Undefined();
}

Napi::Value SignDataAsync(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();

    if (info.Length() < 3) {
        Napi::TypeError::New(env, "Wrong number of arguments").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[0].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"b64data\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[1].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"fingerprint\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }


    std::string b64data = info[0].As<Napi::String>().Utf8Value();
    std::string fingerprint = info[1].As<Napi::String>().Utf8Value();
    Napi::Function callback = info[2].As<Napi::Function>();

    SignDataAsyncWorker* asyncWorker = new SignDataAsyncWorker(callback, b64data, fingerprint);
    asyncWorker->Queue();
    return info.Env().Undefined();
}

Napi::Value ChangeKeyPasswordAsync(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();

    if (info.Length() < 4) {
        Napi::TypeError::New(env, "Wrong number of arguments").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[0].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"keyData\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[1].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"currentPassword\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[2].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"newPassword\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    std::string keyData = info[0].As<Napi::String>().Utf8Value();
    std::string currentPassword = info[1].As<Napi::String>().Utf8Value();
    std::string newPassword = info[2].As<Napi::String>().Utf8Value();
    Napi::Function callback = info[3].As<Napi::Function>();

    ChangeKeyPasswordAsyncWorker* asyncWorker = new ChangeKeyPasswordAsyncWorker(callback, keyData, currentPassword, newPassword);
    asyncWorker->Queue();
    return info.Env().Undefined();
}



