@ECHO OFF

:QUERY

for /f "delims=" %%x in (serviceKey.txt) do set Keytext=%%x

@echo. 
curl.exe -X POST --insecure --header "Content-Type: application/json" --header "Accept: application/json" -d "{\"jsonrpc\": \"2.0\",\"method\": \"query\",\"params\": {\"type\": 1,\"chaincodeID\": {\"name\": \"myHL\"},\"ctorMsg\": {\"function\": \"ballance\",\"args\": [\"admin\"]},\"secureContext\": \"admin\"},\"id\": 1}" "https://localhost:7050/chaincode"
@echo. 

pause

goto :QUERY