ChevronLib Node.JS Wrapper
==========================


## Installation

```bash
npm i @contaquanto/chevronlib
```


## Usage

```javascript
const chevron = require('chevronlib');

const payloadToSign = "HUEBR";


function toBase64(data) {
	const buff = Buffer.from(data);
	return buff.toString('base64');
}

(async() => {
	console.log('Generating key');
	const key = await chevron.generateKey('123457890', 'Test Key', 2048);

	console.log('Loading key');
	const fingerprint = await chevron.loadKey(key);

	console.log(`Unlocking key ${fingerprint}`);
	await chevron.unlockKey(fingerprint, '123457890');

	console.log('Signing data "${payloadToSign}"');
	const signature = await chevron.signData(toBase64(payloadToSign), fingerprint);

	console.log(`Validating signature: ${signature}`);
	const verification = await chevron.verifySignature(toBase64(payloadToSign), signature);
	console.log(`Signature is valid: ${verification}`);
})();

```

## Building for release

TODO