package captainslog

import (
	"fmt"
	"strings"
	"testing"
)

func TestJSONKeyTransformerTransform(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\",\"one.two.three\":\"four.five.six\"}\n")
	original, err := NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	replacer := strings.NewReplacer(".", "_")
	mutator := NewJSONKeyTransformer(replacer)

	mutated, err := mutator.Transform(original)
	if err != nil {
		t.Error(err)
	}

	if want, got := "{\"first_name\":\"captain\",\"one_two_three\":\"four.five.six\"}", mutated.Content; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestJSONKeyTransformerTransformNotCee(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: not a json message\n")
	original, err := NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	replacer := strings.NewReplacer(".", "_")
	mutator := NewJSONKeyTransformer(replacer)

	_, err = mutator.Transform(original)

	if want, got := ErrTransform, err; want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}
}

func TestJSONKeyTransformerTransformMultilevelJSON(t *testing.T) {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"params.name.first\":\"captain\",\"params.name.last\":\"morgan\",\"params.name.ident\":[47],\"sub.arr\":[{\"arr.obj1\":\"val1\"},{\"arr.obj2\":\"val2\"},17],\"sub.obj\":{\"foo.bar\":27,\"bar\":\"baz\"}}\n")
	original, err := NewSyslogMsgFromBytes(b)
	if err != nil {
		t.Error(err)
	}

	replacer := strings.NewReplacer(".", "_")
	mutator := NewJSONKeyTransformer(replacer)

	mutated, err := mutator.Transform(original)
	if err != nil {
		t.Error(err)
	}

	if want, got := "{\"params_name_first\":\"captain\",\"params_name_ident\":[47],\"params_name_last\":\"morgan\",\"sub_arr\":[{\"arr_obj1\":\"val1\"},{\"arr_obj2\":\"val2\"},17],\"sub_obj\":{\"bar\":\"baz\",\"foo_bar\":27}}", mutated.Content; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func ExampleJSONKeyTransformer() {
	b := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\"}\n")
	original, err := NewSyslogMsgFromBytes(b)
	if err != nil {
		panic(err)
	}

	replacer := strings.NewReplacer(".", "_")
	mutator := NewJSONKeyTransformer(replacer)

	mutated, err := mutator.Transform(original)
	if err != nil {
		panic(err)
	}

	fmt.Printf(mutated.Content)
	// Output: {"first_name":"captain"}
}

func BenchmarkJSONKeyTransformerTransform(b *testing.B) {
	m := []byte("<191>2006-01-02T15:04:05.999999-07:00 host.example.org test: @cee:{\"first.name\":\"captain\"}\n")

	replacer := strings.NewReplacer(".", "_")
	mutator := NewJSONKeyTransformer(replacer)

	for i := 0; i < b.N; i++ {
		original, err := NewSyslogMsgFromBytes(m)
		if err != nil {
			panic(err)
		}

		_, err = mutator.Transform(original)
		if err != nil {
			panic(err)
		}
	}
}
