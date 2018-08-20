use warp::Filter;
use warp::filters::BoxedFilter;
use warp::{path, header};

pub fn bye_handler() -> BoxedFilter<(i32, String)> {
    path::param::<i32>()
        .and(header::<String>("host"))
        .boxed()
}