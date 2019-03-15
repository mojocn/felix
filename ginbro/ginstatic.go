package ginbro

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const nameSourceFile = "gin_static.go"
const importString = `
import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)
`
const functionStrings = `

// file holds unzipped read-only file contents and file metadata.
type file struct {
	os.FileInfo
	data []byte
	fs   *ginBinFs
}

type ginBinFs struct {
	files map[string]file
	dirs  map[string][]string
}

const indexHtml = "index.html"

// Static returns a middleware handler that serves static files in the given directory.
func NewGinStaticBinMiddleware(urlPrefix string) (gin.HandlerFunc, error) {
	fs, err := create(zipData)
	if err != nil {
		return nil, err
	}
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		urlPath := strings.TrimSpace(c.Request.URL.Path)
		if urlPath == urlPrefix {
			urlPath = path.Join(urlPrefix, indexHtml)
		}
		f, err := fs.Open(urlPath)
		if err != nil {
			return
		}
		fi, err := f.Stat()
		if strings.HasSuffix(urlPath,".html"){
			c.Header("Cache-Control","no-cache")
		}
		if err != nil || !fi.IsDir() {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}, nil
}

// New creates a new file system with the registered zip contents data.
// It unzips all files and stores them in an in-memory map.
func create(rawZipString string) (http.FileSystem, error) {
	if zipData == "" {
		return nil, errors.New("statik/fs: no zip data registered")
	}
	zipReader, err := zip.NewReader(strings.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}
	files := make(map[string]file, len(zipReader.File))
	dirs := make(map[string][]string)
	fs := &ginBinFs{files: files, dirs: dirs}
	for _, zipFile := range zipReader.File {
		fi := zipFile.FileInfo()
		f := file{FileInfo: fi, fs: fs}
		f.data, err = unzip(zipFile)
		if err != nil {
			return nil, fmt.Errorf("statik/fs: error unzipping file %q: %s", zipFile.Name, err)
		}
		files["/"+zipFile.Name] = f
	}
	for fn := range files {
		// go up directories recursively in order to care deep directory
		for dn := path.Dir(fn); dn != fn; {
			if _, ok := files[dn]; !ok {
				files[dn] = file{FileInfo: dirInfo{dn}, fs: fs}
			} else {
				break
			}
			fn, dn = dn, path.Dir(dn)
		}
	}
	for fn := range files {
		dn := path.Dir(fn)
		if fn != dn {
			fs.dirs[dn] = append(fs.dirs[dn], path.Base(fn))
		}
	}
	for _, s := range fs.dirs {
		sort.Strings(s)
	}
	return fs, nil
}

var _ = os.FileInfo(dirInfo{})

type dirInfo struct {
	name string
}

func (di dirInfo) Name() string       { return path.Base(di.name) }
func (di dirInfo) Size() int64        { return 0 }
func (di dirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (di dirInfo) ModTime() time.Time { return time.Time{} }
func (di dirInfo) IsDir() bool        { return true }
func (di dirInfo) Sys() interface{}   { return nil }

func unzip(zf *zip.File) ([]byte, error) {
	rc, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

// Open returns a file matching the given file name, or os.ErrNotExists if
// no file matching the given file name is found in the archive.
// If a directory is requested, Open returns the file named "index.html"
// in the requested directory, if that file exists.
func (fs *ginBinFs) Open(name string) (http.File, error) {
	name = strings.Replace(name, "//", "/", -1)
	if f, ok := fs.files[name]; ok {
		return newHTTPFile(f), nil
	}
	return nil, os.ErrNotExist
}

func newHTTPFile(file file) *httpFile {
	if file.IsDir() {
		return &httpFile{file: file, isDir: true}
	}
	return &httpFile{file: file, reader: bytes.NewReader(file.data)}
}

// httpFile represents an HTTP file and acts as a bridge
// between file and http.File.
type httpFile struct {
	file

	reader *bytes.Reader
	isDir  bool
	dirIdx int
}

// Read reads bytes into p, returns the number of read bytes.
func (f *httpFile) Read(p []byte) (n int, err error) {
	if f.reader == nil && f.isDir {
		return 0, io.EOF
	}
	return f.reader.Read(p)
}

// Seek seeks to the offset.
func (f *httpFile) Seek(offset int64, whence int) (ret int64, err error) {
	return f.reader.Seek(offset, whence)
}

// Stat stats the file.
func (f *httpFile) Stat() (os.FileInfo, error) {
	return f, nil
}

// IsDir returns true if the file location represents a directory.
func (f *httpFile) IsDir() bool {
	return f.isDir
}

// Readdir returns an empty slice of files, directory
// listing is disabled.
func (f *httpFile) Readdir(count int) ([]os.FileInfo, error) {
	var fis []os.FileInfo
	if !f.isDir {
		return fis, nil
	}
	di, ok := f.FileInfo.(dirInfo)
	if !ok {
		return nil, fmt.Errorf("failed to read directory: %q", f.Name())
	}

	// If count is positive, the specified number of files will be returned,
	// and if negative, all remaining files will be returned.
	// The reading position of which file is returned is held in dirIndex.
	fnames := f.file.fs.dirs[di.name]
	flen := len(fnames)

	// If dirIdx reaches the end and the count is a positive value,
	// an io.EOF error is returned.
	// In other cases, no error will be returned even if, for example,
	// you specified more counts than the number of remaining files.
	start := f.dirIdx
	if start >= flen && count > 0 {
		return fis, io.EOF
	}
	var end int
	if count < 0 {
		end = flen
	} else {
		end = start + count
	}
	if end > flen {
		end = flen
	}
	for i := start; i < end; i++ {
		fis = append(fis, f.file.fs.files[path.Join(di.name, fnames[i])].FileInfo)
	}
	f.dirIdx += len(fis)
	return fis, nil
}

func (f *httpFile) Close() error {
	return nil
}
`

// mtimeDate holds the arbitrary mtime that we assign to files when
// flagNoMtime is set.
var mtimeDate = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

func RunGinStatic(flagSrc, flagDest, flagTags, flagPkg, flagPkgCmt string, flagNoMtime, flagNoCompress, flagForce bool) {

	file, err := generateSource(flagSrc, flagTags, flagPkgCmt, flagPkg, flagNoMtime, flagNoCompress)
	if err != nil {
		exitWithError(err)
	}

	destDir := path.Join(flagDest, flagPkg)
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		exitWithError(err)
	}

	err = rename(file.Name(), path.Join(destDir, nameSourceFile), flagForce)
	if err != nil {
		exitWithError(err)
	}
}

// rename tries to os.Rename, but fall backs to copying from src
// to dest and unlink the source if os.Rename fails.
func rename(src, dest string, flagForce bool) error {
	// Try to rename generated source.
	if err := os.Rename(src, dest); err == nil {
		return nil
	}
	// If the rename failed (might do so due to temporary file residing on a
	// different device), try to copy byte by byte.
	rc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		rc.Close()
		os.Remove(src) // ignore the error, source is in tmp.
	}()

	if _, err = os.Stat(dest); !os.IsNotExist(err) {
		if flagForce {
			if err = os.Remove(dest); err != nil {
				return fmt.Errorf("file %q could not be deleted", dest)
			}
		} else {
			return fmt.Errorf("file %q already exists; use -f to overwrite", dest)
		}
	}

	wc, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer wc.Close()

	if _, err = io.Copy(wc, rc); err != nil {
		// Delete remains of failed copy attempt.
		os.Remove(dest)
	}
	return err
}

// Walks on the source path and generates source code
// that contains source directory's contents as zip contents.
// Generates source registers generated zip contents data to
// be read by the statik/fs HTTP file system.
func generateSource(srcPath, flagTags, flagPkgCmt, flagPkg string, flagNoMtime, flagNoCompress bool) (file *os.File, err error) {
	var (
		buffer    bytes.Buffer
		zipWriter io.Writer
	)

	zipWriter = &buffer
	f, err := ioutil.TempFile("", flagPkg)
	if err != nil {
		return
	}

	zipWriter = io.MultiWriter(zipWriter, f)
	defer f.Close()

	w := zip.NewWriter(zipWriter)
	if err = filepath.Walk(srcPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore directories and hidden files.
		// No entry is needed for directories in a zip file.
		// Each file is represented with a path, no directory
		// entities are required to build the hierarchy.
		if fi.IsDir() || strings.HasPrefix(fi.Name(), ".") {
			return nil
		}
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fHeader, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		if flagNoMtime {
			// Always use the same modification time so that
			// the output is deterministic with respect to the file contents.
			// Do NOT use fHeader.Modified as it only works on go >= 1.10
			fHeader.Modified = mtimeDate
		}
		fHeader.Name = filepath.ToSlash(relPath)
		if !flagNoCompress {
			fHeader.Method = zip.Deflate
		}
		f, err := w.CreateHeader(fHeader)
		if err != nil {
			return err
		}
		_, err = f.Write(b)
		return err
	}); err != nil {
		return
	}
	if err = w.Close(); err != nil {
		return
	}

	var tags string
	if flagTags != "" {
		tags = "\n// +build " + flagTags + "\n"
	}

	var comment string
	if flagPkgCmt != "" {
		comment = "\n" + commentLines(flagPkgCmt)
	}

	// then embed it as a quoted string
	var qb bytes.Buffer
	fmt.Fprintf(&qb, `// Code generated by felix ginbin command. DO NOT EDIT.
%s%s
package %s
%s
const zipData = "`, tags, comment, flagPkg, importString)

	fprintZipData(&qb, buffer.Bytes())

	fmt.Fprintf(&qb, `"%s`, functionStrings)

	if err = ioutil.WriteFile(f.Name(), qb.Bytes(), 0644); err != nil {
		return
	}
	return f, nil
}

// fprintZipData converts zip binary contents to a string literal.
func fprintZipData(dest *bytes.Buffer, zipData []byte) {
	for _, b := range zipData {
		if b == '\n' {
			dest.WriteString(`\n`)
			continue
		}
		if b == '\\' {
			dest.WriteString(`\\`)
			continue
		}
		if b == '"' {
			dest.WriteString(`\"`)
			continue
		}
		if (b >= 32 && b <= 126) || b == '\t' {
			dest.WriteByte(b)
			continue
		}
		fmt.Fprintf(dest, "\\x%02x", b)
	}
}

// comment lines prefixes each line in lines with "// ".
func commentLines(lines string) string {
	lines = "// " + strings.Replace(lines, "\n", "\n// ", -1)
	return lines
}

// Prints out the error message and exists with a non-success signal.
func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
