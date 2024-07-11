package controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/controller"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userObjectID := primitive.NewObjectID()
		userID := userObjectID.Hex()

		mockTaskUsecase := new(mocks.TaskUsecase)

		mockTaskUsecase.On("Create", mock.Anything, mock.Anything).Return(nil)

		gin := gin.Default()

		rec := httptest.NewRecorder()

		tc := &controller.TaskController{
			TaskUsecase: mockTaskUsecase,
		}

		gin.Use(setUserID(userID))
		gin.POST("/task", tc.Create)

		body, err := json.Marshal(domain.SuccessResponse{Message: "Task created successfully"})
		assert.NoError(t, err)

		bodyString := string(body)

		data := url.Values{}
		data.Set("title", "Test Task")
		req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Equal(t, bodyString, rec.Body.String())

		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		userObjectID := primitive.NewObjectID()
		userID := userObjectID.Hex()

		mockTaskUsecase := new(mocks.TaskUsecase)

		customErr := errors.New("Unexpected")

		mockTaskUsecase.On("Create", mock.Anything, mock.Anything).Return(customErr)

		gin := gin.Default()

		rec := httptest.NewRecorder()

		tc := &controller.TaskController{
			TaskUsecase: mockTaskUsecase,
		}

		gin.Use(setUserID(userID))
		gin.POST("/task", tc.Create)

		body, err := json.Marshal(domain.ErrorResponse{Message: customErr.Error()})
		assert.NoError(t, err)

		bodyString := string(body)

		data := url.Values{}
		data.Set("title", "Test Task")
		req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		assert.Equal(t, bodyString, rec.Body.String())

		mockTaskUsecase.AssertExpectations(t)
	})
}

func TestFetchTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userObjectID := primitive.NewObjectID()
		userID := userObjectID.Hex()

		mockTasks := []domain.Task{
			{
				ID:     primitive.NewObjectID(),
				Title:  "Test Task",
				UserID: userObjectID,
			},
			{
				ID:     primitive.NewObjectID(),
				Title:  "Test Task2",
				UserID: userObjectID,
			},
		}

		mockTaskUsecase := new(mocks.TaskUsecase)

		mockTaskUsecase.On("FetchByUserID", mock.Anything, userID).Return(mockTasks, nil)

		gin := gin.Default()

		rec := httptest.NewRecorder()

		tc := &controller.TaskController{
			TaskUsecase: mockTaskUsecase,
		}

		gin.Use(setUserID(userID))
		gin.GET("/task", tc.Fetch)

		body, err := json.Marshal(mockTasks)
		assert.NoError(t, err)

		bodyString := string(body)

		req := httptest.NewRequest(http.MethodGet, "/task", nil)
		gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Equal(t, bodyString, rec.Body.String())

		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		userObjectID := primitive.NewObjectID()
		userID := userObjectID.Hex()

		mockTaskUsecase := new(mocks.TaskUsecase)

		customErr := errors.New("Unexpected")

		mockTaskUsecase.On("FetchByUserID", mock.Anything, userID).Return(nil, customErr)

		gin := gin.Default()

		rec := httptest.NewRecorder()

		tc := &controller.TaskController{
			TaskUsecase: mockTaskUsecase,
		}

		gin.Use(setUserID(userID))
		gin.GET("/task", tc.Fetch)

		body, err := json.Marshal(domain.ErrorResponse{Message: customErr.Error()})
		assert.NoError(t, err)

		bodyString := string(body)

		req := httptest.NewRequest(http.MethodGet, "/task", nil)
		gin.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		assert.Equal(t, bodyString, rec.Body.String())

		mockTaskUsecase.AssertExpectations(t)
	})
}
