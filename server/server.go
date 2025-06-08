package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"syscall"

	"github.com/bhaski-1234/redis-db/config"
	"github.com/bhaski-1234/redis-db/internal/processor"
	"github.com/bhaski-1234/redis-db/protocol"
	diskstorage "github.com/bhaski-1234/redis-db/storage/diskStorage"
	"golang.org/x/sys/unix"
)

type Server struct {
	epollFd     int
	listener    net.Listener
	connections map[int]net.Conn
	mu          sync.RWMutex
	diskstorage *diskstorage.DiskStorage
}

func NewServer() *Server {
	return &Server{
		connections: make(map[int]net.Conn),
		diskstorage: diskstorage.NewDiskStorage(),
	}
}

func (s *Server) Start() error {
	// Create listener
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	fmt.Printf("Server is running on %s:%d\n", config.Host, config.Port)

	// Load data from disk
	err = s.diskstorage.Load("dump")
	if errors.Is(err, os.ErrNotExist) {
		// Handle the "file not found" case specifically
	}

	// Create epoll instance
	s.epollFd, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		s.listener.Close()
		return fmt.Errorf("failed to create epoll: %w", err)
	}

	// Add listener to epoll
	if err := s.addListenerToEpoll(); err != nil {
		s.Close()
		return err
	}

	return s.eventLoop()
}

func (s *Server) addListenerToEpoll() error {
	file, err := s.listener.(*net.TCPListener).File()
	if err != nil {
		return fmt.Errorf("failed to get listener file: %w", err)
	}
	defer file.Close()

	fd := int(file.Fd())
	event := unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(fd),
	}

	if err := unix.EpollCtl(s.epollFd, unix.EPOLL_CTL_ADD, fd, &event); err != nil {
		return fmt.Errorf("failed to add listener to epoll: %w", err)
	}

	// Store a special marker for the listener
	s.mu.Lock()
	s.connections[fd] = nil // nil indicates listener
	s.mu.Unlock()

	return nil
}

func (s *Server) eventLoop() error {
	events := make([]unix.EpollEvent, 100)

	for {
		n, err := unix.EpollWait(s.epollFd, events, -1)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			return fmt.Errorf("epoll wait failed: %w", err)
		}

		for i := 0; i < n; i++ {
			fd := int(events[i].Fd)

			s.mu.RLock()
			conn := s.connections[fd]
			s.mu.RUnlock()

			if conn == nil {
				// This is the listener
				s.handleNewConnection()
			} else {
				// This is a client connection
				s.handleClientData(fd, conn)
			}
		}
	}
}

func (s *Server) handleNewConnection() {
	conn, err := s.listener.Accept()
	if err != nil {
		fmt.Printf("Error accepting connection: %v\n", err)
		return
	}

	tcpConn := conn.(*net.TCPConn)
	file, err := tcpConn.File()
	if err != nil {
		fmt.Printf("Error getting connection file: %v\n", err)
		conn.Close()
		return
	}

	fd := int(file.Fd())
	fmt.Println(fd)
	defer file.Close() // Close the duplicate file descriptor

	// Set non-blocking mode
	if err := syscall.SetNonblock(fd, true); err != nil {
		fmt.Printf("Error setting non-blocking: %v\n", err)
		conn.Close()
		return
	}

	// Add to epoll
	event := unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLET,
		Fd:     int32(fd),
	}

	if err := unix.EpollCtl(s.epollFd, unix.EPOLL_CTL_ADD, fd, &event); err != nil {
		fmt.Printf("Error adding connection to epoll: %v\n", err)
		conn.Close()
		return
	}

	// Store connection
	s.mu.Lock()
	s.connections[fd] = conn
	s.mu.Unlock()

	fmt.Printf("New connection accepted: %v\n", conn.RemoteAddr())
}

func (s *Server) handleClientData(fd int, conn net.Conn) {
	buf := make([]byte, 4096)

	for {
		n, _ := conn.Read(buf)

		// Process data
		resp, _ := processor.Process(buf[:n])
		respEncoded := protocol.EncodeResponse(resp)

		// Echo back for now
		if _, err := conn.Write([]byte(respEncoded)); err != nil {
			s.removeConnection(fd, conn)
			return
		}
	}
}

func (s *Server) removeConnection(fd int, conn net.Conn) {
	// Remove from epoll
	unix.EpollCtl(s.epollFd, unix.EPOLL_CTL_DEL, fd, nil)

	// Remove from map
	s.mu.Lock()
	delete(s.connections, fd)
	s.mu.Unlock()

	// Close connection
	conn.Close()

	fmt.Printf("Connection closed: %v\n", conn.RemoteAddr())
}

func (s *Server) Close() {
	if s.epollFd > 0 {
		unix.Close(s.epollFd)
	}

	if s.listener != nil {
		s.listener.Close()
	}

	s.mu.Lock()
	for _, conn := range s.connections {
		if conn != nil {
			conn.Close()
		}
	}
	s.mu.Unlock()
}

func runTCPServer() {
	server := NewServer()
	if err := server.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
