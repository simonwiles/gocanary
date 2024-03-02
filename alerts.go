package main

import (
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type alertRule struct {
	Rule    string
	Program *vm.Program
}

func compileAlerts(alertRules *[]string) ([]alertRule, error) {
	alerts := make([]alertRule, len(*alertRules))
	for i, ruleStr := range *alertRules {
		program, err := expr.Compile(ruleStr, expr.Env(response{}))
		if err != nil {
			return nil, err
		}
		alerts[i] = alertRule{
			Rule:    ruleStr,
			Program: program,
		}
	}
	return alerts, nil
}
