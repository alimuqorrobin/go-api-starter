package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-api-starter/internal/domain"
	"go-api-starter/internal/http/response"
	"go-api-starter/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	product, err := h.productService.CreateProduct(r.Context(), &req)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Created(w, product)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "Invalid product ID")
		return
	}

	product, err := h.productService.GetProduct(r.Context(), uint(id))
	if err != nil {
		if err == domain.ErrProductNotFound {
			response.NotFound(w, "Product not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, product)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit := 10
	page := 1

	if l := r.URL.Query().Get("limit"); l != "" {
		if pl, err := strconv.Atoi(l); err == nil && pl > 0 && pl <= 100 {
			limit = pl
		}
	}

	if p := r.URL.Query().Get("page"); p != "" {
		if pp, err := strconv.Atoi(p); err == nil && pp > 0 {
			page = pp
		}
	}

	offset := (page - 1) * limit

	products, total, err := h.productService.GetAllProducts(r.Context(), limit, offset)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Paginated(w, products, page, limit, total)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "Invalid product ID")
		return
	}

	var req domain.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	product, err := h.productService.UpdateProduct(r.Context(), uint(id), &req)
	if err != nil {
		if err == domain.ErrProductNotFound {
			response.NotFound(w, "Product not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "Invalid product ID")
		return
	}

	if err := h.productService.DeleteProduct(r.Context(), uint(id)); err != nil {
		if err == domain.ErrProductNotFound {
			response.NotFound(w, "Product not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Product deleted successfully"})
}
