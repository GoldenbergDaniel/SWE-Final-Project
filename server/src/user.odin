package main

import "core:encoding/json"

User :: struct
{
  name: string,
}

user_from_json :: proc(s: string) -> User
{

}

json_from_user :: proc(user: User) -> string
{

}
