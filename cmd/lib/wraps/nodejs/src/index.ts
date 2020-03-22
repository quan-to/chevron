const os = require('os');

const lib = require('bindings')('chevron');

// Ensure that the library has been loaded
lib.__loadnative(__dirname );

const getKeyFingerprints = (keyData: string) => lib.getKeyFingerprints(keyData).split(",");

const loadKey = async function(data: string) {
    return new Promise((resolve, reject) => {
        lib.loadKey(data, (error: string|void, result: any|void) => {
            if (error) {
                reject(error);
            } else {
                const fps = getKeyFingerprints(data);
                resolve(fps[0]);
            }
        });
    });
}

const signData = async function(data: string, fingerprint: string) : Promise<string> {
    return new Promise((resolve, reject) => {
        lib.signData(data, fingerprint, (error: string|void, result: any|void) => {
            if (error) {
                reject(error);
            } else {
                resolve(result);
            }
        });
    });
}

const verifySignature = async function(data: string, signature: string) : Promise<boolean|void> {
    return new Promise((resolve, reject) => {
        lib.verifySignature(data, signature, (error: string|void, result: boolean|void) => {
            if (error) {
                reject(error);
            } else {
                resolve(result);
            }
        });
    });
};

const unlockKey = async function(fingerprint: string, password: string) {
    return new Promise((resolve, reject) => {
        lib.unlockKey(fingerprint, password, (error: string|void, result: any|void) => {
            if (error) {
                reject(error)
            } else {
                resolve(result);
            }
        });
    });
}

export {
	verifySignature,
	signData,
	getKeyFingerprints,
	loadKey,
	unlockKey,
};
