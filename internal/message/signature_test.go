package message

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerify(t *testing.T) {
	verifier := &Verifier{Key: "Spoon!"}
	for _, tc := range []struct {
		filename  string
		signature string
		expected  bool
	}{{
		filename:  "testdata/message1.json",
		signature: "55c7cf6d920e8150c1d04f708f796b17dbd65a0684402a6be39f7a2baef6f998",
		expected:  true,
	}, {
		filename:  "testdata/message2.json",
		signature: "c0b4e60909b8e638c36803f9cd618d66a97507296fc09dca473d6c08351cba8d",
		expected:  true,
	}, {
		filename:  "testdata/message3.json",
		signature: "ca077808b54c2dfeebbab806e6c32e8dd231f74e130f4becea9c8a00274fbc3f",
		expected:  true,
	}, {
		filename:  "testdata/message4.json",
		signature: "66cfdd1e9616a71f80c430857b0d5d67480d11a3669effb1cbca063072d5bb3f",
		expected:  true,
	}, {
		filename:  "testdata/message5.json",
		signature: "f8abd2d59599d1ff763d0562e21e2541d6c0ec8bc93003f3b31dfabd4eb6a3e9",
		expected:  true,
	}, {
		filename:  "testdata/message6.json",
		signature: "dc10a0dc0995093866846fe2974b1d97c18c73c96295c806f6ecd18ee6f1e467",
		expected:  true,
	}, {
		filename:  "testdata/message1.json",
		signature: "dc10a0dc0995093866846fe2974b1d97c18c73c96295c806f6ecd18ee6f1e467",
		expected:  false,
	}, {
		filename:  "testdata/message2.json",
		signature: "f8abd2d59599d1ff763d0562e21e2541d6c0ec8bc93003f3b31dfabd4eb6a3e9",
		expected:  false,
	}, {
		filename:  "testdata/message3.json",
		signature: "66cfdd1e9616a71f80c430857b0d5d67480d11a3669effb1cbca063072d5bb3f",
		expected:  false,
	}, {
		filename:  "testdata/message4.json",
		signature: "ca077808b54c2dfeebbab806e6c32e8dd231f74e130f4becea9c8a00274fbc3f",
		expected:  false,
	}, {
		filename:  "testdata/message5.json",
		signature: "c0b4e60909b8e638c36803f9cd618d66a97507296fc09dca473d6c08351cba8d",
		expected:  false,
	}, {
		filename:  "testdata/message6.json",
		signature: "55c7cf6d920e8150c1d04f708f796b17dbd65a0684402a6be39f7a2baef6f998",
		expected:  false,
	}} {
		basename := filepath.Base(tc.filename)
		t.Run(fmt.Sprintf("%s(valid=%v)", basename, tc.expected), func(t *testing.T) {
			require := require.New(t)

			data, err := ioutil.ReadFile(tc.filename)
			require.NoError(err)

			message := strings.TrimSpace(string(data))
			actual, err := verifier.VerifySignature(message, tc.signature)
			require.NoError(err)

			require.Equal(tc.expected, actual)
		})
	}
}
