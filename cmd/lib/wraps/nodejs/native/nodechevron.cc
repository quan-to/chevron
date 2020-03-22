#include <napi.h>
#include "chevronwrap.h"
#include "asyncworkers.h"
#include "synccalls.h"


Napi::Value LoadNative(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();

    if (info.Length() < 1) {
        Napi::TypeError::New(env, "Wrong number of arguments").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[0].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"script-path\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    std::string path = info[0].As<Napi::String>().Utf8Value();

   if (!loadChevron(path.c_str())) {
        Napi::TypeError::New(env, "Error loading ChevronLib!") .ThrowAsJavaScriptException();
        return Napi::Boolean::New(env, FALSE);
    }

    return Napi::Boolean::New(env, TRUE);
}


Napi::Object init(Napi::Env env, Napi::Object exports) {
    exports.Set(Napi::String::New(env, "__loadnative"),         Napi::Function::New(env, LoadNative));
    exports.Set(Napi::String::New(env, "generateKey"),          Napi::Function::New(env, GenerateKeyAsync));
    exports.Set(Napi::String::New(env, "loadKey"),              Napi::Function::New(env, LoadKeyAsync));
    exports.Set(Napi::String::New(env, "unlockKey"),            Napi::Function::New(env, UnlockKeyAsync));
    exports.Set(Napi::String::New(env, "getKeyFingerprints"),   Napi::Function::New(env, GetKeyFingerprintsSync));
    exports.Set(Napi::String::New(env, "getPublicKey"),         Napi::Function::New(env, GetPublicKeySync));
    exports.Set(Napi::String::New(env, "verifySignature"),      Napi::Function::New(env, VerifySignatureAsync));
    exports.Set(Napi::String::New(env, "signData"),             Napi::Function::New(env, SignDataAsync));
    exports.Set(Napi::String::New(env, "changeKeyPassword"),    Napi::Function::New(env, ChangeKeyPasswordAsync));

    return exports;
};

NODE_API_MODULE(chevron, init);