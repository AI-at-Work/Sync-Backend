package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func (dataBase *Database) AddSession(ctx context.Context, userId string, sessionId string, modelId string, sessionName string) error {
	var query string
	var err error

	// Start a transaction
	tx, err := dataBase.Db.BeginTxx(ctx, nil) // Notice the use of BeginTxx for better context support
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare the SQL query using named parameters
	query = `INSERT INTO Session_Details (Session_Id, User_Id, Model_Id, Session_Name) VALUES (:session_id, :user_id, :model_id, :session_name)`
	params := map[string]interface{}{
		"session_id":   sessionId,
		"user_id":      userId,
		"model_id":     modelId,
		"session_name": sessionName,
	}

	// Execute the query
	result, err := tx.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	// Check how many rows were affected
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return errors.New("no rows were affected, possible invalid user_id or session_id")
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (dataBase *Database) AddChat(ctx context.Context, sessionId string, prompt string, chats string, chatsSummary string, isNew string) error {
	var query string
	var rows sql.Result
	var err error = nil

	// Start a transaction
	tx, err := dataBase.Db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.New("unable to begin transaction")
	}

	// Use defer to roll back transaction if anything goes wrong before commit.
	defer func() {
		if err != nil {
			log.Println("Doing RollBack : ", err)
			tx.Rollback()
		}
	}()

	chatJson := chats

	log.Println("IS NEW : ", isNew)
	log.Println("Session Id : ", sessionId)
	log.Println("chatsSummary : ", chatsSummary)
	log.Println("Vector : ", chatJson)

	if isNew == "new" {
		// Inserting new chat details
		query = `INSERT INTO Chat_Details (Session_Id, Session_Prompt, Chats, Chats_Summary) VALUES ($1, $2, $3::JSONB, $4)`
		rows, err = tx.ExecContext(ctx, query, sessionId, prompt, chatJson, chatsSummary)
	} else {
		// Updating existing chat details
		query = `UPDATE Chat_Details SET Chats = $2::JSONB, Chats_Summary = $3 WHERE Session_Id = $1`
		rows, err = tx.ExecContext(ctx, query, sessionId, chatJson, chatsSummary)
	}

	fmt.Println(chats)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error executing query: %v", err)
	}

	// Check if the operation affected any rows
	affected, err := rows.RowsAffected()
	if err != nil {
		return errors.New("unable to get affected rows")
	}
	if affected == 0 {
		return errors.New("no rows were affected, check session ID")
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return errors.New("unable to commit the transaction")
	}

	return nil
}
