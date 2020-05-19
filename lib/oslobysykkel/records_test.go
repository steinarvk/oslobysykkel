package oslobysykkel

import (
	"testing"
)

func TestReadStatic(t *testing.T) {
	_, err := NewStaticDataSource("../../testdata/station_status.json", "../../testdata/station_information.json")
	if err != nil {
		t.Error(err)
	}
}
