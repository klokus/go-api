# go-api
This is a simple, fast, and easy to use API written in golang

usage: ./goapiv1 0.0.0.0 [webserver_port] [api_key]

request format: http://host/attack?key=[api_key]&host=[host]&port=[port]&time=[time]&method=[method]

request format with custom port: http://host:[webserver_port]/attack?key=[api_key]&host=[host]&port=[port]&time=[time]&method=[method]
