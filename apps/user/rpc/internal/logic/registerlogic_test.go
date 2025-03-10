package logic

import (
	"context"
	"go-chat/apps/user/rpc/internal/svc"
	"go-chat/apps/user/rpc/user"
	"reflect"
	"testing"
)

func TestNewRegisterLogic(t *testing.T) {
	type args struct {
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	tests := []struct {
		name string
		args args
		want *RegisterLogic
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRegisterLogic(tt.args.ctx, tt.args.svcCtx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRegisterLogic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterLogic_Register(t *testing.T) {

	type args struct {
		in *user.RegisterReq
	}
	tests := []struct {
		name      string
		args      args
		wantPrint bool
		wantErr   bool
	}{
		{
			name: "1", args: args{in: &user.RegisterReq{
				Phone:    "12300000003",
				Nickname: "test",
				Password: "123456",
				Sex:      0,
				Avatar:   "",
			}}, wantPrint: true, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRegisterLogic(context.Background(), svcCtx)
			got, err := l.Register(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantPrint {
				t.Logf(tt.name, got)
			}
		})
	}
}
