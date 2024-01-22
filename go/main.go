package main

import (
	"flag"
	"fmt"
	"io"
	"opg-github-actions/cmd/branchname"
	"opg-github-actions/cmd/createtag"
	"opg-github-actions/cmd/latesttag"
	"opg-github-actions/cmd/nexttag"
	"opg-github-actions/cmd/safestring"
	"opg-github-actions/cmd/terraformversion"
	"os"
	"slices"

	"log/slog"

	"github.com/k0kubun/pp"
)

// Log setup options
var (
	logLevels = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	logAsChoices = []string{"text", "json"}
	logToChoices = []string{"stdout", "file"}
	logLevel     = flag.String("log-level", "error", "Logging level to use, must be one of (debug, info, warn, error). Default: error")
	logAs        = flag.String("log-as", "text", "Log as this format type, must be either (text, json). Default: text")
	logTo        = flag.String("log-to", "stdout", "Log to either (stdout, file). Default: stdout")
	logFile      *os.File
)

func logSetup() {
	var (
		level          string               = *logLevel
		as             string               = *logAs
		to             string               = *logTo
		handlerOptions *slog.HandlerOptions = &slog.HandlerOptions{AddSource: true, Level: slog.LevelError}
		validAsChoice  bool                 = slices.Contains(logAsChoices, as)
		validToChoice  bool                 = slices.Contains(logToChoices, to)
		out            io.Writer            = os.Stdout
		log            *slog.Logger
	)
	// setup log level
	if l, ok := logLevels[level]; ok {
		handlerOptions.Level = l
	}
	// if chosen to change output to file, open the file and adjust out
	if validToChoice && to == "file" {
		logFile, _ = os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		out = logFile
	}

	if validAsChoice && as == "json" {
		log = slog.New(slog.NewJSONHandler(out, handlerOptions))
	} else {
		log = slog.New(slog.NewTextHandler(out, handlerOptions))
	}
	slog.SetDefault(log)

}

func main() {
	var (
		err     error
		results map[string]string
	)

	flag.Parse()
	cmd := flag.Arg(0)

	// configure log output
	logSetup()
	if logFile != nil {
		defer logFile.Close()
	}

	slog.Debug("flag parsed command:" + cmd)
	switch cmd {
	case branchname.Name:
		results, err = branchname.Run(flag.Args()[1:])
	case safestring.Name:
		results, err = safestring.Run(flag.Args()[1:])
	case terraformversion.Name:
		results, err = terraformversion.Run(flag.Args()[1:])
	case latesttag.Name:
		results, err = latesttag.Run(flag.Args()[1:])
	case nexttag.Name:
		results, err = nexttag.Run(flag.Args()[1:])
	case createtag.Name:
		results, err = createtag.Run(flag.Args()[1:])

	default:
		err = fmt.Errorf("Sub command [%s] not recognised", cmd)
	}

	if err != nil {
		slog.Error(err.Error())
	}

	pp.Println(results)
}