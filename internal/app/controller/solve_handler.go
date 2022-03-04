package controller

import (
	"net/http"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kannman/waterbucketchallenge/internal/app/waterjug"
)

// FindSolutionHandler godoc
// @Summary      Find watcher bucket challenge solution
// @Tags  		 v1
// @Description  some description
// @Accept       json
// @Produce      json
// @Param        Buckets body FindSolutionRequest true "Buckets"
// @Success      200 {object} FindSolutionResponse
// @Failure      400 {object} ErrResponse "Bad request"
// @Router       /v1/search-solution [post]
func FindSolutionHandler(w http.ResponseWriter, r *http.Request) {
	var (
		request  FindSolutionRequest
		response FindSolutionResponse
	)

	if err := Bind(r, &request); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	solution, found := waterjug.SearchSolution(request.X, request.Y, request.Z)
	response.Solved = found
	for _, node := range solution {
		response.Steps = append(response.Steps, SolveStep{
			X:  node.X,
			Y:  node.Y,
			Op: node.Op,
		})
	}
	render.JSON(w, r, response)
}

type FindSolutionRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

func (s FindSolutionRequest) Bind(*http.Request) error { return nil }

func (s *FindSolutionRequest) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.X, validation.Required, validation.Min(1)),
		validation.Field(&s.Y, validation.Required, validation.Min(1)),
		validation.Field(&s.Z, validation.Required, validation.Min(1)),
	)
}

type FindSolutionResponse struct {
	Steps  []SolveStep `json:"steps"`
	Solved bool        `json:"solved"`
}

type SolveStep struct {
	X  int    `json:"x"`
	Y  int    `json:"y"`
	Op string `json:"op"`
}
