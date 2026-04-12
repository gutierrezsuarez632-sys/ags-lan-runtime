package runtime

import (
	"errors"
	"testing"
)

func TestRuntimeRun_AGSv01(t *testing.T) {
	tests := []struct {
		name          string
		script        string
		wantErr       bool
		expectedCodes []ErrorCode
	}{
		{
			name: "valid fluent script",
			script: `
ctx = createContext("Orders")
sales = ctx["module"]("Sales")
agg = sales.aggregate("Order")
agg.vo("OrderId", "uuid")
agg.id("OrderId")
agg.behavior("create").emits("OrderCreated")
ctx.command("PlaceOrder").handler("Order.create")
`,
			wantErr: false,
		},
		{
			name: "duplicate context",
			script: `
createContext("Orders")
createContext("Orders")
`,
			wantErr:       true,
			expectedCodes: []ErrorCode{ErrDuplicateSymbol},
		},
		{
			name: "aggregate id vo not found",
			script: `
ctx = createContext("Orders")
sales = ctx["module"]("Sales")
agg = sales.aggregate("Order")
agg.id("OrderId")
agg.behavior("create")
ctx.command("PlaceOrder").handler("Order.create")
`,
			wantErr:       true,
			expectedCodes: []ErrorCode{ErrSymbolNotFound},
		},
		{
			name: "chain invalid command after handler",
			script: `
ctx = createContext("Orders")
ctx.command("PlaceOrder").handler("Order.create").middleware("auth")
`,
			wantErr:       true,
			expectedCodes: []ErrorCode{ErrChainInvalid},
		},
		{
			name: "chain invalid provider params after returns",
			script: `
ctx = createContext("Orders")
sales = ctx["module"]("Sales")
sales.provider("PriceProvider").method("getPrice").returns("uuid").params("ProductId")
`,
			wantErr:       true,
			expectedCodes: []ErrorCode{ErrChainInvalid},
		},
		{
			name: "handler unresolved",
			script: `
ctx = createContext("Orders")
sales = ctx["module"]("Sales")
agg = sales.aggregate("Order")
agg.vo("OrderId", "uuid")
agg.id("OrderId")
agg.behavior("create")
ctx.command("PlaceOrder").handler("Order.unknown")
`,
			wantErr:       true,
			expectedCodes: []ErrorCode{ErrHandlerUnresolved},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := NewRuntime()
			err := rt.Run(tt.script)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			if tt.wantErr {
				codes := collectErrorCodes(err)
				for _, expected := range tt.expectedCodes {
					if !containsCode(codes, expected) {
						t.Fatalf("expected code %s in %v, got %v (err=%v)", expected, tt.expectedCodes, codes, err)
					}
				}
			}
		})
	}
}

func collectErrorCodes(err error) []ErrorCode {
	if err == nil {
		return nil
	}

	codes := []ErrorCode{}

	var single *AGSError
	if errors.As(err, &single) {
		codes = append(codes, single.Code)
	}

	var multi *AGSMultiError
	if errors.As(err, &multi) {
		for _, item := range multi.Errors {
			codes = append(codes, item.Code)
		}
	}

	return codes
}

func containsCode(codes []ErrorCode, expected ErrorCode) bool {
	for _, code := range codes {
		if code == expected {
			return true
		}
	}
	return false
}
