package util

import (
    "github.com/labstack/echo"
    "github.com/spf13/viper"
    "github.com/sirupsen/logrus"
    "net/http"
    "strings"
    "fmt"
    "runtime"
)

func SendToDiscord(message string) error {
    content := fmt.Sprintf(`{"content": "%s"}`, message)

    req, err := http.NewRequest("POST", viper.GetString("discord.url") , strings.NewReader(content))
    if err != nil {
        return fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("failed to send message to Discord: Status %s", resp.Status)
    }

    return nil
}

func SendToTeams(message string) error {
    content := fmt.Sprintf(`{"text": "%s"}`, message)

    req, err := http.NewRequest("POST", viper.GetString("teams.url"), strings.NewReader(content))
    if err != nil {
        return fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("failed to send message to Discord: Status %s", resp.Status)
    }

    return nil
}



func ErrorHandlingMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            defer func() {
                if err := recover(); err != nil {
                    pc, file, line, ok := runtime.Caller(1)
                    if !ok {
                        logger.Error("Failed to retrieve caller information")
                    }
                    funcName := runtime.FuncForPC(pc).Name()

                    logger.WithFields(logrus.Fields{
                        "func_name": funcName,
                        "file":      file,
                        "line":      line,
                        "error":     err,
                    }).Error("Recovered from panic")

                    c.JSON(http.StatusInternalServerError, map[string]string{
                        "message": "Internal Server Error...",
                    })
                }
            }()

            err := next(c)
            if err != nil {
                pc, file, line, ok := runtime.Caller(1)
                if !ok {
                    logger.Error("Failed to retrieve caller information")
                }
                funcName := runtime.FuncForPC(pc).Name()

                logger.WithFields(logrus.Fields{
                    "func_name": funcName,
                    "file":      file,
                    "line":      line,
                    "error":     err.Error(),
                }).Error("HTTP Error")

                // Return error message from handler
                return c.JSON(http.StatusInternalServerError, map[string]string{
                    "message": err.Error(),
                })
            }
            return nil
        }
    }
}
