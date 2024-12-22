package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Pop(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func Contains(lst []string, s string) bool {
	for _, el := range lst {
		if el == s {
			return true
		}
	}
	return false
}

func EvaluateLst(lst_expr []string, primary string) []string {
	for i, el := range lst_expr {
		if strings.Contains(primary, el) {
			var result float64
			left, _ := strconv.ParseFloat(lst_expr[i-1], 64)
			right, _ := strconv.ParseFloat(lst_expr[i+1], 64)
			switch el {
			case "+":
				result = left + right
			case "-":
				result = left - right
			case "*":
				result = left * right
			case "/":
				result = left / right
			}
			lst_expr[i] = fmt.Sprintf("%f", result)
			lst_expr = Pop(lst_expr, i+1)
			lst_expr = Pop(lst_expr, i-1)
			break
		}
	}
	for _, pr := range primary {
		if Contains(lst_expr, string(pr)) {
			return EvaluateLst(lst_expr, primary)
		}
	}
	return lst_expr
}

func GetBrackets(expression string) []string {
	var brackets []string
	var s strings.Builder
	bracket_count := 0
	for _, el := range expression {
		if el == '(' {
			bracket_count = bracket_count + 1
			if bracket_count == 1 {
				continue
			}
		}
		if el == ')' {
			bracket_count = bracket_count - 1
			if bracket_count == 0 {
				brackets = append(brackets, s.String())
				s.Reset()
				continue
			}
		}
		if bracket_count == 0 {
			continue
		}
		s.WriteRune(el)
	}
	return brackets
}

func ContainsValidBrackets(expression string) (bool, error) {
	if strings.Contains(expression, "(") || strings.Contains(expression, ")") {
		bracket_flag := 0
		for _, el := range expression {
			if el == '(' {
				bracket_flag = bracket_flag + 1
			} else if el == ')' {
				bracket_flag = bracket_flag - 1
			}
		}
		if bracket_flag != 0 {
			return false, errors.New("incorrect brackets")
		}
		return true, nil
	} else {
		return false, nil
	}
}

func EvaluateExpr(expression string) (string, error) {
	var lst []string
	if strings.HasPrefix(expression, "*") ||
		strings.HasPrefix(expression, "+") ||
		strings.HasPrefix(expression, "-") ||
		strings.HasPrefix(expression, "/") ||
		strings.HasSuffix(expression, "*") ||
		strings.HasSuffix(expression, "+") ||
		strings.HasSuffix(expression, "-") ||
		strings.HasSuffix(expression, "/") {
		return "", errors.New("leading/ending operator")
	}
	flag := "number"
	var dot_flag bool
	var s strings.Builder
	for _, sym := range expression {
		if strings.Contains("+-*/", string(sym)) {
			if flag == "number" {
				flag = "operator"
				lst = append(lst, s.String())
				s.Reset()
			} else {
				return "", errors.New("double operator")
			}
		} else if sym == '.' {
			if dot_flag {
				return "", errors.New("attempting to add another dot")
			}
			dot_flag = true
		} else {
			if flag == "operator" {
				lst = append(lst, s.String())
				s.Reset()
			}
			flag = "number"
			dot_flag = false
		}
		s.WriteRune(sym)
	}
	lst = append(lst, s.String())
	primary := "*/"
	new_lst := EvaluateLst(lst, primary)
	primary = "+-"
	res := EvaluateLst(new_lst, primary)[0]
	return res, nil
}

func EvaluateBrackets(expression string) (string, error) {
	var expr_containg_brackets bool
	var err error
	var brackets []string
	var bracket_res string
	var res string
	expr_containg_brackets, err = ContainsValidBrackets(expression)
	if err != nil {
		return "", err
	}
	if expr_containg_brackets {
		brackets = GetBrackets(expression)
		for _, bracket := range brackets {
			bracket_res, err = EvaluateBrackets(bracket)
			if err != nil {
				return "", err
			}
			expression = strings.Replace(expression, fmt.Sprintf("(%s)", bracket), bracket_res, 1)
		}
	}
	res, err = EvaluateExpr(expression)
	if err != nil {
		return "", err
	}
	return res, nil
}

func AssertExpressionContainsCorrectSymbols(expression string) bool {
	for _, char := range expression {
		if unicode.IsSpace(char) {
			continue
		}
		if unicode.IsDigit(char) {
			continue
		}
		switch char {
		case '+', '-', '*', '/', '(', ')':
			continue
		default:
			return false
		}
	}
	return true
}

func Calc(expression string) (float64, error) {
	var new_expr string
	var err error
	if expression == "" {
		return 0, errors.New("empty expression")
	}
	if !AssertExpressionContainsCorrectSymbols(expression) {
		return 0, errors.New("invalid expression")
	}
	expression = strings.ReplaceAll(expression, " ", "")
	new_expr, err = EvaluateBrackets(expression)
	if err != nil {
		return 0, err
	}
	var res string
	res, err = EvaluateExpr(new_expr)
	if err != nil {
		return 0, err
	}
	result, _ := strconv.ParseFloat(res, 64)
	return result, nil
}
