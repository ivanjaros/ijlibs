package wizard

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func New() ConfigWizard {
	return ConfigWizard{r: bufio.NewReader(os.Stdin)}
}

type ConfigWizard struct {
	r *bufio.Reader
}

func (w ConfigWizard) write(str string) {
	_, _ = os.Stdout.WriteString(str)
}

func (w *ConfigWizard) read() string {
	text, err := w.r.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.Trim(text, "\r\n\t ")
}

func (w ConfigWizard) error(e error) {
	w.err(e.Error())
}

func (w ConfigWizard) err(e string) {
	_, _ = os.Stderr.WriteString(e + "\n")
}

func (w *ConfigWizard) prompt(question string, ask bool, defaultValue ...string) string {
	if len(defaultValue) > 0 {
		question += "[" + defaultValue[0] + "]"
	}
	if ask {
		question += "?"
	} else {
		question += ":"
	}
	question += " "
	w.write(question)
	a := w.read()
	if a == `''` || a == `""` || a == "``" {
		return ""
	}
	if a == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return a
}

func (w *ConfigWizard) Ask(question string, defVal ...bool) bool {
	def := []string{"y/N"}
	if len(defVal) > 0 && defVal[0] {
		def[0] = "Y/n"
	}

	res := w.prompt(question, true, def...)

	if res == "y" || res == "Y" || res == "Y/n" || res == "true" || strings.ToLower(res) == "yes" || res == "1" || res == "t" {
		return true
	}

	return false
}

func (w *ConfigWizard) Get(question string, defVal ...string) string {
	return w.prompt(question, false, defVal...)
}

func (w *ConfigWizard) GetInt(question string, defVal ...int) int {
	var def []string
	if len(defVal) > 0 {
		def = []string{strconv.Itoa(defVal[0])}
	}

	res := w.prompt(question, false, def...)

	if v, err := strconv.Atoi(res); err == nil {
		return v
	}

	return 0
}

func (w *ConfigWizard) Option(question string, options []string, required bool, defVal ...string) string {
	var def []string
	if len(defVal) > 0 {
		def = []string{defVal[0]}
	}

	res := w.prompt(question+"("+strings.Join(options, "|")+")", false, def...)

	if res == "" && required == false {
		return res
	}

	for _, opt := range options {
		if res == opt {
			return res
		}
	}

	w.write("Invalid answer\n")
	return w.Option(question, options, required, defVal...)
}

func (w *ConfigWizard) GetValid(validator func(string) error, question string, defVal ...string) string {
	res := w.Get(question, defVal...)
	err := validator(res)

	for err != nil {
		w.write(err.Error() + "\n")
		res = w.Get(question, defVal...)
		err = validator(res)
	}

	return res
}
