package eureka

import "testing"

func TestGetApplications(t *testing.T) {
	for {
		a, b := GetApplications("http://127.0.0.1:8761")
		t.Logf("%v---%v", a, b)
	}
}
