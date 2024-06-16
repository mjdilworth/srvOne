package apploader

import (
	"context"
	"fmt"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/mjdilworth/srvOne/internal/server"
)

// Helper functions
func returnEmptyContextWithCancel() context.Context {
	ctx, _ := context.WithCancel(context.Background())
	return ctx
}

func returnCancelFunc() context.CancelFunc {
	_, cancel := context.WithCancel(context.Background())
	return cancel
}

func newTestLoader() *AppLoader {
	// An instance of the application will only run for max 100ms for test purposes
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100)*time.Millisecond)

	return &AppLoader{
		app:    testApp{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func sendSignalToRun(t *testing.T, termSignal syscall.Signal, millisecBeforeKill int) bool {
	l := newTestLoader()
	defer l.cancel()

	go func() {
		_ = l.Run() // using mock app
	}()

	exitChan := make(chan bool)

	go func() {
		for {
			select {
			case <-l.ctx.Done():
				if l.ctx.Err() == context.DeadlineExceeded {
					// The application was still running after 5s
					t.Log("Context deadline exceeded")
					exitChan <- false
					return
				}
				// termSignal caused Run to exit!
				t.Log("Application correctly terminated")
				exitChan <- true
				return
			case <-time.After(time.Duration(10) * time.Millisecond):
				t.Log("Waiting for the application to terminate...")
			}
		}
	}()

	time.Sleep(time.Duration(millisecBeforeKill) * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), termSignal)

	return <-exitChan
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancelCause(ctx)
	time.AfterFunc(time.Second, func() { cancel(fmt.Errorf("cancelled from testing")) })

	err := Run(ctx, "7777")
	if err != nil {
		fmt.Printf("expected no error, but got: %s", err)
	}
}

// New()
func TestNew(t *testing.T) {

	a := server.New("8080")
	loader := New(a)

	typeOf := reflect.TypeOf(loader)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if field.Type.Kind() == reflect.Func || (field.Type.Kind() == reflect.Interface) {
			fieldValue := reflect.ValueOf(loader).FieldByName(field.Name)
			if fieldValue.IsNil() {
				t.Errorf("Missing value for field: %s ", field.Name)
			}

		}
	}

	//assert.IsType(t, &AppLoader{}, loader, "Loader not created")
	//assert.IsType(t, returnCancelFunc(), loader.cancel, "CancelFunc not created when loader gets created")
}

/*
// New()
func TestNew(t *testing.T) {
	loader := New()

	want := &AppLoader{}

	assert.IsType(t, &AppLoader{}, loader, "Loader not created")
	assert.IsType(t, returnCancelFunc(), loader.cancel, "CancelFunc not created when loader gets created")
}


// Run()
func TestRunWithWrongMode(t *testing.T) {
	err := NewLoader().Run("WrongMode")
	assert.Error(t, err, "Passing mode 'simon' should have caused an error", err)
}

func TestRunWithSIGTERM(t *testing.T) {
	assert.True(t, sendSignalToRun(t, syscall.SIGTERM, 30), "SIGTERM didn't stop the process after 30ms")
}

func TestRunWithSIGINT(t *testing.T) {
	assert.True(t, sendSignalToRun(t, syscall.SIGINT, 30), "SIGINT didn't stop the process after 30ms")
}

func TestRunWithSIGWINCH(t *testing.T) {
	assert.False(t, sendSignalToRun(t, syscall.SIGWINCH, 30), "SIGWINCH caused the process to stop after 30ms")
}

func TestRunWithNoSignal(t *testing.T) {
	assert.False(t, sendSignalToRun(t, syscall.Signal(0), 0), "Application doesn't seem to run")
}
*/
