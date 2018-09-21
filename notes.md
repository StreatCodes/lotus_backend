## Reserved urls
 - `GET /` Obviously the home page
 - `POST /authorize`
 - `POST/GET /api/*` The api route
 - `GET /admin` Admin login interface
 - `GET /admin/*` Admin/site configuration pages

### TODO 
 - refactor page build into it's own function
 - Add GZIP 

### Caution / Needs review
 - Ordering pages will require a recompute of all order values with a common parent (Our first REST bottle neck!)
 - Deleting a page that is a parent (should be handled by reference (double check the is the same as a foriegn key))