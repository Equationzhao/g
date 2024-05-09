package content

import (
	"testing"
)

func TestParseSize(t *testing.T) {
	testHelper := func(sizeStr string, sizeAmount float64, unit SizeUnit) {
		size, err := ParseSize(sizeStr)
		if err != nil {
			t.Errorf("%s:ParseSize() error = %v", sizeStr, err)
		}
		expected := uint64(sizeAmount * float64(CountBytes(unit)))
		if size.Bytes != expected {
			t.Errorf("%s:ParseSize() got = %v, want %v", sizeStr, size.Bytes, expected)
		}
	}

	_, err := ParseSize("unknown")
	// todo define error struct
	t.Log(err)

	_, err = ParseSize("-1kb")
	t.Log(err)

	testHelper("123byte", 123, Byte)
	testHelper("1,234.5KB", 1_234.5, KB)
	testHelper("1,234.5KiB", 1_234.5, KiB)
	testHelper("1,234.5K", 1_234.5, KB)
	testHelper("1,234.5k", 1_234.5, KB)
	testHelper("1,234.5kb", 1_234.5, KB)
	testHelper("4.5GB", 4.5, GB)
	testHelper("4.5GiB", 4.5, GiB)
	testHelper("4.5G", 4.5, GB)
	testHelper("4.5g", 4.5, GB)
	testHelper("4.5gb", 4.5, GB)
	testHelper("200MB", 200, MB)
	testHelper("200MiB", 200, MiB)
	testHelper("200M", 200, MB)
	testHelper("200mb", 200, MB)
	testHelper("200m", 200, MB)
	testHelper("3,123,432.321TB", 3_123_432.321, TB)
	testHelper("3,123,432.321TiB", 3_123_432.321, TiB)
	testHelper("3,123,432.321T", 3_123_432.321, TB)
	testHelper("3,123,432.321tb", 3_123_432.321, TB)
	testHelper("3,123,432.321t", 3_123_432.321, TB)
}
