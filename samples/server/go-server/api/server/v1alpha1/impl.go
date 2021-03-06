// Copyright 2017 The Grafeas Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// package v1alpha1 is an implementation of the v1alpha1 version of Grafeas.
package v1alpha1

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grafeas/grafeas/samples/server/go-server/api/server/name"
	"golang.org/x/net/context"
	server "github.com/grafeas/grafeas/server-go"
	pb "github.com/grafeas/grafeas/v1alpha1/proto"
	opspb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Grafeas is an implementation of the Grafeas API, which should be called by handler methods for verification of logic
// and storage.
type Grafeas struct {
	S server.Storager
}

// CreateNote validates that a note is valid and then creates a note in the backing datastore.
func (g *Grafeas) CreateNote(ctx context.Context, req *pb.CreateNoteRequest) (*pb.Note, error) {
	n := req.Note
	if req == nil {
		log.Print("Note must not be empty.")
		return nil, status.Error(codes.InvalidArgument, "Note must not be empty")
	}
	if n.Name == "" {
		log.Printf("Invalid note name: %v", n.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid note name")
	}
	// TODO: Validate that operation exists if it is specified when get methods are implmented
	return n, g.S.CreateNote(n)
}

// CreateOccurrence validates that a note is valid and then creates an occurrence in the backing datastore.
func (g *Grafeas) CreateOccurrence(ctx context.Context, req *pb.CreateOccurrenceRequest) (*pb.Occurrence, error) {
	o := req.Occurrence
	if req == nil {
		log.Print("Occurrence must not be empty.")
		return nil, status.Error(codes.InvalidArgument, "Occurrence must not be empty")
	}
	if o.Name == "" {
		log.Printf("Invalid occurrence name: %v", o.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid occurrence name")
	}
	if o.NoteName == "" {
		log.Print("No note is associated with this occurrence")
	}
	pID, nID, err := name.ParseNote(o.NoteName)
	if err != nil {
		log.Printf("Invalid note name: %v", o.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid note name")
	}
	if n, err := g.S.GetNote(pID, nID); n == nil || err != nil {
		log.Printf("Unable to getnote %v, err: %v", n, err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Note %v not found", o.NoteName))
	}
	// TODO: Validate that operation exists if it is specified
	return o, g.S.CreateOccurrence(o)
}

// CreateOperation validates that a note is valid and then creates an operation note in the backing datastore.
func (g *Grafeas) CreateOperation(ctx context.Context, req *pb.CreateOperationRequest) (*opspb.Operation, error) {
	o := req.Operation
	if o.Name == "" {
		log.Printf("Invalid operation name: %v", o.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid operation name")
	}
	return o, g.S.CreateOperation(o)
}

// DeleteOccurrence deletes an occurrence from the datastore.
func (g *Grafeas) DeleteOccurrence(ctx context.Context, req *pb.DeleteOccurrenceRequest) (*empty.Empty, error) {
	pID, oID, err := name.ParseOccurrence(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid occurrence name")
	}
	return &empty.Empty{}, g.S.DeleteOccurrence(pID, oID)
}

// DeleteNote deletes a note from the datastore.
func (g *Grafeas) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*empty.Empty, error) {
	pID, nID, err := name.ParseNote(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid note name")
	}
	// TODO: Check for occurrences tied to this note, and return an error if there are any before deletion.
	return &empty.Empty{}, g.S.DeleteNote(pID, nID)
}

// DeleteOperation deletes an operation from the datastore.
func (g *Grafeas) DeleteOperation(ctx context.Context, req *opspb.DeleteOperationRequest) (*empty.Empty, error) {
	pID, oID, err := name.ParseOperation(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid Operation name")
	}
	// TODO: Check for occurrences and notes tied to this operation, and return an error if there are any before deletion.
	return &empty.Empty{}, g.S.DeleteOperation(pID, oID)
}

// GetNote gets a note from the datastore.
func (g *Grafeas) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.Note, error) {
	pID, nID, err := name.ParseNote(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid Note name")
	}
	return g.S.GetNote(pID, nID)
}

// GetOccurrence gets a occurrence from the datastore.
func (g *Grafeas) GetOccurrence(ctx context.Context, req *pb.GetOccurrenceRequest) (*pb.Occurrence, error) {
	pID, oID, err := name.ParseOccurrence(req.Name)
	if err != nil {
		log.Print("Could note parse name %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Could note parse name")
	}
	return g.S.GetOccurrence(pID, oID)
}

// GetOperation gets a occurrence from the datastore.
func (g *Grafeas) GetOperation(ctx context.Context, req *opspb.GetOperationRequest) (*opspb.Operation, error) {
	pID, oID, err := name.ParseOperation(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid Operation name")
	}
	return g.S.GetOperation(pID, oID)
}

// GetOccurrenceNote gets a the note for the provided occurrence from the datastore.
func (g *Grafeas) GetOccurrenceNote(ctx context.Context, req *pb.GetOccurrenceNoteRequest) (*pb.Note, error) {
	pID, oID, err := name.ParseOccurrence(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid occurrence name")
	}
	o, gErr := g.S.GetOccurrence(pID, oID)
	if gErr != nil {
		return nil, gErr
	}
	npID, nID, err := name.ParseNote(o.NoteName)
	if err != nil {
		log.Printf("Invalid note name: %v", o.Name)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Invalid note name: %v", o.NoteName))
	}
	return g.S.GetNote(npID, nID)
}

func (g *Grafeas) UpdateNote(ctx context.Context, req *pb.UpdateNoteRequest) (*pb.Note, error) {
	pID, nID, err := name.ParseNote(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid Note name")
	}
	// get existing note
	existing, gErr := g.S.GetNote(pID, nID)
	if gErr != nil {
		return nil, err
	}
	// verify that name didnt change
	if req.Note.Name != existing.Name {
		log.Printf("Cannot change note name: %v", req.Note.Name)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot change note name: %v", req.Note.Name))
	}

	// update note
	if gErr = g.S.UpdateNote(pID, nID, req.Note); err != nil {
		log.Printf("Cannot update note : %v", gErr)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot change note name: %v", req.Note.Name))
	}
	return g.S.GetNote(pID, nID)
}

func (g *Grafeas) UpdateOccurrence(ctx context.Context, req *pb.UpdateOccurrenceRequest) (*pb.Occurrence, error) {
	pID, oID, err := name.ParseOccurrence(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid occurrence name")
	}
	// get existing Occurrence
	existing, gErr := g.S.GetOccurrence(pID, oID)
	if gErr != nil {
		return nil, gErr
	}

	// verify that name didnt change
	if req.Name != existing.Name {
		log.Printf("Cannot change occurrence name: %v", req.Occurrence.Name)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot change occurrence name: %v", req.Occurrence.Name))
	}
	// verify that if note name changed, it still exists
	if req.Occurrence.NoteName != existing.NoteName {
		npID, nID, err := name.ParseNote(req.Occurrence.NoteName)
		if err != nil {
			return nil, err
		}
		if newN, err := g.S.GetNote(npID, nID); newN == nil || err != nil {
			return nil, err
		}
	}

	// update Occurrence
	if gErr = g.S.UpdateOccurrence(pID, oID, req.Occurrence); gErr != nil {
		log.Printf("Cannot update occurrence : %v", req.Occurrence.Name)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot update Occurrences: %v", err))
	}
	return g.S.GetOccurrence(pID, oID)
}

func (g *Grafeas) UpdateOperation(ctx context.Context, req *pb.UpdateOperationRequest) (*opspb.Operation, error) {
	pID, oID, err := name.ParseOperation(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid Operation name")
	}
	// get existing operation
	existing, gErr := g.S.GetOperation(pID, oID)
	if gErr != nil {
		return nil, gErr
	}

	// verify that operation isn't marked done
	if req.Operation.Done != existing.Done && existing.Done {
		log.Printf("Trying to update a done operation")
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot update operation in status done: %v", req.Name))
	}

	// verify that name didnt change
	if req.Operation.Name != existing.Name {
		log.Printf("Cannot change operation name: %v", req.Operation.Name)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot change operation name: %v", req.Name))
	}

	// update operation
	if gErr = g.S.UpdateOperation(pID, oID, req.Operation); gErr != nil {
		log.Printf("Cannot update operation : %v", req.Operation.Name)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot update Opreation: %v"))
	}
	return g.S.GetOperation(pID, oID)
}
func (g *Grafeas) ListOperations(ctx context.Context, req *opspb.ListOperationsRequest) (*opspb.ListOperationsResponse, error) {
	pID, err := name.ParseProject(req.Name)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid Project name")
	}
	// TODO: support filters
	ops := g.S.ListOperations(pID, req.Filter)
	return &opspb.ListOperationsResponse{Operations: ops}, nil
}

func (g *Grafeas) ListNotes(ctx context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	pID, err := name.ParseProject(req.Parent)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Parent)
		return nil, status.Error(codes.InvalidArgument, "Invalid Project name")
	}

	// TODO: support filters
	ns := g.S.ListNotes(pID, req.Filter)
	return &pb.ListNotesResponse{Notes: ns}, nil

}

func (g *Grafeas) ListOccurrences(ctx context.Context, req *pb.ListOccurrencesRequest) (*pb.ListOccurrencesResponse, error) {
	pID, err := name.ParseProject(req.Parent)
	if err != nil {
		log.Printf("Error parsing name: %v", req.Parent)
		return nil, err
	}

	// TODO: support filters - prioritizing resource url
	os := g.S.ListOccurrences(pID, req.Filter)
	return &pb.ListOccurrencesResponse{Occurrences: os}, nil
}

func (g *Grafeas) ListNoteOccurrences(ctx context.Context, req *pb.ListNoteOccurrencesRequest) (*pb.ListNoteOccurrencesResponse, error) {
	pID, nID, err := name.ParseNote(req.Name)
	if err != nil {
		log.Printf("Invalid note name: %v", req.Name)
		return nil, status.Error(codes.InvalidArgument, "Invalid note name")
	}
	// TODO: support filters - prioritizing resource url
	os, gErr := g.S.ListNoteOccurrences(pID, nID, req.Filter)
	if gErr != nil {
		return nil, gErr
	}
	return &pb.ListNoteOccurrencesResponse{Occurrences: os}, nil
}

func (g *Grafeas) CancelOperation(context.Context, *opspb.CancelOperationRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "Currently Unimplemented")
}
