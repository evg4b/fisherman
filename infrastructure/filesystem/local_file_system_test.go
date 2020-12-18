package filesystem

import (
	"io"
	"os"
	"os/user"
	"reflect"
	"testing"
)

func TestNewLocalFileSystem(t *testing.T) {
	tests := []struct {
		name string
		want *LocalFileSystem
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLocalFileSystem(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLocalFileSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalFileSystem_Exist(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		f    *LocalFileSystem
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Exist(tt.args.path); got != tt.want {
				t.Errorf("LocalFileSystem.Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalFileSystem_Read(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		f       *LocalFileSystem
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Read(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalFileSystem.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LocalFileSystem.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalFileSystem_Reader(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		f       *LocalFileSystem
		args    args
		want    io.Reader
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Reader(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalFileSystem.Reader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalFileSystem.Reader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalFileSystem_Delete(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		f       *LocalFileSystem
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Delete(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("LocalFileSystem.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalFileSystem_Write(t *testing.T) {
	type args struct {
		path    string
		content string
	}
	tests := []struct {
		name    string
		f       *LocalFileSystem
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Write(tt.args.path, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("LocalFileSystem.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalFileSystem_Chmod(t *testing.T) {
	type args struct {
		path string
		mode os.FileMode
	}
	tests := []struct {
		name    string
		f       *LocalFileSystem
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Chmod(tt.args.path, tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("LocalFileSystem.Chmod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalFileSystem_Chown(t *testing.T) {
	type args struct {
		path string
		user *user.User
	}
	tests := []struct {
		name    string
		f       *LocalFileSystem
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Chown(tt.args.path, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("LocalFileSystem.Chown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
