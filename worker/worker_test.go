package worker

import (
	"reflect"
	"testing"
)

func TestNewWorker(t *testing.T) {
	type args struct {
		name     string
		handle   func() error
		replicas int
	}
	tests := []struct {
		name string
		args args
		want *Worker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWorker(tt.args.name, tt.args.handle, tt.args.replicas); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunWorkers(t *testing.T) {
	type args struct {
		workers     []*Worker
		handleError func(w *Worker, err error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunWorkers(tt.args.workers, tt.args.handleError); (err != nil) != tt.wantErr {
				t.Errorf("RunWorkers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorker_Run(t *testing.T) {

}