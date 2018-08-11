#[macro_use]
extern crate warp;
extern crate postgres;

use std::fs::File;
use std::io::prelude::*;
use postgres::{ Connection, TlsMode };
use warp::Filter;

fn main() {
	println!("Waiting for postgres");
	std::thread::sleep(std::time::Duration::from_millis(800));

	let conn = Connection::connect("postgres://postgres:computer@db", TlsMode::None)
		.expect("Could not connect to postgres database");
	
	println!("Setting up tables");
	let mut f = File::open("./tables.sql").expect("file not found");

	let mut sql_setup = String::new();
	f.read_to_string(&mut sql_setup).expect("Couldn't read tables.sql");

	conn.batch_execute(&sql_setup).expect("Couldn't run setup sql");

    // GET /hello/warp => 200 OK with body "Hello, warp!"
    let hello = path!("hello" / String)
        .map(|name| format!("Hello, {}!", name));

    warp::serve(hello)
        .run(([127, 0, 0, 1], 3030));
}