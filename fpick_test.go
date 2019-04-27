package fpick_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sawadashota/fpick"
)

func TestFilenameExtractMatch(t *testing.T) {
	type args struct {
		filename string
	}
	type match struct {
		filename string
		want     bool
	}
	cases := map[string]struct {
		args  args
		match match
	}{
		"same filename": {
			args: args{
				filename: "target.txt",
			},
			match: match{
				filename: "target.txt",
				want:     true,
			},
		},
		"different filename": {
			args: args{
				filename: "target.txt",
			},
			match: match{
				filename: "bar",
				want:     false,
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			match := fpick.FilenameExtractMatch(c.args.filename)
			if match(c.match.filename) != c.match.want {
				t.Errorf(`FilenameExtractMatch returns unexpected value. "%s" == "%s" => %v`, c.args.filename, c.match.filename, c.match.want)
			}
		})
	}
}

func TestFilenameRegexMatch(t *testing.T) {
	type args struct {
		regex string
	}
	type match struct {
		filename string
		want     bool
	}
	cases := map[string]struct {
		args    args
		wantErr bool
		match   match
	}{
		"match filename": {
			args: args{
				regex: `\Atarget\.txt\z`,
			},
			wantErr: false,
			match: match{
				filename: "target.txt",
				want:     true,
			},
		},
		"un-match filename": {
			args: args{
				regex: `target\.txt`,
			},
			wantErr: false,
			match: match{
				filename: "bar",
				want:     false,
			},
		},
		"broken regex": {
			args: args{
				regex: `(.*?`,
			},
			wantErr: true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			match, err := fpick.FilenameRegexMatch(c.args.regex)
			if err != nil {
				if !c.wantErr {
					t.Error(err)
				}
				return
			}

			if c.wantErr {
				t.Errorf(`want regex error of FilenameRegexMatch but no errors occurred. regex: "%s"`, c.args.regex)
			}

			if match(c.match.filename) != c.match.want {
				t.Errorf(`FilenameRegexMatch returns unexpected value. "%s" ~= "%s" => %v`, c.args.regex, c.match.filename, c.match.want)
			}
		})
	}
}

func TestClient_FileList(t *testing.T) {
	type fields struct {
		rootDir string
		outDir  string
	}
	type args struct {
		match fpick.FileMatcher
	}
	cases := map[string]struct {
		fields     fields
		wantNewErr bool
		args       args
		want       []*fpick.File
		wantErr    bool
	}{
		"normal": {
			fields: fields{
				rootDir: "testdata",
				outDir:  ".out",
			},
			wantNewErr: false,
			args: args{
				match: fpick.FilenameExtractMatch("target.txt"),
			},
			want: []*fpick.File{
				{
					Path: "testdata/target.txt",
					Perm: 0644,
				},
				{
					Path: "testdata/foo/target.txt",
					Perm: 0644,
				},
				{
					Path: "testdata/foo/bar/target.txt",
					Perm: 0644,
				},
			},
			wantErr: false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			client, err := fpick.New(c.fields.rootDir, c.fields.outDir)
			if err != nil {
				if !c.wantNewErr {
					t.Error(err)
				}
				return
			}
			if c.wantNewErr {
				t.Errorf(`want regex error of New but no errors occurred. src: "%s", dst: "%s"`, c.fields.rootDir, c.fields.outDir)
			}

			files, err := client.FileList(c.args.match)
			if (err != nil) != c.wantErr {
				t.Errorf("Client.FileList() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			checkFilesEqual(t, c.want, files)
		})
	}
}

func checkFilesEqual(t *testing.T, expected, actual []*fpick.File) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("Client.FileList() returns %d files but want %d files", len(actual), len(expected))
	}
	for _, f := range expected {
		if !includeFile(f, actual) {
			t.Errorf("Client.FileList() doesn't include %v", *f)
		}
	}

}

func includeFile(needle *fpick.File, haystack []*fpick.File) bool {
	for _, f := range haystack {
		if f.Path == needle.Path && f.Perm == needle.Perm {
			return true
		}
	}
	return false
}

func TestClient_Pick(t *testing.T) {
	type fields struct {
		src string
		dst string
	}
	type args struct {
		match fpick.FileMatcher
		opts  []fpick.OutputOption
	}
	type wantFile struct {
		path string
		body string
	}
	cases := map[string]struct {
		fields    fields
		args      args
		wantErr   bool
		wantFiles []wantFile
	}{
		"mirror output": {
			fields: fields{
				src: "testdata",
				dst: ".out",
			},
			args: args{
				match: fpick.FilenameExtractMatch("target.txt"),
			},
			wantErr: false,
			wantFiles: []wantFile{
				{
					path: ".out/target.txt",
					body: "I'm testdata/target.txt",
				},
				{
					path: ".out/foo/target.txt",
					body: "I'm testdata/foo/target.txt",
				},
				{
					path: ".out/foo/bar/target.txt",
					body: "I'm testdata/foo/bar/target.txt",
				},
			},
		},
		"output flat directory": {
			fields: fields{
				src: "testdata",
				dst: ".out",
			},
			args: args{
				match: fpick.FilenameExtractMatch("target.txt"),
				opts: []fpick.OutputOption{
					fpick.OutputFlatDirOption,
				},
			},
			wantErr: false,
			wantFiles: []wantFile{
				{
					path: ".out/target.txt",
					body: "I'm testdata/target.txt",
				},
				{
					path: ".out/foo/target.txt",
					body: "I'm testdata/foo/target.txt",
				},
				{
					path: ".out/foo__bar/target.txt",
					body: "I'm testdata/foo/bar/target.txt",
				},
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client, err := fpick.New(c.fields.src, c.fields.dst)
			if err != nil {
				t.Error(err)
				return
			}

			if err := client.Pick(c.args.match, c.args.opts...); (err != nil) != c.wantErr {
				t.Errorf("Client.Pick() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			defer os.RemoveAll(c.fields.dst)

			info, err := os.Stat(c.fields.dst)
			if err != nil || !info.IsDir() {
				t.Errorf("expect to created directory %s but not exists", err)
				return
			}

			for _, wf := range c.wantFiles {
				body, err := ioutil.ReadFile(wf.path)
				if err != nil {
					t.Error(err)
					return
				}
				if string(body) != wf.body {
					t.Errorf(`not same body at %s, expect "%s", actual "%s"`, wf.path, wf.body, string(body))
				}
			}
		})
	}
}
