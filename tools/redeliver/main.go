// Command redeliver menjalankan ulang task Asynq yang sudah masuk daftar
// archived (melewati batas retry). Dipakai saat incident: perbaiki penyebabnya,
// lalu jalankan task yang tertahan tanpa menunggu operator membuka dashboard.
//
// Contoh:
//
//	go run ./tools/redeliver -queue default -limit 50
//	go run ./tools/redeliver -queue default -limit 50 -yes
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hibiken/asynq"

	"github.com/vyerru/lapor-kampus/config"
)

func main() {
	queue := flag.String("queue", "default", "nama queue")
	limit := flag.Int("limit", 20, "jumlah task yang dijalankan ulang")
	yes := flag.Bool("yes", false, "lewati konfirmasi")
	flag.Parse()

	config.LoadEnv()
	inspector := asynq.NewInspector(config.RedisOpt())
	defer inspector.Close()

	tasks, err := inspector.ListArchivedTasks(*queue,
		asynq.PageSize(*limit), asynq.Page(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "gagal membaca archived task: %v\n", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("tidak ada task archived di queue", *queue)
		return
	}

	fmt.Printf("ditemukan %d task archived di queue %s:\n", len(tasks), *queue)
	for _, task := range tasks {
		fmt.Printf("  - %s id=%s gagal_terakhir=%s\n",
			task.Type, task.ID, task.LastFailedAt.Format(time.RFC3339))
	}

	if !*yes {
		fmt.Print("jalankan ulang semua task di atas? [y/N] ")
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil || (answer != "y" && answer != "Y") {
			fmt.Println("dibatalkan")
			return
		}
	}

	var succeeded, failed int
	for _, task := range tasks {
		if err := inspector.RunTask(*queue, task.ID); err != nil {
			fmt.Fprintf(os.Stderr, "gagal menjalankan %s: %v\n", task.ID, err)
			failed++
			continue
		}
		succeeded++
	}

	fmt.Printf("selesai: %d dijalankan ulang, %d gagal\n", succeeded, failed)
}
