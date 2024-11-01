package utils

import (
	"testing"
)

func TestParseLine(t *testing.T) {

	tests := []struct {
		name          string
		line          string
		wantNamespace string
		wantName      string
	}{
		{
			name:          "TEST1",
			line:          "NAMESPACE           NAME                                      READY   UP-TO-DATE   AVAILABLE   AGE",
			wantNamespace: "NAMESPACE",
			wantName:      "NAME",
		},
		{
			name:          "TEST2",
			line:          "adam                adam-asset                                1/1     1            1           2y218d",
			wantNamespace: "adam",
			wantName:      "adam-asset",
		},
		{
			name:          "TEST3",
			line:          "cattle-prometheus-p-6n24k   grafana-project-monitoring                0/0     0            0           495d",
			wantNamespace: "cattle-prometheus-p-6n24k",
			wantName:      "grafana-project-monitoring",
		},
		{
			name:          "TEST3",
			line:          "openplatform                ingress-df4186ef1375c28f664e233567db3e57          ClusterIP      10.43.41.69     <none>                      3001/TCP             145d",
			wantNamespace: "openplatform",
			wantName:      "ingress-df4186ef1375c28f664e233567db3e57",
		},
		{
			name:          "TEST3",
			line:          "habitat              habitat.develop2.github.com     habitat2.develop.github.com                        10.200.2.90,10.200.2.91,10.200.2.92,10.200.2.93,10.200.2.95   80        391d",
			wantNamespace: "habitat",
			wantName:      "habitat.develop2.github.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := ParseLine(tt.line)
			gotNamespace := ss[0]
			gotName := ss[1]
			if gotNamespace != tt.wantNamespace {
				t.Errorf("ParseLine() gotNamespace = %v, want %v", gotNamespace, tt.wantNamespace)
			}
			if gotName != tt.wantName {
				t.Errorf("ParseLine() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
