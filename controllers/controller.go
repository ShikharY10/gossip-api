package controllers

import (
	"github.com/ShikharY10/gbAPI/config"
	v1 "github.com/ShikharY10/gbAPI/controllers/v1"
	"github.com/ShikharY10/gbAPI/handler"
	"github.com/ShikharY10/gbAPI/logger"
	"github.com/ShikharY10/gbAPI/middleware"
)

type API_V3 struct {
	Handler    *handler.Handler
	Middleware *middleware.Middleware
	Cloudinary *config.Cloudinary
}

func GetController(handler *handler.Handler, middleware *middleware.Middleware, cloudinary *config.Cloudinary, logging *logger.Logger) (*v1.UserController, *v1.PartnerController, *v1.PostController) {
	userController := v1.UserController{
		Handler:    handler,
		Middleware: middleware,
		Cloudinary: cloudinary,
		Logging:    logging,
	}

	partnerController := v1.PartnerController{
		Handler:    handler,
		Middleware: middleware,
		Cloudinary: cloudinary,
		Logging:    logging,
	}

	postController := v1.PostController{
		Handler:    handler,
		Middleware: middleware,
		Cloudinary: cloudinary,
		Logging:    logging,
	}

	return &userController, &partnerController, &postController
}
