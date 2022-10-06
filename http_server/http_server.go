package http_server

import (
	"database/sql"
	"fmt"
	"github.com/go-http-utils/logger"
	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"scheduler0/cluster"
	"scheduler0/constants"
	"scheduler0/db"
	"scheduler0/fsm"
	controllers2 "scheduler0/http_server/controllers"
	"scheduler0/http_server/middlewares"
	"scheduler0/process"
	repository2 "scheduler0/repository"
	service2 "scheduler0/service"
	"scheduler0/utils"
)

// Start this will start the http server
func Start() {
	ctx := context.Background()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(fmt.Errorf("Fatal error getting working dir: %s \n", err))
	}
	dbFilePath := fmt.Sprintf("%v/%v", dir, constants.SqliteDbFileName)

	sqliteDb := db.NewSqliteDbConnection(dbFilePath)
	conn, err := sqliteDb.OpenConnection()
	if err != nil {
		log.Fatal("Failed to open connection", err)
	}

	dbConnection := conn.(*sql.DB)
	configs := utils.GetScheduler0Configurations()

	err = dbConnection.Ping()
	if err != nil {
		log.Fatalln(fmt.Errorf("ping error: restore failed to create db: %v", err))
	}

	utils.MakeDirIfNotExist(constants.RaftDir)

	dirPath := fmt.Sprintf("%v/%v", constants.RaftDir, configs.NodeId)
	dirPath, exists := utils.MakeDirIfNotExist(dirPath)

	cls := cluster.NewCluster()
	fsmStr := fsm.NewFSMStore(sqliteDb, dbConnection)

	if exists {
		cluster.RecoverRaftStore(cls)
	}

	rft, rfErr := cluster.NewRaft(fsmStr, cls)
	fsmStr.Raft = rft
	if rfErr != nil {
		log.Fatal("failed to create raft store", rfErr)
	}

	if configs.Bootstrap == "true" && !exists {
		err = cluster.BootstrapNode(rft)
		if err != nil {
			log.Fatal("failed to bootstrap node:", err)
		}
	}

	//repository
	credentialRepo := repository2.NewCredentialRepo(fsmStr)
	jobRepo := repository2.NewJobRepo(fsmStr)
	executionRepo := repository2.NewExecutionRepo(fsmStr)
	projectRepo := repository2.NewProjectRepo(fsmStr, jobRepo)

	//services
	credentialService := service2.NewCredentialService(credentialRepo, ctx)
	jobService := service2.NewJobService(jobRepo, ctx)
	executionService := service2.NewExecutionService(executionRepo)
	projectService := service2.NewProjectService(projectRepo)

	// SetupDB logging
	log.SetFlags(0)
	log.SetOutput(new(utils.LogWriter))
	jobProcessor := process.NewJobProcessor(dbConnection, jobRepo, projectRepo)

	// StartJobs process to execute cron-server jobs
	go jobProcessor.StartJobs()

	// HTTP router setup
	router := mux.NewRouter()

	// Security middleware
	secureMiddleware := secure.New(secure.Options{FrameDeny: true})

	// Initialize controllers
	executionController := controllers2.NewExecutionsController(executionService)
	jobController := controllers2.NewJoBHTTPController(jobService, *jobProcessor)
	projectController := controllers2.NewProjectController(projectService)
	credentialController := controllers2.NewCredentialController(credentialService)
	healthCheckController := controllers2.NewHealthCheckController(rft)

	// Mount middleware
	middleware := middlewares.MiddlewareType{}

	router.Use(secureMiddleware.Handler)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middleware.ContextMiddleware)
	router.Use(middleware.AuthMiddleware(credentialService))

	// Executions Endpoint
	router.HandleFunc("/executions", executionController.ListExecutions).Methods(http.MethodGet)

	// Credentials Endpoint
	router.HandleFunc("/credentials", credentialController.CreateOneCredential).Methods(http.MethodPost)
	router.HandleFunc("/credentials", credentialController.ListCredentials).Methods(http.MethodGet)
	router.HandleFunc("/credentials/{uuid}", credentialController.GetOneCredential).Methods(http.MethodGet)
	router.HandleFunc("/credentials/{uuid}", credentialController.UpdateOneCredential).Methods(http.MethodPut)
	router.HandleFunc("/credentials/{uuid}", credentialController.DeleteOneCredential).Methods(http.MethodDelete)

	// Job Endpoint
	router.HandleFunc("/job", jobController.CreateOneJob).Methods(http.MethodPost)
	router.HandleFunc("/jobs", jobController.BatchCreateJobs).Methods(http.MethodPost)
	router.HandleFunc("/jobs", jobController.ListJobs).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{uuid}", jobController.GetOneJob).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{uuid}", jobController.UpdateOneJob).Methods(http.MethodPut)
	router.HandleFunc("/jobs/{uuid}", jobController.DeleteOneJob).Methods(http.MethodDelete)

	// Projects Endpoint
	router.HandleFunc("/projects", projectController.CreateOneProject).Methods(http.MethodPost)
	router.HandleFunc("/projects", projectController.ListProjects).Methods(http.MethodGet)
	router.HandleFunc("/projects/{uuid}", projectController.GetOneProject).Methods(http.MethodGet)
	router.HandleFunc("/projects/{uuid}", projectController.UpdateOneProject).Methods(http.MethodPut)
	router.HandleFunc("/projects/{uuid}", projectController.DeleteOneProject).Methods(http.MethodDelete)

	// Healthcheck Endpoint
	router.HandleFunc("/healthcheck", healthCheckController.HealthCheck).Methods(http.MethodGet)

	router.PathPrefix("/api-docs/").Handler(http.StripPrefix("/api-docs/", http.FileServer(http.Dir("./server/http_server/api-docs/"))))

	log.Println("Server is running on port", configs.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", configs.Port), logger.Handler(router, os.Stdout, logger.CombineLoggerType))
	utils.CheckErr(err)
}
