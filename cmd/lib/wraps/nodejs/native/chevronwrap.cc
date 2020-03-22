#include <napi.h>
#include <dlfcn.h>
#include "chevronwrap.h"

UnlockKey_t 					*chevronlib_unlockkey;
LoadKey_t 						*chevronlib_loadkey;
VerifySignature_t 				*chevronlib_verifysignature;
VerifyBase64DataSignature_t 	*chevronlib_verifybase64datasignature;
SignData_t 						*chevronlib_signdata;
SignBase64Data_t 				*chevronlib_signbase64data;
GetKeyFingerprints_t 			*chevronlib_getkeyfingerprints;
ChangeKeyPassword_t 			*chevronlib_changekeypassword;
GetPublicKey_t 					*chevronlib_getpublickey;
GenerateKey_t					*chevronlib_generatekey;

void *tryLoad(const char *path, const char *name) {
	std::string fullpath = std::string(path) + "/" + std::string(name);
	return dlopen(fullpath.c_str(), RTLD_LAZY);
}

void *loadChevronLibDL(const char *path) {
	void *handler;

	// Try linux
	#ifdef ENVIRONMENT64
	handler = tryLoad(path, "chevron.so");
	if (handler != NULL) {
		return handler;
	}
	#else
	handler = tryLoad(path, "chevron32.so");
	if (handler != NULL) {
		return handler;
	}
	#endif

	// Try MacOSX (only 64 bit)
	handler = tryLoad(path, "chevron.dylib");
	if (handler != NULL) {
		return handler;
	}

	// Try Windows
	#ifdef ENVIRONMENT64
	handler = tryLoad(path, "chevron.dll");
	if (handler != NULL) {
		return handler;
	}
	#else
	handler = tryLoad(path, "chevron32.dll");
	if (handler != NULL) {
		return handler;
	}
	#endif

	return NULL;
}

void loadChevronCalls(void *handler) {
    chevronlib_loadkey 						= (LoadKey_t*) 						dlsym( handler, "LoadKey" );
    chevronlib_unlockkey					= (UnlockKey_t*) 					dlsym( handler, "UnlockKey" );
    chevronlib_verifysignature 				= (VerifySignature_t*) 				dlsym( handler, "VerifySignature" );
    chevronlib_verifybase64datasignature 	= (VerifyBase64DataSignature_t*) 	dlsym( handler, "VerifyBase64DataSignature" );
    chevronlib_signdata 					= (SignData_t*) 					dlsym( handler, "SignData" );
    chevronlib_signbase64data 				= (SignBase64Data_t*) 				dlsym( handler, "SignBase64Data" );
    chevronlib_getkeyfingerprints 			= (GetKeyFingerprints_t*) 			dlsym( handler, "GetKeyFingerprints" );
    chevronlib_changekeypassword 			= (ChangeKeyPassword_t*) 			dlsym( handler, "ChangeKeyPassword" );
    chevronlib_getpublickey 				= (GetPublicKey_t*) 				dlsym( handler, "GetPublicKey" );
    chevronlib_generatekey 					= (GenerateKey_t*) 					dlsym( handler, "GenerateKey" );
}

int loadChevron(const char *path) {
	void *handler = loadChevronLibDL(path);
	if (handler == NULL) {
		return FALSE;
	}
	loadChevronCalls(handler);

	return TRUE;
}