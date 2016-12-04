@ECHO OFF

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

set arg1=admin
set url=http://localhost:7050/chaincode
set chaincodename=myHL
@echo. 
curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"jsonrpc\": \"2.0\",\"method\": \"query\",\"params\": {\"type\": 1,\"chaincodeID\": {\"name\": \"%chaincodename%\"},\"ctorMsg\": {\"function\": \"securities\",\"args\": [\"a\"]},\"secureContext\": \"admin\"},\"id\": 1}" "%url%"
@echo. 

pause

goto :QUERY