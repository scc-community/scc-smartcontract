//应用启动入口
var express=require('express')
//var bodyParser     =         require("body-parser"); 
//创建app应用 =》NodeJS Http.createServer();
var app=express();
var queryScc = require('../queryScc');
var invokeScc = require('../invokeScc');
var chaincodeName = "scc"
var channelName = "sccchannel"

var account = require('../account/account');
var Transaction = require('../account/transaction');


//监听http请求
app.listen(8081);
// app.get('/createPubAndPriKey', async function (req,res) {
// 	var args = [  ]
// 	const a= async ()=> {  
// 		return queryScc.query(chaincodeName, 'createPubAndPriKey', args)
// 	}  
	  
// 	var result = await a().then((info)=>{  
// 		return info
// 	})  
// 	if (typeof(result)=='undefined') {
// 		result = {
// 			"code" : "FAIL",
// 			"msg" : "FAIL"
// 	        }
// 	}
//     res.send(result)
// })
app.get('/createAccount', async function (req,res) {
	// var account = req.query.account
	// var amt = req.query.amt
	// var args = [ account, amt ]
	var password = req.query.password
	if (!_checkParam(password)) {
		res.send({
			"code" : "FAIL",
			"msg" : "param fail!"
	    })
	}
	var accountInfo = account.create(password);
	var address = accountInfo.address
	var args = [ address ]
	const a = async ()=> {
		return invokeScc.invoke(chaincodeName, 'createAccount', args, channelName)
	}  
	  
	var result = await a().then((info)=>{  
		return info
	})
	if (typeof(result)=='undefined') {
		result = {
			"code" : "FAIL",
			"msg" : "call cc fail!"
	        }
	}
	if (result["code"] == 'SUCCESS') {
		for (var key in accountInfo) {
			result[key] = accountInfo[key]
		} 
	}
    res.send(result)
})
app.get('/queryAccount', async function (req,res) {
	var address = req.query.address
	if (!_checkParam(address)) {
		res.send({
			"code" : "FAIL",
			"msg" : "param fail!"
	    })
	}
	var args = [ address ]
	const a= async ()=> {  
		return queryScc.query(chaincodeName, 'queryAccount', args)
	}  
	  
	var result = await a().then((info)=>{  
		return info
	})  
	if (typeof(result)=='undefined') {
		result = {
			"code" : "FAIL",
			"msg" : "FAIL"
	        }
	}
    res.send(result)
})
app.get('/trading', async function (req,res) {
	var from = req.query.from
	var to = req.query.to
	var amt = req.query.amt
	// var sign = req.query.sign
	var password = req.query.password
	var keystore = req.query.keystore	// json
	if (!_checkParam(from) || !_checkParam(to) || !_checkParam(amt) /*|| !_checkParam(sign)*/
		 || !_checkParam(password) || !_checkParam(keystore)) {
		res.send({
			"code" : "FAIL",
			"msg" : "param fail!"
	    })
	}

	// var privateKey = account.recover(password, JSON.stringify(keystore));
	var privateKey = account.recover(password, keystore);
	var tx = new Transaction();
	tx.from = from;
	tx.to = to;
	tx.amount = amt;
	tx.sign(privateKey);
	if (!tx.verify()) {
		res.send({
			"code" : "FAIL",
			"msg" : "sign check fail!"
	    })
	}

	var args = [ from, to, amt ]
	const a= async ()=> {  
		return invokeScc.invoke(chaincodeName, 'trading', args, channelName)
	}  
	  
	var result = await a().then((info)=>{  
		return info
	})  
	if (typeof(result)=='undefined') {
		result = {
			"code" : "FAIL",
			"msg" : "FAIL"
	        }
	}
    res.send(result)
})
app.get('/reward', async function (req,res) {
	var to = req.query.to
	var amt = req.query.amt
	var args = [ to, amt ]
	const a= async ()=> {  
		return invokeScc.invoke(chaincodeName, 'reward', args, channelName)
	}  
	  
	var result = await a().then((info)=>{  
		return info
	})  
	if (typeof(result)=='undefined') {
		result = {
			"code" : "FAIL",
			"msg" : "FAIL"
	    }
	}
    res.send(result)
})

function _checkParam(param) {
	if (param === null || param === undefined || param === '') {
        return false;
    }
    return true;
}
