package base

import (
	//"fmt"
)

func (cxt *cxContext) GetCurrentModule () *cxModule {
	if cxt.CurrentModule != nil {
		return cxt.CurrentModule
	} else {
		return nil
	}
	
}

func (cxt *cxContext) GetCurrentStruct () *cxStruct {
	if cxt.CurrentModule != nil &&
		cxt.CurrentModule.CurrentStruct != nil {
		return cxt.CurrentModule.CurrentStruct
	} else {
		return nil
	}
	
}

func (cxt *cxContext) GetCurrentFunction () *cxFunction {
	if cxt.CurrentModule != nil &&
		cxt.CurrentModule.CurrentFunction != nil {
		return cxt.CurrentModule.CurrentFunction
	} else {
		return nil
	}
	
}

func (cxt *cxContext) GetCurrentExpression () *cxExpression {
	if cxt.CurrentModule != nil &&
		cxt.CurrentModule.CurrentFunction != nil &&
		cxt.CurrentModule.CurrentFunction.CurrentExpression != nil {
		return cxt.CurrentModule.CurrentFunction.CurrentExpression
	} else {
		return nil
	}
}
