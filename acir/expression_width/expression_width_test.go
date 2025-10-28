package expression_width

// import (
// 	"os"
// 	"testing"
// )

// func TestExpressionWidthUnmarshalReaderUnbounded(t *testing.T) {
// 	file, err := os.Open("../../binaries/expression/expression_width/expression_width_unbounded.bin")
// 	if err != nil {
// 		t.Fatalf("failed to open file: %v", err)
// 	}

// 	var width ExpressionWidth
// 	if err := width.UnmarshalReader(file); err != nil {
// 		t.Fatalf("failed to unmarshal ExpressionWidth: %v", err)
// 	}

// 	expectedWidth := &ExpressionWidth{
// 		Kind:  ACIRExpressionWidthUnbounded,
// 		Width: nil,
// 	}

// 	if !width.Equals(expectedWidth) {
// 		t.Errorf("expected %v, got %v", expectedWidth, &width)
// 	}

// 	defer file.Close()
// }

// func TestExpressionWidthUnmarshalReaderBounded(t *testing.T) {
// 	file, err := os.Open("../../binaries/expression/expression_width/expression_width_bounded.bin")
// 	if err != nil {
// 		t.Fatalf("failed to open file: %v", err)
// 	}

// 	var width ExpressionWidth
// 	if err := width.UnmarshalReader(file); err != nil {
// 		t.Fatalf("failed to unmarshal ExpressionWidth: %v", err)
// 	}

// 	expectedWidthValue := uint64(10)
// 	expectedWidth := &ExpressionWidth{
// 		Kind:  ACIRExpressionWidthBounded,
// 		Width: &expectedWidthValue,
// 	}

// 	if !width.Equals(expectedWidth) {
// 		t.Errorf("expected %v, got %v", expectedWidth, &width)
// 	}

// 	defer file.Close()
// }
