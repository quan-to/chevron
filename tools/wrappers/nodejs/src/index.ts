const lib = require('bindings')('chevron');

// Ensure that the library has been loaded
lib.__loadnative(__dirname );

/**
 * Checks if the data string is a base64 encoded payload
 *
 * @param {string} data - A string to be tested
 * @returns {boolean} - True if the string is a valid base64 format
 */
const isBase64 = (data: string) : boolean => {
  const base64regex = /^([0-9a-zA-Z+/]{4})*(([0-9a-zA-Z+/]{2}==)|([0-9a-zA-Z+/]{3}=))?$/;
  return base64regex.test(data);
};

/**
 * Returns all fingerprints contained in the specified key
 * @param {string} asciiArmoredKey - The private / public key you want the fingerprints in ASCII Armored Format
 * @returns {string[]} the fingerprints in the specified keys
 */
const getKeyFingerprints = (asciiArmoredKey: string) => lib.getKeyFingerprints(asciiArmoredKey).split(",");

/**
 * Loads the specified key into memory store for later use
 *
 * Both public and private keys can be loaded using loadKey
 * @param {string} asciiArmoredKey - The private / public key you want the fingerprints in ASCII Armored Format
 * @returns {Promise<string>} the first fingerprint of the loaded key
 */
const loadKey = async function(asciiArmoredKey: string) {
    return new Promise((resolve, reject) => {
        lib.loadKey(asciiArmoredKey, (error: string|void) => {
            if (error) {
                reject(error);
            } else {
                const fps = getKeyFingerprints(asciiArmoredKey);
                resolve(fps[0]);
            }
        });
    });
};

/**
 * Signs the specified data using a pre-loaded and pre-unlocked key specified by fingerprint
 *
 * The key should be previously loaded with loadKey and unlocked with unlockKey
 * The data should be always encoded as base64
 *
 * @param {string} data - Base64 Encoded Data to be signed
 * @param {string} fingerprint - Fingerprint of the key used to sign data
 * @returns {Promise<string>} - The ASCII Armored PGP Signature
 */
const signData = async function(data: string, fingerprint: string) : Promise<string> {
    return new Promise((resolve, reject) => {
        if (!isBase64(data)) {
            return reject('Expected a base64 encoded data');
        }
        lib.signData(data, fingerprint, (error: string|void, result: any|void) => {
            if (error) {
                return reject(error);
            }

            return resolve(result);
        });
    });
};

/**
 * Verifies the signature of a payload and returns true if it's valid.
 *
 * The data should be always encoded as base64
 * The signature field can be in ASCII Armored Format or base64 encoded binary PGP Signature
 *
 * @param {string} data - Base64 Encoded Data to be signed
 * @param {string} signature - A ASCII Armored Format or Base64 Encoded Binary PGP Signature
 * @returns {Promise<boolean>}
 */
const verifySignature = async function(data: string, signature: string) : Promise<boolean|void> {
    return new Promise((resolve, reject) => {
        if (!isBase64(data)) {
            return reject('Expected a base64 encoded data');
        }
        lib.verifySignature(data, signature, (error: string|void, result: boolean|void) => {
            if (error) {
                if (error.indexOf('invalid signature') > -1) {
                    return resolve(false);
                }
                return reject(error);
            }

            return resolve(result);
        });
    });
};

/**
 * Unlocks a pre-loaded private key with the specified password
 *
 * The private key should be pre-loaded with loadKey function
 *
 * @param {string} fingerprint - Fingerprint of the private key that should be unlocked for use
 * @param {string} password - Password of the private key
 * @returns {Promise}
 */
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
};

/**
 * Generates a new PGP Private Key by the specified password, identifier and bits (key-length)
 *
 * The Identifier should be in one of the following formats:
 *  - "Name"
 *  - "Name <email>"
 * The key-length (bits parameter) should not be less than 2048
 * The key is not automatically loaded into memory after generation
 *
 * @param {string} password - The password to encrypt the private key
 * @param {string} identifier - The identifier of the key
 * @param {number} bits - Number of bits of the RSA Key (recommended 3072)
 * @returns {Promise<string>} - The generated private key
 */
const generateKey = async function(password: string, identifier: string, bits: number) : Promise<string|void> {
    return new Promise((resolve, reject) => {
        lib.generateKey(password, identifier, bits, (error: string|void, result: string|void) => {
            if (error) {
                reject(error);
            } else {
                resolve(result);
            }
        });
    });
};

/**
 * Returns a ASCII Armored public key of a pre-loaded key
 *
 * The public/private key should be pre-loaded with loadKey
 *
 * @param {string} fingerprint - Fingerprint to fetch the public key
 * @returns {string} - The public key
 */
const getPublicKey = (fingerprint: string) : string => lib.getPublicKey(fingerprint);


/**
 * Signs the specified data using a pre-loaded and pre-unlocked key specified by fingerprint
 * and returns in Quanto Signature Format
 *
 * The key should be previously loaded with loadKey and unlocked with unlockKey
 * The data should be always encoded as base64
 *
 * @param {string} data - Base64 Encoded Data to be signed
 * @param {string} fingerprint - Fingerprint of the key used to sign data
 * @returns {Promise<string>} - The ASCII Armored PGP Signature
 */
const quantoSignData = async function(data: string, fingerprint: string) : Promise<string> {
    return new Promise((resolve, reject) => {
        if (!isBase64(data)) {
            return reject('Expected a base64 encoded data');
        }
        lib.quantoSignData(data, fingerprint, (error: string|void, result: any|void) => {
            if (error) {
                return reject(error);
            }

            return resolve(result);
        });
    });
};

/**
 * Verifies the signature in Quanto Signature Format of a payload and returns true if it's valid.
 *
 * The data should be always encoded as base64
 * The signature field can be in ASCII Armored Format or base64 encoded binary PGP Signature
 *
 * @param {string} data - Base64 Encoded Data to be signed
 * @param {string} signature - A ASCII Armored Format or Base64 Encoded Binary PGP Signature
 * @returns {Promise<boolean>}
 */
const quantoVerifySignature = async function(data: string, signature: string) : Promise<boolean|void> {
    return new Promise((resolve, reject) => {
        if (!isBase64(data)) {
            return reject('Expected a base64 encoded data');
        }
        lib.quantoVerifySignature(data, signature, (error: string|void, result: boolean|void) => {
            if (error) {
                if (error.indexOf('invalid signature') > -1) {
                    return resolve(false);
                }
                return reject(error);
            }

            return resolve(result);
        });
    });
};

/**
 * Changes a private key password
 *
 * @param {string} keyData - The private key
 * @param {string} currentPassword - The current password of the key
 * @param {string} newPassword - The new password for the key
 * @returns {Promise<string>} the same private key encrypted with the newPassword
 */
const changeKeyPassword = async function(keyData: string, currentPassword: string, newPassword: string): Promise<string|void> {
    return new Promise((resolve, reject) => {
        lib.changeKeyPassword(keyData, currentPassword, newPassword, (error: string|void, result: string|void) => {
            if (error) {
                return reject(error);
            }

            return resolve(result);
        });
    });
};

export {
    verifySignature,
    signData,
    getKeyFingerprints,
    loadKey,
    unlockKey,
    generateKey,
    getPublicKey,
    isBase64,
    quantoSignData,
    quantoVerifySignature,
    changeKeyPassword,
};
