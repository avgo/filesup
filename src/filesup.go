package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"strings"
	"time"

)

type options_s struct {
	filesup_prefix string
	uploaded_dir string
}

var options options_s

func file_load_handler(bindir, part_formname, form_filename, cur_up_dir, up_filename, up_handlers_dir string) {
	file_prefix := "handler-"
	files, err := ioutil.ReadDir(bindir)
	if err != nil {
		filesup_die("file_load_handler() error: \n")
	}
	no_files_with_prefix := true
	for _, cur_file := range files {
		func () {
			cur_file_name := cur_file.Name()
			if !strings.HasPrefix(cur_file_name, file_prefix) {
				return
			}
			fmt.Fprintf(os.Stderr, "file pref: '%s'\n", cur_file_name)
			no_files_with_prefix = false

			up_handler_dir := up_handlers_dir + "/" + cur_file_name
			err := os.Mkdir(up_handler_dir, 0755)
			if err != nil {
				filesup_die("mkdir handler_dir error: " + err.Error() + "\n")
			}

			up_log_dir := up_handler_dir + "/log"
			err = os.Mkdir(up_log_dir, 0755)
			if err != nil {
				filesup_die("mkdir log_dir error: " + err.Error() + "\n")
			}

			up_tmp_dir := up_handler_dir + "/tmp"
			err = os.Mkdir(up_tmp_dir, 0755)
			if err != nil {
				filesup_die("mkdir tmp_dir error: " + err.Error() + "\n")
			}

			dst_stdout, err := os.OpenFile(up_log_dir + "/stdout.log", os.O_CREATE | os.O_WRONLY, 0644)
			if err != nil {
				filesup_die("Can't open file for stdout\n")
			}
			defer dst_stdout.Close()

			dst_stderr, err := os.OpenFile(up_log_dir + "/stderr.log", os.O_CREATE | os.O_WRONLY, 0644)
			if err != nil {
				filesup_die("Can't open file for stderr\n")
			}
			defer dst_stderr.Close()

			cmd := exec.Command(bindir + "/" + cur_file_name, cur_up_dir, up_filename, up_tmp_dir)

			src_stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: stdout " + err.Error() + "\n")
			}

			src_stderr, err := cmd.StderrPipe()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: stderr " + err.Error() + "\n")
			}

			err = cmd.Start()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Start() error: " + err.Error() + "\n")
			}

			_, err = io.Copy(dst_stdout, src_stdout)
			if err != nil {
				filesup_die("err copy stdout!\n")
			}

			_, err = io.Copy(dst_stderr, src_stderr)
			if err != nil {
				filesup_die("err copy stderr!\n")
			}
			cmd.Wait()
		}()
	}
	if no_files_with_prefix {
		fmt.Fprintf(os.Stderr, "files with prefix is not finded!\n")
	}
}

func filesup_die(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

func filesup_save_to_file(part *multipart.Part, filename string) {
	dst, err := os.OpenFile(filename, os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		filesup_die("err2!\n")
	}
	defer dst.Close()
	_, err = io.Copy(dst, part)
	if err != nil {
		filesup_die("err3!\n")
	}
}

func is_multipart(r *http.Request) bool {
	v := r.Header.Get("Content-Type")
	if v == "" {
		return false
	}
	d, _, err := mime.ParseMediaType(v)
	if err != nil || d != "multipart/form-data" {
		return false
	}
	return true
}

func load_options() {
	var envok bool
	options.filesup_prefix, envok = os.LookupEnv("FILESUP_PREFIX")
	if !envok {
		filesup_die("FILESUP_PREFIX env must be set!\n")
	}
	options.uploaded_dir, envok = os.LookupEnv("FILESUP_UPLOADED_DIR")
	if !envok {
		filesup_die("FILESUP_UPLOADED_DIR env must be set!\n")
	}
}

func main() {
	load_options()
	r, err := cgi.Request()
	if err != nil {
		filesup_die("err!\n")
	}
	if !is_multipart(r) {
		fmt.Printf("Content-type: text/html\n\n<a href=\"/\">HOME</a><br>No multipart!\n")
		os.Exit(0)
	}
	read_form, err := r.MultipartReader()
	if err != nil {
		filesup_die("err!\n")
	}
	type rec struct {
		Formname string
		Filename string
	}
	doc := struct {
		Files_str []rec
		Envs []string
	}{
		Files_str: []rec{},
		Envs: os.Environ(),
	}

	t := time.Now()
	cur_up_dir := options.uploaded_dir + "/" +
			fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second());
	err = os.Mkdir(cur_up_dir, 0755)
	if err != nil {
		filesup_die("mkdir up error: " + err.Error() + "\n")
	}
	cur_up_data_dir := cur_up_dir + "/data";
	err = os.Mkdir(cur_up_data_dir, 0755)
	if err != nil {
		filesup_die("mkdir up_data error: " + err.Error() + "\n")
	}
	cur_up_handlers_dir := cur_up_dir + "/handlers";
	err = os.Mkdir(cur_up_handlers_dir, 0755)
	if err != nil {
		filesup_die("mkdir up_handlers error: " + err.Error() + "\n")
	}

	for {
		if func () int {
		part, err := read_form.NextPart()
		if err == io.EOF {
			return 1
		}
		if err != nil {
			filesup_die("err1!\n")
		}
		form_filename := part.FileName()
		if form_filename == "" {
			form_filename = "empty"
		} else {
			up_filename := cur_up_data_dir + "/" + form_filename
			filesup_save_to_file(part, up_filename)
			cur_up_handlers_f_dir := cur_up_handlers_dir + "/" + form_filename;
			err = os.Mkdir(cur_up_handlers_f_dir, 0755)
			if err != nil {
				filesup_die("mkdir up_handlers_f_dir error: " + err.Error() + "\n")
			}
			file_load_handler(options.filesup_prefix + "/bin", part.FormName(), form_filename, cur_up_dir, up_filename, cur_up_handlers_f_dir)
		}
		doc.Files_str = append(doc.Files_str, rec{ Filename: form_filename, Formname: part.FormName() })
		return 0
		}() == 1 {
			break
		}
	}
	tmpl, err := template.New("main.html").ParseFiles(options.filesup_prefix + "/templates/main.html")
	if err != nil {
		filesup_die("err4!\n")
	}
	fmt.Printf("Content-type: text/html\n\n")
	// doc.Files_str = []rec{}
	err = tmpl.Execute(os.Stdout, doc)
	if err != nil {
		filesup_die(err.Error())
	}
}
