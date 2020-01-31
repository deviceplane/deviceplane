package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestChainOrder(t *testing.T) {
	orders := []int{}
	set := []bool{}

	middlewareMaker := func(count int) func(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
		return func(with *FetchObject, hf http.HandlerFunc) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if with.Project == nil {
					set = append(set, false)
					with.Project = &models.Project{}
				} else {
					set = append(set, true)
				}

				orders = append(orders, count)
				hf.ServeHTTP(w, r)
			})
		}
	}
	handler := func(with *FetchObject) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if with.Project == nil {
				set = append(set, false)
			} else {
				set = append(set, true)
			}

			orders = append(orders, -1)
			w.WriteHeader(200)
			return
		})
	}

	first := middlewareMaker(1)
	second := middlewareMaker(2)
	third := middlewareMaker(3)

	s := Service{}
	chainedHandlers := s.initWith(first, second, third)(handler)

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chainedHandlers.ServeHTTP(recorder, req)

	assert.Equal(t, []int{1, 2, 3, -1}, orders)
	assert.Equal(t, []bool{false, true, true, true}, set)
}
