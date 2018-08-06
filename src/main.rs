extern crate postgres;
extern crate warp;

use postgres::{ Connection, TlsMode };
use std::fs::File;
use std::path::Path;
use std::io::prelude::*;

fn main() {
	println!("Howdy!!");
	let conn = Connection::connect("postgres://postgres@db:5433", TlsMode::None)
		.expect("Could not connect to postgres database");
	let mut f = File::open("./tables.sql").expect("file not found");

	let mut sql_setup = String::new();
	f.read_to_string(&mut sql_setup).expect("Couldn't read tables.sql");

	conn.batch_execute(&sql_setup).expect("Couldn't run setup sql");
	println!("test");
}