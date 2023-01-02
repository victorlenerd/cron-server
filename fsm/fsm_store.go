package fsm

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"os"
	"scheduler0/constants"
	"scheduler0/db"
	"scheduler0/models"
	"scheduler0/protobuffs"
	"scheduler0/utils"
	"sync"
	"time"
)

type Store struct {
	rwMtx                   sync.RWMutex
	DataStore               *db.DataStore
	logger                  *log.Logger
	Raft                    *raft.Raft
	QueueJobsChannel        chan []interface{}
	JobExecutionLogsChannel chan models.CommitJobStateLog
	StopAllJobs             chan bool

	raft.BatchingFSM
}

type Response struct {
	Data  []interface{}
	Error string
}

const (
	JobQueuesTableName        = "job_queues"
	JobQueueIdColumn          = "id"
	JobQueueNodeIdColumn      = "node_id"
	JobQueueLowerBoundJobId   = "lower_bound_job_id"
	JobQueueUpperBound        = "upper_bound_job_id"
	JobQueueDateCreatedColumn = "date_created"
	JobQueueVersion           = "version"

	ExecutionsUnCommittedTableName    = "job_executions_uncommitted"
	ExecutionsTableName               = "job_executions_committed"
	ExecutionsIdColumn                = "id"
	ExecutionsUniqueIdColumn          = "unique_id"
	ExecutionsStateColumn             = "state"
	ExecutionsNodeIdColumn            = "node_id"
	ExecutionsLastExecutionTimeColumn = "last_execution_time"
	ExecutionsNextExecutionTime       = "next_execution_time"
	ExecutionsJobIdColumn             = "job_id"
	ExecutionsDateCreatedColumn       = "date_created"
	ExecutionsJobQueueVersion         = "job_queue_version"
	ExecutionsVersion                 = "execution_version"

	JobQueuesVersionTableName     = "job_queue_versions"
	JobNumberOfActiveNodesVersion = "number_of_active_nodes"
)

var _ raft.FSM = &Store{}

func NewFSMStore(db *db.DataStore, logger *log.Logger) *Store {
	return &Store{
		DataStore:               db,
		QueueJobsChannel:        make(chan []interface{}, 1),
		JobExecutionLogsChannel: make(chan models.CommitJobStateLog, 1),
		StopAllJobs:             make(chan bool, 1),
		logger:                  logger,
	}
}

func (s *Store) Apply(l *raft.Log) interface{} {
	s.rwMtx.Lock()
	defer s.rwMtx.Unlock()

	return ApplyCommand(
		s.logger,
		l,
		s.DataStore,
		true,
		s.QueueJobsChannel,
		s.StopAllJobs,
	)
}

func (s *Store) ApplyBatch(logs []*raft.Log) []interface{} {
	s.rwMtx.Lock()
	defer s.rwMtx.Unlock()

	results := []interface{}{}

	for _, l := range logs {
		result := ApplyCommand(
			s.logger,
			l,
			s.DataStore,
			true,
			s.QueueJobsChannel,
			s.StopAllJobs,
		)
		results = append(results, result)
	}

	return results
}

func ApplyCommand(
	logger *log.Logger,
	l *raft.Log,
	db *db.DataStore,
	useQueues bool,
	queue chan []interface{},
	stopAllJobsQueue chan bool) interface{} {

	logPrefix := logger.Prefix()
	logger.SetPrefix(fmt.Sprintf("%s[apply-raft-command] ", logPrefix))
	defer logger.SetPrefix(logPrefix)

	if l.Type == raft.LogConfiguration {
		return nil
	}

	command := &protobuffs.Command{}

	marsherErr := proto.Unmarshal(l.Data, command)
	if marsherErr != nil {
		logger.Fatal("failed to unmarshal command", marsherErr.Error())
	}
	switch command.Type {
	case protobuffs.Command_Type(constants.CommandTypeDbExecute):
		params := []interface{}{}
		err := json.Unmarshal(command.Data, &params)
		if err != nil {
			return Response{
				Data:  nil,
				Error: err.Error(),
			}
		}
		ctx := context.Background()

		db.ConnectionLock.Lock()

		tx, err := db.Connection.BeginTx(ctx, nil)
		if err != nil {
			logger.Println("failed to execute sql command", err.Error())
			return Response{
				Data:  nil,
				Error: err.Error(),
			}
		}

		if utils.MonitorMemoryUsage(logger) {
			return Response{
				Data:  nil,
				Error: "out of memory",
			}
		}

		exec, err := tx.Exec(command.Sql, params...)
		if err != nil {
			logger.Println("failed to execute sql command", err.Error())
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				return Response{
					Data:  nil,
					Error: err.Error(),
				}
			}
		}

		err = tx.Commit()
		if err != nil {
			return Response{
				Data:  nil,
				Error: err.Error(),
			}
		}
		db.ConnectionLock.Unlock()

		lastInsertedId, err := exec.LastInsertId()
		if err != nil {
			logger.Println("failed to get last ", err.Error())
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				return Response{
					Data:  nil,
					Error: rollBackErr.Error(),
				}
			}
			return Response{
				Data:  nil,
				Error: err.Error(),
			}
		}
		rowsAffected, err := exec.RowsAffected()
		if err != nil {
			logger.Println(err.Error())
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				return Response{
					Data:  nil,
					Error: rollBackErr.Error(),
				}
			}
			return Response{
				Data:  nil,
				Error: err.Error(),
			}
		}
		data := []interface{}{lastInsertedId, rowsAffected}

		return Response{
			Data:  data,
			Error: "",
		}
	case protobuffs.Command_Type(constants.CommandTypeJobQueue):
		jobIds := []interface{}{}
		err := json.Unmarshal(command.Data, &jobIds)
		if err != nil {
			return Response{
				Data:  nil,
				Error: err.Error(),
			}
		}
		lowerBound := jobIds[0].(float64)
		upperBound := jobIds[1].(float64)
		lastVersion := jobIds[2].(float64)

		serverNodeId, err := utils.GetNodeIdWithRaftAddress(logger, raft.ServerAddress(command.ActionTarget))

		db.ConnectionLock.Lock()

		insertBuilder := sq.Insert(JobQueuesTableName).Columns(
			JobQueueNodeIdColumn,
			JobQueueLowerBoundJobId,
			JobQueueUpperBound,
			JobQueueVersion,
			JobQueueDateCreatedColumn,
		).Values(
			serverNodeId,
			lowerBound,
			upperBound,
			lastVersion,
			time.Now().UTC(),
		).RunWith(db.Connection)

		_, err = insertBuilder.Exec()
		if err != nil {
			logger.Fatalln("failed to insert new job queues", err.Error())
		}
		db.ConnectionLock.Unlock()
		if useQueues {
			queue <- []interface{}{command.Sql, int64(lowerBound), int64(upperBound)}
		}
		break
	case protobuffs.Command_Type(constants.CommandTypeJobExecutionLogs):
		jobState := models.CommitJobStateLog{}
		err := json.Unmarshal(command.Data, &jobState)
		if err != nil {
			return Response{
				Data:  nil,
				Error: err.Error(),
			}
		}

		if len(jobState.Logs) < 1 {
			return Response{
				Data: nil,
			}
		}

		db.ConnectionLock.Lock()

		batchInsert := func(jobExecutionLogs []models.JobExecutionLog) {
			query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES ",
				ExecutionsTableName,
				ExecutionsUniqueIdColumn,
				ExecutionsStateColumn,
				ExecutionsNodeIdColumn,
				ExecutionsLastExecutionTimeColumn,
				ExecutionsNextExecutionTime,
				ExecutionsJobIdColumn,
				ExecutionsDateCreatedColumn,
				ExecutionsJobQueueVersion,
				ExecutionsVersion,
			)

			query += "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
			params := []interface{}{
				jobExecutionLogs[0].UniqueId,
				jobExecutionLogs[0].State,
				jobExecutionLogs[0].NodeId,
				jobExecutionLogs[0].LastExecutionDatetime,
				jobExecutionLogs[0].NextExecutionDatetime,
				jobExecutionLogs[0].JobId,
				jobExecutionLogs[0].DataCreated,
				jobExecutionLogs[0].JobQueueVersion,
				jobExecutionLogs[0].ExecutionVersion,
			}

			for _, executionLog := range jobExecutionLogs[1:] {
				params = append(params,
					executionLog.UniqueId,
					executionLog.State,
					executionLog.NodeId,
					executionLog.LastExecutionDatetime,
					executionLog.NextExecutionDatetime,
					executionLog.JobId,
					executionLog.DataCreated,
					executionLog.JobQueueVersion,
					executionLog.ExecutionVersion,
				)
				query += ",(?, ?, ?, ?, ?, ?, ?, ?, ?)"
			}

			query += ";"

			ctx := context.Background()
			tx, err := db.Connection.BeginTx(ctx, nil)
			if err != nil {
				logger.Fatalln("failed to create transaction for batch insertion", err)
			}

			_, err = tx.Exec(query, params...)
			if err != nil {
				trxErr := tx.Rollback()
				if trxErr != nil {
					logger.Fatalln("failed to rollback update transition", trxErr)
				}
				logger.Fatalln("failed to update committed status of executions", err)
			}
			err = tx.Commit()
			if err != nil {
				logger.Fatalln("failed to commit transition", err)
			}
		}

		batchDelete := func(jobExecutionLogs []models.JobExecutionLog) {
			paramPlaceholder := "?"
			params := []interface{}{
				jobExecutionLogs[0].UniqueId,
			}

			for _, jobExecutionLog := range jobExecutionLogs[1:] {
				paramPlaceholder += ",?"
				params = append(params, jobExecutionLog.UniqueId)
			}

			ctx := context.Background()
			tx, err := db.Connection.BeginTx(ctx, nil)
			if err != nil {
				logger.Fatalln("failed to create transaction for batch insertion", err)
			}
			query := fmt.Sprintf("DELETE FROM %s WHERE %s IN (%s)", ExecutionsUnCommittedTableName, ExecutionsUniqueIdColumn, paramPlaceholder)
			_, err = tx.Exec(query, params...)
			if err != nil {
				trxErr := tx.Rollback()
				if trxErr != nil {
					logger.Fatalln("failed to rollback update transition", trxErr)
				}
				logger.Fatalln("failed to update committed status of executions", err)
			}
			err = tx.Commit()
			if err != nil {
				logger.Fatalln("failed to commit transition", err)
			}
		}

		batchInsert(jobState.Logs)
		batchDelete(jobState.Logs)

		db.ConnectionLock.Unlock()

		//logger.Println(fmt.Sprintf("received %v jobs execution logs from %s", len(jobState.Logs), jobState.Address))
		break
	case protobuffs.Command_Type(constants.CommandTypeStopJobs):
		if useQueues {
			stopAllJobsQueue <- true
		}
	}

	return nil
}

func (s *Store) Snapshot() (raft.FSMSnapshot, error) {
	logPrefix := s.logger.Prefix()
	s.logger.SetPrefix(fmt.Sprintf("%s[snapshot-fsm] ", logPrefix))
	defer s.logger.SetPrefix(logPrefix)
	fmsSnapshot := NewFSMSnapshot(s.DataStore)
	s.logger.Println("took snapshot")
	return fmsSnapshot, nil
}

func (s *Store) Restore(r io.ReadCloser) error {
	logPrefix := s.logger.Prefix()
	s.logger.SetPrefix(fmt.Sprintf("%s[restoring-snapshot] ", logPrefix))
	defer s.logger.SetPrefix(logPrefix)
	s.logger.Println("restoring snapshot")

	b, err := utils.BytesFromSnapshot(r)
	if err != nil {
		return fmt.Errorf("restore failed: %s", err.Error())
	}
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Fatal error getting working dir: %s \n", err)
	}
	dbFilePath := fmt.Sprintf("%v/%v", dir, constants.SqliteDbFileName)
	if err := os.Remove(dbFilePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if b != nil {
		if err := ioutil.WriteFile(dbFilePath, b, os.ModePerm); err != nil {
			return err
		}
	}

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return fmt.Errorf("restore failed to create db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("ping error: restore failed to create db: %v", err)
	}

	s.DataStore.ConnectionLock.Lock()
	s.DataStore.Connection = db
	s.DataStore.ConnectionLock.Unlock()

	return nil
}
