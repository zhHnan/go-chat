package bitmap

import (
	"reflect"
	"testing"
)

func TestBitmap_Export(t *testing.T) {
	type fields struct {
		bits []byte
		size int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bitmap{
				bits: tt.fields.bits,
				size: tt.fields.size,
			}
			if got := b.Export(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Export() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitmap_IsSet(t *testing.T) {
	type fields struct {
		bits []byte
		size int
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bitmap{
				bits: tt.fields.bits,
				size: tt.fields.size,
			}
			if got := b.IsSet(tt.args.id); got != tt.want {
				t.Errorf("IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitmap_Set(t *testing.T) {
	b := NewBitmap(100)
	b.Set("pppp")
	b.Set("222")
	b.Set("pppp")
	b.Set("ccc")

	for _, bit := range b.bits {
		t.Logf("%b, %v", bit, bit)
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want *Bitmap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Load(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBitmap(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want *Bitmap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBitmap(tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBitmap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hash(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hash(tt.args.id); got != tt.want {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
