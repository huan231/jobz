package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/huan231/jobz/internal/cronjob"
	"github.com/huan231/jobz/internal/health"
	"github.com/huan231/jobz/internal/job"
	"github.com/huan231/jobz/pkg/events"
	_ "github.com/mattn/go-sqlite3"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func main() {
	migrations := flag.String("migrations", "", "path to the migrations")

	flag.Parse()

	db, err := sql.Open("sqlite3", os.Getenv("DATABASE_FILE_PATH"))

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+*migrations, os.Getenv("DATABASE_FILE_PATH"), driver)

	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))

	if err != nil {
		log.Fatal(err)
	}

	clientset := kubernetes.NewForConfigOrDie(config)
	factory := informers.NewSharedInformerFactory(clientset, 0)

	cronJobLister := factory.Batch().V1().CronJobs().Lister()
	jobLister := factory.Batch().V1().Jobs().Lister()

	cronJobInformer := factory.Batch().V1().CronJobs().Informer()
	eventInformer := factory.Core().V1().Events().Informer()

	stop := make(chan struct{})

	defer close(stop)

	go factory.Start(stop)

	for _, ok := range factory.WaitForCacheSync(stop) {
		if !ok {
			log.Fatal(fmt.Errorf("timed out waiting for caches to sync"))
		}
	}

	hub := events.NewHub()
	defer hub.Close()

	cronJobService := cronjob.NewService(cronJobLister)
	jobService := job.NewService(job.NewRepository(db))

	healthController := health.NewController()
	cronJobController := cronjob.NewController(cronJobService)
	jobController := job.NewController(jobService)
	eventsController := events.NewController(hub)

	cronJobEventHandler := cronjob.NewEventHandler(cronJobService, hub)
	jobEventHandler := job.NewEventHandler(jobService, hub, cronJobLister, jobLister)

	cronJobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			cronJobEventHandler.OnCronJobAdd(obj.(*batchv1.CronJob))
		},
		UpdateFunc: func(oldObj, newObj any) {
			cronJobEventHandler.OnCronJobUpdate(oldObj.(*batchv1.CronJob), newObj.(*batchv1.CronJob))
		},
		DeleteFunc: func(obj any) {
			jobEventHandler.OnCronJobDelete(obj.(*batchv1.CronJob))
			cronJobEventHandler.OnCronJobDelete(obj.(*batchv1.CronJob))
		},
	})
	eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			jobEventHandler.OnEventAdd(obj.(*corev1.Event))
		},
	})

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(cors.Default())

	router.GET("/livez", healthController.Liveness)
	router.GET("/cronjobs", cronJobController.List)
	router.GET("/jobs", jobController.List)
	router.GET("/events", eventsController.Stream)

	if err = router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
