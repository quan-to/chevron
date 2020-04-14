#ifndef CHEVRON_SYNC_H_
#define CHEVRON_SYNC_H_

#include <napi.h>

Napi::Value GetKeyFingerprintsSync(const Napi::CallbackInfo& info);
Napi::Value GetPublicKeySync(const Napi::CallbackInfo& info);

#endif  // CHEVRON_ASYNC_H_