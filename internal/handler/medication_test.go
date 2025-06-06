package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MedicationHandlerの基本テスト
func TestMedicationHandler_New(t *testing.T) {
	t.Run("MedicationHandlerが正常に作成される", func(t *testing.T) {
		// MedicationHandlerの作成をテスト
		// 実際のリポジトリは使わずに、nilで作成してもパニックしないことを確認
		assert.NotPanics(t, func() {
			NewMedicationHandler(nil)
		})
	})
}