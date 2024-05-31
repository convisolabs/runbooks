package functions

import "sync"

var lock = &sync.Mutex{}

var iFunctions IFunctions

func GetFunctionsSingletonInstance() IFunctions {
	if iFunctions == nil {
		lock.Lock()
		defer lock.Unlock()
		if iFunctions == nil {

			iFunctions = FunctionsNew()
		}
	}
	return iFunctions
}
