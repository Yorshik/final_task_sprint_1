package calc

import (
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		want        float64
		expectError bool
	}{
		// 1. Простые арифметические операции
		{"Addition", "1+1", 2, false},
		{"Subtraction", "5-3", 2, false},
		{"Multiplication", "4*2", 8, false},
		{"Division", "10/2", 5, false},

		// 2. Операции с приоритетами и скобками
		{"PriorityWithoutParentheses", "2+2*2", 6, false},
		{"PriorityWithParentheses", "(2+2)*2", 8, false},

		// 3. Некорректные выражения
		{"DivisionByZero", "10/0", 0, true},
		{"IncompleteExpression", "2+2*", 0, true},
		{"InvalidCharacters", "2+2*a", 0, true},
		{"UnbalancedParentheses", "(2+3", 0, true},
		{"EmptyExpression", "", 0, true},

		// 4. Сложные выражения (с примерами вашей функции)
		{"Complex expression 1", "2+2*2", 6, false},
		{"Complex expression 2", "(2+2)*2", 8, false},
		{"Complex expression 3", "10/2", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calc(tt.expression)
			if (err != nil) != tt.expectError {
				t.Errorf("Calc() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}
