package main

import "core:fmt"
import "core:strings"
import sql "ext:sqlite3"

Database :: sql.sqlite3

db_open :: proc(path: cstring, db: ^^Database)
{
	sql.open(path, db)
}

db_close :: proc(db: ^Database)
{
	sql.close(db)
}

db_create_user_table :: proc(db: ^Database, name: string)
{
	schema: cstring = "CREATE TABLE Users(id int, name varchar(255));"
	stmt: ^sql.Stmt
	res := sql.prepare_v2(db, cast(^u8) schema, cast(i32) len(schema), &stmt, nil)
	db_check_result(db, res)
	
	res = sql.step(stmt)
	db_check_result(db, res)

	sql.finalize(stmt)
}

db_select_user :: proc(db: ^Database, name: string) -> (result: User)
{
	schema: cstring = "SELECT * FROM Users;"
	stmt: ^sql.Stmt
	res := sql.prepare_v2(db, cast(^u8) schema, cast(i32) len(schema), &stmt, nil)
	db_check_result(db, res)

	res = sql.step(stmt)
	db_check_result(db, res)

	id_text := strings.clone_from_cstring(sql.column_text(stmt, 0))
	name_text := strings.clone_from_cstring(sql.column_text(stmt, 1))
	result.name = strings.concatenate({id_text, name_text})

	sql.finalize(stmt)

	return
}

db_insert_user :: proc(db: ^Database, user: User)
{
	schema: cstring = "INSERT INTO Users (id, name) VALUES (3, 'Daniel');"
	stmt: ^sql.Stmt
	res := sql.prepare_v2(db, cast(^u8) schema, cast(i32) len(schema), &stmt, nil)
	db_check_result(db, res)
	
	res = sql.step(stmt)
	db_check_result(db, res)

	sql.finalize(stmt)
}

@(private="file")
db_check_result :: proc(db: ^Database, res: sql.ResultCode, loc := #caller_location)
{
	fmt.printf("Result: %s\n", res)

	if res == .ERROR || res == .CONSTRAINT || res == .MISUSE
  {
		text := fmt.tprintf("%s %v %s", res, sql.errmsg(db))
		panic(text, loc)
	}
}

db_test :: proc(db: ^Database)
{
	fmt.println(" - CREATE -")
	db_create_user_table(db, "X")
	fmt.println(" - INSERT - ")
	db_insert_user(db, User{})
	fmt.println(" - SELECT - ")
	result := db_select_user(db, "X")
	fmt.println("Result:", result)
}
