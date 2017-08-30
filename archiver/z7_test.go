package archiver

import (
	"testing"

	"pkg.re/essentialkaos/ek.v9/fsutil"

	check "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { check.TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type Z7Suite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = check.Suite(&Z7Suite{})

// ////////////////////////////////////////////////////////////////////////////////// //
func (s *Z7Suite) TestCheck(c *check.C) {
	ok, err := Check("testdata/test-max-compression.7z")

	c.Assert(ok, check.Equals, true)
	c.Assert(err, check.IsNil)

	ok, err = Check("testdata/test-broken.7z")

	c.Assert(ok, check.Equals, false)
	c.Assert(err, check.NotNil)

	ok, err = Check("testdata/test-not-exist.7z")

	c.Assert(ok, check.Equals, false)
	c.Assert(err, check.NotNil)
}

func (s *Z7Suite) TestAdd(c *check.C) {
	resultFile := c.MkDir() + "/test.7z"

	_, err := Add(&Props{File: resultFile}, "testdata/test")

	c.Assert(err, check.IsNil)

	c.Assert(fsutil.IsExist(resultFile), check.Equals, true)
	c.Assert(fsutil.IsReadable(resultFile), check.Equals, true)
	c.Assert(fsutil.IsNonEmpty(resultFile), check.Equals, true)

	ok, err := Check(resultFile)

	c.Assert(ok, check.Equals, true)
	c.Assert(err, check.IsNil)

	resultFile = c.MkDir() + "/test1.7z"

	_, err = Add(resultFile, "testdata/test")

	c.Assert(err, check.IsNil)

	c.Assert(fsutil.IsExist(resultFile), check.Equals, true)
	c.Assert(fsutil.IsReadable(resultFile), check.Equals, true)
	c.Assert(fsutil.IsNonEmpty(resultFile), check.Equals, true)

	ok, err = Check(resultFile)

	c.Assert(ok, check.Equals, true)
	c.Assert(err, check.IsNil)
}

func (s *Z7Suite) TestExtract(c *check.C) {
	outputDir := c.MkDir()

	_, err := Extract(
		&Props{
			File:      "testdata/test-max-compression.7z",
			OutputDir: outputDir,
		},
	)

	c.Assert(err, check.IsNil)

	c.Assert(fsutil.CheckPerms("DR", outputDir+"/test"), check.Equals, true)
	c.Assert(fsutil.CheckPerms("DR", outputDir+"/test/dir1"), check.Equals, true)
	c.Assert(fsutil.CheckPerms("FRS", outputDir+"/test/dir1/file1.log"), check.Equals, true)
	c.Assert(fsutil.CheckPerms("FRS", outputDir+"/test/file1.log"), check.Equals, true)
	c.Assert(fsutil.CheckPerms("FRS", outputDir+"/test/file2.log"), check.Equals, false)
	c.Assert(fsutil.CheckPerms("FRS", outputDir+"/test/file3.log"), check.Equals, true)
}

func (s *Z7Suite) TestDelete(c *check.C) {
	testArchive := c.MkDir() + "/test.7z"

	fsutil.CopyFile("testdata/test-max-compression.7z", testArchive)

	_, err := Delete(testArchive, "test/file2.log", "test/file3.log")

	c.Assert(err, check.IsNil)

	outputDir := c.MkDir()

	_, err = Extract(
		&Props{
			File:      testArchive,
			OutputDir: outputDir,
		},
	)

	c.Assert(err, check.IsNil)

	c.Assert(fsutil.CheckPerms("DR", outputDir+"/test"), check.Equals, true)
	c.Assert(fsutil.CheckPerms("DR", outputDir+"/test/dir1"), check.Equals, true)
	c.Assert(fsutil.CheckPerms("FR", outputDir+"/test/dir1/file1.log"), check.Equals, true)
	c.Assert(fsutil.CheckPerms("FR", outputDir+"/test/file1.log"), check.Equals, true)
	c.Assert(fsutil.CheckPerms("FR", outputDir+"/test/file2.log"), check.Equals, false)
	c.Assert(fsutil.CheckPerms("FR", outputDir+"/test/file3.log"), check.Equals, false)
}
