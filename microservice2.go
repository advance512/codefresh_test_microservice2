package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	"github.com/labstack/echo/middleware"

	"github.com/caarlos0/env"
	"github.com/kr/pretty"
	raml "gopkg.in/raml.v0"
)

func verifyRAML(c echo.Context) error {
	fmt.Println("Handler for /v1/verifyRAML called.")

	requestBody, err := ioutil.ReadAll(c.Request().Body())

	if err != nil {
		c.Error(err)
		return err
	}

	tmpfile, err := ioutil.TempFile("", "ramltmp")
	if err != nil {
		c.Error(err)
		return err
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(requestBody); err != nil {
		c.Error(err)
		return err
	}

	if err := tmpfile.Close(); err != nil {
		c.Error(err)
		return err
	}

	if apiDefinition, err := raml.ParseFile(tmpfile.Name()); err != nil {
		return c.String(http.StatusOK, err.Error()) // Application error, not a runtime error
	} else {
		result := fmt.Sprintf("Successfully parsed RAML file:\n\n%s", pretty.Sprintf("%s", apiDefinition))
		return c.String(http.StatusOK, result)
	}
}

type Config struct {
	Ms1Port int `env:"MS1PORT" envDefault:"3000"`
	Ms2Port int `env:"PORT" envDefault:"3000"`
}

// Global configuration
var _configuration = Config{}

func main() {
	e := echo.New()
	env.Parse(&_configuration)

	fmt.Println("Configuration: ", _configuration)

	e.Use(middleware.Recover())

	e.POST("/v1/verifyRAML", verifyRAML)

	fmt.Printf("Microservice 2 is now listening on port %d...", _configuration.Ms2Port)

	e.Run(standard.New(fmt.Sprintf(":%d", _configuration.Ms2Port)))

}
