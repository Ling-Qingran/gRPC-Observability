// main.go
package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/Ling-Qingran/gRPC-Observability/user"
	"google.golang.org/grpc"
)

type userServiceServer struct {
	user.UnimplementedUserServiceServer
}

var userInputs = []user.User{
	{Name: "Jack", Age: 21, CommuteMethod: "Bike", College: "Boston University", Hobbies: "Golf"},
	{Name: "David", Age: 21, CommuteMethod: "Bike", College: "Boston University", Hobbies: "Golf"},
	{Name: "Austin", Age: 21, CommuteMethod: "Bike", College: "Boston University", Hobbies: "Golf"},
}

func (*userServiceServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	name := req.GetName()
	userInput, err := getUserByName(name)

	if err != nil {
		return nil, err
	}

	return userInput, nil
}

func (*userServiceServer) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	name := req.GetName()
	index, err := getUserByNameWithIndex(name)

	if err != nil {
		return nil, err
	}

	updatedUser := req.GetUser()
	// Update the user record with the new data
	userInputs[index] = *updatedUser

	return updatedUser, nil
}

func (*userServiceServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	name := req.GetName()
	index, err := getUserByNameWithIndex(name)

	if err != nil {
		return &user.DeleteUserResponse{Success: false}, err
	}

	// Remove the user from the slice
	userInputs = append(userInputs[:index], userInputs[index+1:]...)

	return &user.DeleteUserResponse{Success: true}, nil
}

func (*userServiceServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	newUser := req.GetUser()
	userInputs = append(userInputs, *newUser)

	return newUser, nil
}

func getUserByName(name string) (*user.User, error) {
	for i, t := range userInputs {
		if t.Name == name {
			return &userInputs[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func getUserByNameWithIndex(name string) (int, error) {
	for i, t := range userInputs {
		if t.Name == name {
			return i, nil
		}
	}
	return -1, errors.New("user not found")
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &userServiceServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
