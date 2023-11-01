package routes
// If you keep the name of package same as name of folder, it becomes easy
// for you to use the functions of this package in other packages

import (
	"github.com/23piyush/go-redis-fiber-url-shortener/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// ResolveURL ...
func ResolveURL(c *fiber.Ctx) error {
	// get the short from the url
	url := c.Params("url")
	// query the db to find the original URL, if a match is found
	// increment the redirect counter and redirect to the original URL
	// else return error message
	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result() // redis is a key-value pair database, corresponding to key - "url", you will get the value
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found on database", // connected to database, but this url was not found
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to DB",  // couldnot even connect to database
		})
	}
	// increment the counter
	rInr := database.CreateClient(1)
	defer rInr.Close()
	_ = rInr.Incr(database.Ctx, "counter")
	// redirect to original URL
	return c.Redirect(value, 301) // with status 301, redirect to original URL
}
// We are storing both original URL and shortened URL in the redis database
// Once we get shortened URL to resolve, find the original URL of this shortened URL and redirect user to the original URL