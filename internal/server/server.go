package flare_server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	FlareModel "github.com/soulteary/flare/config/model"
	FlareState "github.com/soulteary/flare/config/state"
	FlareLogger "github.com/soulteary/flare/internal/logger"

	FlareAuth "github.com/soulteary/flare/internal/auth"
	FlareIcons "github.com/soulteary/flare/internal/icons"
	FlareAssets "github.com/soulteary/flare/internal/resources/assets"
	FlareTemplates "github.com/soulteary/flare/internal/resources/templates"

	FlareDeprecated "github.com/soulteary/flare/internal/misc/deprecated"
	FlareHealth "github.com/soulteary/flare/internal/misc/health"
	FlareRedir "github.com/soulteary/flare/internal/misc/redir"
	FlareEditor "github.com/soulteary/flare/internal/pages/editor"
	FlareGuide "github.com/soulteary/flare/internal/pages/guide"
	FlareHome "github.com/soulteary/flare/internal/pages/home"
	FlareSettings "github.com/soulteary/flare/internal/settings"
	FlareAppearance "github.com/soulteary/flare/internal/settings/appearance"
	FlareOthers "github.com/soulteary/flare/internal/settings/others"
	FlareSearch "github.com/soulteary/flare/internal/settings/search"
	FlareTheme "github.com/soulteary/flare/internal/settings/theme"
	FlareWeather "github.com/soulteary/flare/internal/settings/weather"
)

func StartDaemon(AppFlags *FlareModel.Flags) {

	if !AppFlags.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := FlareLogger.GetLogger(AppFlags.LogLevel)
	router := gin.New()
	router.Use(gin.Recovery())
	if AppFlags.DebugMode {
		router.Use(FlareLogger.New(log))
	}

	if !AppFlags.DisableLoginMode {
		FlareAuth.RequestHandle(router)
	}

	if !AppFlags.DebugMode {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	FlareState.Init()
	FlareAssets.RegisterRouting(router)

	FlareIcons.Init()
	FlareIcons.RegisterRouting(router)

	FlareTemplates.RegisterRouting(router)
	FlareAppearance.RegisterRouting(router)
	FlareDeprecated.RegisterRouting(router)
	FlareHealth.RegisterRouting(router)
	FlareWeather.RegisterRouting(router)
	FlareHome.RegisterRouting(router)
	FlareOthers.RegisterRouting(router)
	FlareRedir.RegisterRouting(router)
	FlareSearch.RegisterRouting(router)
	FlareSettings.RegisterRouting(router)
	FlareTheme.RegisterRouting(router)

	if AppFlags.EnableEditor {
		FlareEditor.RegisterRouting(router)
		log.Info("在线编辑模块启用，可以访问 " + FlareState.RegularPages.Editor.Path + " 来进行数据编辑。")
	}

	if AppFlags.EnableGuide {
		FlareGuide.RegisterRouting(router)
		log.Info("向导模块启用，可以访问 " + FlareState.RegularPages.Guide.Path + " 来获取程序使用帮助。")
	}

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(AppFlags.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("程序启动出错: %s\n", err)
			os.Exit(1)
		}
	}()
	log.Info("程序已启动完毕 🚀")

	<-ctx.Done()

	stop()
	log.Info("程序正在关闭中，如需立即结束请按 CTRL+C")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("程序强制关闭: ", slog.Any("error", err))
	}

	log.Info("期待与你的再次相遇 ❤️")
}
