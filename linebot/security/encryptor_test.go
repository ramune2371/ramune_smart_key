package security_test

import (
	"fmt"
	"linebot/props"
	"linebot/security"
	"linebot/testutil"
	"testing"
)

func TestSalHash(t *testing.T) {
	tests := []struct {
		salt   string
		target string
		expect string
	}{
		{
			salt:   "TEST",
			target: "test",
			expect: "047feec539d0d561765a913ac3d3ad0b7188df00eb32cddae34077afd15c8af9",
		},
		{
			salt:   "TEST",
			target: "TEST",
			expect: "fc344ccc84c06056a9a06ffa48de548abc309913722b9da21c8c68ccc31daf49",
		},
		{
			salt:   "test",
			target: "test",
			expect: "45b79a422b2b93633983ee16aeb90f36ad054174e491e2a2d528daab75a263a4",
		},
		{
			salt:   "test",
			target: "TEST",
			expect: "47eb4b196f412b206309b4d5391a3e667f4366fb156682e00bfd0da1dc332c4a",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test SaltHash(Salt: %s, Input: %s", test.salt, test.target), func(t *testing.T) {
			props.Salt = test.salt
			ret := security.EncryptorImpl{}.SaltHash(test.target)
			if ret != test.expect {
				t.Errorf(testutil.STRING_TEST_MSG_FMT, "", test.expect, ret)
			}
		})
	}
}
