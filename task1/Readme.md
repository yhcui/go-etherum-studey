go mod init github.com/local/go-etherum-studey/task1
go get github.com/ethereum/go-ethereum
go get github.com/ethereum/go-ethereum/rpc


go get github.com/ethereum/go-ethereum/accounts/keystore@v1.16.7
solcjs --bin Counter.sol
solcjs --abi Counter.sol

abigen --bin=Counter_sol_Counter.bin --abi=Counter_sol_Counter.abi --pkg=counter --out=counter.go