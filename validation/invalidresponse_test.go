package validation

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestGeneratesaProperJsonMessage(t *testing.T) {
  InitInvalidPool(1, 1024)
  spec := gspec.New(t)
  response := newResponse(map[string][]*Definition {
    "usern\"ame": []*Definition {
      Define("un.req").Field("username").Message("username is required"),
      Define("un.dupe").Field("username").Message("username is already taken"),
    },
    "password": []*Definition {
      Define("pw.len").Field("password").Message("is too\", short"),
    },
  })
  spec.Expect(string(response.Body())).ToEqual(`{"usern\"ame":["username is required","username is already taken"],"password":["is too\", short"]}`)
  spec.Expect(response.Header().Get("Content-Length")).ToEqual("98")
}


func TestCloseReleasesTheBuffer(t *testing.T) {
  InitInvalidPool(2, 1024)
  spec := gspec.New(t)
  response := newResponse(map[string][]*Definition {
    "password": []*Definition {
      Define("pw.len").Field("password").Message("is too\", short"),
    },
  })
  spec.Expect(invalidBytePool.Len()).ToEqual(1)
  response.Close()
  spec.Expect(invalidBytePool.Len()).ToEqual(2)
}
