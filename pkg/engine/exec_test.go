package engine

import (
	"github.com/storm-blue/rubick/pkg/modifier/action"
	"testing"
)

func TestExec(t *testing.T) {
	type args struct {
		ctx     action.Context
		yaml    string
		scripts string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TEST1",
			args: args{
				ctx: action.NewContext(nil),
				yaml: `apiVersion: v1
kind: Service
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"name":"dev-java-user-auth-core-service-svc","namespace":"java-dev"},"spec":{"ports":[{"port":8080,"protocol":"TCP","targetPort":8080}],"selector":{"app":"dev-java-user-auth-core-service"}}}
  creationTimestamp: "2023-09-25T06:30:18Z"
  labels:
    appName: user-auth-core-service
    devops.aliyun.com/app-name: user-auth-core-service
    devops.aliyun.com/env-name: dev
    devops.aliyun.com/org-id: 626253aaa60d8a4bbe1753f0
    devops.aliyun.com/version: 20241019152519-752
    group: dev
  name: dev-java-user-auth-core-service-svc
  namespace: java-dev
  resourceVersion: "2785940354"
  uid: dbc5e023-b4be-4f4b-910a-548d05157916
spec:
  clusterIP: 10.100.72.127
  clusterIPs:
  - 10.100.72.127
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: dev-java-user-auth-core-service
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}
`,
				scripts: `
DELETE(metadata.annotations.(kubectl.kubernetes.io/last-applied-configuration))
DELETE(metadata.creationTimestamp)
DELETE(metadata.resourceVersion)
DELETE(metadata.uid)
DELETE(status)
IF VALUE_OF(kind)=="Service" THEN DELETE(spec.clusterIP)
IF VALUE_OF(kind)=="Service" THEN DELETE(spec.clusterIPs)
IF VALUE_OF(kind)=="Service" THEN SET(spec.ports[port=8080].port, 80)
IF VALUE_OF(kind)=="Service" THEN SET_WITH_VALUE_OF(metadata.name, metadata.labels.(devops.aliyun.com/app-name))
`,
			},
			want: `apiVersion: v1
kind: Service
metadata:
  annotations: {}
  labels:
    appName: user-auth-core-service
    devops.aliyun.com/app-name: user-auth-core-service
    devops.aliyun.com/env-name: dev
    devops.aliyun.com/org-id: 626253aaa60d8a4bbe1753f0
    devops.aliyun.com/version: 20241019152519-752
    group: dev
  name: user-auth-core-service
  namespace: java-dev
spec:
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: dev-java-user-auth-core-service
  sessionAffinity: None
  type: ClusterIP
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exec(tt.args.ctx, tt.args.yaml, tt.args.scripts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exec()\n----------------- got -----------------\n%v\n----------------- want -----------------\n%v\n", got, tt.want)
			}
		})
	}
}
