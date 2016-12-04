$(document).ready(function() {
	$('#askbidresponse').hide()
	$("#mainWindow").hide()
	$('#askorbid').hide()
	$('#availableSecurities').hide()
    $('#dateRangePicker')
        .datepicker({
            format: 'mm/dd/yyyy',
            startDate: '01/01/2010',
            endDate: '12/30/2020'
        })
    
})    

var post="POST"
var get="GET"
var enrollId ="admin";
var enrollSecret = "Xurw3yU9zI0l"; 
var chainCodeId = 'myHL' 
	var url='http://localhost:7050'
var urlExtension

var userName;   

 


function login() {
	userName = document.getElementById("input4783297498217").value;
	var xhttp = new XMLHttpRequest();

//	var params ='{"enrollId":"user_type1_0","enrollSecret":"7cd5f37f09"}'

	var params ='{"enrollId":"'+document.getElementById('input4783297498217').value+'","enrollSecret":"'+document.getElementById('input79843205234523532').value+'"}'
	if (!window.XMLHttpRequest) {
	    // code for old IE browsers
		xhttp = new ActiveXObject("Microsoft.XMLHTTP");
	} 
	xhttp.onreadystatechange = function() { 
		document.getElementById("response").innerHTML=this.responseText+"\n" +document.getElementById("response").innerHTML
	    if (this.readyState != 4){   
	    	return false;
	    }else if(this.status == 200) {
	    	if(this.statusText=="OK"){
	    		MainView();
	    			return false
	    	}
	    }else 
	    if (this.status == 401) {
	    	update("Unautorized user");
	    }
	    return false;
	}
  xhttp.open("POST", url+'/registrar', true);
  xhttp.send(params);
 // xhttp.send();
	return false;
}; 

/*
function postFunction(name, params){
	params = '{"Secret":{"enrollId":"'+document.getElementById('input4783297498217') +'","enrollSecret":"'+document.getElementById('input79843205234523532') +'"}}';
    method =  "post"; // Set method to post by default if not specified.
    path = 'https://d48e4ff54a324330ac90c7ed2a4ddaa2-vp1.us.blockchain.ibm.com:444/registrar'
    // The rest of this code assumes you are not using a library.
    // It can be made less wordy if you use one.
    var form = document.createElement("form");
    form.setAttribute("method", method);
    form.setAttribute("action", path);

}





 
function post() {
	params = '{"Secret":{"enrollId":"'+document.getElementById('input4783297498217') +'","enrollSecret":"'+document.getElementById('input79843205234523532') +'"}}';
    method =  "post"; // Set method to post by default if not specified.
    path = 'https://d48e4ff54a324330ac90c7ed2a4ddaa2-vp1.us.blockchain.ibm.com:444/registrar'
    // The rest of this code assumes you are not using a library.
    // It can be made less wordy if you use one.
    var form = document.createElement("form");
    form.setAttribute("method", method);
    form.setAttribute("action", path);

    for(var key in params) {
        if(params.hasOwnProperty(key)) {
            var hiddenField = document.createElement("d3424322");
            hiddeField.setAttribute("type", "hidden");
            hiddenField.setAttribute("name", key);
            hiddenField.setAttribute("value", params[key]);

            form.appendChild(hiddenField);
         }
    }

    document.body.appendChild(form);
    form.submit();
    return false;
}
*/
var uuid=0
function query(isQuery, functionName, args,handlerFunction){ 
	uuid=uuid+1
	 var argString='['
		 for(var i=0;i<args.length;i++){
			 argString+='"'+args[i]+'"'
			 if(i<args.length-1){
				 argString+=","
			 }
		 }
	 argString+=']'
		var xhttp = new XMLHttpRequest();
	 var query;
	 if(isQuery==true){
		 query ='query'
	 }else{
		 query = 'invoke'
	 }
		//var qString ='{"jsonrpc": "2.0",  "method": "'+query+'",  "params": {  "type": 1,  "chaincodeID": {  "name": "'+ chainCodeId +'" },  "ctorMsg": {"function": "'+functionName+'",   "args": '+argString+' }, "secureContext":"'+ userName +'" },  "id": '+uuid+Date.now()+'}'
		var qString ='{"jsonrpc": "2.0",  "method": "'+query+'",  "params": {  "type": 1,  "chaincodeID": {  "name": "'+ chainCodeId +'" },  "ctorMsg": {"function": "'+functionName+'",   "args": '+argString+' }, "secureContext":"'+ userName +'" },  "id": '+uuid+Date.now()+'}'
	  	
		if (!window.XMLHttpRequest) {
		    // code for old IE browsers
			xhttp = new ActiveXObject("Microsoft.XMLHTTP");
		} 
		xhttp.onreadystatechange = function() {
			document.getElementById("response").innerHTML=this.responseText +"\n" +document.getElementById("response").innerHTML
		    if (xhttp.readyState != 4){
		    	return false;
		    }else if(xhttp.status == 200) { 
		    	if(xhttp.statusText=="OK"){ 
		    		handlerFunction(xhttp.response)
		    		return false
		    	}
		    }else 
		    if (xhttp.status == 401) {
		    	update("Unautorized user");
		    }
		    return false;
		}
	  xhttp.open("POST", url+'/chaincode', true);
	  xhttp.send(qString);
	  console.log(qString)
	  return false
 }
 




//  XHHTP response handlers

// update the cash on response recieved
function cashHandler(response){  
	var result = JSON.parse(response)
	document.getElementById("balance").innerHTML=result.result.message 
	return false
}


// Update holdings on response recieved
function holdingHandler(response){ 
	
	$('.holdingsId').show()
	$('.holdingsCount').show()
	var result = JSON.parse(response) 
	var holds =result.result.message.replace("[","").replace("]","").trim()
	var test = holds.split("},")
	document.getElementById('holdingsId').innerHTML = ""
	for (var i =0; i< test.length-1;i++){
		var s  =   test[i].trim()
		if(i<test.length-1){
			s=s.concat('}')  
		} 
		s = ("{\"security\":" +s+"}")
		test[i] = JSON.parse(s)
		document.getElementById('holdingsId').innerHTML = document.getElementById('holdingsId').innerHTML +'<b>'+ test[i].security.securityid + ' : '  +test[i].security.units +'</b></br>'
		}
	
	
	
	return false
}

// Update securuties on response recieved
function securitiesHandler(response){
	var result = JSON.parse(response).result.message
	result = result.substr(1,result.length )
	result = result.substr(0,result.length-1)
	var e=result.split(",") 	
	document.getElementById('availableSecurities').innerHTML=""
	for(var temp=0; temp <e.length;temp++){
		document.getElementById('availableSecurities').innerHTML+='<div class=\'btn btn-custom1\' value=\'this.innerHTML\' onclick=\'askorbid(this.innerHTML,selectedSecurity,0)\'>'+e[temp].substr(1,e[temp].trim().length-2)+'</div>'
		}
	return false
}


function askHandler(){
	
}
function bidHandler(response){
	$('#askbidresponse').show() 
	$('#askorbid').hide()
	var result = JSON.parse(response).result 
	document.getElementById('askbidresponse').innerHTML = result.message
}




var selectedSecurity
function askorbid(security, amount){
	$('#askorbid').show()
	$('#askbidresponse').hide()
	document.getElementById('selectedSecurity').value = security
	selectedSecurity = security
	if(amount!=0){
		
	}
	// 	$('#selectedSecurity').value=security
}
function ask(){
	var price =$('#price').value
	var units =$('#units').value
	var exdate=$('#expirationdate').value
	var params = ['ask',userName,selectedSecurity,price,units,exdate]
	query(false,'ask',params,bidHandler) 
}

function bid(){
	var price =$('#price').value
	var units =$('#units').value
	var exdate=$('#expirationdate').value
	var params = ['ask',userName,selectedSecurity,price,units,exdate]
	query(false,'bid',params,bidHandler) 
}
// tradeType, userid, security, price, units, expiry)")

$('#holdinglist').hide()
function showHoldings(){ 
	if($('#holdinglist').visible){
		$('#holdinglist').hide()
		return false
	}else{
		$('#holdinglist').show() 
		$('#askorbid').hide() 
		$('#availableSecurities').hide() 
		var args =[userName]
		var q =query
		q(true, 'holdings',args,holdingHandler)
	return false
	}
}
 
function showAllSecurities(){  
	if( $('#availableSecurities').visible){
		 $('#availableSecurities').hide() 
		 $('#askorbid').hide()  
	}else{ 
	 $('#availableSecurities').toggle()
	 
		 $('#holdinglist').hide()
		var args =[]
		var q =query
		q(true,'securities',args,securitiesHandler)
	}
	return false	 
}


 


function MainView(){
	$('#loginMenu').hide()
	$("#mainWindow").show()
	var args =[userName]
	document.getElementById('userName').innerHTML = userName	
	document.getElementById("response").innerHTML =document.getElementById("response").innerHTML + "\nLogged as "+userName
	var xhttp = new XMLHttpRequest();
	var q = query
	q(true,'ballance',args,cashHandler) 
	 
  return false
}