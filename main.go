package main

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type config struct {
	Timeout       int     `env:"ENV_TIMEOUT" envDefault:"5"`
	URL           url.URL `env:"ENV_URL"`
	VmConnand     string  `env:"ENV_COMMAND" envDefault:"reboot"`
	VmID          string  `env:"ENV_VMID"`
	VmCommandPath string  `env:"ENV_COMMAND_PATH"`
	RetryWaitTime int     `env:"ENV_WAITTIME"`
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	log.Info("Starting programm")
	defer log.Info("Exit programm")
	// Loading the environment variables from '.env' file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("unable to load .env file: %e", slog.String("error", err.Error()))
	}
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Error("faild to load config", slog.String("error", err.Error()))
	}
	cfg.URL.Scheme = "https"
	log.Info("Starting check with", slog.String("url", cfg.URL.String()), slog.Int("timeout", cfg.Timeout), slog.Int("loop time", cfg.RetryWaitTime))
	retry := 0
	for {
		connection, err := checkConnection(&cfg, log)
		if err != nil {
			if retry >= 3 {
				break
			}
			retry++
		}
		if err == nil {
			if !connection {
				retry++
			} else {
				break
			}
		}
		log.Info("retry", slog.Int("retryValue:", retry))
		if retry == 3 {
			break
		}
		time.Sleep(time.Duration(cfg.RetryWaitTime) * time.Second)
	}
	if retry >= 3 {
		cmd := exec.Command(cfg.VmCommandPath, cfg.VmConnand, cfg.VmID)
		log.Debug("console command used", slog.String("command:", cmd.String()))
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Error("Error", slog.Any("error:", err))
			return
		}
	}
}

func checkConnection(cfg *config, log *slog.Logger) (bool, error) {
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout * int(time.Second)),
	}

	resp, err := client.Get(cfg.URL.String())
	if err != nil {
		log.Error("Error:", slog.Any("error:", err))
		return false, err
	}
	log.Info("Code:", slog.Int("code:", resp.StatusCode))
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	return true, nil
}
