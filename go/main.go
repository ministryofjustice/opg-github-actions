package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"opg-github-actions/cmd/branchname"
	"opg-github-actions/cmd/createtag"
	"opg-github-actions/cmd/latesttag"
	"opg-github-actions/cmd/nexttag"
	"opg-github-actions/cmd/safestring"
	"opg-github-actions/cmd/terraformversion"
	"os"
	"slices"
	"strings"

	"log/slog"
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
		level          string     = *logLevel
		as             string     = *logAs
		to             string     = *logTo
		validAsChoice  bool       = slices.Contains(logAsChoices, as)
		validToChoice  bool       = slices.Contains(logToChoices, to)
		out            io.Writer  = os.Stdout
		logLevel       slog.Level = slog.LevelError
		handlerOptions *slog.HandlerOptions
		log            *slog.Logger
	)
	// setup log level
	if l, ok := logLevels[level]; ok {
		logLevel = l
	}
	handlerOptions = &slog.HandlerOptions{AddSource: true, Level: logLevel}
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

var (
	workspace string = "GITHUB_WORKSPACE"
	githubout string = "GITHUB_OUTPUT"
)

// Out provides a wrapper to echo results information out to the stdout
func Out(results map[string]string) {
	var outFile string
	var f *os.File
	// github output via os.environ['GITHUB_OUTPUT']
	if os.Getenv(workspace) != "" {
		slog.Info("github workspace found, writting to GITHUB_OUTPUT")
		outFile = os.Getenv(githubout)
		f, _ = os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
	} else {
		slog.Info("github workspace not found...")
	}

	for k, v := range results {
		str := fmt.Sprintf("%s=%s\n", k, v)
		fmt.Printf(str)
		if outFile != "" && f != nil {
			slog.Debug("writting to GITHUB_OUTPUT")
			f.WriteString(str)
		}
	}
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

	slog.Info("flag parsed command:" + cmd)
	slog.Error(fmt.Sprintf("arguments: [%s]", strings.Join(flag.Args(), " ")))

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
		log.Fatal(err.Error())
	}

	Out(results)

}
