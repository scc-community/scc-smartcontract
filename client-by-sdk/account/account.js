'use strict'

var keythereum = require('keythereum');

module.exports = {
    createParams: {
        keyBytes: 32,
        ivBytes: 16
    },
    keystoreOption: {
        kdf: "pbkdf2",
        cipher: "aes-128-ctr",
        kdfparams: {
            c: 262144,
            dklen: 32,
            prf: "hmac-sha256"
        }
    },
    create: function (password) {
        var dk = keythereum.create(this.createParams);
        var privateKey = dk.privateKey;
        var address = keythereum.privateKeyToAddress(privateKey);
        var keystore = keythereum.dump(password, dk.privateKey, dk.salt, dk.iv, this.keystoreOption);
        var result = {
            privateKey: privateKey.toString("hex"),
            address: address.toString("hex"),
            keystore: keystore
        }
        return result;
    },
    recover: function (password, keystore) {
        var keyObject = JSON.parse(keystore);
        return keythereum.recover(password, keyObject).toString("hex");
    }
}




