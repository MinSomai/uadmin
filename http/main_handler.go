package http

import (
	"context"
	"fmt"
	authapi "github.com/uadmin/uadmin/blueprint/auth/api"
	userapi "github.com/uadmin/uadmin/blueprint/user/api"
	settingsapi "github.com/uadmin/uadmin/blueprint/settings/api"
	"github.com/uadmin/uadmin/metrics"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/utils"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// mainHandler is the main handler for the admin
func mainHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckRateLimit(r) {
		w.Write([]byte("Slow down. You are going too fast!"))
		return
	}
	if !authapi.ValidateIP(r, preloaded.AllowedIPs, preloaded.BlockedIPs) {
		if r.Form == nil {
			r.Form = url.Values{}
		}
		r.Form.Set("err_msg", "Your IP Address ("+r.RemoteAddr+") is not Allowed to Access this Page")
		r.Form.Set("err_code", "403")
		PageErrorHandler(w, r, nil)
		return
	}
	r.URL.Path = strings.TrimPrefix(r.URL.Path, preloaded.RootURL)
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	URLParts := strings.Split(r.URL.Path, "/")

	if URLParts[0] == "resetpassword" {
		userapi.PasswordResetHandler(w, r)
		return
	}

	// Authentecation
	// This session is preloaded with a user
	session := authapi.IsAuthenticated(r)
	if session == nil {
		userapi.LoginHandler(w, r)
		return
	}

	// Check remote access
	if !(utils.IsLocal(r.RemoteAddr) || session.User.RemoteAccess) {
		if r.Form == nil {
			r.Form = url.Values{}
		}
		r.Form.Set("err_msg", "Remote Access Denied")
		PageErrorHandler(w, r, nil)
		return
	}

	if r.URL.Path == "" {
		homeHandler(w, r, session)
		return
	}
	if len(URLParts) == 1 {
		if URLParts[0] == "logout" {
			userapi.LogoutHandler(w, r, session)
			return
		}
		if URLParts[0] == "export" {
			exportHandler(w, r, session)
			return
		}
		if URLParts[0] == "cropper" {
			cropImageHandler(w, r, session)
			return
		}
		if URLParts[0] == "profile" {
			userapi.ProfileHandler(w, r, session)
			return
		}
		if URLParts[0] == "settings" {
			settingsapi.SettingsHandler(w, r, session)
			return
		}
		listHandler(w, r, session)
		return
	} else if len(URLParts) == 2 {
		formHandler(w, r, session)
		return
	}
	PageErrorHandler(w, r, session)
}

// Handler is a function that takes an http handler function and returns an http handler function
// that has extra functionality including logging
func Handler(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
			Prepare log message. Valid place holders:
			A perfect list (not fully inplemented): http://httpd.apache.org/docs/current/mod/mod_log_config.html
			 %a: Client IP address
			 %{remote}p: Client port
			 %A: Server hostname/IP
			 %{local}p: Server port
			 %U: Path
			 %c: All coockies
			 %{NAME}c: Cookie named 'NAME'
			 %{GET}f: GET request parameters
			 %{POST}f: POST request parameters
			 %B: Respnse length
			 %>s: Response code
			 %D: Time taken in microseconds
			 %T: Time taken in seconds
			 %I: Request length
			 %i: All headers
			 %{NAME}i: header named 'NAME'
		*/
		HTTP_LOG_MSG := preloaded.HTTPLogFormat
		host, port, _ := net.SplitHostPort(r.RemoteAddr)
		HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%a", host, -1)
		HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%{remote}p", port, -1)
		host, port, _ = net.SplitHostPort(r.Host)
		HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%A", host, -1)
		HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%{local}p", port, -1)
		HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%U", r.URL.Path, -1)
		HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%I", fmt.Sprint(r.ContentLength), -1)

		// Process cookies
		if strings.Contains(HTTP_LOG_MSG, "%c") {
			v := []string{}
			for _, c := range r.Cookies() {
				v = append(v, c.Name+"="+c.Value)
			}
			HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%c", strings.Join(v, "&"), -1)
		}
		re := regexp.MustCompile(`%{[^ ;,}]*}c`)
		cookies := re.FindAll([]byte(HTTP_LOG_MSG), -1)
		for _, cookie := range cookies {
			cookieName := strings.TrimPrefix(string(cookie), "%{")
			cookieName = strings.TrimSuffix(cookieName, "}c")
			c, _ := r.Cookie(cookieName)
			if c != nil {
				HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, string(cookie), c.Name+"="+c.Value, -1)
			} else {
				HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, string(cookie), c.Name+"=", -1)
			}
		}

		// Process headers
		if strings.Contains(HTTP_LOG_MSG, "%i") {
			v := []string{}
			for k := range r.Header {
				v = append(v, k+"="+r.Header.Get(k))
			}
			HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%i", strings.Join(v, "&"), -1)
		}
		re = regexp.MustCompile(`%{[^ ;,}]*}i`)
		headers := re.FindAll([]byte(HTTP_LOG_MSG), -1)
		for _, header := range headers {
			headerName := strings.TrimPrefix(string(header), "%{")
			headerName = strings.TrimSuffix(headerName, "}i")
			h := r.Header.Get(headerName)
			HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, string(header), headerName+"="+h, -1)
		}

		// Process GET/POST parameters
		if strings.Contains(HTTP_LOG_MSG, "%{GET}f") {
			HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%{GET}f", r.URL.RawQuery, -1)
		}
		if strings.Contains(HTTP_LOG_MSG, "%{POST}f") {
			v := []string{}
			err := r.ParseMultipartForm(32 << 20)
			if err != nil {
				r.ParseForm()
			}
			for key, val := range r.PostForm {
				v = append(v, key+"=["+strings.Join(val, ",")+"]")
			}
			HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%{POST}f", strings.Join(v, "&"), -1)
		}

		// Add context with stime
		ctx := context.WithValue(r.Context(), preloaded.CKey("start"), time.Now())
		r = r.WithContext(ctx)
		res := responseWriter{
			w: w,
		}

		// Execute the actual handler
		f(&res, r)

		// add etime
		ctx = context.WithValue(r.Context(), preloaded.CKey("end"), time.Now())
		r = r.WithContext(ctx)

		// Add post execution context
		// response counter
		HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%B", fmt.Sprint(res.GetCounter()), -1)

		// response code
		HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%>s", fmt.Sprint(res.GetCode()), -1)

		// time taken
		sTime := r.Context().Value(preloaded.CKey("start")).(time.Time)
		eTime := r.Context().Value(preloaded.CKey("end")).(time.Time)
		if strings.Contains(HTTP_LOG_MSG, "%D") {
			HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%D", fmt.Sprint(eTime.Sub(sTime).Nanoseconds()/1000), -1)
		}
		if strings.Contains(HTTP_LOG_MSG, "%T") {
			HTTP_LOG_MSG = strings.Replace(HTTP_LOG_MSG, "%T", fmt.Sprintf("%0.3f", float64(eTime.Sub(sTime).Nanoseconds())/1000000000), -1)
		}

		// Log Metrics
		metrics.SetMetric("uadmin/http/responsetime", float64(eTime.Sub(sTime).Nanoseconds()/1000000))
		metrics.IncrementMetric("uadmin/http/requestrate")

		go func() {
			if preloaded.LogHTTPRequests {
				// Send log to syslog
				utils.Syslogf(utils.INFO, HTTP_LOG_MSG)
			}
		}()
	}
}

type responseWriter struct {
	w       http.ResponseWriter
	counter uint64
	code    int
}

func (w *responseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.counter += uint64(len(b))
	return w.w.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.w.WriteHeader(statusCode)
}

func (w *responseWriter) GetCode() int {
	if w.code == 0 {
		return 200
	}
	return w.code
}

func (w *responseWriter) GetCounter() uint64 {
	return w.counter
}
