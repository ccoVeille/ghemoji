package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	pkg = flag.String("package", "ghemoji", "package name to output")
	arg = ""
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [options] <file>:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	arg = flag.Arg(0)
}

const templ = `// Code generated by "ghupdate"; DO NOT EDIT

package {{.Package}}

var emoji = map[string]string{
{{range $key, $val := .EmojiData}}{{if $val.Emoji}}{{range .Aliases}}	"{{.}}": "{{$val.Emoji}}",
{{end}}{{end}}{{end}}}
`

type emojiStruct []struct {
	Emoji          string   `json:"emoji,omitempty"`
	Description    string   `json:"description,omitempty"`
	Category       string   `json:"category,omitempty"`
	Aliases        []string `json:"aliases"`
	Tags           []string `json:"tags"`
	UnicodeVersion string   `json:"unicode_version,omitempty"`
	IosVersion     string   `json:"ios_version,omitempty"`
}

type templateData struct {
	Package   string
	EmojiData emojiStruct
}

func main() {
	emoji := emojiStruct{}

	res, err := http.Get("https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	x := json.NewDecoder(res.Body)

	err = x.Decode(&emoji)
	if err != nil {
		log.Fatal(err.Error())
	}

	td := templateData{
		Package:   *pkg,
		EmojiData: emoji,
	}

	tmpl := template.New("code template")
	tmpl, err = tmpl.Parse(templ)
	if err != nil {
		log.Fatal(err.Error())
	}

	f, err := os.Create(arg)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer f.Close()

	err = tmpl.Execute(f, td)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("success")
}
