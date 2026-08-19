package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LYH2263/go-apigatex"
	"github.com/LYH2263/go-apigatex/internal/api"
)

func main() {
	addr := flag.String("addr", ":8094", "HTTP 监听地址")
	web := flag.String("web", "web", "静态管理页目录")
	persist := flag.String("persist", "", "路由快照 JSON 路径（可选）")
	token := flag.String("token", "dev-token", "默认 Bearer token")
	timeout := flag.Duration("proxy-timeout", 15*time.Second, "上游超时")
	flag.Parse()

	opts := []apigatex.Option{
		apigatex.WithDefaultToken(*token),
		apigatex.WithProxyTimeout(*timeout),
		apigatex.WithAllowInsecureUpstream(true),
	}
	if *persist != "" {
		opts = append(opts, apigatex.WithPersistPath(*persist))
	}

	g := apigatex.New(opts...)
	defer g.Close()
	if *persist != "" {
		if err := g.LoadPersist(); err != nil {
			log.Printf("load persist: %v", err)
		}
	}

	srv := api.New(g, api.Options{WebDir: *web, AllowCORS: true})
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("gated 监听 %s", *addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	_ = httpSrv.Close()
}
