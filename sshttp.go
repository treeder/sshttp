package main

import (
	"flag"
	"github.com/bmizerany/pat"
	"github.com/spf13/viper"
	"github.com/treeder/sshttp/system"
	"gopkg.in/inconshreveable/log15.v2"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const VERSION = "0.1"

func main() {
	Log := log15.New("svc", "sshttp", "v", VERSION)

	viper.AutomaticEnv()
	viper.SetDefault("verbosity", "info")
	viper.SetDefault("port", 8022)
	viper.SetDefault("ssl", false)

	var (
		verbosity = flag.String("v", "info", "Verbosity (options: debug, info, warn, error, crit)")
		port      = flag.Int("p", 8022, "Port to expose HTTP service")
		ssl       = flag.Bool("ssl", false, "Enable TLS")
		token     = flag.String("t", "", "Auth token to require in HTTP requests")
	)

	flag.Parse()

	if *verbosity != "" {
		viper.Set("verbosity", verbosity)
	}
	if *port != 0 {
		viper.Set("port", port)
	}
	if *ssl {
		viper.Set("ssl", ssl)
	}
	if *token != "" {
		viper.Set("token", token)
	}

	lvl, err := log15.LvlFromString(viper.GetString("verbosity"))
	if err != nil {
		Log.Error("Invalid verbosity selected (-v <verbosity>)")
		os.Exit(1)
	}
	Log.Info("Verbosity level is " + viper.GetString("verbosity"))
	Log.SetHandler(log15.LvlFilterHandler(lvl, log15.StdoutHandler))

	if viper.GetString("token") != "" {
		SetToken(viper.GetString("token"))
	} else {
		Log.Error("Token argument must be provided (-t <token>)")
		os.Exit(1)
	}

	Log.Info("Starting sshttp on port", viper.GetInt("port"))

	Log.Debug("Configuration", log15.Ctx{
		"port":  viper.GetInt("port"),
		"ssl":   viper.GetBool("ssl"),
		"token": viper.GetString("token"),
	})

	buildServer(Log)

	portStr := ":" + strconv.Itoa(viper.GetInt("port"))

	if viper.GetBool("ssl") {
		Log.Debug("Serving on https://localhost" + portStr)
		err := http.ListenAndServeTLS(portStr, "cert.pem", "key.pem", nil)
		if err != nil {
			Log.Error("ListenAndServeTLS: " + err.Error())
			os.Exit(1)
		}
	} else {
		Log.Warn("Server is not using SSL")
		Log.Debug("Serving on http://localhost" + portStr)
		err := http.ListenAndServe(portStr, nil)
		if err != nil {
			Log.Error("ListenAndServe: " + err.Error())
			os.Exit(1)
		}
	}
}

func buildServer(mainLog log15.Logger) {
	base := "/v1"

	// Switching to a mux lib to handle not found errors, etc. We were returning blank pages if method wasn't allowed for instance.
	m := pat.New()
	m.Get("/", http.HandlerFunc(hi))
	m.Post(base+"/shell", middlewareAuth(mainLog, system.ShellHandler))
	http.Handle("/", m)

	routes := map[string]http.HandlerFunc{
		"/system":           system.SystemHandler,
		"/system/ram":       system.RamHandler,
		"/system/load":      system.LoadHandler,
		"/system/host":      system.HostHandler,
		"/system/disk":      system.DiskHandler,
		"/system/cpuinfo":   system.CpuHandler,
		"/system/processes": system.ProcessesHandler,
		"/files":            system.FHandler,
	}

	for path, handler := range routes {
		http.HandleFunc(base+path, middlewareAuth(mainLog, handler))
	}
}

func hi(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hi")
}

type authHandlers struct {
	l    log15.Logger
	next http.HandlerFunc
}

func (a *authHandlers) checkAuthHandler(res http.ResponseWriter, req *http.Request) {
	if userToken := getTokenString(req); userToken != token {
		http.Error(res, "Not Authorized", http.StatusUnauthorized)
		return
	}
	a.next(res, req)
}

// Authentication middleware
func middlewareAuth(l log15.Logger, f http.HandlerFunc) http.HandlerFunc {
	return (&authHandlers{l, f}).checkAuthHandler
}

var token string

func SetToken(setToken string) {
	token = setToken
}

func getTokenString(r *http.Request) string {
	tokenStr := r.URL.Query().Get("oauth")
	if tokenStr == "" {
		authHeader := r.Header.Get("Authorization")
		authFields := strings.Fields(authHeader)
		if len(authFields) == 2 && authFields[0] == "OAuth" {
			tokenStr = authFields[1]
		}
	}
	if tokenStr == "" {
		tokenStr = r.FormValue("token")
	}
	return tokenStr
}
