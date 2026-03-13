package examples

import "log/slog"

func examples() {
	slog.Info("Starting server")       // want `log message should start with a lowercase letter`
	slog.Info("запуск сервера")        // want `log message should be in English only: found: .з.`
	slog.Info("server started!")       // want `log message should not contain special character: found: .!.`
	slog.Info("user password exposed") // want `log message may expose sensitive data: found: .password.`

	slog.Info("starting server")    // OK
	slog.Info("server started")     // OK
	slog.Info("user authenticated") // OK
}
