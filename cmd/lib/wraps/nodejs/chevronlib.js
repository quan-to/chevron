const lib = require('bindings')('chevron');

const getKeyFingerprints = (keyData) => lib.getKeyFingerprints(key).split(",");

const loadKey = async function(data) {
    return new Promise((resolve, reject) => {
        lib.loadKey(data, (error, result) => {
            if (error) {
                reject(error);
            } else {
                const fps = getKeyFingerprints(data);
                resolve(fps[0]);
            }
        });
    });
}

const signData = async function(data, fingerprint) {
    return new Promise((resolve, reject) => {
        lib.signData(data, fingerprint, (error, result) => {
            if (error) {
                reject(error);
            } else {
                resolve(result);
            }
        });
    });
}

const verifySignature = async function(data, signature) {
    return new Promise((resolve, reject) => {
        lib.verifySignature(payloadToSign, signature, (error, result) => {
            if (error) {
                reject(error);
            } else {
                resolve(result);
            }
        });
    });
};

const unlockKey = async function(fingerprint, password) {
    return new Promise((resolve, reject) => {
        lib.unlockKey(fingerprint, password, (error, result) => {
            if (error) {
                reject(error)
            } else {
                resolve(result);
            }
        });
    });
}

module.exports.verifySignature = verifySignature;
module.exports.signData = signData;
module.exports.getKeyFingerprints = getKeyFingerprints;
module.exports.loadKey = loadKey;
module.exports.unlockKey = unlockKey;