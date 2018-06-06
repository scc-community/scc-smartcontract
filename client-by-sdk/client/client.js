//应用启动入口
var express=require('express')
//var bodyParser     =         require("body-parser"); 
//创建app应用 =》NodeJS Http.createServer();
var app=express();
var queryScc = require('../queryScc');
var invokeScc = require('../invokeScc');
var chaincodeName = "scc"
var channelName = "mychannel"

//监听http请求
app.listen(8081);
app.get('/test',function (req,res) {
    	console.log(req.query.a)
	console.log(req.query.b)
    	res.send("test")
})
app.get('/queryAccount', async function (req,res) {
	var account = req.query.account
	var args = [ account ]
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
app.get('/createAccount', async function (req,res) {
	var account = req.query.account
	var amt = req.query.amt
	var args = [ account, amt ]
	const a= async ()=> {  
		return invokeScc.invoke(chaincodeName, 'createAccount', args, channelName)
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
