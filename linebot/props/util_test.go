package props_test

import (
	"linebot/props"
	"linebot/testutil"
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("String Environment", func(t *testing.T) {
		tests := []struct {
			description string
			target      map[string]string
			expect      map[string]string
		}{
			{
				description: "empty",
				target: map[string]string{
					// key_server
					"KEY_SERVER_URL": "",

					// line
					"CHANNEL_SECRET": "",
					"CHANNEL_TOKEN":  "",

					// security
					"SALT": "",
				},
				expect: map[string]string{
					"KEY_SERVER_URL": "http://localhost:9999/",
					"CHANNEL_SECRET": "channelSecret",
					"CHANNEL_TOKEN":  "channelToken",
					"SALT":           "test",
				},
			},
			{
				description: "not empty",
				target: map[string]string{
					// key_server
					"KEY_SERVER_URL": "http://test.test",

					// line
					"CHANNEL_SECRET": "testValue",
					"CHANNEL_TOKEN":  "testValue",

					// security
					"SALT": "testValue",
				},
				expect: map[string]string{
					"KEY_SERVER_URL": "http://test.test",
					"CHANNEL_SECRET": "testValue",
					"CHANNEL_TOKEN":  "testValue",
					"SALT":           "testValue",
				},
			},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				for k, v := range test.target {
					t.Setenv(k, v)
				}
				props.LoadEnv()
				// KEY_SERVER_URL
				if props.KeyServerURL != test.expect["KEY_SERVER_URL"] {
					t.Errorf("KEY_SERVER_URL:"+testutil.STRING_TEST_MSG_FMT, test.description, test.expect["KEY_SERVER_URL"], props.KeyServerURL)
				}

				// CHANNEL_SECRET check
				if props.ChannelSecret != test.expect["CHANNEL_SECRET"] {
					t.Errorf("CHANNEL_SECRET:"+testutil.STRING_TEST_MSG_FMT, "", test.expect["CHANNEL_SECRET"], props.ChannelSecret)
				}

				// CHANNEL_TOKEN check
				if props.ChannelToken != test.expect["CHANNEL_TOKEN"] {
					t.Errorf("CHANNEL_TOKEN:"+testutil.STRING_TEST_MSG_FMT, "", test.expect["CHANNEL_TOKEN"], props.ChannelToken)
				}

				// SALT check
				if props.Salt != test.expect["SALT"] {
					t.Errorf("SALT:"+testutil.STRING_TEST_MSG_FMT, "", test.expect["SALT"], props.Salt)
				}

			})
		}

		t.Run("load key server url exit test", func(t *testing.T) {
			defaultOsExit := props.OsExit
			defer func() { props.OsExit = defaultOsExit }()

			props.OsExit = func(value int) {
				if value != 1 {
					t.Errorf(testutil.INT_TEST_MSG_FMT, "failed key server url load exit status", 1, value)
				}
			}
			t.Setenv("KEY_SERVER_URL", "fizz")
			props.LoadEnv()
		})
	})

	t.Run("Boolean environment", func(t *testing.T) {
		tests := []struct {
			description string
			target      map[string]string
			expect      map[string]bool
		}{
			{
				description: "empty",
				target: map[string]string{
					"SIGNATURE_CHECK": "",
				},
				expect: map[string]bool{
					"SIGNATURE_CHECK": true,
				},
			},
			{
				description: "not empty(true)",
				target: map[string]string{
					"SIGNATURE_CHECK": "true",
				},
				expect: map[string]bool{
					"SIGNATURE_CHECK": true,
				},
			},
			{
				description: "not empty(false)",
				target: map[string]string{
					"SIGNATURE_CHECK": "false",
				},
				expect: map[string]bool{
					"SIGNATURE_CHECK": false,
				},
			},
			{
				description: "invalid value",
				target: map[string]string{
					"SIGNATURE_CHECK": "fizz",
				},
				expect: map[string]bool{
					"SIGNATURE_CHECK": true,
				},
			},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				for k, v := range test.target {
					t.Setenv(k, v)
					props.LoadEnv()

					if props.IsSignatureCheck != test.expect["SIGNATURE_CHECK"] {
						t.Errorf(testutil.BOOL_TEST_MSG_FMT, test.description, test.expect["SIGNATURE_CHECK"], props.IsSignatureCheck)
					}
				}
			})
		}
	})
}
