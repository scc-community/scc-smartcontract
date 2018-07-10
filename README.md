# fabric-chaincode
* 将chaincode目录下所有文件拷贝到chaincode工作目录
* 进入cli容器`docker exec -it cli bash`，在容器内依次执行如下命令
	* `go get -u github.com/shopspring/decimal`（如果已下载过该依赖库则无需再执行）
	* `go get -u github.com/ethereum/go-ethereum`（如果已下载过该依赖库则无需再执行）
	* `go get -u github.com/btcsuite/btcd`（如果已下载过该依赖库则无需再执行）
	* `peer chaincode install -n scc -p github.com/chaincode_scc_forum/go -v xx`
	* `peer chaincode upgrade -o orderer:7050 -C sccchannel -n scc -c '{"Args":["init",""]}' -v xx`
* 验证
	* 启动服务器：`node client-by-sdk/client/client.js`
	* 浏览器上打开如下地址检查各个接口是否ok
		* 创建一个新的账户：http://xxxx:8081/createAccount?password=xxxxxx
		* 查询一个指定账户：http://xxxx:8081/queryAccount?address=0x12345
		* 交易，地址from向to转token：http://xxxx:8081/trading?from=0x12345&to=0x67890&amt=33.3&timestamp=1531192164708&version=1&privateKey=xxxxxxx
		* 奖励（后续待去掉）：http://xxxx:8081/reward?to=abc&amt=8.2
