'use strict'

var account = require('../account/account');
var Transaction = require('../account/transaction');

var password = '12345678';

var result = account.create(password);
var recoverKey = account.recover(password, JSON.stringify(result.keystore));

console.log(result);
console.log(recoverKey);

var privateKey = Buffer.from(result.privateKey, 'hex');
var from = result.address;
console.log(privateKey);
var tx = new Transaction();
tx.from = from;
tx.to = '0x1111';
tx.amount = 123;

tx.sign(privateKey);
console.log(tx.verify());
