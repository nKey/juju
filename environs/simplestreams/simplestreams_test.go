// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package simplestreams_test

import (
	"bytes"
	"testing"

	gc "launchpad.net/gocheck"

	"launchpad.net/juju-core/environs/simplestreams"
	sstesting "launchpad.net/juju-core/environs/simplestreams/testing"
)

func Test(t *testing.T) {
	registerSimpleStreamsTests()
	gc.Suite(&signingSuite{})
	gc.TestingT(t)
}

func registerSimpleStreamsTests() {
	gc.Suite(&simplestreamsSuite{
		LocalLiveSimplestreamsSuite: sstesting.LocalLiveSimplestreamsSuite{
			BaseURL:       "test:",
			RequireSigned: false,
			ValidConstraint: sstesting.NewTestConstraint(simplestreams.LookupParams{
				CloudSpec: simplestreams.CloudSpec{
					Region:   "us-east-1",
					Endpoint: "https://ec2.us-east-1.amazonaws.com",
				},
				Series: "precise",
				Arches: []string{"amd64", "arm"},
			}),
		},
	})
}

type simplestreamsSuite struct {
	sstesting.TestDataSuite
	sstesting.LocalLiveSimplestreamsSuite
}

func (s *simplestreamsSuite) SetUpSuite(c *gc.C) {
	s.LocalLiveSimplestreamsSuite.SetUpSuite(c)
	s.TestDataSuite.SetUpSuite(c)
}

func (s *simplestreamsSuite) TearDownSuite(c *gc.C) {
	s.TestDataSuite.TearDownSuite(c)
	s.LocalLiveSimplestreamsSuite.TearDownSuite(c)
}

func (s *simplestreamsSuite) TestGetProductsPath(c *gc.C) {
	indexRef, err := s.GetIndexRef(sstesting.Index_v1)
	c.Assert(err, gc.IsNil)
	path, err := indexRef.GetProductsPath(s.ValidConstraint)
	c.Assert(err, gc.IsNil)
	c.Assert(path, gc.Equals, "streams/v1/image_metadata.json")
}

func (s *simplestreamsSuite) TestMetadataCatalog(c *gc.C) {
	metadata := s.AssertGetMetadata(c)
	c.Check(len(metadata.Products), gc.Equals, 2)
	c.Check(len(metadata.Aliases), gc.Equals, 1)
	metadataCatalog := metadata.Products["com.ubuntu.cloud:server:12.04:amd64"]
	c.Check(len(metadataCatalog.Items), gc.Equals, 2)
	c.Check(metadataCatalog.Series, gc.Equals, "precise")
	c.Check(metadataCatalog.Version, gc.Equals, "12.04")
	c.Check(metadataCatalog.Arch, gc.Equals, "amd64")
	c.Check(metadataCatalog.RegionName, gc.Equals, "au-east-1")
	c.Check(metadataCatalog.Endpoint, gc.Equals, "https://somewhere")
}

func (s *simplestreamsSuite) TestItemCollection(c *gc.C) {
	ic := s.AssertGetItemCollections(c, "20121218")
	c.Check(ic.RegionName, gc.Equals, "au-east-2")
	c.Check(ic.Endpoint, gc.Equals, "https://somewhere-else")
	c.Assert(len(ic.Items) > 0, gc.Equals, true)
	ti := ic.Items["usww2he"].(*sstesting.TestItem)
	c.Check(ti.Id, gc.Equals, "ami-442ea674")
	c.Check(ti.Storage, gc.Equals, "ebs")
	c.Check(ti.VType, gc.Equals, "hvm")
	c.Check(ti.RegionName, gc.Equals, "us-east-1")
	c.Check(ti.Endpoint, gc.Equals, "https://ec2.us-east-1.amazonaws.com")
}

func (s *simplestreamsSuite) TestDenormalisationFromCollection(c *gc.C) {
	ic := s.AssertGetItemCollections(c, "20121218")
	ti := ic.Items["usww1pe"].(*sstesting.TestItem)
	c.Check(ti.RegionName, gc.Equals, ic.RegionName)
	c.Check(ti.Endpoint, gc.Equals, ic.Endpoint)
}

func (s *simplestreamsSuite) TestDenormalisationFromCatalog(c *gc.C) {
	metadata := s.AssertGetMetadata(c)
	metadataCatalog := metadata.Products["com.ubuntu.cloud:server:12.04:amd64"]
	ic := metadataCatalog.Items["20111111"]
	ti := ic.Items["usww3pe"].(*sstesting.TestItem)
	c.Check(ti.RegionName, gc.Equals, metadataCatalog.RegionName)
	c.Check(ti.Endpoint, gc.Equals, metadataCatalog.Endpoint)
}

func (s *simplestreamsSuite) TestDealiasing(c *gc.C) {
	metadata := s.AssertGetMetadata(c)
	metadataCatalog := metadata.Products["com.ubuntu.cloud:server:12.04:amd64"]
	ic := metadataCatalog.Items["20121218"]
	ti := ic.Items["usww3he"].(*sstesting.TestItem)
	c.Check(ti.RegionName, gc.Equals, "us-west-3")
	c.Check(ti.Endpoint, gc.Equals, "https://ec2.us-west-3.amazonaws.com")
}

var testSigningKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: GnuPG v1.4.10 (GNU/Linux)

lQHYBE2rFNoBBADFwqWQIW/DSqcB4yCQqnAFTJ27qS5AnB46ccAdw3u4Greeu3Bp
idpoHdjULy7zSKlwR1EA873dO/k/e11Ml3dlAFUinWeejWaK2ugFP6JjiieSsrKn
vWNicdCS4HTWn0X4sjl0ZiAygw6GNhqEQ3cpLeL0g8E9hnYzJKQ0LWJa0QARAQAB
AAP/TB81EIo2VYNmTq0pK1ZXwUpxCrvAAIG3hwKjEzHcbQznsjNvPUihZ+NZQ6+X
0HCfPAdPkGDCLCb6NavcSW+iNnLTrdDnSI6+3BbIONqWWdRDYJhqZCkqmG6zqSfL
IdkJgCw94taUg5BWP/AAeQrhzjChvpMQTVKQL5mnuZbUCeMCAN5qrYMP2S9iKdnk
VANIFj7656ARKt/nf4CBzxcpHTyB8+d2CtPDKCmlJP6vL8t58Jmih+kHJMvC0dzn
gr5f5+sCAOOe5gt9e0am7AvQWhdbHVfJU0TQJx+m2OiCJAqGTB1nvtBLHdJnfdC9
TnXXQ6ZXibqLyBies/xeY2sCKL5qtTMCAKnX9+9d/5yQxRyrQUHt1NYhaXZnJbHx
q4ytu0eWz+5i68IYUSK69jJ1NWPM0T6SkqpB3KCAIv68VFm9PxqG1KmhSrQIVGVz
dCBLZXmIuAQTAQIAIgUCTasU2gIbAwYLCQgHAwIGFQgCCQoLBBYCAwECHgECF4AA
CgkQO9o98PRieSoLhgQAkLEZex02Qt7vGhZzMwuN0R22w3VwyYyjBx+fM3JFETy1
ut4xcLJoJfIaF5ZS38UplgakHG0FQ+b49i8dMij0aZmDqGxrew1m4kBfjXw9B/v+
eIqpODryb6cOSwyQFH0lQkXC040pjq9YqDsO5w0WYNXYKDnzRV0p4H1pweo2VDid
AdgETasU2gEEAN46UPeWRqKHvA99arOxee38fBt2CI08iiWyI8T3J6ivtFGixSqV
bRcPxYO/qLpVe5l84Nb3X71GfVXlc9hyv7CD6tcowL59hg1E/DC5ydI8K8iEpUmK
/UnHdIY5h8/kqgGxkY/T/hgp5fRQgW1ZoZxLajVlMRZ8W4tFtT0DeA+JABEBAAEA
A/0bE1jaaZKj6ndqcw86jd+QtD1SF+Cf21CWRNeLKnUds4FRRvclzTyUMuWPkUeX
TaNNsUOFqBsf6QQ2oHUBBK4VCHffHCW4ZEX2cd6umz7mpHW6XzN4DECEzOVksXtc
lUC1j4UB91DC/RNQqwX1IV2QLSwssVotPMPqhOi0ZLNY7wIA3n7DWKInxYZZ4K+6
rQ+POsz6brEoRHwr8x6XlHenq1Oki855pSa1yXIARoTrSJkBtn5oI+f8AzrnN0BN
oyeQAwIA/7E++3HDi5aweWrViiul9cd3rcsS0dEnksPhvS0ozCJiHsq/6GFmy7J8
QSHZPteedBnZyNp5jR+H7cIfVN3KgwH/Skq4PsuPhDq5TKK6i8Pc1WW8MA6DXTdU
nLkX7RGmMwjC0DBf7KWAlPjFaONAX3a8ndnz//fy1q7u2l9AZwrj1qa1iJ8EGAEC
AAkFAk2rFNoCGwwACgkQO9o98PRieSo2/QP/WTzr4ioINVsvN1akKuekmEMI3LAp
BfHwatufxxP1U+3Si/6YIk7kuPB9Hs+pRqCXzbvPRrI8NHZBmc8qIGthishdCYad
AHcVnXjtxrULkQFGbGvhKURLvS9WnzD/m1K2zzwxzkPTzT9/Yf06O6Mal5AdugPL
VrM0m72/jnpKo04=
=zNCn
-----END PGP PRIVATE KEY BLOCK-----
`

var validClearsignInput = `
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA1

Hello world
line 2
`

var invalidClearsignInput = `
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA1

Invalid
`

var testSig = `-----BEGIN PGP SIGNATURE-----
Version: GnuPG v1.4.10 (GNU/Linux)

iJwEAQECAAYFAk8kMuEACgkQO9o98PRieSpMsAQAhmY/vwmNpflrPgmfWsYhk5O8
pjnBUzZwqTDoDeINjZEoPDSpQAHGhjFjgaDx/Gj4fAl0dM4D0wuUEBb6QOrwflog
2A2k9kfSOMOtk0IH/H5VuFN1Mie9L/erYXjTQIptv9t9J7NoRBMU0QOOaFU0JaO9
MyTpno24AjIAGb+mH1U=
=hIJ6
-----END PGP SIGNATURE-----
`

type signingSuite struct {
	origKey string
}

func (s *signingSuite) SetUpSuite(c *gc.C) {
	s.origKey = simplestreams.SetSigningKey(testSigningKey)
}

func (s *signingSuite) TearDownSuite(c *gc.C) {
	simplestreams.SetSigningKey(s.origKey)
}

func (s *signingSuite) TestDecodeCheckValidSignature(c *gc.C) {
	r := bytes.NewReader([]byte(validClearsignInput + testSig))
	txt, err := simplestreams.DecodeCheckSignature(r)
	c.Assert(err, gc.IsNil)
	c.Assert(txt, gc.DeepEquals, []byte("Hello world\nline 2\n"))
}

func (s *signingSuite) TestDecodeCheckInvalidSignature(c *gc.C) {
	r := bytes.NewReader([]byte(invalidClearsignInput + testSig))
	_, err := simplestreams.DecodeCheckSignature(r)
	c.Assert(err, gc.Not(gc.IsNil))
	_, ok := err.(*simplestreams.NotPGPSignedError)
	c.Assert(ok, gc.Equals, false)
}

func (s *signingSuite) TestDecodeCheckMissingSignature(c *gc.C) {
	r := bytes.NewReader([]byte("foo"))
	_, err := simplestreams.DecodeCheckSignature(r)
	_, ok := err.(*simplestreams.NotPGPSignedError)
	c.Assert(ok, gc.Equals, true)
}
