package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *appDependencies) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.notAllowedResponse)

	// router.HandlerFunc(http.MethodGet, "/api/v1/books", a.getBooks)
	// something else here
	router.HandlerFunc(http.MethodPost, "/api/v1/books", a.postBook)

	// router.HandlerFunc(http.MethodPost, "/createProduct", a.createProduct)
	// router.HandlerFunc(http.MethodGet, "/displayProduct/:id", a.displayProduct)
	// router.HandlerFunc(http.MethodDelete, "/deleteProduct/:id", a.deleteProduct)
	// router.HandlerFunc(http.MethodGet, "/displayAllProducts", a.displayAllProducts)
	// router.HandlerFunc(http.MethodPatch, "/updateProduct/:id", a.updateProduct)

	// router.HandlerFunc(http.MethodPost, "/product/:id/createReview", a.createReview)
	// router.HandlerFunc(http.MethodGet, "/product/:id/getReview/:rid", a.getReview)
	// router.HandlerFunc(http.MethodPatch, "/product/:id/updateReview/:rid", a.updateReview)
	// router.HandlerFunc(http.MethodDelete, "/product/:id/deleteReview/:rid", a.deleteReview)

	// router.HandlerFunc(http.MethodGet, "/reviews", a.GetAllReviews)

	return a.recoverPanic(a.rateLimit(router))
}
