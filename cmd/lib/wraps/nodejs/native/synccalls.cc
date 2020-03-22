#include <napi.h>
#include "chevronwrap.h"



Napi::Value GetKeyFingerprintsSync(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();

    if (info.Length() < 1) {
        Napi::TypeError::New(env, "Wrong number of arguments").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[0].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"keyData\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    std::string keyData = info[0].As<Napi::String>().Utf8Value();
    char resultArray[_CHEVRON_BUFFER_SIZE];
    memset(resultArray, 0x00, _CHEVRON_BUFFER_SIZE);

    if (chevronlib_getkeyfingerprints((char *)keyData.c_str(), resultArray, _CHEVRON_BUFFER_SIZE) == ERROR) {
        Napi::TypeError::New(env, resultArray).ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    return Napi::String::New(env, resultArray);
}

Napi::Value GetPublicKeySync(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();

    if (info.Length() < 1) {
        Napi::TypeError::New(env, "Wrong number of arguments").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    if (!info[0].IsString()) {
        Napi::TypeError::New(env, "Expected first argument \"fingerPrint\" to be string.").ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    std::string fingerprint = info[0].As<Napi::String>().Utf8Value();
    char resultArray[_CHEVRON_BUFFER_SIZE];
    memset(resultArray, 0x00, _CHEVRON_BUFFER_SIZE);

    if (chevronlib_getpublickey((char *)fingerprint.c_str(), resultArray, _CHEVRON_BUFFER_SIZE) == ERROR) {
        Napi::TypeError::New(env, resultArray).ThrowAsJavaScriptException();
        return Napi::String::New(env, "");
    }

    return Napi::String::New(env, resultArray);
}
