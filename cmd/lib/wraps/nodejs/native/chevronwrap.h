#ifndef __CHEVRON_WRAP__
#define __CHEVRON_WRAP__

#define TRUE 1
#define FALSE 0
#define OK TRUE
#define ERROR -1
#define _CHEVRON_BUFFER_SIZE 16384

// Check windows
#if _WIN32 || _WIN64
#   if _WIN64
#       define ENVIRONMENT64
#   else
#       define ENVIRONMENT32
#   endif
#endif

// Check GCC
#if __GNUC__
#   if __x86_64__ || __ppc64__
#       define ENVIRONMENT64
#   else
#       define ENVIRONMENT32
#   endif
#endif

struct LoadKey_return {
    int err;
    int loadedPrivateKeys;
};

typedef int                     UnlockKey_t(char* fingerprint, char* password, char* result, int resultLen);
typedef struct LoadKey_return   LoadKey_t(char* keyData, char* result, int resultLen);
typedef int                     VerifySignature_t(char* data, int dataLen, char* signature, char* result, int resultLen);
typedef int                     VerifyBase64DataSignature_t(char* b64data, char* signature, char* result, int resultLen);
typedef int                     SignData_t(char* data, int dataLen, char* fingerprint, char* result, int resultLen);
typedef int                     SignBase64Data_t(char* b64data, char* fingerprint, char* result, int resultLen);
typedef int                     GetKeyFingerprints_t(char* keyData, char* result, int resultLen);
typedef int                     ChangeKeyPassword_t(char* keyData, char* currentPassword, char* newPassword, char* result, int resultLen);
typedef int                     GetPublicKey_t(char* fingerprint, char* result, int resultLen);
typedef int                     GenerateKey_t(char* password, char* identifier, int bits, char* result, int resultLen);

extern UnlockKey_t                      *chevronlib_unlockkey;
extern LoadKey_t                        *chevronlib_loadkey;
extern VerifySignature_t                *chevronlib_verifysignature;
extern VerifyBase64DataSignature_t      *chevronlib_verifybase64datasignature;
extern SignData_t                       *chevronlib_signdata;
extern SignBase64Data_t                 *chevronlib_signbase64data;
extern GetKeyFingerprints_t             *chevronlib_getkeyfingerprints;
extern ChangeKeyPassword_t              *chevronlib_changekeypassword;
extern GetPublicKey_t                   *chevronlib_getpublickey;
extern GenerateKey_t                    *chevronlib_generatekey;

int loadChevron(const char *path);

#endif