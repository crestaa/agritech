package api

import (
	"agritech/server/database"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Serve() {
	r := gin.Default()

	r.Static("/web", "./static")

	apiHandler(r)

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func apiHandler(r *gin.Engine) {
	apiGroup := r.Group("/api")

	apiGroup.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "API test")
	})

	apiGroup.GET("/fields", func(c *gin.Context) {
		r, err := returnCampi()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, r)
	})
	apiGroup.GET("/fields/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		ret, err := returnCampoByID(id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, ret)
	})
	apiGroup.GET("/fields/:id/readings", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		r, err := returnFieldReadings(id, -1)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, r)
	})
	apiGroup.GET("/fields/:id/readings/:type", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		t, err := strconv.Atoi(c.Param("type"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		r, err := returnFieldReadings(id, t)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, r)
	})
	apiGroup.GET("/fields/:id/sensors/", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		r, err := returnFieldSensors(id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, r)
	})

	apiGroup.GET("/sensors", func(c *gin.Context) {
		r, err := returnSensori()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, r)
	})
	apiGroup.GET("/sensors/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		ret, err := returnSensoreByID(id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, ret)
	})

}

func returnCampoByID(id int) (string, error) {
	campo, err := database.GetCampo(id)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	ret, err := json.Marshal(campo)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return bytes.NewBuffer(ret).String(), nil
}

func returnCampi() (string, error) {
	campi, err := database.GetCampi()
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	ret, err := json.Marshal(campi)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return bytes.NewBuffer(ret).String(), nil
}

func returnSensori() (string, error) {
	sen, err := database.GetSensori()

	if err != nil {
		fmt.Print(err)
		return "", err
	}
	ret, err := json.Marshal(sen)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return bytes.NewBuffer(ret).String(), nil
}

func returnSensoreByID(id int) (string, error) {
	sen, err := database.GetSensorByID(id)

	if err != nil {
		fmt.Print(err)
		return "", err
	}
	ret, err := json.Marshal(sen)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return bytes.NewBuffer(ret).String(), nil
}

func returnFieldReadings(id int, t int) (string, error) {
	reads, err := database.GetReadings(id, t)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	ret, err := json.Marshal(reads)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return bytes.NewBuffer(ret).String(), nil
}

func returnFieldSensors(id int) (string, error) {
	reads, err := database.GetFieldSensors(id)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	ret, err := json.Marshal(reads)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return bytes.NewBuffer(ret).String(), nil
}
