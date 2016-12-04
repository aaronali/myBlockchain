@ECHO OFF

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

@echo
set arg1=admin
set url=http://localhost:7050/chaincode
set chaincodename=myHL
curl -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"jsonrpc\": \"2.0\",\"method\": \"query\",\"params\": {\"type\": 1,\"chaincodeID\": {\"name\": \"%chaincodename%\"},\"ctorMsg\": {\"function\": \"holdings\",\"args\": [\"Aaron\"]},\"secureContext\": \"%arg1%\"},\"id\": 1}" "%url%"
@echo. 

pause

goto :QUERY