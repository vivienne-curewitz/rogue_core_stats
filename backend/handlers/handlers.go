package handlers

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	saveparser "github.com/vivienne-curewitz/rogue_core_stats/saveParser"
	"github.com/vivienne-curewitz/rogue_core_stats/types"
)

func StartHandlers(wg *sync.WaitGroup) {
	dataPipe := make(chan types.SaveDataTask, 100)
	ctx := context.Background()
	statuses := make(map[uuid.UUID]types.UploadStatus)
	// start save parser thread
	go saveparser.SaveDataPipe(dataPipe, statuses, ctx)
	// stats handle funcs here
	http.HandleFunc("/saveUpload", UploadHandlerFactory(statuses, dataPipe))
	http.HandleFunc("/saveUploadStatus", UploadStatusHandlerFactory(statuses))
	http.HandleFunc("/overview", OverviewHandlerFactory())
	http.ListenAndServe("0.0.0.0:8080", nil)
	log.Println("Server shutting Down")
	wg.Done()
}
