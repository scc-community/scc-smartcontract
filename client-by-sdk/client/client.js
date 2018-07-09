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
var interval = 10000;	// 10s

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
	    return
	}
	// var accountInfo = account.create(password);
	// var address = accountInfo.address
	var args = [ password ]
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
	// if (result["code"] == 'SUCCESS') {
	// 	for (var key in accountInfo) {
	// 		result[key] = accountInfo[key]
	// 		if(key === "keystore"){
	// 			console.log(JSON.stringify(accountInfo[key]))
	// 		}
	// 	} 
	// }
    res.send(result)
})
app.get('/queryAccount', async function (req,res) {
	var address = req.query.address
	if (!_checkParam(address)) {
		res.send({
			"code" : "FAIL",
			"msg" : "param fail!"
	    })
	    return
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
			"msg" : "call cc fail"
	    }
	}
    res.send(result)
})
app.get('/trading', async function (req,res) {
	var from = req.query.from
	var to = req.query.to
	var amt = req.query.amt
	var r = req.query.r
	var v = req.query.v
	var s = req.query.s
	var password = req.query.password
	var keystore = req.query.keystore	// json
	var timestamp = req.query.timestamp
	var currentTimestamp = Date.now();
    if(!_checkParam(timestamp) || Math.abs(currentTimestamp - timestamp) > interval) {
        res.send({
			"code" : "FAIL",
			"msg" : "param fail! Timestamp Error!"
	    })
		return
    }
	if (!_checkParam(from) || !_checkParam(to) || !_checkParam(amt) 
		 || !_checkParam(r) || !_checkParam(v) || !_checkParam(s)
		 || !_checkParam(password) || !_checkParam(keystore)) {
		res.send({
			"code" : "FAIL",
			"msg" : "param fail!"
	    })
		return
	}

	try {
		// var privateKey = account.recover(password, JSON.stringify(keystore));
		// var privateKey = account.recover(password, keystore)
		// var privateKeyBuf = Buffer.from(privateKey, 'hex')
		var tx = new Transaction()
		tx.from = from
		tx.to = to
		tx.amount = parseFloat(amt)
		tx.timestamp = parseInt(timestamp)
		tx.r = Buffer.from(r, 'hex')
		tx.v = parseInt(v)
		tx.s = Buffer.from(s, 'hex')
		// tx.sign(privateKeyBuf);
		if (!tx.verify()) {
			res.send({
				"code" : "FAIL",
				"msg" : "sign check fail!"
		    })
		    return
		}
	} catch(err) {
		console.error(err)
		res.send({
			"code" : "FAIL",
			"msg" : "sign verify catch error"
	    })
		return
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
