package validation

import (
  "regexp"
  "testing"
  "github.com/viki-org/gspec"
)

func TestErrorWhenPatterDoesNotMatch(t *testing.T) {
  spec := gspec.New(t)
  spec.Expect(Pattern(regexp.MustCompile("over 9000!?")).Verify("over 9001")).ToEqual(false)
}

func TestOkWhenPatterMatches(t *testing.T) {
  spec := gspec.New(t)
  spec.Expect(Pattern(regexp.MustCompile("over 9000!?")).Verify("over 9000")).ToEqual(true)
  spec.Expect(Pattern(regexp.MustCompile("over 9000!?")).Verify("over 9000!")).ToEqual(true)
}

func TestInvalidEmails(t *testing.T) {
  spec := gspec.New(t)
  for _, email := range []string{"", "a", "a@", "a@a"} {
    spec.Expect(Email().Verify(email)).ToEqual(false)
  }
}

func TestValidEmails(t *testing.T) {
  spec := gspec.New(t)
  for _, email := range []string{"p@d.g", "leto@dune.gov", "Irulan.Corrino@salusa-secundus.museum"} {
    spec.Expect(Email().Verify(email)).ToEqual(true)
  }
}
