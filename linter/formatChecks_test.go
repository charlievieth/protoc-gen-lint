package linter

import "testing"

func TestIsCamelCase(t *testing.T) {
	var stringsToTest = []struct {
		test string
		want bool
	}{
		{"hello_world", false},
		{"HELLO_WORLD", false},
		{"helloWorld", false},
		{"helloworld", false},
		{"HELLOWORLD", false},
		{"HelloWorld", true},
		{"ETA", true},
	}

	for _, v := range stringsToTest {
		if got := isCamelCase(v.test); got != v.want {
			t.Errorf("%s: Expected %t, Received %t", v.test, v.want, got)
		}
	}
}

func BenchmarkIsCamelCase_LongCamel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isCamelCase("OneSuperLongCamelCaseTestName")
	}
}

func BenchmarkIsCamelCase_LongSnake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isCamelCase("platform_version_unknown")
	}
}

func BenchmarkIsLowerUnderscore_LongSnake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isLowerUnderscore("platform_version_unknown")
	}
}

func TestIsLowerUnderscore(t *testing.T) {
	var stringsToTest = []struct {
		test string
		want bool
	}{
		{"hello_world", true},
		{"HELLO_WORLD", false},
		{"helloWorld", false},
		{"helloworld", true},
		{"hello_world", true},
		{"_hello_world", false},
		{"hello_world_", false},
		{"HELLOWORLD", false},
		{"HelloWorld", false},
	}

	for _, v := range stringsToTest {
		if got := isLowerUnderscore(v.test); got != v.want {
			t.Errorf("Expected %t, Received %t", v.want, got)
		}
	}
}

func TestIsUpperUnderscore(t *testing.T) {
	var stringsToTest = []struct {
		test string
		want bool
	}{
		{"hello_world", false},
		{"HELLO_WORLD", true},
		{"_HELLO_WORLD", false},
		{"HELLO_WORLD_", false},
		{"helloWorld", false},
		{"helloworld", false},
		{"HELLOWORLD", true},
		{"HelloWorld", false},
	}

	for _, v := range stringsToTest {
		if got := isUpperUnderscore(v.test); got != v.want {
			t.Errorf("Expected %t, Received %t", v.want, got)
		}
	}
}
