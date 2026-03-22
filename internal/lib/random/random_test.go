package random

import (
	"testing"
	"fmt"
)

func TestNewRandomString(t *testing.T) {
	tests:= []struct{
		length int
		want string
	}{
		{length: 10, want: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"},
	}
	for _, test:=range tests{
		t.Run(fmt.Sprintf("length %d", test.length), func(t *testing.T) {
			got:= NewRandomString(test.length)
			if len(got)!= test.length {
				t.Errorf("NewRandomString(%d) = %s, want %s", test.length, got, test.want)
			}
		})
	}
}
