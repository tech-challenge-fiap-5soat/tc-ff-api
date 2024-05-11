package valueobject

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

func TestCPF(t *testing.T) {
	t.Run("should return true when cpf is valid", func(t *testing.T) {
		cpf := CPF("19119119100")

		assert.True(t, cpf.IsValid())
	})
	t.Run("should return false when cpf is invalid", func(t *testing.T) {
		cpf := CPF("12345678910")

		assert.False(t, cpf.IsValid())
	})

	t.Run("should return false when cpf has eleven digits equal", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			cpf := CPF(strings.Repeat(strconv.Itoa(i), 11))

			assert.False(t, cpf.IsValid())
		}

	})
}
