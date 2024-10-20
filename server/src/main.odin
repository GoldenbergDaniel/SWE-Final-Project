package main

import "core:fmt"
import "core:net"
// import "core:encoding/json"

import "ext:http"
import "src:basic/mem"

global: Global

main :: proc()
{
	// Test
	{
		db: ^Database
		db_open("data.db", &db)
		db_test(db)
		
		if true do return
	}

	mem.init_arena_growing(&global.temp_arena)

	server: http.Server
	http.server_shutdown_on_interrupt(&server)

	router: http.Router
	http.router_init(&router)
	http.route_get(&router, "/api", http.handler(get_something))
  http.route_post(&router, "/signin", http.handler(post_signin))
  http.route_post(&router, "/signup", http.handler(post_signup))
	routed := http.router_handler(&router)

	fmt.println("Listening on http://localhost:5174")
	err := http.listen_and_serve(&server, routed, net.Endpoint{net.IP4_Loopback, 5174})
  if err != nil
  {
    fmt.eprintln("Server stopped with error:", err)
  }
}

get_something :: proc(req: ^http.Request, res: ^http.Response)
{
	if err := http.respond_json(res, req.line); err != nil
  {
		fmt.eprintf("could not respond with JSON: %s", err)
	}
}

post_signin :: proc(req: ^http.Request, res: ^http.Response)
{
  enable_cors(&res.headers)

	http.body(req, -1, res, proc(
		res: rawptr, 
		body: http.Body, 
		err: http.Body_Error)
	{
		res := cast(^http.Response) res

		if err != nil
    {
			http.respond(res, http.body_error_status(err))
      fmt.eprintln("Error with response:", err)
			return
		}

		http.respond_plain(res, "sign-in: hello from the server")
	})
}

post_signup :: proc(req: ^http.Request, res: ^http.Response)
{
  enable_cors(&res.headers)
  
	http.body(req, -1, res, proc(res: rawptr, body: http.Body, err: http.Body_Error) {
		res := cast(^http.Response) res

		if err != nil
    {
			http.respond(res, http.body_error_status(err))
      fmt.eprintln("Error with response:", err)
			return
		}

		http.respond_plain(res, "sign-up: hello from the server")
	})
}

enable_cors :: proc(headers: ^http.Headers)
{
  http.headers_set(headers, "Access-Control-Allow-Origin", "*")
  http.headers_set(headers, "Access-Control-Allow-Methods", "GET, POST, PUT, PATCH")
  http.headers_set(headers, "Access-Control-Allow-Headers", "Content-Type")
}

Global :: struct
{
	temp_arena: mem.Arena,
}

// User //////////////////////////////////////////////////////////////////////////////////

User_Data :: struct
{
	
}

User_Name  :: distinct string
