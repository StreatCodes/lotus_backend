## Install
 - Create a .env with the following:
```
LOTUS_ENV="dev"
LOTUS_HTTP_ADDR="0.0.0.0"
LOTUS_HTTP_PORT="3001"
LOTUS_DB_USER="lotus"
LOTUS_DB_PASS="lotus"
LOTUS_DB_ADDR="localhost"
LOTUS_DB_NAME="lotus"
```
 - Create the corresponding database and database users.
 - Create a link to the admin frontend `ln -s ~/src/node/lotus_frontend admin`
 - Build and run `go build; ./lotus_backend`

## Reserved urls
 - `GET /` Obviously the home page
 - `POST /authorize`
 - `POST/GET /api/*` The api route
 - `GET /admin` Admin login interface
 - `GET /admin/*` Admin/site configuration pages

### TODO 
 - Create and admin interface where we can create, delete update and sort pages
 - Create and interface to add, remove, sort and update components
 - Clean up
 - refactor page build into it's own function

### Caution / Needs review
 - Ordering pages will require a recompute of all order values with a common parent (Our first REST bottle neck!)
 - Deleting a page that is a parent (should be handled by reference (double check the is the same as a foriegn key))