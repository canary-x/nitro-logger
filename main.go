package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/mdlayher/vsock"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logFileName := flag.String("file", "/var/log/nitro-logger/logfile.log", "Log file name")
	maxSize := flag.Int("max-size", 50, "Max size of log file in megabytes")
	maxBackups := flag.Int("max-backups", 7, "Max number of backups")
	maxAge := flag.Int("max-age", 28, "Max age of log file in days")
	compress := flag.Bool("compress", true, "Compress rotated log files with gzip")
	port := flag.Int("port", 8090, "Port number")
	flag.Parse()

	if *maxSize <= 0 {
		log.Fatalln("Max size must be greater than 0")
	}
	if *maxBackups <= 0 {
		log.Fatalln("Max backups must be greater than 0")
	}
	if *maxAge <= 0 {
		log.Fatalln("Max age must be greater than 0")
	}
	if *port <= 0 {
		log.Fatalln("Port must be greater than 0")
	}

	logger := &lumberjack.Logger{
		Filename:   *logFileName,
		MaxSize:    *maxSize,
		MaxBackups: *maxBackups,
		MaxAge:     *maxAge,
		Compress:   *compress,
	}
	if err := run(logger, *port); err != nil {
		log.Fatalf("Fatal error: %+v\n", err)
	}
}

func run(logger *lumberjack.Logger, port int) error {
	ln, err := listen(port)
	if err != nil {
		return fmt.Errorf("listening on socket: %w", err)
	}
	defer ln.Close()

	log.Println("Listening on port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, logger)
	}
}

func handleConnection(conn net.Conn, logger *lumberjack.Logger) {
	log.Println("Nitro logger: received connection, will start streaming logs...")
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if _, err := logger.Write(scanner.Bytes()); err != nil {
			log.Printf("Nitro logger: error writing to log file: %v\n", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Nitro logger: error reading from connection: %v\n", err)
	}
	log.Println("Nitro logger: connection closed, waiting for further connections...")
}

func listen(port int) (net.Listener, error) {
	listenTCP := func(port int) (net.Listener, error) {
		return net.Listen("tcp", fmt.Sprintf(":%d", port))
	}
	contextID, err := vsock.ContextID()
	if err != nil {
		log.Printf("OS does not support vsock (error on getting CID: %v): falling back to regular TCP socket\n", err)
		return listenTCP(port)
	}

	ln, err := vsock.ListenContextID(contextID, uint32(port), nil)
	if err != nil && strings.Contains(err.Error(), "vsock: not implemented") {
		log.Println("OS does not support vsock: falling back to regular TCP socket")
		return listenTCP(port)
	}
	log.Println("Vsock connected on CID", contextID)
	return ln, err
}
