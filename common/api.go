package common

import (
	"strconv"
	"errors"
	"encoding/json"
)

//Host host
var Host = "54.188.118.102:8732"

//CurrentLevelAt func
func CurrentLevelAt(hash string) CurrentLevelType {
	var curr CurrentLevelType
	route := fmt.Sprintf("/chains/main/blocks/%s/helpers/current_level", hash)
	body := Get(Host, route)
	json.Unmarshal(body, &curr)
	return curr
}

//CurrentLevel func
func CurrentLevel() CurrentLevelType {
	return CurrentLevelAt("head")
}

//HeaderAt func
func HeaderAt(hash string) BlockHeaderType {
	var blockheader BlockHeaderType
	route := fmt.Sprintf("/chains/main/blocks/%s/header", hash)
	body := Get(Host, route)
	json.Unmarshal(body, &blockheader)
	return blockheader
}

//Header func
func Header() BlockHeaderType {
	return HeaderAt("head")
}

//CycleInfo func
func CycleInfo(hash string, cycle int) CycleInfoType {
	var cycleInfo CycleInfoType
	route := fmt.Sprintf("/chains/main/blocks/%s/context/raw/json/cycle/%d", hash, cycle)
	body := Get(Host, route)
	json.Unmarshal(body, &cycleInfo)
	return cycleInfo
}

//DelegatedContracts func
func DelegatedContracts(hash string, delegate string) []string {
	var contracts []string
	route := fmt.Sprintf("/chains/main/blocks/%s/context/delegates/%s/delegated_contracts", hash, delegate)
	body := Get(Host, route)
	json.Unmarshal(body, &contracts)
	return contracts
}

//BalanceAt func
func BalanceAt(hash string, delegate string) string {
	var balance string
	route := fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/balance", hash, delegate)
	body := Get(Host, route)
	json.Unmarshal(body, &balance)
	return balance
}

//DelegateBalanceAt func
func DelegateBalanceAt(hash string, delegate string) string {
	var balance string
	route := fmt.Sprintf("/chains/main/blocks/%s/context/delegates/%s/balance", hash, delegate)
	body := Get(Host, route)
	json.Unmarshal(body, &balance)
	return balance
}

//FrozenBalanceByCycle func
func FrozenBalanceByCycle(hash string, delegate string) []FrozenBalanceByCycleType {
	var balance []FrozenBalanceByCycleType
	route := fmt.Sprint("/chains/main/blocks/%s/context/delegates/%s/frozen_balance_by_cycle", hash, delegate)
	body := Get(Host, route)
	json.Unmarshal(body, &balance)
	return balance
}

//StakingBalanceAt func
func StakingBalanceAt(hash string, delegate string) string {
	var balance string
	route := fmt.Sprint("/chains/main/blocks/%s/context/delegates/%s/staking_balance", hash, delegate)
	body := Get(Host, route)
	json.Unmarshal(body, &balance)
	return balance
}

//BakingRightsFor func
func BakingRightsFor(hash string, delegate string, cycle int) []BakingRightType {
	var bakingRights []BakingRightType
	route := fmt.Sprint("/chains/main/blocks/%s/helpers/baking_rights?delegate=%s&cycle=%d", hash, delegate, cycle)
	body := Get(Host, route)
	json.Unmarshal(body, &bakingRights)
	return bakingRights
}

//EndorsingRightsFor func
func EndorsingRightsFor(hash string, delegate string, cycle int) []EndorsingRightType {
	var endorsingRights []EndorsingRightType
	route := fmt.Sprint("/chains/main/blocks/%s/helpers/endorsing_rights?delegate=%s&cycle=%d", hash, delegate, cycle)
	body := Get(Host, route)
	json.Unmarshal(body, &endorsingRights)
	return endorsingRights
}

//Metadata func
func Metadata(hash string) BlockMetadataType {
	var metadata BlockMetadataType
	route := fmt.Sprint("/chains/main/blocks/%s/metadata", hash)
	body := Get(Host, route)
	json.Unmarshal(body, &metadata)
	return metadata
}

//Operations func
func Operations(hash string) OperationType {
	var operations OperationType
	route := fmt.Sprint("/chains/main/blocks/%s/operations", hash)
	body := Get(Host, route)
	json.Unmarshal(body, &operations)
	return operations
}

//Counter func
func Counter(hash string, contract string) int {
	var value string
	route := fmt.Sprint("/chains/main/blocks/%s/context/contracts/%s/counter", hash, contract)
	body := Get(Host, route)
	json.Unmarshal(body, &value)
	counter, _ := strconv.Atoi(value)
	return counter
}

//RunOperation func
func RunOperation(hash string, data string) (OperationContentsAndResultsType, error) {
	var result OperationContentsAndResultsType
	route := fmt.Sprint("/chains/main/blocks/%s/helpers/scripts/run_operation", head)
	body := Post(Host, route, data)
	if json.Unmarshal(body, &result) == nil {
		return result, nil
	}
	return result, errors.New(string(body))
}

//ForgeOperations func
func ForgeOperations(data string) (string, error) {
	var result string
	route := fmt.Sprint("/chains/main/blocks/%s/helpers/forge/operations", head)
	body := Post(Host, route, data)
	if json.Unmarshal(body, &result) == nil {
		return result, nil
	}
	return result, errors.New(string(body))
}

//PreapplyOperations func
func PreapplyOperations(data string) ([]PreapplyResultType, error) {
	var result []PreapplyResultType
	route := fmt.Sprint("/chains/main/blocks/%s/helpers/preapply/operations", head)
	body := Post(Host, route, data)
	if json.Unmarshal(body, &result) == nil {
		return result, nil
	}
	return result, errors.New(string(body))
}

//Injection func
func Injection(data string) (string, error) {
	var result string
	route := "/injection/operation?chain=main"
	body := Post(Host, route, data)
	if json.Unmarshal(body, &result) == nil {
		return result, nil
	}
	return result, errors.New(string(body))
}
