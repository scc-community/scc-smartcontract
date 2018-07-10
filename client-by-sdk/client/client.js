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
	var privateKey = req.query.privateKey
	var timestamp = Date.now()
	var version = 1
	var sign = ""
	if (!_checkParam(from) || !_checkParam(to) || !_checkParam(amt) || !_checkParam(privateKey)) {
		res.send({
			"code" : "FAIL",
			"msg" : "param fail!"
	    })
		return
	}

	var tx = new Transaction()
	try {

		tx.from = from
		tx.to = to
		tx.amount = parseFloat(amt)
		tx.timestamp = parseInt(timestamp)
		tx.version = parseInt(version)
		tx.sign(privateKey)
		var suf
		if (tx.v == 27) {
			suf = "00"
		} else {
			suf = "01"
		}
		sign = "0x" + tx.r.toString("hex") + tx.s.toString("hex") + suf

	} catch(err) {
		console.error(err)
		res.send({
			"code" : "FAIL",
			"msg" : "build tx catch error"
	    })
		return
	}

	var args = [ tx.from, tx.to, tx.amount, tx.timestamp, tx.version, sign ]
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
