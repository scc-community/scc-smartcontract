'use strict'

var util = require('ethereumjs-util');
var crypto = require('crypto');

var hashAlgo = "sha256";

function Transaction() {
    var v;
    var s;
    var r;
    var from;
    var to;
    var amount;
    var version;
    var txHash;
}

function _getMsgHash(tx) {
    var msg = {
        from: tx.from,
        to: tx.to,
        amount: tx.amount,
        version: tx.version != null ? tx.version : 1,
    }
    return crypto.createHash(hashAlgo).update(JSON.stringify(msg)).digest();
}

function _verifyTxHash(txHash) {
    if(txHash && txHash.length == 32) {
        return true;
    }
    return false;
}

Transaction.prototype.generateTxHash = function () {
    if(_verifyTxHash(this.txHash)) {
        return this;
    }
    this.txHash = crypto.createHash(hashAlgo).update(JSON.stringify(this)).digest();
    return this;
};

Transaction.prototype.sign = function (privateKey) {
    var msgHash = _getMsgHash(this);
    var result = util.ecsign(msgHash, privateKey);
    if(result instanceof Error) {
        throw result;
    }
    this.r = result.r;
    this.v = result.v;
    this.s = result.s;
    return this;
};

Transaction.prototype.verify = function () {
    // compare from address
    var publicKey = util.ecrecover(_getMsgHash(this), this.v, this.r, this.s);
    var address = '0x' + util.publicToAddress(publicKey).toString('hex');
    return address === this.from;
};

Transaction.prototype.serialize = function () {

};

module.exports = Transaction;
