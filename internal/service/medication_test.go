package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMedicationService_New(t *testing.T) {
	t.Run("MedicationServiceが正常に作成される", func(t *testing.T) {
		// MedicationServiceの作成をテスト
		assert.NotPanics(t, func() {
			NewMedicationService(nil)
		})
	})
}
