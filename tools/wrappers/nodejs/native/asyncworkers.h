#ifndef CHEVRON_ASYNC_H_
#define CHEVRON_ASYNC_H_

#include <napi.h>

Napi::Value GenerateKeyAsync(const Napi::CallbackInfo& info);
Napi::Value LoadKeyAsync(const Napi::CallbackInfo& info);
Napi::Value UnlockKeyAsync(const Napi::CallbackInfo& info);


Napi::Value VerifySignatureAsync(const Napi::CallbackInfo& info);
Napi::Value SignDataAsync(const Napi::CallbackInfo& info);
Napi::Value ChangeKeyPasswordAsync(const Napi::CallbackInfo& info);
Napi::Value QuantoVerifySignatureAsync(const Napi::CallbackInfo& info);
Napi::Value QuantoSignDataAsync(const Napi::CallbackInfo& info);

#endif  // CHEVRON_ASYNC_H_
