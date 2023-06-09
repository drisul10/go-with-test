package reflection

import (
	"reflect"
	"testing"
)

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func TestWalk(t *testing.T) {
	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"Struct with one string field",
			struct {
				Name string
			}{"Salma"},
			[]string{"Salma"},
		},
		{
			"Struct with two string fields",
			struct {
				Name string
				City string
			}{"Salma", "Bandung"},
			[]string{"Salma", "Bandung"},
		},
		{
			"Struct with integer field",
			struct {
				Name string
				Age  int
			}{"Salma", 23},
			[]string{"Salma"},
		},
		{
			"Struct with nested fields",
			Person{
				"Salma",
				Profile{23, "Bandung"},
			},
			[]string{"Salma", "Bandung"},
		},
		{
			"Pointer to things",
			&Person{
				"Salma",
				Profile{23, "Bandung"},
			},
			[]string{"Salma", "Bandung"},
		},
		{
			"Slices",
			[]Profile{
				{23, "Bandung"},
				{25, "Klaten"},
			},
			[]string{"Bandung", "Klaten"},
		},
		{
			"Arrays",
			[2]Profile{
				{23, "Bandung"},
				{25, "Klaten"},
			},
			[]string{"Bandung", "Klaten"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %q, want %q", got, test.ExpectedCalls)
			}
		})
	}

	t.Run("With maps", func(t *testing.T) {
		aMap := map[string]string{
			"Foo": "Bar",
			"Baz": "Boz",
		}

		var got []string
		walk(aMap, func(input string) {
			got = append(got, input)
		})

		assertContains(t, got, "Bar")
		assertContains(t, got, "Boz")
	})

	t.Run("With channels", func(t *testing.T) {
		aChan := make(chan Profile)

		go func() {
			aChan <- Profile{23, "Bandung"}
			aChan <- Profile{25, "Klaten"}
			close(aChan)
		}()

		var got []string
		want := []string{"Bandung", "Klaten"}

		walk(aChan, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("With function", func(t *testing.T) {
		aFunc := func() (Profile, Profile) {
			return Profile{23, "Bandung"}, Profile{25, "Klaten"}
		}

		var got []string
		want := []string{"Bandung", "Klaten"}

		walk(aFunc, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func assertContains(t testing.TB, haystack []string, needle string) {
	t.Helper()
	contains := false

	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}

	if !contains {
		t.Errorf("expected %v to contain %q, but it didn't", haystack, needle)
	}
}
