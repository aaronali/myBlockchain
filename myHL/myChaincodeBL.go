/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/



//todo:  add constants for all string litterals
//todo:  need to make consitent status.  need better way to take them out of the process when closed
//todo: data abstraction layer, abstract persistance
//todo: add security to get user names



package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const separator = 		"."
const userIndex =		"UserIndex" + separator
const tradeIndex =		"TradeIndex" + separator
const securityIndex =	"SecurityIndex" + separator
const holdingIndex =	"HoldingIndex" + separator
const initialCash =		1000
const payout =			5
const defaultPrice =	5


type ChaincodeBusinessLayer struct {
	userRep			UserRepository 
	holdingsRep		HoldingsRepository
	securitiesRep	SecurityRepository
	tradeRep		TradeRepository
	stub shim.ChaincodeStubInterface
	logger Logger
}

func (t *ChaincodeBusinessLayer) initObjects(stub shim.ChaincodeStubInterface) error {
	logger = shim.NewLogger("ChaincodeBusinessLayer")
	logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL")) 
	t.stub = stub
	t.writeOut("in init obbjects")
	logger.debugf("func (t *ChaincodeBusinessLayer) initObjects(stub shim.ChaincodeStubInterface) erro");
	
	//initialize our repositories
	t.userRep.init(stub)
	t.holdingsRep.init(stub)
	t.securitiesRep.init(stub)
	t.tradeRep.init(stub)
	
	return nil
}





//********************************************************************************************************
//****                        Debug function inplimentations                                          ****
//********************************************************************************************************


const debug =			true


func (t *ChaincodeBusinessLayer) writeOut(out string) ([]byte, error) {
	if debug {
		curOutByteA,err := t.stub.GetState("currentOutput")
		outByteA := []byte(string(curOutByteA) + ":::" + out)
		err = t.stub.PutState("currentOutput", outByteA)
		return nil, err
	}
	
	return nil, nil
}



func (t *ChaincodeBusinessLayer) readOut() string {
	if debug {
		curOutByteA, err := t.stub.GetState("currentOutput")
		if err != nil {
			return "error"
		}
		
		return string(curOutByteA)
	}
	
	return ""
}




//********************************************************************************************************
//****                        Query function inplimentations                                          ****
//********************************************************************************************************


//to do.  actually return the regersitered entities not the hard coded ones
func (t *ChaincodeBusinessLayer) securities() ([]byte, error) {
	var sec Security
	var out string
	
	numberSecurities, err := t.securitiesRep.getLastIndex()
	if err != nil {
		return nil, err
	}
	
	for b := 1; b <= numberSecurities; b++{
		sec, err = t.securitiesRep.getSecurityByPostion(b)
		if err != nil {
			return nil, err
		}
		
		if out == "" {
			out = "\"" + sec.SecurityID + "\""
		} else {
			out = out + ",\"" + sec.SecurityID + "\""
		}
	}
	
	return []byte("[" + out + "]"), nil
}



func (t *ChaincodeBusinessLayer) users() ([]byte, error) {
	var out string

	for user, err := t.userRep.getFirstUser(); err==nil && user.UserID != ""; user, err = t.userRep.getNextUser(){
		if out == "" {
			out = "\"" + user.UserID + "\""
		} else {
			out = out + ",\"" + user.UserID + "\""
		}
	}
	
	return []byte("[" + out + "]"), nil
}



func (t *ChaincodeBusinessLayer) holdings(userID string) ([]byte, error) {
	fmt.Printf("Running holdings")
	var out string
	
	for holding, err := t.holdingsRep.getFirstHolding(); err==nil && holding.UserID != ""; holding, err = t.holdingsRep.getNextHolding(){
		if holding.UserID == userID {
			outByteA, err := json.Marshal(holding)
			if err != nil {
				return nil, err
			}
			
			if out == "" {
				out = string(outByteA)
			} else {
				out = out + "," + string(outByteA)
			}
		}
	}
	
	return []byte("[" + out + "]"), nil
}



func (t *ChaincodeBusinessLayer) ballance(stub shim.ChaincodeStubInterface, userID string) ([]byte, error) {
	fmt.Printf("Running ballance")
	
	user, err := t.userRep.getUser(userID)
	if err != nil {
		return nil, err
	}
	
	return []byte(strconv.Itoa(user.getBallance())), nil
}



//********************************************************************************************************
//****                        Invoke function implimentations                                         ****
//********************************************************************************************************



// initial public offering for a square
func (t *ChaincodeBusinessLayer) registerTrade(tradeType string, userID string, securityID string, price float64, units int, expiry string) ([]byte, error) {
	fmt.Printf("Running registerTrade")
	
	var trade Trade
	var err error
	
	securityPointer, err := t.securitiesRep.getSecurityPosition(securityID)
	if err != nil {
		return nil, err
	}
	
	if securityPointer == 0 {
		return nil, errors.New("Security Not Found.")
	}
	
	trade.init(userID, securityID, securityPointer, tradeType, price, units, expiry, "Active", 0)
	
	_, err = t.tradeRep.newTrade(trade)
	if err != nil {
		return nil, err
	}
	
	_, err = t.exchange()
	if err != nil {
		return nil, err
	}
	
	return nil, nil  //what should this function return?
}



// initial public offering for a square
func (t *ChaincodeBusinessLayer) registerSecurity(securityID string, desc string) ([]byte, error) {
	fmt.Printf("Running registerSecurity")
	
	var security Security
	
	security.init(securityID, desc, "Active")
	_, err := t.securitiesRep.newSecurity(security)
	if err != nil {
		return nil, err
	}
	
	return nil, nil  //what should this function return?
}



// called by the moderator watson?  to specify that an event happened pay it out
//todo: need to make dividends payout for each share not just once if there are holdings
func (t *ChaincodeBusinessLayer) dividend(securityID string, amount int) ([]byte, error) {

	fmt.Printf("Running dividend")
	t.writeOut("in dividend")
	holding, _ := t.holdingsRep.getFirstHolding()
	t.writeOut("in dividend: user is " + holding.UserID)
	
	
	//For each holding
	for holding, err := t.holdingsRep.getFirstHolding(); err==nil && holding.UserID != ""; holding, err = t.holdingsRep.getNextHolding(){
		t.writeOut("in dividend Holding for loop")
		if holding.SecurityID == securityID {
				t.writeOut("in dividend- found matching userID.  need to pay out")
			//payout user
			user, err := t.userRep.getUser(holding.UserID)
			if err != nil {
				return nil, err
			}
			
			user.Ballance = user.Ballance + amount
			_,err = t.userRep.updateUser(user)
			if err != nil {
				return nil, err
			}
		}
		
		err = errors.New("Time to exit")
	}
	
	
	t.writeOut("in dividend: before return")
	
	return nil,nil
}



// run on a schedule to execute any pending trades. matching asks with bids and updating the ledger
// first iteration will:
//		only match buyer and seller based on ticker and not on bid and ask prices.  this will simplify and elimiate items we are not trying to prove
//		assume one one share per trade.  this will elimate having to match the number of trades from buy with the sell
//		ignore expiry
//		ignore if the counterparties have the security
//		or if user is active
func (t *ChaincodeBusinessLayer) exchange() ([]byte, error) {
	fmt.Printf("Running exchange")
	t.writeOut("in exchange")
	
	var buyTrade	Trade
	var sellTrade	Trade
	
	numberTrades, err := t.tradeRep.getLastIndex()
	if err != nil {
		return nil, err
	}
	
	t.writeOut("in exchange: before matching loop")
	//trade matching loop
	for b := 1; b <= numberTrades; b++{
		buyTrade, err = t.tradeRep.getTradeByPosition(b)
		if err != nil {
			return nil, err
		}
		
		for s := 1; s <= numberTrades; s++ {
			sellTrade, err = t.tradeRep.getTradeByPosition(s)
			if err != nil {
				return nil, err
			}
			
			if sellTrade.Status == "Active" && buyTrade.Status == "Active" && sellTrade.TransType == "ask" && buyTrade.TransType == "bid" && sellTrade.SecurityID == buyTrade.SecurityID {
				t.writeOut("in exchange: before executeTrade")
				
				_, err := t.executeTrade(b, buyTrade, s, sellTrade)
				
				if err != nil {
					return nil, err
				}
				
				
			}
		}
	}
	
	t.writeOut("in exchange: before return")
	return nil, nil
}



//actually make the trade.  does not vlidate anything.  this should be added at some point:
	//no roll back / transation management
	//doesnt check holdings
	//or users are valid
	//or expiry date etc...
func (t *ChaincodeBusinessLayer) executeTrade(buyTradeIndex int, buyTrade Trade, sellTradeIndex int, sellTrade Trade) ([]byte, error) {
	fmt.Printf("Running exchange")
	
	var buyUser		User
	var sellUser	User
	var buyHolding	Holding
	var sellHolding	Holding
	var err			error
	
	
	t.writeOut("buyTrade.userID = " + buyTrade.getUserID())
	
	//get counterparties (need to change "user" to accounts or counterparties)
	buyUser, err = t.userRep.getUser(buyTrade.getUserID()) 
	if err != nil {
		return nil, err
	}
	
	t.writeOut("buyUser.UserID = " + buyUser.UserID)

	
	sellUser, err = t.userRep.getUser(sellTrade.getUserID()) 
	if err != nil {
		return nil, err
	}
	
	
	//get holdings
	buyHolding, err = t.holdingsRep.getHolding(buyTrade.getUserID(), buyTrade.SecurityID) 
	if err != nil {
		return nil, errors.New("this is the error3")
	}
		
	if buyHolding.UserID == "" {
		//buyer does not already have a holding
		buyHolding.init(buyUser.UserID, buyTrade.SecurityID, 0)
	}
t.writeOut("buyHolding.userID = " + buyUser.UserID)

	if sellTrade.UserID != "BANK" {
		sellHolding, err = t.holdingsRep.getHolding(sellTrade.getUserID(), sellTrade.SecurityID) 
		if err != nil {
			return nil, err
		}
		
		if sellHolding.UserID == "" {  //should not have to check if the seller has the holding but we have not implimented that level of checks yet
			//sell does not already have a holding
			sellHolding.init(sellUser.UserID, sellTrade.SecurityID, 0)
		}
	}

	//transfers funds and closes the trades
	//no transaction rolling back etc...  dont know how best to handle
	buyUser.Ballance = buyUser.Ballance - defaultPrice
	sellUser.Ballance = sellUser.Ballance + defaultPrice
	buyTrade.Status = "Closed"
	sellTrade.Status = "Closed"
	buyHolding.Units = buyHolding.Units + buyTrade.Units  //need to actually determine the units that are being traded
	
	if sellTrade.UserID != "BANK" {
		sellHolding.Units = sellHolding.Units - buyTrade.Units
	}

	//Saves the changes to the trades
	_,err = t.tradeRep.updateTrade(buyTradeIndex, buyTrade)
	if err != nil {
		return nil, err
	}

	_,err = t.tradeRep.updateTrade(sellTradeIndex, sellTrade)
	if err != nil {
		return nil, err
	}
	
	//Saves the changes to the seller and buyer
	_,err = t.userRep.updateUser(buyUser)
	if err != nil {
		return nil, err
	}
	
	_,err = t.userRep.updateUser(sellUser)
	if err != nil {
		return nil, err
	}
	
	t.writeOut("OriginKey: " + t.holdingsRep.LinkedList.originKey + "        userID: " + buyHolding.UserID + buyHolding.SecurityID + strconv.Itoa(buyHolding.Units))
	
	//Save changes to the holdings
	_,err = t.holdingsRep.updateHolding(buyHolding)
	if err != nil {
		return nil, err
	}
	
	t.writeOut("OriginKey: " + t.holdingsRep.LinkedList.originKey + "        userID: " + buyHolding.UserID)
	/*
	if sellTrade.UserID != "BANK" {
		_,err = t.holdingsRep.updateHolding(sellHolding)
		if err != nil {
			return nil, err
		}
	}
	*/
	
	t.writeOut("in execute trade: before return")
	return nil, nil
}



// register user
func (t *ChaincodeBusinessLayer) registerUser(userID string) ([]byte, error) {
	fmt.Printf("Running registerUser")
	//need to make sure the user is not already registered
	newCash := initialCash
	
	if userID == "BANK" {
		newCash = 100000
	} 
	
	index, err := t.userRep.newUser(userID, newCash, "Active")
	if err != nil {
		return nil, err
	}
	
	return []byte(index), nil
}



//   curently not used but should be used in place of taking the user id via the interface.  user id should come from the security model
func (t *ChaincodeBusinessLayer) getUserID(args []string) ([]byte, error) {
	//returns the user's ID 
	
	return nil, nil  //dont know how to get the current user
}










//********************************************************************************************************
//****                              Security implimentations                                          ****
//********************************************************************************************************





func (t *ChaincodeBusinessLayer) isCaller(certificate []byte) (bool, error) {
	// In order to enforce access control, we require that the
	// metadata contains the signature under the signing key corresponding
	// to the verification key inside certificate of
	// the payload of the transaction (namely, function name and args) and
	// the transaction binding (to avoid copying attacks)

	// Verify \sigma=Sign(certificate.sk, tx.Payload||tx.Binding) against certificate.vk
	// \sigma is in the metadata

	sigma, err := t.stub.GetCallerMetadata()
	if err != nil {
		return false, errors.New("Failed getting metadata")
	}
	payload, err := t.stub.GetPayload()
	if err != nil {
		return false, errors.New("Failed getting payload")
	}
	binding, err := t.stub.GetBinding()
	if err != nil {
		return false, errors.New("Failed getting binding")
	}

	t.writeOut("passed certificate " + string(certificate))
	t.writeOut("passed sigma " + string(sigma))
	t.writeOut("passed payload " + string(payload))
	t.writeOut("passed binding " + string(binding))

	ok, err := t.stub.VerifySignature(certificate, sigma, append(payload, binding...), )
	if err != nil {
		t.writeOut("Failed checking signature: ") // + string(err.text))  dont knwo how to convert error ro string
		return ok, err
	}
	
	if !ok {
		t.writeOut("Invalid signature")
	}

	t.writeOut("Check caller...Verified!")

	return ok, err
}