## Reserved urls
 - `GET /` Obviously the home page
 - `POST /authorize`
 - `POST/GET /api/*` The api route
 - `GET /admin` Admin login interface
 - `GET /admin/*` Admin/site configuration pages

### TODO 
 - Create a working login
 - Create and admin interface where we can create, delete update and sort pages
 - Create and interface to add, remove, sort and update components
 - Clean up
 - Create some
 - refactor page build into it's own function

### Caution / Needs review
 - Ordering pages will require a recompute of all order values with a common parent (Our first REST bottle neck!)
 - Deleting a page that is a parent (should be handled by reference (double check the is the same as a foriegn key))