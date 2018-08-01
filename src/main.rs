extern crate rusqlite;

use rusqlite::Connection;
use std::fs::File;
use std::path::Path;
use std::io::prelude::*;

fn main() {
	if Path::new("lotus.db").exists() {
		std::fs::remove_file("lotus.db").expect("Couldn't delete database");
	}

	let conn = Connection::open("lotus.db").expect("Couldn't open lotus.db");
	let mut f = File::open("tables.sql").expect("file not found");

	let mut sql_setup = String::new();
	f.read_to_string(&mut sql_setup).expect("Couldn't read tables.sql");

	conn.execute_batch(&sql_setup).expect("Couldn't run setup sql");
	println!("test");
}